package main

import (
	"log"
	"os"

	"github.com/ryotarai/smock/pkg/command"
)

func main() {
	if err := command.New().App().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
