package cli

import (
	"io"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
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

func stringSlice(c *cli.Context, name string) []string {
	var results []string
	for _, v := range c.StringSlice(name) {
		results = append(results, strings.Split(v, ",")...)
	}
	return results
}

func read(path string) ([]byte, error) {
	if path == "-" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(path)
}

func writer(path string) (io.WriteCloser, error) {
	if path == "-" {
		return os.Stdout, nil
	}
	return os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0660)
}
