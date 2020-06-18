package main

import (
	"log"
	"os"

	"github.com/ryotarai/smock/pkg/command"
)

func main() {
	app, err := command.New().App()
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
