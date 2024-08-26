// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright © 2024 Christian Hering

package gomonerolight

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

type LoginRequest struct {
	Address          string `json:"address"`
	ViewKey          string `json:"view_key"`
	CreateAccount    bool   `json:"create_account"`
	GeneratedLocally bool   `json:"generated_locally"`
}

type LoginResponse struct {
	NewAddress       bool   `json:"new_address"`
	GeneratedLocally bool   `json:"generated_locally"`
	StartHeight      uint64 `json:"start_height"`
}

var ErrorLoginRequestEncode = errors.New("failed to encode login request using data from 'client' and 'request'")

func (c *Client) Login(request *LoginRequest) (*LoginResponse, error) {
	const path = "/login"

	b := new(bytes.Buffer)

	request.Address = c.address
	request.ViewKey = c.viewKey

	err := json.NewEncoder(b).Encode(request)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "In call to Login(), failed to encode the following data:\n%#v\n", *request)
		if err != nil {
			panic(err)
		}

		return &LoginResponse{}, ErrorLoginRequestEncode
	}

	url, err := url.JoinPath(c.serverURL, path)
	if err != nil {
		printToStderr("Failed to join " + c.serverURL + " and " + path + " in call to Login().")

		return &LoginResponse{}, err
	}

	retries := 0

POST_REQUEST:

	resp, err := c.client.Post(url, "application/json", bytes.NewReader(b.Bytes()))
	if err != nil {
		printToStderr("Failed to post the following data to " + url + ": " + b.String())

		return &LoginResponse{}, err
	}

	if resp.StatusCode == http.StatusServiceUnavailable {
		if retries < c.retryCount {
			time.Sleep(c.retryTime)

			goto POST_REQUEST
		}

		return &LoginResponse{}, ErrorServiceUnavailable
	} else if resp.StatusCode != http.StatusOK {
		return &LoginResponse{}, ErrorStatusCodeNotOK
	}

	var response = &LoginResponse{}

	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		printToStderr("Failed to decode /login response.")

		return &LoginResponse{}, err
	}

	return response, nil
}
