// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright © 2024 Christian Hering

package gomonerolight

import (
	"errors"
	"time"
)

var ErrorStandardRequestEncode = errors.New("failed to encode standard request using data from 'client'")

type StandardRequest struct {
	Address string `json:"address"`
	ViewKey string `json:"view_key"`
}

type Transaction struct {
	ID            uint64    `json:"id"`
	Hash          []byte    `json:"hash"`
	Timestamp     time.Time `json:"timestamp"`
	TotalReceived string    `json:"total_received"`
	TotalSent     string    `json:"total_sent"`
	UnlockTime    uint64    `json:"unlock_time"`
	Height        uint64    `json:"height"`
	SpentOutputs  []Spend   `json:"spent_outputs"`
	PaymentID     []byte    `json:"payment_id"`
	Coinbase      bool      `json:"coinbase"`
	Mempool       bool      `json:"mempool"`
	Mixin         uint64    `json:"mixin"`
}

type Spend struct {
	Amount      string `json:"amount"`
	KeyImage    []byte `json:"key_image"`
	TxPublicKey []byte `json:"tx_pub_key"`
	OutIndex    uint16 `json:"out_index"`
	Mixin       uint32 `json:"mixin"`
}
