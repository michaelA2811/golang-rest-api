package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"test-api/config"
	"test-api/entity"
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order entity.Orders

	db := config.Connect()

	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&order)
	if err != nil {
		responseOrderReturn(501, "Error Encoding JSON", []entity.Orders{}, w)
	}

	err = db.Create(&order).Error
	if err != nil {
		responseOrderReturn(500, "Error Create Order", []entity.Orders{}, w)
	}

	responseOrderReturn(200, "Order Created", []entity.Orders{order}, w)

}

func GetOrder(w http.ResponseWriter, r *http.Request) {
	order := []entity.Orders{}
	db := config.Connect()
	result := db.Model(&entity.Orders{}).Preload("Items").Find(&order).Error
	if result != nil {
		responseOrderReturn(500, "Order not Found", []entity.Orders{}, w)
	}
	fmt.Println(result)
	// res, _ := json.Marshal(order)
	responseOrderReturn(200, "Order Created", order, w)
}

func UpdateOrder(w http.ResponseWriter, r *http.Request) {

}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {

}
func responseOrderReturn(status int, msg string, data []entity.Orders, w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(entity.OrderResponse{Status: status, Message: msg, Data: data})
}
