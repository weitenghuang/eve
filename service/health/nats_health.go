package health

import (
	"github.com/nats-io/go-nats"
	"github.com/scipian/eve"
	"github.com/scipian/eve/pkg/config"
	"strconv"
	"time"
)

type NatsChecker struct {
	Url            string
	AllowReconnect bool
	MaxReconnect   int
	ReconnectWait  time.Duration
	Timeout        time.Duration
}

func NewNatsChecker() *NatsChecker {
	nConfig := config.NewNatsConfig()
	return &NatsChecker{
		Url:            nConfig.Url,
		AllowReconnect: false,
		MaxReconnect:   1,
		ReconnectWait:  1 * time.Second,
		Timeout:        1 * time.Second,
	}
}

func (nChecker *NatsChecker) Ping() *eve.Error {
	opts := nats.Options{
		AllowReconnect: nChecker.AllowReconnect,
		MaxReconnect:   nChecker.MaxReconnect,
		ReconnectWait:  nChecker.ReconnectWait,
		Timeout:        nChecker.Timeout,
		Url:            nChecker.Url,
	}
	if _, err := opts.Connect(); err != nil {
		meta := nChecker.natsOptsMeta()
		return &eve.Error{
			Type:        "NATS message system",
			Description: "Ping error",
			Metadata:    meta,
			Error:       err.Error(),
		}
	}
	return nil
}

func (nChecker *NatsChecker) natsOptsMeta() map[string]string {
	return map[string]string{
		"Url":            nChecker.Url,
		"AllowReconnect": strconv.FormatBool(nChecker.AllowReconnect),
		"MaxReconnect":   strconv.Itoa(nChecker.MaxReconnect),
		"ReconnectWait":  nChecker.ReconnectWait.String(),
		"Timeout":        nChecker.Timeout.String(),
	}
}
