package main

import (
	"log"
	"os"

	"github.com/moshebe/gtrace/internal/cli"
	cliv2 "github.com/urfave/cli/v2"
)

func main() {
	app := &cliv2.App{
		Name:      "gtrace",
		Version:   "v1.0.0",
		HelpName:  "gtrace",
		Usage:     "Google Cloud Trace CLI tool",
		UsageText: "Simple command-line tool for query and fetch tracing information from Cloud Trace API.\n   Find more information at: https://cloud.google.com/trace/docs",
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
