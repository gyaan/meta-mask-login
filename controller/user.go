package controller

import (
	"net/http"
	"fmt"
	"github.com/go-chi/chi"
	"encoding/json"
	"github.com/gyaan/meta-mask-login/models"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"

)

func GetUserFiles(s * r.Session) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		fmt.Printf("Testing the accesstoken validation")
		//todo get the loggedin user files and return or something like that
        writer.WriteHeader(http.StatusOK)
	}
}

func GetUserDetails(s * r.Session) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		publicAddress := chi.URLParam(request, "publicAddress")
		fmt.Printf(publicAddress)

		cursor, err := r.DB("block_chain").Table("public_addresses").Filter(r.Row.Field("public_address").Eq(publicAddress)).Run(s)
		if err != nil {
			writer.WriteHeader(http.StatusNotFound)
			return
		}

		//create the random nonce store it and send it with response
        var user models.User
        err = cursor.One(&user)

		cursor.Close()
		fmt.Printf("%+v\n", user)

        if err == r.ErrEmptyResult {

        	//todo put  flag for this
			user.Nonce = "123x" //todo generate with some complex logic
			user.PublicAddress = publicAddress
			_, err := r.DB("block_chain").Table("public_addresses").Insert(user).RunWrite(s)
			if err != nil{
				fmt.Println(err)
			}
		}

		res, _ := json.Marshal(user)
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		fmt.Fprintf(writer, "%s", res)

	}
}

func CreateUser(s * r.Session) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		user := models.User{}
		json.NewDecoder(request.Body).Decode(&user)
        user.Nonce = "123x" //todo generate with some complex logic
        result, err := r.DB("block_chain").Table("public_addresses").Insert(user).RunWrite(s)

        if err != nil{
        	fmt.Println(err)
		}

		fmt.Println(result)
		res, _ := json.Marshal(user)

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		fmt.Fprintf(writer, "%s", res)
	}
}
