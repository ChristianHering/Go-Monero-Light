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

type GetAddressTxsResponse struct {
	TotalReceived      string        `json:"total_received"`
	ScannedHeight      uint64        `json:"scanned_height"`
	ScannedBlockHeight uint64        `json:"scanned_block_height"`
	StartHeight        uint64        `json:"start_height"`
	BlockchainHeight   uint64        `json:"blockchain_height"`
	Transactions       []Transaction `json:"transactions"`
}

func (c *Client) GetAddressTxs() (*GetAddressTxsResponse, error) {
	const path = "/get_address_txs"

	b := new(bytes.Buffer)

	reqObj := &StandardRequest{
		Address: c.address,
		ViewKey: c.viewKey,
	}

	err := json.NewEncoder(b).Encode(reqObj)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "In call to GetAddressTxsResponse(), failed to encode the following data:\n%#v\n", *reqObj)
		if err != nil {
			panic(err)
		}

		return &GetAddressTxsResponse{}, ErrorStandardRequestEncode
	}

	url, err := url.JoinPath(c.serverURL, path)
	if err != nil {
		printToStderr("Failed to join " + c.serverURL + " and " + path + " in call to GetAddressTxs().")

		return &GetAddressTxsResponse{}, err
	}

	retries := 0

POST_REQUEST:

	resp, err := c.client.Post(url, "application/json", bytes.NewReader(b.Bytes()))
	if err != nil {
		printToStderr("Failed to post the following data to " + url + ": " + b.String())

		return &GetAddressTxsResponse{}, err
	}

	if resp.StatusCode == http.StatusServiceUnavailable {
		if retries < c.retryCount {
			time.Sleep(c.retryTime)

			goto POST_REQUEST
		}

		return &GetAddressTxsResponse{}, ErrorServiceUnavailable
	} else if resp.StatusCode != http.StatusOK {
		return &GetAddressTxsResponse{}, ErrorStatusCodeNotOK
	}

	var response = &GetAddressTxsResponse{}

	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		printToStderr("Failed to decode /get_address_txs response.")

		return &GetAddressTxsResponse{}, err
	}

	return response, nil
}
