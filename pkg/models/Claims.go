package models

import "github.com/dgrijalva/jwt-go"


type Claims struct{
	Email string
	jwt.StandardClaims
}
