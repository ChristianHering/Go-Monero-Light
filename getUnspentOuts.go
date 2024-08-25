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

type GetUnspentOutsRequest struct {
	Address       string `json:"address"`
	ViewKey       string `json:"view_key"` // hex encoded binary
	Amount        string `json:"amount"`
	Mixin         uint32 `json:"mixin"`
	UseDust       bool   `json:"use_dust"`
	DustThreshold string `json:"dust_threshold"`
}

type GetUnspentOutsResponse struct {
	PerByteFee string   `json:"per_byte_fee"`
	FeeMask    string   `json:"fee_mask"`
	Amount     string   `json:"amount"`
	Outputs    []Output `json:"outputs"`
}

type Output struct {
	TxID           uint64   `json:"tx_id"`
	Amount         string   `json:"amount"`
	Index          uint16   `json:"index"`
	GlobalIndex    string   `json:"global_index"`
	RingCT         string   `json:"rct"`              // hex encoded binary
	TxHash         string   `json:"tx_hash"`          // hex encoded binary
	TxPrefixHash   string   `json:"tx_prefix_hash"`   // hex encoded binary
	PublicKey      string   `json:"public_key"`       // hex encoded binary
	TxPublicKey    string   `json:"tx_pub_key"`       // hex encoded binary
	SpendKeyImages []string `json:"spend_key_images"` // hex encoded binary elements
	Timestamp      string   `json:"timestamp"`        // Time in the format: "YYYY-HH-MM-SS.0-00:00"
	Height         uint64   `json:"height"`
}

var ErrorGetUnspentOutsRequestEncode = errors.New("failed to encode GetUnspentOutsRequest using data from 'client' and 'request'")

func (c *Client) GetUnspentOuts(request *GetUnspentOutsRequest) (*GetUnspentOutsResponse, error) {
	const path = "/get_unspent_outs"

	b := new(bytes.Buffer)

	request.Address = c.address
	request.ViewKey = c.viewKey

	err := json.NewEncoder(b).Encode(request)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "In call to GetUnspentOutsResponse(), failed to encode the following data:\n%#v\n", *request)
		if err != nil {
			panic(err)
		}

		return &GetUnspentOutsResponse{}, ErrorGetUnspentOutsRequestEncode
	}

	url, err := url.JoinPath(c.serverURL, path)
	if err != nil {
		printToStderr("Failed to join " + c.serverURL + " and " + path + " in call to GetUnspentOuts().")

		return &GetUnspentOutsResponse{}, err
	}

	retries := 0

POST_REQUEST:

	resp, err := c.client.Post(url, "application/json", bytes.NewReader(b.Bytes()))
	if err != nil {
		printToStderr("Failed to post the following data to " + url + ": " + b.String())

		return &GetUnspentOutsResponse{}, err
	}

	if resp.StatusCode == http.StatusServiceUnavailable {
		if retries < c.retryCount {
			time.Sleep(c.retryTime)

			goto POST_REQUEST
		}

		return &GetUnspentOutsResponse{}, ErrorServiceUnavailable
	} else if resp.StatusCode != http.StatusOK {
		return &GetUnspentOutsResponse{}, ErrorStatusCodeNotOK
	}

	var response = &GetUnspentOutsResponse{}

	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		printToStderr("Failed to decode /get_unspent_outs response.")

		return &GetUnspentOutsResponse{}, err
	}

	return response, nil
}
