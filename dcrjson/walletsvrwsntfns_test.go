// Copyright (c) 2014 The btcsuite developers
// Copyright (c) 2015-2016 The commanderu developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package cdrjson_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/commanderu/cdrd/cdrjson"
)

// TestWalletSvrWsNtfns tests all of the chain server websocket-specific
// notifications marshal and unmarshal into valid results include handling of
// optional fields being omitted in the marshalled command, while optional
// fields with defaults have the default assigned on unmarshalled commands.
func TestWalletSvrWsNtfns(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		newNtfn      func() (interface{}, error)
		staticNtfn   func() interface{}
		marshalled   string
		unmarshalled interface{}
	}{
		{
			name: "accountbalance",
			newNtfn: func() (interface{}, error) {
				return cdrjson.NewCmd("accountbalance", "acct", 1.25, true)
			},
			staticNtfn: func() interface{} {
				return cdrjson.NewAccountBalanceNtfn("acct", 1.25, true)
			},
			marshalled: `{"jsonrpc":"1.0","method":"accountbalance","params":["acct",1.25,true],"id":null}`,
			unmarshalled: &cdrjson.AccountBalanceNtfn{
				Account:   "acct",
				Balance:   1.25,
				Confirmed: true,
			},
		},
		{
			name: "cdrdconnected",
			newNtfn: func() (interface{}, error) {
				return cdrjson.NewCmd("cdrdconnected", true)
			},
			staticNtfn: func() interface{} {
				return cdrjson.NewBtcdConnectedNtfn(true)
			},
			marshalled: `{"jsonrpc":"1.0","method":"cdrdconnected","params":[true],"id":null}`,
			unmarshalled: &cdrjson.BtcdConnectedNtfn{
				Connected: true,
			},
		},
		{
			name: "walletlockstate",
			newNtfn: func() (interface{}, error) {
				return cdrjson.NewCmd("walletlockstate", true)
			},
			staticNtfn: func() interface{} {
				return cdrjson.NewWalletLockStateNtfn(true)
			},
			marshalled: `{"jsonrpc":"1.0","method":"walletlockstate","params":[true],"id":null}`,
			unmarshalled: &cdrjson.WalletLockStateNtfn{
				Locked: true,
			},
		},
		{
			name: "newtx",
			newNtfn: func() (interface{}, error) {
				return cdrjson.NewCmd("newtx", "acct", `{"account":"acct","address":"1Address","category":"send","amount":1.5,"fee":0.0001,"confirmations":1,"txid":"456","walletconflicts":[],"time":12345678,"timereceived":12345876,"vout":789,"otheraccount":"otheracct"}`)
			},
			staticNtfn: func() interface{} {
				result := cdrjson.ListTransactionsResult{
					Account:         "acct",
					Address:         "1Address",
					Category:        "send",
					Amount:          1.5,
					Fee:             cdrjson.Float64(0.0001),
					Confirmations:   1,
					TxID:            "456",
					WalletConflicts: []string{},
					Time:            12345678,
					TimeReceived:    12345876,
					Vout:            789,
					OtherAccount:    "otheracct",
				}
				return cdrjson.NewNewTxNtfn("acct", result)
			},
			marshalled: `{"jsonrpc":"1.0","method":"newtx","params":["acct",{"account":"acct","address":"1Address","amount":1.5,"category":"send","confirmations":1,"fee":0.0001,"time":12345678,"timereceived":12345876,"txid":"456","vout":789,"walletconflicts":[],"otheraccount":"otheracct"}],"id":null}`,
			unmarshalled: &cdrjson.NewTxNtfn{
				Account: "acct",
				Details: cdrjson.ListTransactionsResult{
					Account:         "acct",
					Address:         "1Address",
					Category:        "send",
					Amount:          1.5,
					Fee:             cdrjson.Float64(0.0001),
					Confirmations:   1,
					TxID:            "456",
					WalletConflicts: []string{},
					Time:            12345678,
					TimeReceived:    12345876,
					Vout:            789,
					OtherAccount:    "otheracct",
				},
			},
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Marshal the notification as created by the new static
		// creation function.  The ID is nil for notifications.
		marshalled, err := cdrjson.MarshalCmd("1.0", nil, test.staticNtfn())
		if err != nil {
			t.Errorf("MarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !bytes.Equal(marshalled, []byte(test.marshalled)) {
			t.Errorf("Test #%d (%s) unexpected marshalled data - "+
				"got %s, want %s", i, test.name, marshalled,
				test.marshalled)
			continue
		}

		// Ensure the notification is created without error via the
		// generic new notification creation function.
		cmd, err := test.newNtfn()
		if err != nil {
			t.Errorf("Test #%d (%s) unexpected NewCmd error: %v ",
				i, test.name, err)
		}

		// Marshal the notification as created by the generic new
		// notification creation function.    The ID is nil for
		// notifications.
		marshalled, err = cdrjson.MarshalCmd("1.0", nil, cmd)
		if err != nil {
			t.Errorf("MarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !bytes.Equal(marshalled, []byte(test.marshalled)) {
			t.Errorf("Test #%d (%s) unexpected marshalled data - "+
				"got %s, want %s", i, test.name, marshalled,
				test.marshalled)
			continue
		}

		var request cdrjson.Request
		if err := json.Unmarshal(marshalled, &request); err != nil {
			t.Errorf("Test #%d (%s) unexpected error while "+
				"unmarshalling JSON-RPC request: %v", i,
				test.name, err)
			continue
		}

		cmd, err = cdrjson.UnmarshalCmd(&request)
		if err != nil {
			t.Errorf("UnmarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !reflect.DeepEqual(cmd, test.unmarshalled) {
			t.Errorf("Test #%d (%s) unexpected unmarshalled command "+
				"- got %s, want %s", i, test.name,
				fmt.Sprintf("(%T) %+[1]v", cmd),
				fmt.Sprintf("(%T) %+[1]v\n", test.unmarshalled))
			continue
		}
	}
}
