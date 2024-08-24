// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright © 2024 Christian Hering

package gomonerolight

import "errors"

var ErrorServiceUnavailable = errors.New(`server responded with HTTP error "Service Unavailible" (status code 503) too many times`)
