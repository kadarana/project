package main

import (
	"log"
	"project/internal/db"
)

func main() {
	connStr := "user=orders_user password=orders_pass host=localhost port=5432 dbname=orders_service sslmode=disable"

	database, err := db.NewDB(connStr)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer database.Close()
}
