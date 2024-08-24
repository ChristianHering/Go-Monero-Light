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

func TestGetAddressInfo(t *testing.T) {
	tryCount := 1 //Number of times to send HTTP Service Unavailable

	response := GetAddressInfoResponse{
		LockedFunds:        "1",
		TotalReceived:      "2",
		TotalSent:          "3",
		ScannedHeight:      4,
		ScannedBlockHeight: 5,
		StartHeight:        6,
		TransactionHeight:  7,
		BlockchainHeight:   8,
		SpentOutputs: []Spend{{
			Amount:      "1",
			KeyImage:    []byte{},
			TxPublicKey: []byte{},
			OutIndex:    2,
			Mixin:       3,
		}},
		ExchangeRates: Rates{
			AUD: 4.1,
			EUR: 4.2,
			GBP: 4.3,
			USD: 4.4,
			RUB: 4.5,
		},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		var request = &StandardRequest{}

		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil {
			t.Error("GetAddressInfo() made an invalid request: ", err)
		}

		if request.Address != "xmr_address" {
			t.Error("GetAddressInfo() mangled our XMR address: " + request.Address + " != xmr_address")
		}
		if request.ViewKey != "xmr_view_key" {
			t.Error("GetAddressInfo() mangled our view key: " + request.ViewKey + " != xmr_view_key")
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

	resp, err := client.GetAddressInfo()
	if err != nil {
		t.Error("GetAddressInfo() returned the error: ", err)
	}

	if reflect.DeepEqual(*resp, response) != true {
		t.Error("Response struct didn't match the original data.")
	}
}
