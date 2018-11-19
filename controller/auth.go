package controller

import (
	"net/http"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"encoding/json"
	"github.com/gyaan/meta-mask-login/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gyaan/meta-mask-login/services"
)

func Auth(s *r.Session) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		fmt.Printf("Testing buddy!!")


		//ss := s.MongoDb.Copy()
		//defer ss.Close()

		signature := models.EtherSignature{}
		json.NewDecoder(request.Body).Decode(&signature)

		//first check if public address is there and get the nonce of the public address
		//err := ss.DB(os.Getenv("DATABASE_NAME")).C("users").Find(bson.M{"public_address": signature.PublicAddress}).One(&user)
		cursor, err := r.DB("block_chain").Table("public_addresses").Filter(r.Row.Field("public_address").Eq(signature.PublicAddress)).Run(s)
		if err!=nil{
			writer.WriteHeader(http.StatusNotFound)
			return
		}

		//create the random nonce store it and send it with response
		var user models.User
		err = cursor.One(&user)

		cursor.Close()
		fmt.Printf("%+v\n", user)

		//cursor.Close()
		//public address is there lets verify it
		// get the public address and signature and verify it
		//get the user details and verify end provide the access token

		msg:= []byte("I am signing my one-time nonce: "+user.Nonce) //get this from db

		fromAddress:= common.HexToAddress(signature.PublicAddress)
		sig:= hexutil.MustDecode(signature.Signature)

		if len(sig) < 60{  //todo put some basic validation for signature
			writer.WriteHeader(http.StatusBadRequest)
		}
		if sig[64]!=27 && sig[64]!=28{
			writer.WriteHeader(http.StatusUnauthorized) //todo add the error message as well
			return
		}

		sig[64] -=27

		pubKey, err := crypto.SigToPub(signHash(msg), sig)

		if err !=nil{
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		recoveredAddress:= crypto.PubkeyToAddress(*pubKey)

		//recovered address same as signed public address
        if fromAddress != recoveredAddress{
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		te, err:=  services.GenerateToken(user); if err!=nil{
			writer.WriteHeader(http.StatusUnauthorized)
		}

		AccessToken:= models.AccessToken{}
		AccessToken.Token = te.Token
		AccessToken.User = user

		res, _ := json.Marshal(AccessToken)
		writer.Header().Set("Content-type", "application/json")
		writer.WriteHeader(http.StatusOK)
		fmt.Fprintf(writer, "%s",res)
	}
}

func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	//fmt.Println(msg)
	return crypto.Keccak256([]byte(msg))
}
