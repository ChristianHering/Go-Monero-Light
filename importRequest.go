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

// ImportRequestResponse returns the result of our account rescan,
// its status, and more if the server charged for the wallet service.
//
// PaymentAddress, PaymentID, and ImportFee are optional responses and are
// typically returned if the client needs to pay to complete the request.
type ImportRequestResponse struct {
	PaymentAddress   string `json:"payment_address"`
	PaymentID        string `json:"payment_id"` // hex encoded binary
	ImportFee        string `json:"import_fee"`
	NewRequest       bool   `json:"new_request"`
	RequestFulfilled bool   `json:"request_fulfilled"`
	Status           string `json:"status"`
}

// ImportRequest requests a rescan for our
// account's address since Monero's genesis block.
func (c *Client) ImportRequest() (*ImportRequestResponse, error) {
	const path = "/import_request"

	b := new(bytes.Buffer)

	request := &StandardRequest{
		Address: c.address,
		ViewKey: c.viewKey,
	}

	err := json.NewEncoder(b).Encode(request)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to encode:\n\n%#v\n\nwith error:\n%v\n\n", *request, err)

		return &ImportRequestResponse{}, ErrorStandardRequestEncode
	}

	url, err := url.JoinPath(c.serverURL, path)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to join:\n%s\nand\n%s\nwith error:\n\n%v\n\n", c.serverURL, path, err)

		return &ImportRequestResponse{}, ErrorJoinPathFailed
	}

	retries := 0

POST_REQUEST:

	resp, err := c.client.Post(url, "application/json", bytes.NewReader(b.Bytes()))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to post:\n\n%s\n\n to our endpoint at:\n\n%s\n\nwith error:\n\n%v\n\n", b.String(), url, err)

		return &ImportRequestResponse{}, ErrorPostRequestFailed
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
		_, _ = fmt.Fprintf(os.Stderr, "failed to decode:\n\n%#v\n\nwith error:\n%v\n", response, err)

		return &ImportRequestResponse{}, ErrorResponseUnmarshalFailed
	}

	return response, nil
}
