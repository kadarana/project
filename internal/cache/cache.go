package cache

import "project/internal/model"

type Cache struct {
	data map[string]model.Order
}

func NewCache() *Cache {
	return &Cache{data: make(map[string]model.Order)}
}

func (c *Cache) Set(order model.Order) {
	c.data[order.OrderUID] = order
}

func (c *Cache) Get(id string) (model.Order, bool) {
	order, ok := c.data[id]
	return order, ok
}

func (c *Cache) LoadFromDB(orders []model.Order) {
	for _, o := range orders {
		c.data[o.OrderUID] = o
	}
}
