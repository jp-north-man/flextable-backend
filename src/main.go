package main

import (
	"net/http"

	"example.com/packages/models"
)

func main() {
	models.ConnectDB()
	http.ListenAndServe(":8080", nil)
}
