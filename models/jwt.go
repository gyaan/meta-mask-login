package models

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	PublicAddress string `json:"public_address"`
	jwt.StandardClaims
}

type TokenEndExpire struct {
	Token string
	Expire string
}

type Key int

//MyKey jwt handdler
const (
	JwtKey    Key = 100000
	DbKey     Key = 200000
	UserKey   Key = 300000
	ProjKey   Key = 400000
	SecretKey     = "My-temporary-SECRETKey-2017"
)