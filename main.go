package main

import (
	"log"
	"os"

	"github.com/moshebe/gtrace/internal/cli"
	"github.com/moshebe/gtrace/internal/version"
)

func main() {
	if err := cli.App(version.Name()).Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
