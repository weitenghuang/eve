package main

import (
	"github.com/concur/rohr/cmd/eve/command"
	"log"
	"os"
)

func main() {
	if err := command.Execute(); err != nil {
		log.Fatalln(err)
		os.Exit(-1)
	}
}
