package main

import (
	"log"
	"os"

	"github.com/joaovicdsantos/rip/internal"
)

func main() {
	ripCli := internal.CreateRipCli()
	if err := ripCli.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
