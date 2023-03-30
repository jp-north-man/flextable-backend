package main

import (
	"log"
	"net/http"

	"example.com/packages/handler"
	"example.com/packages/models"
)

func main() {
	models.CreateDB()
	err := models.DeleteAllData()
	if err != nil {
		log.Fatal("Failed to delete all data: ", err)
	}
	http.HandleFunc("/create", handler.CreateHandler)
	http.HandleFunc("/add_row", handler.AddRowHandler)
	http.HandleFunc("/list", handler.GetFlexTablesHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
