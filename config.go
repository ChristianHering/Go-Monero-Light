// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright © 2024 Christian Hering

package gomonerolight

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

var ErrorBadConfig = errors.New("configuration options passed to NewClient were invalid")

type Config struct {
	Address    string        // Your XMR address
	HTTPClient *http.Client  // For setting custom cookies, etc. Likely to remain unused.
	RetryCount int           // The number of times to retry a method call before giving up
	RetryTime  time.Duration // The time to wait in between retry requests
	ServerURL  string        // The URL of the API server (eg. https://api.mymonero.com)
	ViewKey    string        // Your XMR private view key
}

func checkConfig(cfg *Config) error {
	if cfg.Address == "" {
		// TODO: Generate a new, random address (and viewkey) if one is not provided

		printToStderr("no XMR address was passed to NewClient() call")

		return ErrorBadConfig
	}

	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{}
	}

	if cfg.ServerURL == "" {
		cfg.ServerURL = "https://api.mymonero.com" //Default to using MyMonero
	}

	if cfg.ViewKey == "" {
		_, err := fmt.Fprintln(os.Stderr, "No viewkey was passed to NewClient call")
		if err != nil {
			panic(err)
		}

		return ErrorBadConfig
	}

	return nil
}
