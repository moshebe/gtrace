package main

import (
	"log"
	"os"
	"time"

	"github.com/moshebe/gtrace/internal/cli"
	cliv2 "github.com/urfave/cli/v2"
)

func main() {
	app := &cliv2.App{
		Name:      "gtrace",
		Version:   "v1.0.0",
		Compiled:  time.Now(),
		HelpName:  "gtrace",
		Usage:     "demonstrate available API",
		UsageText: "gtrace - demonstrating the available API",
		ArgsUsage: "[args and such]",
		Commands: []*cliv2.Command{
			cli.GetCommand,
			cli.ListCommand,
			cli.URLCommand,
			cli.FormatCommand,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
