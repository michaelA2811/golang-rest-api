package entity

import "time"

type User struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Age       int    `json:"age"`
	CreatedAt time.Time
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []User
}
