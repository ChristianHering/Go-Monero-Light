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

type SubmitRawTxRequest struct {
	Tx string `json:"tx"`
}

type SubmitRawTxResponse struct {
	Status string `json:"status"`
}

var ErrorSubmitRawTxRequestEncode = errors.New("failed to encode SubmitRawTxRequest using data from 'request'")

func (c *Client) SubmitRawTx(request *SubmitRawTxRequest) (*SubmitRawTxResponse, error) {
	const path = "/submit_raw_tx"

	b := new(bytes.Buffer)

	err := json.NewEncoder(b).Encode(request)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "In call to SubmitRawTxResponse(), failed to encode the following data:\n%#v\n", *request)
		if err != nil {
			panic(err)
		}

		return &SubmitRawTxResponse{}, ErrorSubmitRawTxRequestEncode
	}

	url, err := url.JoinPath(c.serverURL, path)
	if err != nil {
		printToStderr("Failed to join " + c.serverURL + " and " + path + " in call to SubmitRawTx().")

		return &SubmitRawTxResponse{}, err
	}

	retries := 0

POST_REQUEST:

	resp, err := c.client.Post(url, "application/json", bytes.NewReader(b.Bytes()))
	if err != nil {
		printToStderr("Failed to post the following data to " + url + ": " + b.String())

		return &SubmitRawTxResponse{}, err
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
		printToStderr("Failed to decode /submit_raw_tx response.")

		return &SubmitRawTxResponse{}, err
	}

	return response, nil
}
