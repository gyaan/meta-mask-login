package dbs

import (
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"log"
)

func GetRethinkSession() *r.Session {
	session, err := r.Connect(r.ConnectOpts{
		Address: "db:28015",
	})
	if err != nil {
		log.Fatalln(err)
	}

	//check db is there or not if not create the database
	err = r.TableList().Contains("public_addresses").Do(r.Branch(r.Row, r.Expr(nil), r.Do(func() r.Term {
		return r.TableCreate("public_addresses")
	}))).Exec(session)

	return session
}
