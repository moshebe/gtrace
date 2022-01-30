package main

import (
	"log"
	"os"

	"github.com/moshebe/gtrace/internal/cli"
)

func main() {
	err := cli.App.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
