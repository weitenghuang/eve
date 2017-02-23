package nats

import (
	log "github.com/Sirupsen/logrus"
	"github.com/concur/rohr/pkg/config"
	"github.com/nats-io/go-nats"
	"math"
	"time"
)

func EncodedConn() (*nats.EncodedConn, error) {
	natsConfig := config.NewNatsConfig()
	opts := &nats.Options{
		AllowReconnect: natsConfig.AllowReconnect,
		MaxReconnect:   natsConfig.MaxReconnect,
		ReconnectWait:  natsConfig.ReconnectWait,
		Timeout:        natsConfig.Timeout,
		Url:            natsConfig.Url,
	}

	// n maxRetry will produce nth partial sum for total re-try's wait time
	maxRetry := int(math.Floor(math.Sqrt(float64(opts.MaxReconnect * 2))))
	nc, err := connect(opts, 0, maxRetry)
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

func connect(opts *nats.Options, retry int, maxRetry int) (*nats.Conn, error) {
	nc, err := (*opts).Connect()
	if err != nil {
		if err == nats.ErrNoServers && retry < maxRetry {
			retry++
			log.Println("Wait for", retry, "seconds to re-try nats queue server connection.")
			time.Sleep(time.Duration(retry) * time.Second)
			nc, err = connect(opts, retry, maxRetry)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return nc, nil
}
