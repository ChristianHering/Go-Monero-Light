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

// LoginRequest holds the information needed for calling /login.
//
// However, the elements "Address" and "ViewKey" are passed in
// from client and are not needed in calls to Login().
//
// CreateAccount and GeneratedLocally define
type LoginRequest struct {
	Address          string `json:"address"`
	ViewKey          string `json:"view_key"` // hex encoded binary
	CreateAccount    bool   `json:"create_account"`
	GeneratedLocally bool   `json:"generated_locally"`
}

// LoginResponse is what you get back from a call to /login.
//
// NewAddress lets you know if you're a new address to the server
// and if so, you *may* receive HTTP 403 (Forbidden) status codes
// until your account is manually reviewed.
//
// GeneratedLocally and StartHeight are optional.
type LoginResponse struct {
	NewAddress       bool   `json:"new_address"`
	GeneratedLocally bool   `json:"generated_locally"`
	StartHeight      uint64 `json:"start_height"`
}

var ErrorLoginRequestEncode = errors.New("failed to encode login request using data from 'client' and 'request'")

// Login checks for the existance of an account
// on the server or, optionally, creates one.
//
// When calling Login, pass a LoginRequest struct*
// that doesn't have Address/ViewKey fields set.
// They will be overwritten with the values set
// for your client 'c'.
func (c *Client) Login(request *LoginRequest) (*LoginResponse, error) {
	const path = "/login"

	b := new(bytes.Buffer)

	request.Address = c.address
	request.ViewKey = c.viewKey

	err := json.NewEncoder(b).Encode(request)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to encode:\n\n%#v\n\nwith error:\n%v\n\n", *request, err)

		return &LoginResponse{}, ErrorLoginRequestEncode
	}

	url, err := url.JoinPath(c.serverURL, path)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to join:\n%s\nand\n%s\nwith error:\n\n%v\n\n", c.serverURL, path, err)

		return &LoginResponse{}, ErrorJoinPathFailed
	}

	retries := 0

POST_REQUEST:

	resp, err := c.client.Post(url, "application/json", bytes.NewReader(b.Bytes()))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to post:\n\n%s\n\n to our endpoint at:\n\n%s\n\nwith error:\n\n%v\n\n", b.String(), url, err)

		return &LoginResponse{}, ErrorPostRequestFailed
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
		_, _ = fmt.Fprintf(os.Stderr, "failed to decode:\n\n%#v\n\nwith error:\n%v\n", response, err)

		return &LoginResponse{}, ErrorResponseUnmarshalFailed
	}

	return response, nil
}
