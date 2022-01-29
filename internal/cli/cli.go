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
