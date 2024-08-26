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

func TestLogin(t *testing.T) {
	tryCount := 1 //Number of times to send HTTP Service Unavailable

	request := &LoginRequest{
		Address: "xmr_address",
		ViewKey: "xmr_view_key",
		CreateAccount:    false,
		GeneratedLocally: true,
	}

	response := LoginResponse{
		NewAddress:       true,
		GeneratedLocally: false,
		StartHeight:      3223243,
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		var req = &LoginRequest{}

		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			t.Error("Login() made an invalid request: ", err)
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

	resp, err := client.Login(request)
	if err != nil {
		t.Error("Login() returned the error: ", err)
	}

	if reflect.DeepEqual(*resp, response) != true {
		t.Error("response struct didn't match the original data")
	}
}
