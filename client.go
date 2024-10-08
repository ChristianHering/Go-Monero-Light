// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright © 2024 Christian Hering

package gomonerolight

import (
	"net/http"
	"time"
)

type Client struct {
	address    string
	client     *http.Client
	retryCount int
	retryTime  time.Duration
	serverURL  string
	viewKey    string
}

// NewClient creates a new client using the
// given Config 'cfg'. After calling NewClient()
// and getting a client 'c', call c.Login()
// and then the subsequent methods you need.
func NewClient(cfg Config) (*Client, error) {
	c := &Client{}

	err := checkConfig(&cfg)
	if err != nil {
		return nil, err
	}

	c.address = cfg.Address
	c.client = cfg.HTTPClient
	c.retryCount = cfg.RetryCount
	c.retryTime = cfg.RetryTime
	c.serverURL = cfg.ServerURL
	c.viewKey = cfg.ViewKey

	return c, nil
}
