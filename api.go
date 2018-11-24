package main

import (
	"gopkg.in/mgo.v2"
	"net/http"
	"log"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"strings"
)

type User struct {
	Id            bson.ObjectId `bson:"_id" json:"id"`
	PublicAddress string        `bson:"public_address" json:"public_address"`
	Nonce         int           `bson:"nonce" json:"nonce"`
	UserName      string        `bson:"user_name" json:"user_name"`
}
type DatabaseObj struct {
	Server   string
	Database string
}

type AuthReq struct {
	PublicAddress string `bson:"public_address" json:"public_address"`
	Signature     string `bson:"signature" json:"signature"`
}

func (m *DatabaseObj) connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

var dbObj = DatabaseObj{}
var db *mgo.Database

const collection = "users"

func init() {
	dbObj.Database = "test"
	dbObj.Server = "localhost"
	dbObj.connect()
}

func main() {
	r := mux.NewRouter()
	/**Routes*/
	r.HandleFunc("/users", getUserDetails).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/authentication", AuthenticateUser).Methods("POST", "OPTIONS")

	/*testing jwt*/

	//r.HandleFunc("get-token",GetJwtToken).Methods("GET")

	if err := http.ListenAndServe(":1335", r);
		err != nil {
		log.Fatal(err)
	}
}

func AuthenticateUser(writer http.ResponseWriter, request *http.Request) {

	if request.Method == "OPTIONS" {
		OptionRequestCros(writer)
		return
	} else {
		enableCors(&writer)
	}

	decoder := json.NewDecoder(request.Body)
	var t AuthReq
	err := decoder.Decode(&t)
	if err != nil {
		log.Fatal(err)
	}

	msg := []byte("I am signing my one-time nonce: 1234")

	fromAddr := common.HexToAddress(strings.TrimSpace(t.PublicAddress))
	sig := hexutil.MustDecode(strings.TrimSpace(t.Signature))

	if sig[64] != 27 && sig[64] != 28 {
		fmt.Println("something wrong with the address!!")
	}
	sig[64] -= 27

	//get the public address from signature
	pubKey, err := crypto.SigToPub(signHash(msg), sig)
	if err != nil {
		fmt.Println("something wrong with public key")
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	responseMessage1 := map[string]string{}

	if fromAddr == recoveredAddr {
		responseMessage1 = map[string]string{"error": "", "data": "someAccessToken"}
	}
	response, _ := json.Marshal(responseMessage1)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(200)
	writer.Write(response)
}


func createUser(writer http.ResponseWriter, request *http.Request) {
	enableCors(&writer)
	defer request.Body.Close()
	var user User
	if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
		fmt.Println("problem with request")
	}
	user.Id = bson.NewObjectId()
	user.Nonce = 1234

	if err := db.C(collection).Insert(user); err != nil {
		fmt.Println("problem with insertion")
	}
	response, _ := json.Marshal(user)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(200)
	writer.Write(response)
}

func getUserDetails(writer http.ResponseWriter, request *http.Request) {
	enableCors(&writer)
	PublicAddress := request.URL.Query().Get("publicAddress")
	fmt.Println(PublicAddress)
	user := User{}
	err := db.C(collection).Find(bson.M{"public_address": PublicAddress}).One(&user)
	if err != nil {
		//panic(err)
	}

	if (User{}) == user { //no user found
		user.Id = bson.NewObjectId()
		user.PublicAddress = PublicAddress
		user.Nonce = 1234
		if err := db.C(collection).Insert(&user); err != nil {
			fmt.Println("problem with insertion")
		}
	}
	response, _ := json.Marshal(user)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(200)
	writer.Write(response)
}
func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	//fmt.Println(msg)
	return crypto.Keccak256([]byte(msg))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Origin, Accept, Content-Type, Content-Length,X-Requested-With, Accept-Encoding, X-CSRF-Token, Authorization")
	(*w).Header().Set("Access-Control-Allow-Accept", "application/json, text/html")
	(*w).Header().Set("Access-Control-Allow-Content-type", "application/json, text/html")
}
func OptionRequestCros(writer http.ResponseWriter) {
	// Add a generous access-control-allow-origin header for CORS requests
	writer.Header().Add("Access-Control-Allow-Origin", "*")
	// Only GET/POST Methods are supported
	writer.Header().Set("Access-Control-Allow-Methods", "GET, POST")
	// Only these headers can be set
	(writer).Header().Set("Access-Control-Allow-Headers", "Origin, Accept, Content-Type, Content-Length,X-Requested-With, Accept-Encoding, X-CSRF-Token, Authorization")
	(writer).Header().Set("Access-Control-Allow-Accept", "application/json, text/html")
	(writer).Header().Set("Access-Control-Allow-Content-type", "application/json, text/html")
	// Indicate that no content will be returned
	writer.WriteHeader(204)
}

func GetJwtToken(writer http.ResponseWriter, requst http.Request) {



}
