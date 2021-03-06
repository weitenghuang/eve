package health

import (
	"crypto/tls"
	"github.com/concur/eve"
	"github.com/concur/eve/pkg/config"
	r "gopkg.in/gorethink/gorethink.v3"
	"strconv"
	"strings"
	"time"
)

type RethinkdbChecker struct {
	Addresses    []string
	DatabaseName string
	InitialCap   int
	MaxOpen      int
	Timeout      time.Duration
	TLSConfig    *tls.Config
}

func NewRethinkdbChecker() *RethinkdbChecker {
	rethinkConfig := config.NewRethinkDbConfig()
	// Health page will use rethink cluster loadbalacer directly as connection address
	return &RethinkdbChecker{
		Addresses:    rethinkConfig.Addresses,
		DatabaseName: rethinkConfig.DatabaseName,
		InitialCap:   1,
		MaxOpen:      1,
		Timeout:      1 * time.Second,
		TLSConfig:    rethinkConfig.TLSConfig,
	}
}

func (rChecker *RethinkdbChecker) Ping() *eve.Error {
	meta := rChecker.rethinkOptsMeta()

	hostname, port := rChecker.splitAddress()
	host := r.NewHost(hostname, port)
	pool, err := r.NewPool(host, &r.ConnectOpts{
		Addresses:  rChecker.Addresses,
		InitialCap: rChecker.InitialCap,
		MaxOpen:    rChecker.MaxOpen,
		Timeout:    rChecker.Timeout,
		TLSConfig:  rChecker.TLSConfig,
	})

	if err != nil {
		return &eve.Error{
			Type:        "RethinkDB",
			Description: "NewPool error",
			Metadata:    meta,
			Error:       err.Error(),
		}
	}

	if err := pool.Ping(); err != nil {
		return &eve.Error{
			Type:        "RethinkDB",
			Description: "Ping error",
			Metadata:    meta,
			Error:       err.Error(),
		}
	}
	return nil
}

func (rChecker *RethinkdbChecker) DbReady() *eve.Error {
	meta := rChecker.rethinkOptsMeta()

	session, err := r.Connect(r.ConnectOpts{
		Addresses:  rChecker.Addresses,
		InitialCap: rChecker.InitialCap,
		MaxOpen:    rChecker.MaxOpen,
		Timeout:    rChecker.Timeout,
		TLSConfig:  rChecker.TLSConfig,
	})

	if err != nil {
		return &eve.Error{
			Type:        "RethinkDB",
			Description: "Connection error",
			Metadata:    meta,
			Error:       err.Error(),
		}
	}

	cursor, err := r.DBList().Run(session)
	defer cursor.Close()
	if err != nil {
		return &eve.Error{
			Type:        "RethinkDB",
			Description: "DBList error",
			Metadata:    meta,
			Error:       err.Error(),
		}
	}

	var row interface{}
	for cursor.Next(&row) {
		if row.(string) == rChecker.DatabaseName {
			return nil
		}
	}

	return &eve.Error{
		Type:        "RethinkDB",
		Description: "DbReady error: eve db not found",
		Metadata:    meta,
		Error:       "eve db not found",
	}
}

func (rChecker *RethinkdbChecker) TableReady() *eve.Error {
	meta := rChecker.rethinkOptsMeta()

	session, err := r.Connect(r.ConnectOpts{
		Addresses:  rChecker.Addresses,
		InitialCap: rChecker.InitialCap,
		MaxOpen:    rChecker.MaxOpen,
		Timeout:    rChecker.Timeout,
		TLSConfig:  rChecker.TLSConfig,
	})

	if err != nil {
		return &eve.Error{
			Type:        "RethinkDB",
			Description: "Connection error",
			Metadata:    meta,
			Error:       err.Error(),
		}
	}

	cursor, err := r.DB(rChecker.DatabaseName).TableList().Run(session)
	defer cursor.Close()
	if err != nil {
		return &eve.Error{
			Type:        "RethinkDB",
			Description: "TableList error",
			Metadata:    meta,
			Error:       err.Error(),
		}
	}

	tables := map[string]byte{
		"quoin":          0,
		"quoinArchive":   0,
		"infrastructure": 0,
	}

	var row interface{}
	for cursor.Next(&row) {
		name := row.(string)
		if _, ok := tables[name]; ok {
			delete(tables, name)
		}
	}

	if len(tables) > 0 {
		err := "Table:"
		for key := range tables {
			err = strings.Join([]string{err, key, "is not found."}, " ")
		}
		return &eve.Error{
			Type:        "RethinkDB",
			Description: "TableReady error",
			Metadata:    meta,
			Error:       err,
		}
	}

	return nil
}

func (rChecker *RethinkdbChecker) rethinkOptsMeta() map[string]string {
	return map[string]string{
		"Addresses":  strings.Join(rChecker.Addresses, ","),
		"InitialCap": strconv.Itoa(rChecker.InitialCap),
		"MaxOpen":    strconv.Itoa(rChecker.MaxOpen),
		"Timeout":    rChecker.Timeout.String(),
	}
}

func (rChecker *RethinkdbChecker) splitAddress() (hostname string, port int) {
	hostname = "localhost"
	port = 28015
	if len(rChecker.Addresses) > 0 {
		addrParts := strings.Split(rChecker.Addresses[0], ":")
		if len(addrParts) >= 1 {
			hostname = addrParts[0]
		}
		if len(addrParts) >= 2 {
			port, _ = strconv.Atoi(addrParts[1])
		}
	}

	return
}
