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

func CreateTable(tableName string, columns []Column) error {
	db := models.ConnectDB()
	defer db.Close()

	columnsJSON, err := json.Marshal(columns)
	if err != nil {
		return err
	}

	creator := "user" // 実際のユーザー名に置き換えてください。
	_, err = db.Exec("INSERT INTO table_definitions (name, creator, columns) VALUES ($1, $2, $3)", tableName, creator, columnsJSON)
	if err != nil {
		return err
	}

	return nil
}
