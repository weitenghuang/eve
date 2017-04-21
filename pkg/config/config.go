package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
)

const (
	DEFAULT_PORT        = "8088"
	DEFAULT_DNS         = "localhost"
	DEFAULT_SCHEME      = "http"
	DEFAULT_DB_PORT     = "28015"
	DEFAULT_DB_NAME     = "eve"
	DEFAULT_QUEUE_PORT  = "4222"
	DEFAULT_ENVIRONMENT = "DEV"
)

type ApiServerConfig struct {
	Port     string
	DNS      string
	Scheme   string
	Hostname string
	CertFile string
	KeyFile  string
}

type RethinkDbConfig struct {
	Url          string
	DatabaseName string
	InitialCap   int
	MaxOpen      int
	TLSConfig    *tls.Config
}

type NatsConfig struct {
	Url            string
	AllowReconnect bool
	MaxReconnect   int
	ReconnectWait  time.Duration
	Timeout        time.Duration
}

type SystemConfig struct {
	Hostname    string
	Version     string
	Environment string
}

func NewApiServerConfig() *ApiServerConfig {
	port := os.Getenv("EVE_PORT")
	if port == "" {
		port = DEFAULT_PORT
	}
	dns := os.Getenv("EVE_DNS")
	if dns == "" {
		dns = DEFAULT_DNS
	}
	scheme := os.Getenv("EVE_SCHEME")
	if scheme == "" {
		scheme = DEFAULT_SCHEME
	}
	certfile := os.Getenv("EVE_CERT_FILE")
	if certfile == "" {
		certfile = filepath.Join("/opt", "tls", "eve-server.pem")
	}
	keyfile := os.Getenv("EVE_KEY_FILE")
	if keyfile == "" {
		keyfile = filepath.Join("/opt", "tls", "eve-server-key.pem")
	}
	hostname, _ := os.Hostname()
	return &ApiServerConfig{
		Port:     port,
		DNS:      dns,
		Scheme:   scheme,
		Hostname: hostname,
		CertFile: certfile,
		KeyFile:  keyfile,
	}
}

func NewRethinkDbConfig() *RethinkDbConfig {
	url := os.Getenv("EVE_DB_URL")
	if url == "" {
		url = fmt.Sprintf("%s:%s", DEFAULT_DNS, DEFAULT_DB_PORT)
	}

	dbName := os.Getenv("EVE_DB_NAME")
	if dbName == "" {
		dbName = DEFAULT_DB_NAME
	}

	tlsConfig := &tls.Config{}
	rootCAFile := os.Getenv("EVE_DB_CA_CERT")

	if rootCAFile != "" {
		pool, err := LoadCAFile(rootCAFile)
		if err != nil {
			log.Println(err.Error())
		} else {
			tlsConfig.RootCAs = pool
		}
	}

	return &RethinkDbConfig{
		Url:          url,
		DatabaseName: dbName,
		InitialCap:   4,
		MaxOpen:      8,
		TLSConfig:    tlsConfig,
	}
}

func LoadCAFile(caFile string) (*x509.CertPool, error) {
	pool := x509.NewCertPool()

	pem, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("Error loading CA File: %s", err)
	}

	ok := pool.AppendCertsFromPEM(pem)
	if !ok {
		return nil, fmt.Errorf("Error loading CA File: Couldn't parse PEM in: %s", caFile)
	}

	return pool, nil
}

func NewNatsConfig() *NatsConfig {
	url := os.Getenv("EVE_QUEUE_URL")
	if url == "" {
		url = fmt.Sprintf("nats://%s:%s", DEFAULT_DNS, DEFAULT_QUEUE_PORT)
	}
	maxReconnect, err := strconv.Atoi(os.Getenv("EVE_QUEUE_MAX_RECONNECT"))
	if err != nil {
		maxReconnect = 30
	}
	return &NatsConfig{
		AllowReconnect: true,
		MaxReconnect:   maxReconnect,
		ReconnectWait:  1 * time.Second,
		Timeout:        1 * time.Second,
		Url:            url,
	}
}

func NewSystemConfig() *SystemConfig {
	hostname, err := os.Hostname()
	if err != nil {
		log.Println(err)
	}
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = DEFAULT_ENVIRONMENT
	}
	version := os.Getenv("VERSION")
	return &SystemConfig{
		Hostname:    hostname,
		Version:     version,
		Environment: env,
	}
}
