package main

import (
	"log"
	"net/http"

	"example.com/packages/handler"
	"example.com/packages/models"
)

func main() {
	models.CreateDB()
	http.HandleFunc("/create", handler.CreateHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
