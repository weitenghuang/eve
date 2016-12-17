package nats

import (
	"github.com/nats-io/go-nats"
	"log"
	"os"
	"time"
)

const (
	MAX_RETRY int = 60
)

func EncodedConn() (*nats.EncodedConn, error) {
	url := os.Getenv("EVE_QUEUE_URL")
	if url == "" {
		url = nats.DefaultURL
	}

	opts := &nats.Options{
		AllowReconnect: true,
		MaxReconnect:   30,
		ReconnectWait:  5 * time.Second,
		Timeout:        1 * time.Second,
		Url:            url,
	}
	nc, err := Connect(opts, 0)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return c, nil
}

func Connect(opts *nats.Options, retry int) (*nats.Conn, error) {
	nc, err := (*opts).Connect()
	if err != nil {
		if err == nats.ErrNoServers && retry < MAX_RETRY {
			retry++
			log.Println("Wait for", retry, "seconds to re-try nats queue server connection.")
			// Max total wait time = 30 minutes
			time.Sleep(time.Duration(retry) * time.Second)
			nc, err = Connect(opts, retry)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return nc, nil
}
