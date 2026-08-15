package models

import "time"

type OrderItem struct {
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity"`
}
type Order struct {
	ID        int         `json:"id"`
	Customer  string      `json:"customer"`
	Email     string      `json:"email"`
	Items     []OrderItem `json:"items"`
	Total     float64     `json:"total"`
	Status    string      `json:"status"`
	CreatedAt time.Time   `json:"createdAt"`
}
