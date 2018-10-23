package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"fmt"
	"github.com/gyaan/meta-mask-login/routes"
	"github.com/gyaan/meta-mask-login/dbs"
	"net/http"
)

func init() {
	if err:= godotenv.Load();err!=nil{
		log.Printf("Error loading .env file!")
        os.Exit(1)
	}
}

func main() {
    /*
	sessions:= dbs.StartDispatch()
	dbName:=os.Getenv("DATABASE_NAME")
	fmt.Printf(dbName)
	*/

	sessions := dbs.GetRethinkSession()

	err:= http.ListenAndServe(":"+os.Getenv("APPLICATION_PORT"), routes.Router(sessions))

	if err==nil{
		fmt.Println("Problem while creating server")
	}
}