package dbs

import (
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"log"
)

func GetRethinkSession() *r.Session {
	session, err := r.Connect(r.ConnectOpts{
		Address: "localhost:28015",
	})
	if err != nil {
		log.Fatalln(err)
	}
	return session
}
