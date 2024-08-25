// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright © 2024 Christian Hering

package gomonerolight

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestGetRandomOuts(t *testing.T) {
	tryCount := 1 //Number of times to send HTTP Service Unavailable

	response := GetRandomOutsResponse{
		AmountOuts: []RandomOutput{{
			GlobalIndex: "1",
			PublicKey:   "915CA54F54AA5545945567554CD552ACB3CAA946",
			RingCT:      "915CA54F54AA5545945567554CD552ACB3CAA946",
		}},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		var request = &GetRandomOutsRequest{}

		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil {
			t.Error("GetRandomOuts() made an invalid request: ", err)
		}

		if request.Count != 0 {
			t.Error("GetRandomOuts() mangled our mixin count: " + strconv.Itoa(int(request.Count)) + " != 0")
		}
		if len(request.Amounts) != 1 || request.Amounts[0] != "0" {
			_, err := fmt.Fprintf(os.Stderr, `GetRandomOuts() mangled our amounts:\n%#v != []string{"0"}\n`, request.Amounts)
			if err != nil {
				panic(err)
			}
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

	request := &GetRandomOutsRequest{
		Count:   0,
		Amounts: []string{"0"},
	}

	resp, err := client.GetRandomOuts(request)
	if err != nil {
		t.Error("GetRandomOuts() returned the error: ", err)
	}

	if reflect.DeepEqual(*resp, response) != true {
		t.Error("Response struct didn't match the original data.")
	}
}
