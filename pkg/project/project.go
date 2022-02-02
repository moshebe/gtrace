package project

import (
	"context"
	"os"

	"golang.org/x/oauth2/google"
)

// FromEnv attempts to find active project id from environment variables.
// The environment variables are taken from this page: https://cloud.google.com/functions/docs/configuring/env-var
// If no project found, returns empty string.
func FromEnv() string {
	envs := []string{
		"GOOGLE_CLOUD_PROJECT",
		"GCP_PROJECT",
	}

	for _, e := range envs {
		project := os.Getenv(e)
		if project != "" {
			return project
		}
	}

	return ""
}

// FromApplicationDefault attempts to find active project id following the command
// `gcloud auth application-default login`
// https://cloud.google.com/sdk/gcloud/reference/auth/application-default/login
func FromApplicationDefault(ctx context.Context) string {
	credentials, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		return ""
	}

	return credentials.ProjectID
}
