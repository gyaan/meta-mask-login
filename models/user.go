package models

import "time"

type User struct {
	//Id            bson.ObjectId `bson:"_id" json:"id" rethinkdb:"id"`
	PublicAddress string        `bson:"public_address" json:"public_address" rethinkdb:"public_address"`
	Name          string        `bson:"name" json:"name" rethinkdb:"name"`
	Nonce         string        `bson:"nonce" json:"nonce" rethinkdb:"nonce"`
	CreatedAt     time.Time      `bson:"created_at" json:"created_at" rethinkdb:"created_at"`
}
