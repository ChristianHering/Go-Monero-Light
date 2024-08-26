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

// GetAddressTxsResponse holds an array of candidate spend events
// that can be used to get an account's transaction history.
type GetAddressTxsResponse struct {
	TotalReceived      string        `json:"total_received"`
	ScannedHeight      uint64        `json:"scanned_height"`
	ScannedBlockHeight uint64        `json:"scanned_block_height"`
	StartHeight        uint64        `json:"start_height"`
	BlockchainHeight   uint64        `json:"blockchain_height"`
	Transactions       []Transaction `json:"transactions"`
}

// GetAddressTxs returns candidate spends to show transaction history
//
// Your Monero spend key is required to calculate if a candidate spend
// was an actual spend so it only returns candidate spend events and
// leaves the calculation for the client.
func (c *Client) GetAddressTxs() (*GetAddressTxsResponse, error) {
	const path = "/get_address_txs"

	b := new(bytes.Buffer)

	request := &StandardRequest{
		Address: c.address,
		ViewKey: c.viewKey,
	}

	err := json.NewEncoder(b).Encode(request)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to encode:\n\n%#v\n\nwith error:\n%v\n\n", *request, err)

		return &GetAddressTxsResponse{}, ErrorStandardRequestEncode
	}

	url, err := url.JoinPath(c.serverURL, path)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to join:\n%s\nand\n%s\nwith error:\n\n%v\n\n", c.serverURL, path, err)

		return &GetAddressTxsResponse{}, ErrorJoinPathFailed
	}

	retries := 0

POST_REQUEST:

	resp, err := c.client.Post(url, "application/json", bytes.NewReader(b.Bytes()))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to post:\n\n%s\n\n to our endpoint at:\n\n%s\n\nwith error:\n\n%v\n\n", b.String(), url, err)

		return &GetAddressTxsResponse{}, ErrorPostRequestFailed
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
		_, _ = fmt.Fprintf(os.Stderr, "failed to decode:\n\n%#v\n\nwith error:\n%v\n", response, err)

		return &GetAddressTxsResponse{}, ErrorResponseUnmarshalFailed
	}

	return response, nil
}
