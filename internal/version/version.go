package version

import (
	"fmt"
	"strings"
)

var (
	Version = "dev"
	Commit  = "none"
)

func Name() string {
	if !strings.Contains(Version, "-") {
		return Version
	}
	return fmt.Sprintf("%s (%s)", Version, Commit)
}
