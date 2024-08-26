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

func TestImportRequest(t *testing.T) {
	tryCount := 1 //Number of times to send HTTP Service Unavailable

	request := &StandardRequest{
		Address: "xmr_address",
		ViewKey: "xmr_view_key",
	}

	response := ImportRequestResponse{
		PaymentAddress:   "payment_addr",
		PaymentID:        "e8021307f3e7f10a41a712ea7f26d4f9f72f78597011cc5859b11d3eaa97d998",
		ImportFee:        "c2feea824a63c4eba3101136ab065e80b24217bcd46f1adfdf87df1360213fc1",
		NewRequest:       true,
		RequestFulfilled: false,
		Status:           "Custom status message!",
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		var req = &StandardRequest{}

		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			t.Error("ImportRequest() made an invalid request: ", err)
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

	resp, err := client.ImportRequest()
	if err != nil {
		t.Error("ImportRequest() returned the error: ", err)
	}

	if reflect.DeepEqual(*resp, response) != true {
		t.Error("response struct didn't match the original data")
	}
}
