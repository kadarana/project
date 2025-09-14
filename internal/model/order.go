package model

type Order struct {
	OrderUID    string `json:"order_uid"`
	TrackNumber string `json:"track_number"`
	Entry       string `json:"entry"`
	CustomerID  string `json:"customer_id"`
	DateCreated string `json:"date_created"`

	Delivery Delivery `json:"delivery"`
	Payment  Payment  `json:"payment"`
	Item     []Item   `json:"items"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type Payment struct {
	Transaction string `json:"transaction"`
	Currency    string `json:"currency"`
	Amount      string `json:"amount"`
}

type Item struct {
	ID    int     `json:"id,omitempty"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
