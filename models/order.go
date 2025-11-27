package models

import "time"

type OrderStatus string


const (
	StatusPending   OrderStatus = "pending"
	StatusPreparing OrderStatus = "preparing"
	StatusReady     OrderStatus = "ready"
	StatusPickedUp  OrderStatus = "picked-up"
	StatusCancelled OrderStatus = "cancelled"
)


type Order struct {
	ID           string      `json:"id"`
	DrinkID      string      `json:"drinkId"`
	DrinkName    string      `json:"drinkName"`
	Size         string      `json:"size"`
	Extras       []string    `json:"extras"`
	CustomerName string      `json:"customerName"`
	Status       OrderStatus `json:"status"`
	OrderedAt    time.Time   `json:"orderedAt"`
	TotalPrice   float64     `json:"totalPrice"`
}



type OrderUpdate struct {
	Extras       []string `json:"extras"`
}