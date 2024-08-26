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

// SubmitRawTxRequest holds a raw (binary) Monero
// transaction that's been encoded as an ASCII string.
type SubmitRawTxRequest struct {
	Tx string `json:"tx"` // hex encoded binary
}

// SubmitRawTxResponse holds the status of a call to
// the light wallet server's monero daemon at /submit_raw_tx.
// This response is typically from the Monero daemon.
type SubmitRawTxResponse struct {
	Status string `json:"status"`
}

var ErrorSubmitRawTxRequestEncode = errors.New("failed to encode SubmitRawTxRequest using data from 'request'")

// SubmitRawTx sends a raw transaction to our XMR light
// wallet server so it can be relayed on the Monero network.
//
// In order to call it, you must supply a request struct*
// that has a raw transaction encoded into an ASCII string.
func (c *Client) SubmitRawTx(request *SubmitRawTxRequest) (*SubmitRawTxResponse, error) {
	const path = "/submit_raw_tx"

	b := new(bytes.Buffer)

	err := json.NewEncoder(b).Encode(request)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to encode:\n\n%#v\n\nwith error:\n%v\n\n", *request, err)

		return &SubmitRawTxResponse{}, ErrorSubmitRawTxRequestEncode
	}

	url, err := url.JoinPath(c.serverURL, path)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to join:\n%s\nand\n%s\nwith error:\n\n%v\n\n", c.serverURL, path, err)

		return &SubmitRawTxResponse{}, ErrorJoinPathFailed
	}

	retries := 0

POST_REQUEST:

	resp, err := c.client.Post(url, "application/json", bytes.NewReader(b.Bytes()))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to post:\n\n%s\n\n to our endpoint at:\n\n%s\n\nwith error:\n\n%v\n\n", b.String(), url, err)

		return &SubmitRawTxResponse{}, ErrorPostRequestFailed
	}

	if resp.StatusCode == http.StatusServiceUnavailable {
		if retries < c.retryCount {
			time.Sleep(c.retryTime)

			goto POST_REQUEST
		}

		return &SubmitRawTxResponse{}, ErrorServiceUnavailable
	} else if resp.StatusCode != http.StatusOK {
		return &SubmitRawTxResponse{}, ErrorStatusCodeNotOK
	}

	var response = &SubmitRawTxResponse{}

	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to decode:\n\n%#v\n\nwith error:\n%v\n", response, err)

		return &SubmitRawTxResponse{}, ErrorResponseUnmarshalFailed
	}

	return response, nil
}
