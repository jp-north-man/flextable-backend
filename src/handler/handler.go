package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"example.com/packages/models"
)

type Cell struct {
	ColumnNumber int         `json:"columnNumber"`
	ColumnName   string      `json:"columnName"`
	Value        interface{} `json:"value"`
}

type Row struct {
	RowNumber int    `json:"rowNumber"`
	Cells     []Cell `json:"cells"`
}

type Column struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Width string `json:"width"`
	Type  string `json:"type"`
}

type CreateRequest struct {
	Columns   []Column `json:"columns"`
	TableName string   `json:"tableName"`
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var req CreateRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Print(req)
		err = CreateTable(req.TableName, req.Columns) // modelsパッケージ内で実装されたCreateTable関数を呼び出す
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Table created successfully",
		})

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func CreateTable(tableName string, columns []Column) error {
	db := models.ConnectDB()
	defer db.Close()

	columnsJSON, err := json.Marshal(columns)
	if err != nil {
		log.Println("Error marshalling columns:", err)
		return err
	}

	creator := "user" // 実際のユーザー名に置き換えてください。
	var tableID int
	err = db.QueryRow("INSERT INTO test_table_definitions (name, creator, columns) VALUES ($1, $2, $3) RETURNING id", tableName, creator, columnsJSON).Scan(&tableID)
	if err != nil {
		log.Println("INSERT INTO test_table_definitions (name, creator, columns) VALUES ($1, $2, $3) RETURNING id:", err)
		return err
	}

	// 初期データを空の配列として作成します。
	initialData := make([]Row, 0)

	// 初期データをJSON形式に変換します。
	initialDataJSON, err := json.Marshal(initialData)
	if err != nil {
		log.Println("Error marshalling initialData:", err)
		return err
	}

	_, err = db.Exec(`
		INSERT INTO test_data_tables (table_id, data)
		VALUES ($1, $2)
		`, tableID, initialDataJSON)
	if err != nil {
		log.Println("INSERT INTO test_data_tables (table_id, data) VALUES ($1, $2):", err)
		return err
	}

	// 新しく作成されたテーブルのデータを取得します。
	rows, err := db.Query("SELECT * FROM test_table_definitions")
	if err != nil {
		log.Println("Error querying test_table_definitions:", err)
		return err
	}
	defer rows.Close()

	// データをログに出力します。
	for rows.Next() {
		var id int
		var name, creator, columns string

		err := rows.Scan(&id, &name, &creator, &columns)
		if err != nil {
			return err
		}

		log.Println("ID:", id, ", Name:", name, ", Creator:", creator, ", Columns:", columns)
	}

	// test_data_tables からデータを取得します。
	dataRows, err := db.Query("SELECT * FROM test_data_tables")
	if err != nil {
		log.Println("Error querying test_data_tables:", err)
		return err
	}
	defer dataRows.Close()

	// データをログに出力します。
	for dataRows.Next() {
		var id, table_id int
		var data string

		err := dataRows.Scan(&id, &table_id, &data)
		if err != nil {
			return err
		}

		log.Println("ID:", id, ", TableID:", table_id, ", Data:", data)
	}

	return nil
}

type AddRowRequest struct {
	TableID int `json:"table_id"`
}

func AddRowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var req AddRowRequest
		// リクエストボディから AddRowRequest をデコードします。
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Println(req)

		// AddRow 関数を呼び出して、新しい行をデータベースに追加します。
		err = AddRow(req.TableID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 成功メッセージを返します。
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Row added successfully",
		})

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// AddRow は、指定されたテーブルIDに新しい行を追加します。
func AddRow(tableID int) error {
	// データベースに接続します。
	db := models.ConnectDB()
	defer db.Close()

	// 指定されたテーブルIDのカラム情報を取得します。
	var columnsJSON string
	err := db.QueryRow("SELECT columns FROM test_table_definitions WHERE id = $1", tableID).Scan(&columnsJSON)
	if err != nil {
		return err
	}

	// JSONをColumnのスライスにデコードします。
	var columns []Column
	err = json.Unmarshal([]byte(columnsJSON), &columns)
	if err != nil {
		return err
	}

	// 新しい行を作成します。
	newRow := make(map[string]interface{})
	for _, column := range columns {
		newRow[column.Name] = ""
	}

	// 新しい行をJSONに変換します。
	newRowJSON, err := json.Marshal(newRow)
	if err != nil {
		return err
	}
	// データベースに新しい行を追加します。
	_, err = db.Exec(`
		INSERT INTO test_data_tables (table_id, data)
		VALUES ($1, $2)
	`, tableID, newRowJSON)
	if err != nil {
		return err
	}

	return nil
}

type GetFlexTablesResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetFlexTablesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		flexTables, err := GetFlexTables()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(flexTables)

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func GetFlexTables() ([]GetFlexTablesResponse, error) {
	db := models.ConnectDB()
	defer db.Close()

	rows, err := db.Query("SELECT id, name FROM test_table_definitions")
	if err != nil {
		log.Println("Error querying test_table_definitions:", err)
		return nil, err
	}
	defer rows.Close()

	flexTables := make([]GetFlexTablesResponse, 0)

	for rows.Next() {
		var id int
		var name string

		err := rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}

		flexTables = append(flexTables, GetFlexTablesResponse{
			ID:   id,
			Name: name,
		})
	}

	return flexTables, nil
}

type UpdateCellRequest struct {
	TableID      int         `json:"table_id"`
	RowNumber    int         `json:"row_number"`
	ColumnNumber int         `json:"column_number"`
	Value        interface{} `json:"value"`
}

func UpdateCellHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		var req UpdateCellRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Println(req)

		err = UpdateCell(req.TableID, req.RowNumber, req.ColumnNumber, req.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Cell updated successfully",
		})

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func UpdateCell(tableID int, rowNumber int, columnNumber int, value interface{}) error {
	db := models.ConnectDB()
	defer db.Close()

	var dataJSON string
	err := db.QueryRow("SELECT data FROM test_data_tables WHERE table_id = $1", tableID).Scan(&dataJSON)
	if err != nil {
		log.Println("Error querying test_data_tables:", err)
		return err
	}

	var rows []Row
	err = json.Unmarshal([]byte(dataJSON), &rows)
	if err != nil {
		log.Println("Error unmarshalling rows:", err)
		return err
	}

	rows[rowNumber].Cells[columnNumber].Value = value

	updatedDataJSON, err := json.Marshal(rows)
	if err != nil {
		log.Println("Error marshalling updated data:", err)
		return err
	}

	_, err = db.Exec(`
		UPDATE test_data_tables
		SET data = $1
		WHERE table_id = $2
	`, updatedDataJSON, tableID)
	if err != nil {
		log.Println("Error updating test_data_tables:", err)
		return err
	}

	return nil
}

type GetTableDataRequest struct {
	TableID int `json:"id"`
}

type GetTableDataResponse struct {
	TableName string   `json:"table_name"`
	Columns   []Column `json:"columns"`
	Data      []Row    `json:"data"`
}

func GetTableDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var req GetTableDataRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Println("table id", req.TableID)

		tableData, err := GetTableData(req.TableID)
		if err != nil {
			log.Printf("Error decoding request body 2: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tableData)

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func GetTableData(tableID int) (GetTableDataResponse, error) {
	db := models.ConnectDB()
	defer db.Close()

	var tableName, columnsJSON string
	err := db.QueryRow("SELECT name, columns FROM test_table_definitions WHERE id = $1", tableID).Scan(&tableName, &columnsJSON)
	if err != nil {
		log.Printf("Error decoding request body 3: %v", err)
		return GetTableDataResponse{}, err
	}

	var columns []Column
	err = json.Unmarshal([]byte(columnsJSON), &columns)
	if err != nil {
		log.Printf("Error decoding request body 4: %v", err)
		return GetTableDataResponse{}, err
	}

	var dataJSON string
	err = db.QueryRow("SELECT data FROM test_data_tables WHERE table_id = $1", tableID).Scan(&dataJSON)
	if err != nil {
		log.Printf("Error decoding request body 5: %v", err)
		return GetTableDataResponse{}, err
	}

	var data []Row
	err = json.Unmarshal([]byte(dataJSON), &data)
	if err != nil {
		log.Printf("Error decoding request body 6: %v", err)
		return GetTableDataResponse{}, err
	}

	return GetTableDataResponse{
		TableName: tableName,
		Columns:   columns,
		Data:      data,
	}, nil
}
