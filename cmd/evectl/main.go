package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/scipian/eve/cmd/evectl/command"
)

func main() {
	cmd := command.NewRootCommand(os.Stdin, os.Stdout, os.Stderr)
	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
		os.Exit(-1)
	}
}
