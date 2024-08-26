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

// GetAddressInfoResponse holds the information to calculate
// a wallet's balance using our spend key.
//
// ExchangeRates is optional and may not be sent by the server.
type GetAddressInfoResponse struct {
	LockedFunds        string  `json:"locked_funds"`
	TotalReceived      string  `json:"total_received"`
	TotalSent          string  `json:"total_sent"`
	ScannedHeight      uint64  `json:"scanned_height"`
	ScannedBlockHeight uint64  `json:"scanned_block_height"`
	StartHeight        uint64  `json:"start_height"`
	TransactionHeight  uint64  `json:"transaction_height"`
	BlockchainHeight   uint64  `json:"blockchain_height"`
	SpentOutputs       []Spend `json:"spent_outputs"`
	ExchangeRates      Rates   `json:"rates"`
}

// Rates lists XMR exchange rates for common fiat currencies
type Rates struct {
	AUD float32 `json:"AUD"`
	BRL float32 `json:"BRL"`
	BTC float32 `json:"BTC"`
	CAD float32 `json:"CAD"`
	CHF float32 `json:"CHF"`
	CNY float32 `json:"CNY"`
	EUR float32 `json:"EUR"`
	GBP float32 `json:"GBP"`
	HKD float32 `json:"HKD"`
	INR float32 `json:"INR"`
	JPY float32 `json:"JPY"`
	KRW float32 `json:"KRW"`
	MXN float32 `json:"MXN"`
	NOK float32 `json:"NOK"`
	NZD float32 `json:"NZD"`
	SEK float32 `json:"SEK"`
	SGD float32 `json:"SGD"`
	TRY float32 `json:"TRY"`
	USD float32 `json:"USD"`
	RUB float32 `json:"RUB"`
	ZAR float32 `json:"ZAR"`
}

// GetAddressInfo gets the information to calculate a wallet's balance
//
// The server returns candidate spends that can be used to calculate
// a wallet's balance using our spend key.
func (c *Client) GetAddressInfo() (*GetAddressInfoResponse, error) {
	const path = "/get_address_info"

	b := new(bytes.Buffer)

	request := &StandardRequest{
		Address: c.address,
		ViewKey: c.viewKey,
	}

	err := json.NewEncoder(b).Encode(request)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to encode:\n\n%#v\n\nwith error:\n%v\n\n", *request, err)

		return &GetAddressInfoResponse{}, ErrorStandardRequestEncode
	}

	url, err := url.JoinPath(c.serverURL, path)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to join:\n%s\nand\n%s\nwith error:\n\n%v\n\n", c.serverURL, path, err)

		return &GetAddressInfoResponse{}, ErrorJoinPathFailed
	}

	retries := 0

POST_REQUEST:

	resp, err := c.client.Post(url, "application/json", bytes.NewReader(b.Bytes()))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to post:\n\n%s\n\n to our endpoint at:\n\n%s\n\nwith error:\n\n%v\n\n", b.String(), url, err)

		return &GetAddressInfoResponse{}, ErrorPostRequestFailed
	}

	if resp.StatusCode == http.StatusServiceUnavailable {
		if retries < c.retryCount {
			time.Sleep(c.retryTime)

			goto POST_REQUEST
		}

		return &GetAddressInfoResponse{}, ErrorServiceUnavailable
	} else if resp.StatusCode != http.StatusOK {
		return &GetAddressInfoResponse{}, ErrorStatusCodeNotOK
	}

	var response = &GetAddressInfoResponse{}

	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to decode:\n\n%#v\n\nwith error:\n%v\n", response, err)

		return &GetAddressInfoResponse{}, ErrorResponseUnmarshalFailed
	}

	return response, nil
}
