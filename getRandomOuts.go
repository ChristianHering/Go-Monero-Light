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

// GetRandomOutsRequest holds request data for GetRandomOuts()
//
// Amounts represents the XMR amounts that need mixing. Clients
// should take care when making several ring signatures. See:
// https://github.com/monero-project/meta/blob/master/api/lightwallet_rest.md#get_random_outs
type GetRandomOutsRequest struct {
	Count   uint32   `json:"count"`
	Amounts []string `json:"amounts"`
}

// GetRandomOutsResponse
type GetRandomOutsResponse struct {
	AmountOuts []RandomOutputs `json:"amount_outs"`
}

type RandomOutputs struct {
	Amount  string         `json:"amount"`
	Outputs []RandomOutput `json:"outputs"`
}

type RandomOutput struct {
	GlobalIndex string `json:"global_index"`
	PublicKey   string `json:"public_key"`
	RingCT      string `json:"rct"`
}

var ErrorRandomOutsRequestEncode = errors.New("failed to encode random outs request using data from 'request' and 'client'")

// GetRandomOuts selects random outputs to be
// used for a ring signature in a new transaction.
func (c *Client) GetRandomOuts(request *GetRandomOutsRequest) (*GetRandomOutsResponse, error) {
	const path = "/get_random_outs"

	b := new(bytes.Buffer)

	err := json.NewEncoder(b).Encode(request)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to encode:\n\n%#v\n\nwith error:\n%v\n\n", *request, err)

		return &GetRandomOutsResponse{}, ErrorRandomOutsRequestEncode
	}

	url, err := url.JoinPath(c.serverURL, path)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to join:\n%s\nand\n%s\nwith error:\n\n%v\n\n", c.serverURL, path, err)

		return &GetRandomOutsResponse{}, ErrorJoinPathFailed
	}

	retries := 0

POST_REQUEST:

	resp, err := c.client.Post(url, "application/json", bytes.NewReader(b.Bytes()))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to post:\n\n%s\n\n to our endpoint at:\n\n%s\n\nwith error:\n\n%v\n\n", b.String(), url, err)

		return &GetRandomOutsResponse{}, ErrorPostRequestFailed
	}

	if resp.StatusCode == http.StatusServiceUnavailable {
		if retries < c.retryCount {
			time.Sleep(c.retryTime)

			goto POST_REQUEST
		}

		return &GetRandomOutsResponse{}, ErrorServiceUnavailable
	} else if resp.StatusCode != http.StatusOK {
		return &GetRandomOutsResponse{}, ErrorStatusCodeNotOK
	}

	var response = &GetRandomOutsResponse{}

	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to decode:\n\n%#v\n\nwith error:\n%v\n", response, err)

		return &GetRandomOutsResponse{}, ErrorResponseUnmarshalFailed
	}

	return response, nil
}
