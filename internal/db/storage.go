package db

import (
	"database/sql"
	"project/internal/model"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (r *Storage) SaveOrder(order model.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO orders (order_uid, track_number, entry, customer_id, date_created)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (order_uid) DO NOTHING
	`, order.OrderUID, order.TrackNumber, order.Entry, order.CustomerID, order.DateCreated)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO deliveries (order_uid, name, phone, address)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (order_uid) DO NOTHING
	`, order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Address)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO payments (order_uid, transaction, currency, amount)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (order_uid) DO NOTHING
	`, order.OrderUID, order.Payment.Transaction, order.Payment.Currency, order.Payment.Amount)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		_, err = tx.Exec(`
		INSET INTO items (order_uid, name, price)
		VALUES ($1, $2, $3)
		`, order.OrderUID, item.Name, item.Price)
	}

	return tx.Commit()
}

func (r *Storage) GetOrderByID(id string) (model.Order, error) {
	var order model.Order

	err := r.db.QueryRow(`
		SELECT order_uid, track_number, entry, customer_id, date_created 
		FROM orders WHERE order_uid = $1
	`, id).Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.CustomerID, &order.DateCreated)
	if err != nil {
		return order, err
	}

	err = r.db.QueryRow(`
		SELECT order_uid, name, phone, address 
		FROM deliveries WHERE order_uid = $1
	`, id).Scan(&order.OrderUID, &order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Address)
	if err != nil {
		return order, err
	}

	err = r.db.QueryRow(`
		SELECT order_uid, transaction, currency, amount 
		FROM payments WHERE order_uid = $1
	`, id).Scan(&order.OrderUID, &order.Payment.Transaction, &order.Payment.Currency, &order.Payment.Amount)
	if err != nil {
		return order, err
	}

	rows, err := r.db.Query(`
        SELECT id, name, price
        FROM items WHERE order_uid = $1
    `, id)
	if err != nil {
		return order, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Price); err != nil {
			return order, err
		}

		order.Items = append(order.Items, item)
	}

	return order, nil

}

func (r *Storage) GetAllOrders() ([]model.Order, error) {
	orders := []model.Order{}

	// 1. Берём все заказы
	rows, err := r.db.Query(`
        SELECT order_uid, track_number, entry, customer_id, date_created
        FROM orders
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 2. Для каждого заказа достаём связанные данные
	for rows.Next() {
		var order model.Order
		if err := rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.CustomerID, &order.DateCreated); err != nil {
			return nil, err
		}

		// 2.1 Доставка
		err = r.db.QueryRow(`
            SELECT name, phone, address
            FROM deliveries WHERE order_uid = $1
        `, order.OrderUID).Scan(&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Address)
		if err != nil {
			return nil, err
		}

		// 2.2 Оплата
		err = r.db.QueryRow(`
            SELECT transaction, currency, amount
            FROM payments WHERE order_uid = $1
        `, order.OrderUID).Scan(&order.Payment.Transaction, &order.Payment.Currency, &order.Payment.Amount)
		if err != nil {
			return nil, err
		}

		// 2.3 Товары
		itemsRows, err := r.db.Query(`
            SELECT id, name, price
            FROM items WHERE order_uid = $1
        `, order.OrderUID)
		if err != nil {
			return nil, err
		}

		for itemsRows.Next() {
			var item model.Item
			if err := itemsRows.Scan(&item.ID, &item.Name, &item.Price); err != nil {
				itemsRows.Close()
				return nil, err
			}
			order.Items = append(order.Items, item)
		}
		itemsRows.Close()

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil

}
