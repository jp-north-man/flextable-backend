package models

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func ConnectDB() {
	db, err := sql.Open("postgres",
		"host=localhost port=5432 user=user password=password dbname=mydb sslmode=disable") //今回は直接書きます。Please read from .env etc.
	if err != nil {
		log.Println("sql.Open: ", err)
	} else {
		log.Println("sql.Open: ", "ok")
	}

	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE table (
			id SERIAL PRIMARY KEY,
		);
	`)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Tables created successfully")

}
