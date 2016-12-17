package rethinkdb

import (
	r "gopkg.in/gorethink/gorethink.v3"
	"log"
	"os"
	"sync"
)

const (
	DEFAULT_URL     = "localhost:28015"
	DEFAULT_DB_NAME = "eve"
)

type DbSession struct {
	Session *r.Session
	Url     string
	DbName  string
}

var defaultSession *DbSession

func init() {
	buildSession()
}

func buildSession() error {
	var mu sync.Mutex
	url := os.Getenv("EVE_DB_URL")
	if url == "" {
		url = DEFAULT_URL
	}

	dbName := os.Getenv("EVE_DB_NAME")
	if dbName == "" {
		dbName = DEFAULT_DB_NAME
	}
	log.Println("Connecting to database...")
	session, err := r.Connect(r.ConnectOpts{
		Address:    url,
		InitialCap: 4,
		MaxOpen:    8,
	})
	if err != nil {
		return err
	}
	mu.Lock()
	defaultSession = &DbSession{
		Session: session,
		Url:     url,
		DbName:  dbName,
	}
	mu.Unlock()
	return nil
}

func closeDefaultSession() error {
	err := defaultSession.Session.Close()
	if err != nil {
		log.Println(err.Error())
	}
	return nil
}

func DefaultSession() *DbSession {
	if defaultSession == nil {
		if err := buildSession(); err != nil {
			log.Println(err)
		}
	}
	if defaultSession != nil && defaultSession.Session != nil && !defaultSession.Session.IsConnected() {
		if err := defaultSession.Session.Reconnect(); err != nil {
			log.Println(err)
		}
	}
	return defaultSession
}
