package models

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	db, err := sql.Open("postgres",
		"host=localhost port=5432 user=user password=password dbname=mydb sslmode=disable") // 今回は直接書きます。Please read from .env etc.
	if err != nil {
		log.Fatal("sql.Open: ", err)
	} else {
		log.Println("sql.Open: ", "ok")
	}

	return db
}

func CreateDB() {
	db := ConnectDB()
	defer db.Close()

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS table_definitions (
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL,
            creator TEXT NOT NULL,
            columns JSONB NOT NULL
        )
	`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS data_tables (
            id SERIAL PRIMARY KEY,
            table_id INTEGER NOT NULL REFERENCES table_definitions(id),
            data JSONB NOT NULL
        )
    `)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS test_table_definitions (
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL,
            creator TEXT NOT NULL,
            columns JSONB NOT NULL
        )
	`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS test_data_tables (
            id SERIAL PRIMARY KEY,
            table_id INTEGER NOT NULL REFERENCES test_table_definitions(id),
            data JSONB NOT NULL
        )
    `)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Tables created successfully")

}

func DeleteAllData() error {
	db := ConnectDB()
	defer db.Close()

	_, err := db.Exec("DELETE FROM test_data_tables")
	if err != nil {
		return fmt.Errorf("failed to delete data from test_data_tables: %v", err)
	}

	_, err = db.Exec("DELETE FROM test_table_definitions")
	if err != nil {
		return fmt.Errorf("failed to delete data from test_table_definitions: %v", err)
	}

	// シーケンスリセット　→　test_data_tables
	_, err = db.Exec("ALTER SEQUENCE test_data_tables_id_seq RESTART WITH 1")
	if err != nil {
		return fmt.Errorf("failed to reset sequence for test_data_tables: %v", err)
	}

	// シーケンスリセット　→　test_table_definitions
	_, err = db.Exec("ALTER SEQUENCE test_table_definitions_id_seq RESTART WITH 1")
	if err != nil {
		return fmt.Errorf("failed to reset sequence for test_table_definitions: %v", err)
	}

	fmt.Println("Data deleted successfully from both tables")
	return nil
}
