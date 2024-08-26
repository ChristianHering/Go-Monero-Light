// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright © 2024 Christian Hering

package gomonerolight

import "errors"

// Request errors
var ErrorJoinPathFailed = errors.New("failed to join server url with path. Is the server URL in 'client' valid?")
var ErrorPostRequestFailed = errors.New("failed to post data to endpoint")

// Encoding errors
var ErrorStandardRequestEncode = errors.New("failed to encode standard request using data from 'client'")

// Request status errors
var ErrorServiceUnavailable = errors.New(`server responded with HTTP error "Service Unavailible" (status code 503) too many times`)
var ErrorStatusCodeNotOK = errors.New("server responded with a non-OK status code")

// Response errors
var ErrorResponseUnmarshalFailed = errors.New("failed to unmarshal response body from our POST request")
