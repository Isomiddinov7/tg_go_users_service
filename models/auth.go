package models

import "github.com/golang-jwt/jwt"

var JwtKey = []byte("my_secret_key")

type Claims struct {
	Login string `json:"login"`
	jwt.StandardClaims
}
