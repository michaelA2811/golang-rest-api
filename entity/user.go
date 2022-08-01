package entity

import (
	"time"

	"github.com/golang-jwt/jwt"
)

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

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
