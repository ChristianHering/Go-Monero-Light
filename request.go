// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright © 2024 Christian Hering

package gomonerolight

import (
	"time"
)

// StandardRequest is the most common request used
// in Monero's light wallet API and doesn't need to
// be passed to function calls after client creation.
type StandardRequest struct {
	Address string `json:"address"`
	ViewKey string `json:"view_key"`
}

type Transaction struct {
	ID            uint64    `json:"id"`
	Hash          string    `json:"hash"` // hex encoded binary
	Timestamp     time.Time `json:"timestamp"`
	TotalReceived string    `json:"total_received"`
	TotalSent     string    `json:"total_sent"`
	UnlockTime    uint64    `json:"unlock_time"`
	Height        uint64    `json:"height"`
	SpentOutputs  []Spend   `json:"spent_outputs"`
	PaymentID     string    `json:"payment_id"` // hex encoded binary
	Coinbase      bool      `json:"coinbase"`
	Mempool       bool      `json:"mempool"`
	Mixin         uint64    `json:"mixin"`
}

type Spend struct {
	Amount      string `json:"amount"`
	KeyImage    string `json:"key_image"`  // hex encoded binary
	TxPublicKey string `json:"tx_pub_key"` // hex encoded binary
	OutIndex    uint16 `json:"out_index"`
	Mixin       uint32 `json:"mixin"`
}
