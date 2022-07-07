package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"test-api/config"
	"test-api/entity"

	"github.com/gorilla/mux"
)

func AllUser(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	var response entity.Response
	var arrUser []entity.User

	db := config.Connect()
	defer db.Close()

	rows, err := db.Query("SELECT id,username,email,password,age FROM users")
	if err != nil {
		log.Print(err)
	}

	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Age)
		if err != nil {
			log.Fatal(err.Error())
		} else {
			arrUser = append(arrUser, user)
		}
	}

	response.Status = 200
	response.Message = "Success"
	response.Data = arrUser

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	var arrUser []entity.User
	var response entity.Response

	userId := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(userId)

	db := config.Connect()
	defer db.Close()

	rows, err := db.Query("SELECT id,username,email,password,age FROM users WHERE id=?", id)
	if err != nil {
		log.Print(err)
	}

	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Age)
		if err != nil {
			log.Fatal(err.Error())
		} else {
			arrUser = append(arrUser, user)
		}
	}

	response.Status = 200
	response.Message = "Success"
	response.Data = arrUser

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func NewUser(w http.ResponseWriter, r *http.Request) {
	var response entity.Response

	db := config.Connect()
	defer db.Close()

	err := r.ParseMultipartForm(4096)
	if err != nil {
		panic(err)
	}
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	age := r.FormValue("age")

	_, err = db.Exec("INSERT INTO users(username,email,password,age) VALUES(?,?,?,?)", username, email, password, age)

	if err != nil {
		log.Print(err)
		return
	}
	response.Status = 200
	response.Message = "Insert data successfully"
	fmt.Print("Insert data to database")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var response entity.Response
	userId := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(userId)

	db := config.Connect()
	defer db.Close()

	err := r.ParseMultipartForm(4096)

	if err != nil {
		panic(err)
	}

	_, err = db.Exec("DELETE FROM users WHERE id=?", id)

	if err != nil {
		log.Print(err)
		return
	}

	response.Status = 200
	response.Message = "Delete data successfully"
	fmt.Print("Delete data successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var response entity.Response
	userId := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(userId)

	db := config.Connect()
	defer db.Close()

	err := r.ParseMultipartForm(4096)

	if err != nil {
		panic(err)
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	age := r.FormValue("age")

	if username != "" {
		_, err = db.Exec("UPDATE users SET username=? WHERE id=?", username, id)
	} else if email != "" {
		_, err = db.Exec("UPDATE users SET email=? WHERE id=?", email, id)
	} else if password != "" {
		_, err = db.Exec("UPDATE users SET password=? WHERE id=?", password, id)
	} else {
		_, err = db.Exec("UPDATE users SET age=? WHERE id=?", age, id)
	}

	if err != nil {
		log.Print(err)
	}

	response.Status = 200
	response.Message = "Update data successfully"
	fmt.Print("Update data to database")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)

}
