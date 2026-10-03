package main

import (
	"context"
	"log"
	"net/http"

	database "go-task-api/DB"
	"go-task-api/routes"
)

func main() {
	client, collection, err := database.ConnectDB()
	if err != nil {
		log.Fatal("database connection failed: ", err)
	}
	defer client.Disconnect(context.Background())

	// Main wires dependencies and starts the HTTP server.
	log.Println("API listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", routes.NewRouter(collection)))
}
