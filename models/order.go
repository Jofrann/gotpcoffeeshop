package models

import "time"

type OrderStatus string

const (
	StatusPending  OrderStatus = "pending"
	StatusReady    OrderStatus = "ready"
	StatusPickedUp OrderStatus = "picked-up"
)

type Order struct {
	ID         string      `json:"id"`
	DrinkID    string      `json:"drinkId"`
	DrinkName  string      `json:"drinkName"`
	Size       string      `json:"size"`
	Extras     []string    `json:"extras"`
	Status     OrderStatus `json:"status"`
	OrderedAt  time.Time   `json:"orderedAt"`
	TotalPrice float64     `json:"totalPrice"`
}


// Order représente une commande
type OrderInput struct {
	DrinkID      string `json:"drink_id"`
	Size         string `json:"size"`
	CustomerName string `json:"customer_name"`
}
