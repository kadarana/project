package main

import (
	"log"
	"net/http"
	"project/internal/api"
	"project/internal/cache"
	"project/internal/db"
	"project/internal/nats"
)

func main() {
	connStr := "user=orders_user password=orders_pass host=localhost port=5432 dbname=orders_service sslmode=disable"
	database, err := db.NewDB(connStr)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	defer database.Close()

	storage := db.NewStorage(database)

	c := cache.NewCache()

	orders, err := storage.GetAllOrders()
	if err != nil {
		log.Println("Failed to preload cache:", err)
	} else {
		c.LoadFromDB(orders)
		log.Printf("Cache preloaded with %d orders\n", len(orders))
	}

	go func() {
		if err := nats.Subscribe(storage, c); err != nil {
			log.Fatal("NATS subscribe error:", err)
		}
	}()

	srv := api.NewServer(c)
	log.Println("HTTP server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", srv))
}
