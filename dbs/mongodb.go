package dbs

import (
	"gopkg.in/mgo.v2"
	"os"
	"time"
	"log"
)

type MgoSession struct {
	Session *mgo.Session
}

func newMgoSession(s *mgo.Session) *MgoSession {
	return &MgoSession{s}
}

func StartMongoDb(msg string) *MgoSession {
	mongoDbDailInfo := &mgo.DialInfo{
		Addrs:   []string{os.Getenv("DATABASE_HOST")},
		Timeout: 60 * time.Second,
	}

	mongoSession, err := mgo.DialWithInfo(mongoDbDailInfo)

	if err != nil {
		log.Fatalf("[MongoDB] CreateSession: %s\n", err)
	}

	mongoSession.SetMode(mgo.Monotonic, true)
	log.Printf("[MongoDB] connected! %s", msg)

	return newMgoSession(mongoSession)
}
