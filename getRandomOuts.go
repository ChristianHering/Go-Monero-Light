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

type GetRandomOutsRequest struct {
	Count   uint32   `json:"count"`
	Amounts []string `json:"amounts"`
}

type GetRandomOutsResponse struct {
	AmountOuts []RandomOutput `json:"amount_outs"`
}

type RandomOutput struct {
	GlobalIndex string `json:"global_index"`
	PublicKey   string `json:"public_key"`
	RingCT      string `json:"rct"`
}

var ErrorRandomOutsRequestEncode = errors.New("failed to encode random outs request using data from 'request' and 'client'")

func (c *Client) GetRandomOuts(request *GetRandomOutsRequest) (*GetRandomOutsResponse, error) {
	const path = "/get_random_outs"

	b := new(bytes.Buffer)

	err := json.NewEncoder(b).Encode(request)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "In call to GetRandomOuts(), failed to encode the following data:\n%#v\n", *request)
		if err != nil {
			panic(err)
		}

		return &GetRandomOutsResponse{}, ErrorRandomOutsRequestEncode
	}

	url, err := url.JoinPath(c.serverURL, path)
	if err != nil {
		printToStderr("Failed to join " + c.serverURL + " and " + path + " in call to GetRandomOuts().")

		return &GetRandomOutsResponse{}, err
	}

	retries := 0

POST_REQUEST:

	resp, err := c.client.Post(url, "application/json", bytes.NewReader(b.Bytes()))
	if err != nil {
		printToStderr("Failed to post the following data to " + url + ": " + b.String())

		return &GetRandomOutsResponse{}, err
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
		printToStderr("Failed to decode /get_random_outs response.")

		return &GetRandomOutsResponse{}, err
	}

	return response, nil
}
