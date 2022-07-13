package main

import (
	"fmt"
	"log"
	"net/http"

	"test-api/service"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	// myRouter.HandleFunc("/user", service.NewUser).Methods("POST")
	// myRouter.HandleFunc("/user", service.AllUser).Methods("GET")
	// myRouter.HandleFunc("/user/{id}", service.GetUserById).Methods("GET")
	// myRouter.HandleFunc("/user/{id}", service.UpdateUser).Methods("PUT")
	// myRouter.HandleFunc("/user/{id}", service.DeleteUser).Methods("DELETE")
	myRouter.HandleFunc("/orders", service.CreateOrder).Methods("POST")
	myRouter.HandleFunc("/orders", service.GetOrder).Methods("GET")
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {
	handleRequests()
}
