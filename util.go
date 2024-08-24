// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright © 2024 Christian Hering

package gomonerolight

import (
	"fmt"
	"os"
)

func printToStderr(errorString string) {
	_, err := fmt.Fprintln(os.Stderr, errorString)
	if err != nil {
		panic(err)
	}
}
