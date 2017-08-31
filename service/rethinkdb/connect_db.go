package rethinkdb

import (
	log "github.com/Sirupsen/logrus"
	"github.com/scipian/eve/pkg/config"
	r "gopkg.in/gorethink/gorethink.v3"
	"sync"
)

type DbSession struct {
	Session *r.Session
	DbName  string
}

var defaultSession *DbSession

func buildSession() error {
	var mu sync.Mutex
	dbConfig := config.NewRethinkDbConfig()
	log.Println("Connecting to database...")
	session, err := r.Connect(r.ConnectOpts{
		Addresses:     dbConfig.Addresses,
		InitialCap:    dbConfig.InitialCap,
		MaxOpen:       dbConfig.MaxOpen,
		TLSConfig:     dbConfig.TLSConfig,
		DiscoverHosts: dbConfig.DiscoverHosts,
		Timeout:       dbConfig.Timeout,
		ReadTimeout:   dbConfig.ReadTimeout,
		WriteTimeout:  dbConfig.WriteTimeout,
	})
	if err != nil {
		return err
	}
	mu.Lock()
	defaultSession = &DbSession{
		Session: session,
		DbName:  dbConfig.DatabaseName,
	}
	mu.Unlock()
	return nil
}

func closeDefaultSession() error {
	err := defaultSession.Session.Close()
	if err != nil {
		log.Debugf("Failed to close db session, %v", err)
	}
	return nil
}

func DefaultSession() *DbSession {
	if defaultSession == nil || defaultSession.Session == nil {
		if err := buildSession(); err != nil {
			log.Errorf("Failed to build db session, %v", err)
		}
	} else {
		if !defaultSession.Session.IsConnected() {
			if err := defaultSession.Session.Reconnect(); err != nil {
				log.Debugf("Failed to reconnect db session, %v", err)
				closeDefaultSession()
				if err := buildSession(); err != nil {
					log.Errorf("Failed to build db session, %v", err)
				}
			}
		}
	}
	return defaultSession
}
