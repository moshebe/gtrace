package cli

import (
	"os"
	"strings"

	cliv2 "github.com/urfave/cli/v2"
)

const (
	createFileFlags = os.O_RDWR | os.O_CREATE
	createFilePerm  = 0660
)

var App = &cliv2.App{
	Name:      "gtrace",
	Version:   "v1.0.0",
	HelpName:  "gtrace",
	Usage:     "Google Cloud Trace CLI tool",
	UsageText: "Simple command-line tool for query and fetch tracing information from Cloud Trace API.\n   Find more information at: https://cloud.google.com/trace/docs",
	Commands: []*cliv2.Command{
		GetCommand,
		ListCommand,
		URLCommand,
		FormatCommand,
	},
}

func stdio(value string) bool {
	return value == "-" || value == ""
}

func stringSlice(c *cliv2.Context, name string) []string {
	var results []string
	for _, v := range c.StringSlice(name) {
		results = append(results, strings.Split(v, ",")...)
	}
	return results
}
