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

func (c *Client) GetAddressInfo() (*GetAddressInfoResponse, error) {
	const path = "/get_address_info"

	b := new(bytes.Buffer)
	retries := 0

	reqObj := &StandardRequest{
		Address: c.address,
		ViewKey: c.viewKey,
	}

	err := json.NewEncoder(b).Encode(reqObj)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "In call to GetAddressInfoResponse(), failed to encode the following data:\n%#v\n", *reqObj)
		if err != nil {
			panic(err)
		}

		return &GetAddressInfoResponse{}, ErrorStandardRequestEncode
	}

	url, err := url.JoinPath(c.serverURL, path)
	if err != nil {
		printToStderr("Failed to join " + c.serverURL + " and " + path + " in call to GetAddressInfo().")

		return &GetAddressInfoResponse{}, err
	}

POST_REQUEST:

	resp, err := c.client.Post(url, "application/json", bytes.NewReader(b.Bytes()))
	if err != nil {
		printToStderr("Failed to post the following data to " + url + ": " + b.String())

		return &GetAddressInfoResponse{}, err
	}

	if resp.StatusCode == http.StatusServiceUnavailable {
		if retries < c.retryCount {
			time.Sleep(c.retryTime)

			goto POST_REQUEST
		}

		return &GetAddressInfoResponse{}, ErrorServiceUnavailable
	}

	var response = &GetAddressInfoResponse{}

	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		printToStderr("Failed to decode /get_address_info response.")

		return &GetAddressInfoResponse{}, err
	}

	return response, nil
}
