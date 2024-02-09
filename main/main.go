package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"postfix_to_cf/api"
	"postfix_to_cf/config"
	"postfix_to_cf/mail"
)

var (
	// Version will be injected at build time using -ldflags
	Version = "development"
)

func main() {
	cfg := config.New()

	if cfg.ShowVersion {
		fmt.Println(Version)
		return
	}

	if !cfg.ParseOnly && (cfg.Token == "" || cfg.Endpoint == "") {
		fmt.Fprintln(os.Stderr, "Error: Token and Endpoint are required when not in parse-only mode")
		os.Exit(1)
	}

	emailData, err := mail.ParseEmailFromStdin()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing email: %v\n", err)
		os.Exit(1)
	}

	jsonPayload, err := json.Marshal(emailData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling email data to JSON: %v\n", err)
		os.Exit(1)
	}

	if cfg.ParseOnly {
		fmt.Println(string(jsonPayload))
		return
	}

	// Create a standard http.Client
	httpClient := &http.Client{}

	err = api.SendPostRequest(httpClient, cfg.Endpoint, cfg.Token, string(jsonPayload))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to send POST request: %s\n", err)
		os.Exit(1)
	}
}
