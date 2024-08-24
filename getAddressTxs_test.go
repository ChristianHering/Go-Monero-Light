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

func TestGetAddressTxs(t *testing.T) {
	tryCount := 1 //Number of times to send HTTP Service Unavailable

	response := GetAddressTxsResponse{
		TotalReceived:      "31415926535897",
		ScannedHeight:      3222370,
		ScannedBlockHeight: 3222370,
		StartHeight:        3222370,
		BlockchainHeight:   3222370,
		Transactions: []Transaction{{
			ID:            7,
			Hash:          "a70d679d2052f752732659680f27afe54b83686866c906cb5b4d9c91ce65a942",
			Timestamp:     time.Time{},
			TotalReceived: "31415926535897",
			TotalSent:     "31415926535900",
			UnlockTime:    3222370,
			Height:        3222370,
			SpentOutputs: []Spend{{
				Amount:      "31415926535897",
				KeyImage:    "A555554CD552ACB3554CAACA94914F54A5567946",
				TxPublicKey: "a06c7f33eb578148a578167babf3367f87a43ee601d9feda98c803c135c0b506", // Don't send XMR here
				OutIndex:    0,
				Mixin:       4,
			}},
			PaymentID: "",
			Coinbase:  false,
			Mempool:   false,
			Mixin:     0,
		}},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		var request = &StandardRequest{}

		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil {
			t.Error("GetAddressTxs() made an invalid request: ", err)
		}

		if request.Address != "xmr_address" {
			t.Error("GetAddressTxs() mangled our XMR address: " + request.Address + " != xmr_address")
		}
		if request.ViewKey != "xmr_view_key" {
			t.Error("GetAddressTxs() mangled our view key: " + request.ViewKey + " != xmr_view_key")
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

	resp, err := client.GetAddressTxs()
	if err != nil {
		t.Error("GetAddressTxs() returned the error: ", err)
	}

	if reflect.DeepEqual(*resp, response) != true {
		t.Error("Response struct didn't match the original data.")
	}
}
