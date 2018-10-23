package dbs

import "gopkg.in/mgo.v2"

type Dispatch struct {
	MongoDb *mgo.Session
}

func StartDispatch() *Dispatch {
	mongoSession := StartMongoDb("Dispatch Service").Session
	return &Dispatch{MongoDb: mongoSession}
}
