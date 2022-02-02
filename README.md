<p align="center">
    <img src="assets/gtrace.svg" height="250" width="250"/>
</p>

</br>

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/moshebe/gtrace)](https://pkg.go.dev/github.com/moshebe/gtrace)

Unofficial, simple yet effective Google Cloud Trace CLI tool.

</br></br>
# Installation
### [Homebrew](https://brew.sh/) (Linux/macOS)
```shell
brew install moshebe/pkg/gtrace
```
### [Go](https://golang.org) (Linux/Windows/macOS/any other platform supported by Go)
If you have Go 1.16+, you can install latest released version of `gtrace` directly from source by running:
```shell
go install github.com/moshebe/gtrace@latest
```

# Usage
```shell
› gtrace help
NAME:
   gtrace - Google Cloud Trace CLI tool

USAGE:
   Simple command-line tool for query and fetch tracing information from Cloud Trace API.
   Find more information at: https://cloud.google.com/trace/docs

VERSION:
   v1.0.0

COMMANDS:
   get      Get a specific trace by id from one or more projects
   list     Query traces from a project according to the given conditions
   url      Generate a browsable URL for a given trace
   format   Format trace spans according to a given template
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

# Authentication

Google Cloud APIs has a few methods to authenticate:
1. On your machine maka sure you have `gcloud` and that you are logged-in and then run `gcloud auth application-default login`
1. If you are running on GKE you can use Workload Identity
1. Service accounts with JSON keys
1. Pointing to your service account via env `GOOGLE_APPLICATION_CREDENTIALS`

> You can read about it more on: https://cloud.google.com/docs/authentication/getting-started

# Examples

Fetch a specific trace from multiple projects:
```shell
gtrace get --project production-a,production-b 5e26a889fa12da351beee9ea16ce0a65
```

Format trace spans by a specific template:
```shell
gtrace format -f /tmp/trace.json --template "{{ .Name }}, {{ .Duration }}"
```

Query traces by multiple filters from the last 3 hours:
```shell
gtrace list --project dev --limit 10 --since 3h --filter service:api --filter user-id:1234
```

