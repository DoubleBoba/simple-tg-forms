//go:build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/sh"
)

// Deploy function to Google Cloud
func Deploy() error {
	project_id, ok := os.LookupEnv("PROJECT_ID")
	if !ok {
		return fmt.Errorf("PROJECT_ID must be set")
	}
	bot_token, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		return fmt.Errorf("BOT_TOKEN must be set")
	}
	webhook_token, ok := os.LookupEnv("WEBHOOK_TOKEN")
	if !ok {
		return fmt.Errorf("WEBHOOK_TOKEN must be set")
	}

	env_vars_string := fmt.Sprintf(
		"--update-env-vars=PROJECT_ID=%s,BOT_TOKEN=%s,WEBHOOK_SECRET_TOKEN=%s,DEPLOYED=true",
		project_id, bot_token, webhook_token,
	)
	err := sh.RunV(
		"gcloud", "functions", "deploy", "go-http-function",
		"--gen2", "--runtime=go119", "--region=europe-west1",
		"--source=tg-bot/", "--entry-point=process-tg-update", "--trigger-http",
		"--memory=128Mi", "--allow-unauthenticated", env_vars_string,
	)
	return err
}

// Run go tests
func Test() error {
	return sh.RunV("go", "test")
}
