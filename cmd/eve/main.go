package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/scipian/eve/cmd/eve/command"
	"os"
)

func main() {
	if err := command.Execute(); err != nil {
		log.Fatalln(err)
		os.Exit(-1)
	}
}
