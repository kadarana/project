package nats

import (
	"encoding/json"
	"log"
	"project/internal/cache"
	"project/internal/db"
	"project/internal/model"

	stan "github.com/nats-io/stan.go"
)

func Subscribe(storage *db.Storage, c *cache.Cache) error {
	sc, err := stan.Connect("test-cluster", "orders-service-client", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		return err
	}
	defer sc.Close()

	_, err = sc.Subscribe("orders", func(m *stan.Msg) {
		var order model.Order

		if err := json.Unmarshal(m.Data, &order); err != nil {
			log.Println("Invalid JSON:", err)
			return
		}

		if err := storage.SaveOrder(order); err != nil {
			log.Println("Failed to save order", err)
			return
		}

		c.Set(order)

		log.Printf("Order %s saved $ cached\n", order.OrderUID)

	}, stan.DeliverAllAvailable())

	if err != nil {
		return err
	}

	select {}

}
