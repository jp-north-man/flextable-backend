package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"example.com/packages/models"
)

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
		log.Println(req)
		// err = CreateTable(req.TableName, req.Columns) // modelsパッケージ内で実装されたCreateTable関数を呼び出す
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Table created successfully",
		})

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func CreateTable(tableName string, columns []Column) (int, error) {
	db := models.ConnectDB()
	defer db.Close()

	columnsJSON, err := json.Marshal(columns)
	if err != nil {
		return 0, err
	}

	creator := "user" // 実際のユーザー名に置き換えてください。
	var tableID int
	err = db.QueryRow("INSERT INTO table_definitions (name, creator, columns) VALUES ($1, $2, $3) RETURNING id", tableName, creator, columnsJSON).Scan(&tableID)
	if err != nil {
		return 0, err
	}

	// columnsをもとに空のJSONオブジェクトを作成します。
	initialData := make(map[string]interface{})
	for _, column := range columns {
		initialData[column.Name] = ""
	}

	// 初期データをJSON形式に変換します。
	initialDataJSON, err := json.Marshal(initialData)
	if err != nil {
		return 0, err
	}

	_, err = db.Exec(`
		INSERT INTO data_tables (table_id, data)
		VALUES ($1, $2)
	`, tableID, initialDataJSON)
	if err != nil {
		return 0, err
	}

	return tableID, nil
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
	err := db.QueryRow("SELECT columns FROM table_definitions WHERE id = $1", tableID).Scan(&columnsJSON)
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
		INSERT INTO data_tables (table_id, data)
		VALUES ($1, $2)
	`, tableID, newRowJSON)
	if err != nil {
		return err
	}

	return nil
}
