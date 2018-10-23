package middlewares

import (
	"net/http"
	"strings"
	"log"
	"github.com/dgrijalva/jwt-go"
	mid "github.com/gyaan/meta-mask-login/models"
	"context"
	"fmt"
)

func TokenAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

		var tokenString string

		//if jwt is passed in query string
		//tokenString = r.URL.Query().Get("jwt")
		bearer := r.Header.Get("Authorization")

		if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
			tokenString = bearer[7:]
		}

		if tokenString==""{
			log.Printf("No token found!!")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		claims := mid.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return  []byte(mid.SecretKey),nil
		})

		//use this address of the apis to get the user details
		fmt.Printf(claims.PublicAddress)

		if claims, ok := token.Claims.(*mid.Claims); ok && token.Valid{
			log.Printf("Token verified")
			ctx := context.WithValue(r.Context(), mid.JwtKey, *claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors & jwt.ValidationErrorMalformed !=0{
				log.Printf("Token is not valid jwt token")
			} else if ve.Errors & (jwt.ValidationErrorExpired | jwt.ValidationErrorNotValidYet) !=0{
				log.Printf("Token expired or not active")
			} else {
				log.Printf("Can't handle token")
			}
		}else {
			log.Printf("Can't handle token")
		}
		w.WriteHeader(http.StatusUnauthorized)
		})
}