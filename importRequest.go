// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright © 2024 Christian Hering

package gomonerolight

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

type ImportRequestResponse struct {
	PaymentAddress   string `json:"payment_address"`
	PaymentID        string `json:"payment_id"`
	ImportFee        string `json:"import_fee"`
	NewRequest       bool   `json:"new_request"`
	RequestFulfilled bool   `json:"request_fulfilled"`
	Status           string `json:"status"`
}

func (c *Client) ImportRequest() (*ImportRequestResponse, error) {
	const path = "/import_request"

	b := new(bytes.Buffer)

	reqObj := &StandardRequest{
		Address: c.address,
		ViewKey: c.viewKey,
	}

	err := json.NewEncoder(b).Encode(reqObj)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "In call to ImportRequestResponse(), failed to encode the following data:\n%#v\n", *reqObj)
		if err != nil {
			panic(err)
		}

		return &ImportRequestResponse{}, ErrorStandardRequestEncode
	}

	url, err := url.JoinPath(c.serverURL, path)
	if err != nil {
		printToStderr("Failed to join " + c.serverURL + " and " + path + " in call to ImportRequest().")

		return &ImportRequestResponse{}, err
	}

	retries := 0

POST_REQUEST:

	resp, err := c.client.Post(url, "application/json", bytes.NewReader(b.Bytes()))
	if err != nil {
		printToStderr("Failed to post the following data to " + url + ": " + b.String())

		return &ImportRequestResponse{}, err
	}

	if resp.StatusCode == http.StatusServiceUnavailable {
		if retries < c.retryCount {
			time.Sleep(c.retryTime)

			goto POST_REQUEST
		}

		return &ImportRequestResponse{}, ErrorServiceUnavailable
	} else if resp.StatusCode != http.StatusOK {
		return &ImportRequestResponse{}, ErrorStatusCodeNotOK
	}

	var response = &ImportRequestResponse{}

	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		printToStderr("Failed to decode /import_request response.")

		return &ImportRequestResponse{}, err
	}

	return response, nil
}
