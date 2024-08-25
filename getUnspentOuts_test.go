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

func TestGetUnspentOuts(t *testing.T) {
	tryCount := 1 //Number of times to send HTTP Service Unavailable

	request := &GetUnspentOutsRequest{
		Address:       "xmr_address",
		ViewKey:       "xmr_view_key",
		Amount:        "314159",
		Mixin:         312250,
		UseDust:       false,
		DustThreshold: "2",
	}

	response := GetUnspentOutsResponse{
		PerByteFee: "9",
		FeeMask:    "0",
		Amount:     "314159",
		Outputs: []Output{{
			TxID:         7,
			Amount:       "31415",
			Index:        6360,
			GlobalIndex:  "6363",
			RingCT:       "84e1c2349335412e307c518d572526b2f92c7a8d20d0cd108ee97654e3455d5b",
			TxHash:       "1adfdf87df1301136ab065e80b24217bcc2feea824a63c4eba31d46f60213fc1",
			TxPrefixHash: "01136ab065e80b24217bcc2feea824a63c4eba31d46f1adfdf87df1360213fc1",
			PublicKey:    "7aeda98c803c13a0bc7f63eb578148a578143ee601d9f67fabf3367f85c0b509",
			TxPublicKey:  "b119701f3d3eaa97d998a4e8021307785e7f107f26d4f9f72f1cc58591a712ea",
			SpendKeyImages: []string{
				"1c4c466d8d4b6546895dae3b79f2ec97cc1e3e99545191b5e2c799e3ecbabe9e",
				"62afa3a0182853cef04a7953bd191b3e3910a7775bc734fe8081856f5f68f509",
				"7920f681fa4071336774c0ca56546a8090abcbceee3b6f579fb7c337c53788f7",
			},
			Timestamp: "2024-14-19-27.0-00:00",
			Height:    3223048,
		}},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		var req = &GetUnspentOutsRequest{}

		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			t.Error("GetUnspentOuts() made an invalid request: ", err)
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
		address:    request.Address,
		client:     &http.Client{},
		retryCount: tryCount,
		retryTime:  time.Duration(0),
		serverURL:  ts.URL,
		viewKey:    request.ViewKey,
	}

	resp, err := client.GetUnspentOuts(request)
	if err != nil {
		t.Error("GetUnspentOuts() returned the error: ", err)
	}

	if reflect.DeepEqual(*resp, response) != true {
		t.Error("response struct didn't match the original data")
	}
}
