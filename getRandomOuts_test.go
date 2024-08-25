// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright © 2024 Christian Hering

package gomonerolight

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestGetRandomOuts(t *testing.T) {
	tryCount := 1 //Number of times to send HTTP Service Unavailable

	request := &GetRandomOutsRequest{
		Count:   0,
		Amounts: []string{"0"},
	}

	response := GetRandomOutsResponse{
		AmountOuts: []RandomOutput{{
			GlobalIndex: "1",
			PublicKey:   "915CA54F54AA5545945567554CD552ACB3CAA946",
			RingCT:      "915CA54F54AA5545945567554CD552ACB3CAA946",
		}},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		var req = &GetRandomOutsRequest{}

		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			t.Error("GetRandomOuts() made an invalid request: ", err)
		}

		if reflect.DeepEqual(req, request) != true {
			t.Error("req struct didn't match the original data in request")
		}

		if tryCount != 0 {
			tryCount--

			w.WriteHeader(http.StatusServiceUnavailable)

			_, err := w.Write([]byte{})
			if err != nil {
				t.Error("failed to write HTTP status code 503")
			}

			return
		}

		b := new(bytes.Buffer)

		err = json.NewEncoder(b).Encode(response)
		if err != nil {
			t.Error("Failed to marshal our response!")
		}

		w.Write(b.Bytes())
		if err != nil {
			t.Error("Failed to marshal our response!")
		}
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	client := &Client{
		address:    "xmr_address",
		client:     &http.Client{},
		retryCount: tryCount,
		retryTime:  time.Duration(0),
		serverURL:  ts.URL,
		viewKey:    "xmr_view_key",
	}

	resp, err := client.GetRandomOuts(request)
	if err != nil {
		t.Error("GetRandomOuts() returned the error: ", err)
	}

	if reflect.DeepEqual(*resp, response) != true {
		t.Error("response struct didn't match the original data")
	}
}
