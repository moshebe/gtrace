package cli

import (
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

const (
	createFileFlags = os.O_RDWR | os.O_CREATE
	createFilePerm  = 0660
)

func App(version string) *cli.App {
	return &cli.App{
		Name:      "gtrace",
		Version:   version,
		HelpName:  "gtrace",
		Usage:     "Google Cloud Trace CLI tool",
		UsageText: "Simple command-line tool for query and fetch tracing information from Cloud Trace API.\n   Find more information at: https://cloud.google.com/trace/docs",
		Commands: []*cli.Command{
			GetCommand,
			ListCommand,
			URLCommand,
			FormatCommand,
		},
	}
}

func stdio(value string) bool {
	return value == "-" || value == ""
}

func stringSlice(c *cli.Context, name string) []string {
	var results []string
	for _, v := range c.StringSlice(name) {
		results = append(results, strings.Split(v, ",")...)
	}
	return results
}
