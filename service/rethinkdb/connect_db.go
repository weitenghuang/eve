package rethinkdb

import (
	log "github.com/Sirupsen/logrus"
	"github.com/concur/eve/pkg/config"
	r "gopkg.in/gorethink/gorethink.v3"
	"sync"
)

type DbSession struct {
	Session *r.Session
	Url     string
	DbName  string
}

var defaultSession *DbSession

func buildSession() error {
	var mu sync.Mutex
	dbConfig := config.NewRethinkDbConfig()
	log.Println("Connecting to database...")
	session, err := r.Connect(r.ConnectOpts{
		Address:    dbConfig.Url,
		InitialCap: dbConfig.InitialCap,
		MaxOpen:    dbConfig.MaxOpen,
	})
	if err != nil {
		return err
	}
	mu.Lock()
	defaultSession = &DbSession{
		Session: session,
		Url:     dbConfig.Url,
		DbName:  dbConfig.DatabaseName,
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
