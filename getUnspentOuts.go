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

// GetUnspentOutsRequest holds a request for GetUnspentOuts().
//
// It is not required to pass Address or ViewKey, as those are
// derived from 'client'. Additionally, if the value given for
// Amount is greater than the total received outputs for our
// account, the server will respond with HTTP 400 (Bad Request)
type GetUnspentOutsRequest struct {
	Address       string `json:"address"`
	ViewKey       string `json:"view_key"` // hex encoded binary
	Amount        string `json:"amount"`
	Mixin         uint32 `json:"mixin"`
	UseDust       bool   `json:"use_dust"`
	DustThreshold string `json:"dust_threshold"`
}

// GetUnspentOutsResponse is the response from a call to GetUnspentOuts().
//
// It holds the total value of all the outputs in Outputs as
// well as the actual data for each output in our Outputs slice.
type GetUnspentOutsResponse struct {
	PerByteFee string   `json:"per_byte_fee"`
	FeeMask    string   `json:"fee_mask"`
	Amount     string   `json:"amount"`
	Outputs    []Output `json:"outputs"`
}

// Output represents a single monero output.
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

// GetUnspentOuts gets a list of received outputs.
//
// It does not return or distinguish when outputs were spent.
func (c *Client) GetUnspentOuts(request *GetUnspentOutsRequest) (*GetUnspentOutsResponse, error) {
	const path = "/get_unspent_outs"

	b := new(bytes.Buffer)

	request.Address = c.address
	request.ViewKey = c.viewKey

	err := json.NewEncoder(b).Encode(request)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to encode:\n\n%#v\n\nwith error:\n%v\n\n", *request, err)

		return &GetUnspentOutsResponse{}, ErrorGetUnspentOutsRequestEncode
	}

	url, err := url.JoinPath(c.serverURL, path)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to join:\n%s\nand\n%s\nwith error:\n\n%v\n\n", c.serverURL, path, err)

		return &GetUnspentOutsResponse{}, ErrorJoinPathFailed
	}

	retries := 0

POST_REQUEST:

	resp, err := c.client.Post(url, "application/json", bytes.NewReader(b.Bytes()))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to post:\n\n%s\n\n to our endpoint at:\n\n%s\n\nwith error:\n\n%v\n\n", b.String(), url, err)

		return &GetUnspentOutsResponse{}, ErrorPostRequestFailed
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
		_, _ = fmt.Fprintf(os.Stderr, "failed to decode:\n\n%#v\n\nwith error:\n%v\n", response, err)

		return &GetUnspentOutsResponse{}, ErrorResponseUnmarshalFailed
	}

	return response, nil
}
