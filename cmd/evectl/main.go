package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/concur/rohr/cmd/evectl/command"
	"os"
)

func main() {
	if err := command.Execute(); err != nil {
		log.Fatalln(err)
		os.Exit(-1)
	}
}
