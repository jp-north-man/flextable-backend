package main

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

func main() {
	models.CreateDB()
	http.HandleFunc("/create", createHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createHandler(w http.ResponseWriter, r *http.Request) {
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
