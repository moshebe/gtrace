package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/moshebe/gtrace/internal/cli"
)

var (
	version = "dev"
	commit  = "none"
)

func versionName() string {
	if !strings.Contains(version, "-") {
		return version
	}
	return fmt.Sprintf("%s (%s)", version, commit)
}

func main() {
	err := cli.App(versionName()).Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
