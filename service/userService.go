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

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func AllUser(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	var response entity.Response
	var arrUser []entity.User

	db := config.Connect()

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
	// var response entity.Response
	var user entity.User

	db := config.Connect()
	defer db.Close()

	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)

	if err != nil {
		responseReturn(501, "Error Encoding JSON", []entity.User{}, w)
	}

	bs, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	rows, err := db.Exec("INSERT INTO users(username,email,password,age,created_at) VALUES(?,?,?,?,?)", user.Username, user.Email, bs, user.Age, time.Now())
	if err != nil {
		responseReturn(501, "Error Encoding JSON", []entity.User{}, w)
	}

	insertUser, _ := rows.LastInsertId()
	user.Id = int(insertUser)
	user.CreatedAt = time.Now()

	responseReturn(200, "User Created", []entity.User{user}, w)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(userId)

	db := config.Connect()
	defer db.Close()

	_, err := db.Exec("DELETE FROM users WHERE id=?", id)

	if err != nil {
		log.Print(err)
		responseReturn(200, "Error Delete User", []entity.User{}, w)
	}

	responseReturn(200, "Success Delete User", []entity.User{}, w)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
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
		responseReturn(500, "Error Update User", []entity.User{}, w)
	}
	responseReturn(200, "Success Update User", []entity.User{}, w)

}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var creds entity.Credentials
	var user entity.User
	var jwtKey = []byte("my_secret_key")

	db := config.Connect()
	defer db.Close()

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rows, err := db.Query("SELECT username,password FROM users WHERE username=?", creds.Username)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for rows.Next() {
		err = rows.Scan(&user.Username, &user.Password)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	fmt.Println(user.Password)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &entity.Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(tokenString))

	// http.SetCookie(w, &http.Cookie{
	// 	Name:    "token",
	// 	Value:   tokenString,
	// 	Expires: expirationTime,
	// })
}

func responseReturn(status int, msg string, data []entity.User, w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(entity.Response{Status: status, Message: msg, Data: data})
}
