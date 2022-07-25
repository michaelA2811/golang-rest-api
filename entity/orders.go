package entity

import (
	"time"
)

type Items struct {
	ItemId      int    `json:"itemId"`
	ItemCode    string `json:"itemCode"`
	Description string `json:"description"`
	Quantity    int    `json:"qty"`
	OrderId     int    `json:"orderId"`
}

type Orders struct {
	OrderId      int       `json:"orderId"`
	CustomerName string    `json:"customerName"`
	OrderedAt    time.Time `json:"ordered_at"`
	Item         []Items   `json:"items"`
}

type OrderResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []Orders
}
