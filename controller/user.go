package controller

import (
	"net/http"
	"fmt"
	"github.com/go-chi/chi"
	"encoding/json"
	"github.com/gyaan/meta-mask-login/models"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"

	"math/rand"
	"io"
	"os"
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
			fmt.Println("I m here!!!");
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
		randomString := RandomString(7)
		user.Nonce = randomString

        //first check if  user is there
		cursor, err := r.DB("block_chain").Table("public_addresses").Filter(r.Row.Field("public_address").Eq(user.PublicAddress)).Run(s)
		if err != nil {
			writer.WriteHeader(http.StatusNotFound)
			return
		}

		//create the random nonce store it and send it with response
		var existingUser models.User
		err = cursor.One(&existingUser)

		cursor.Close()
		fmt.Printf("%+v\n", existingUser)

		if err == r.ErrEmptyResult {
			//user.PublicAddress = publicAddress
			_, err := r.DB("block_chain").Table("public_addresses").Insert(user).RunWrite(s)
			if err != nil{
				fmt.Println(err)
			}
		}else {
			//generate the new  nonce and update it to db and sent the user details
			//existingUser.Nonce = RandomString(7)
			filter:= map[string]interface{}{"public_address":existingUser.PublicAddress}
			updateData := map[string]interface{}{"nonce": randomString}
			_, err := r.DB("block_chain").Table("public_addresses").Filter(filter).Update(updateData).RunWrite(s)
			if err != nil{
				fmt.Println(err)
			}
			existingUser.Nonce = randomString
			user = existingUser
		}

		res, _ := json.Marshal(user)
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		fmt.Fprintf(writer, "%s", res)
	}
}

func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
func UploadFile(s * r.Session) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		request.ParseMultipartForm(32 << 20)
		file, handler, err := request.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(writer, "%v", handler.Header)
		f, err := os.OpenFile("./"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}