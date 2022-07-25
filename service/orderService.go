package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"test-api/config"
	"test-api/entity"
	"time"

	"github.com/gorilla/mux"
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order entity.Orders

	db := config.Connect()
	defer db.Close()

	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&order)
	if err != nil {
		responseOrderReturn(501, "Error Encoding JSON", []entity.Orders{}, w)
	}

	rows, err := db.Exec("INSERT INTO orders (customer_name,ordered_at) VALUES(?,NOW())", order.CustomerName)
	if err != nil {
		responseOrderReturn(500, "Error Create Order", []entity.Orders{}, w)
	}
	insertOrder, _ := rows.LastInsertId()

	for _, x := range order.Item {
		_, err := db.Exec("INSERT INTO items (item_code,description,quantity,order_id) VALUES (?,?,?,?)", x.ItemCode, x.Description, x.Quantity, insertOrder)
		fmt.Printf("%v", x)
		if err != nil {
			panic(err)
		}
	}

	responseOrderReturn(200, "Order Created", []entity.Orders{order}, w)

}

func GetOrder(w http.ResponseWriter, r *http.Request) {
	orders := []entity.Orders{}
	db := config.Connect()
	defer db.Close()

	result, err := db.Query("SELECT order_id,customer_name,ordered_at FROM orders")
	if err != nil {
		responseOrderReturn(500, "Order not Found", []entity.Orders{}, w)
	}

	for result.Next() {
		var order entity.Orders
		if err := result.Scan(&order.OrderId, &order.CustomerName, &order.OrderedAt); err != nil {
			log.Fatal(err)
		}
		orders = append(orders, order)
	}

	for i, x := range orders {
		rows, err := db.Query("SELECT item_id,item_code,description,quantity,order_id FROM items WHERE order_id=?", x.OrderId)
		if err != nil {
			log.Fatal(err)
		}

		var items []entity.Items
		for rows.Next() {
			var item entity.Items
			if err := rows.Scan(&item.ItemId, &item.ItemCode, &item.Description, &item.Quantity, &item.OrderId); err != nil {
				log.Fatal(err)
			}
			items = append(items, item)
		}
		orders[i].Item = items
	}
	responseOrderReturn(200, "Order Lists", orders, w)
}

func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	var orders entity.Orders
	// var items []entity.Items
	orderId := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(orderId)

	db := config.Connect()
	defer db.Close()

	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&orders)
	if err != nil {
		responseOrderReturn(501, "Error Encoding JSON", []entity.Orders{}, w)
	}

	_, err = db.Exec("UPDATE orders SET customer_name=? ,ordered_at=? WHERE order_id=?", orders.CustomerName, time.Now(), id)
	if err != nil {
		responseOrderReturn(500, "Error Create Order", []entity.Orders{}, w)
	}

	for _, x := range orders.Item {
		_, err := db.Exec("UPDATE items SET quantity=? WHERE order_id=?", x.Quantity, id)

		if err != nil {
			panic(err)
		}
	}

	responseOrderReturn(200, "Success Update Order", []entity.Orders{orders}, w)
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	orderId := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(orderId)
	db := config.Connect()
	defer db.Close()

	_, err := db.Exec("UPDATE items SET order_id=NULL WHERE order_id=?", id)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("DELETE FROM orders WHERE order_id=?", id)
	if err != nil {
		log.Fatal(err)
	}
	w.Write([]byte("Order deleted successfully"))
}
func responseOrderReturn(status int, msg string, data []entity.Orders, w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(entity.OrderResponse{Status: status, Message: msg, Data: data})
}
