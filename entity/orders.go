package entity

import (
	"time"
)

type Items struct {
	ItemId      int    `gorm:"primaryKey" json:"itemId"`
	ItemCode    string `gorm:"itemCode" json:"itemCode"`
	Description string `gorm:"description" json:"description"`
	Quantity    int    `gorm:"quantity" json:"qty"`
	OrderId     int    `gorm:"references:OrderId" json:"orderId"`
}

type Orders struct {
	OrderId      int       `gorm:"primaryKey"`
	CustomerName string    `json:"customerName"`
	OrderedAt    time.Time `gorm:"autoCreateTime"`
	Item         []Items   `gorm:"foreignKey:OrderId" json:"items"`
}

type OrderResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []Orders
}
