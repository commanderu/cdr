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

// TestChainSvrWsCmds tests all of the chain server websocket-specific commands
// marshal and unmarshal into valid results include handling of optional fields
// being omitted in the marshalled command, while optional fields with defaults
// have the default assigned on unmarshalled commands.
func TestChainSvrWsCmds(t *testing.T) {
	t.Parallel()

	testID := int(1)
	tests := []struct {
		name         string
		newCmd       func() (interface{}, error)
		staticCmd    func() interface{}
		marshalled   string
		unmarshalled interface{}
	}{
		{
			name: "authenticate",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("authenticate", "user", "pass")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewAuthenticateCmd("user", "pass")
			},
			marshalled:   `{"jsonrpc":"1.0","method":"authenticate","params":["user","pass"],"id":1}`,
			unmarshalled: &cdrjson.AuthenticateCmd{Username: "user", Passphrase: "pass"},
		},
		{
			name: "notifywinningtickets",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("notifywinningtickets")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewNotifyWinningTicketsCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"notifywinningtickets","params":[],"id":1}`,
			unmarshalled: &cdrjson.NotifyWinningTicketsCmd{},
		},
		{
			name: "notifyspentandmissedtickets",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("notifyspentandmissedtickets")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewNotifySpentAndMissedTicketsCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"notifyspentandmissedtickets","params":[],"id":1}`,
			unmarshalled: &cdrjson.NotifySpentAndMissedTicketsCmd{},
		},
		{
			name: "notifynewtickets",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("notifynewtickets")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewNotifyNewTicketsCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"notifynewtickets","params":[],"id":1}`,
			unmarshalled: &cdrjson.NotifyNewTicketsCmd{},
		},
		{
			name: "notifystakedifficulty",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("notifystakedifficulty")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewNotifyStakeDifficultyCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"notifystakedifficulty","params":[],"id":1}`,
			unmarshalled: &cdrjson.NotifyStakeDifficultyCmd{},
		},
		{
			name: "notifyblocks",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("notifyblocks")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewNotifyBlocksCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"notifyblocks","params":[],"id":1}`,
			unmarshalled: &cdrjson.NotifyBlocksCmd{},
		},
		{
			name: "stopnotifyblocks",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("stopnotifyblocks")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewStopNotifyBlocksCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"stopnotifyblocks","params":[],"id":1}`,
			unmarshalled: &cdrjson.StopNotifyBlocksCmd{},
		},
		{
			name: "notifynewtransactions",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("notifynewtransactions")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewNotifyNewTransactionsCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"notifynewtransactions","params":[],"id":1}`,
			unmarshalled: &cdrjson.NotifyNewTransactionsCmd{
				Verbose: cdrjson.Bool(false),
			},
		},
		{
			name: "notifynewtransactions optional",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("notifynewtransactions", true)
			},
			staticCmd: func() interface{} {
				return cdrjson.NewNotifyNewTransactionsCmd(cdrjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"notifynewtransactions","params":[true],"id":1}`,
			unmarshalled: &cdrjson.NotifyNewTransactionsCmd{
				Verbose: cdrjson.Bool(true),
			},
		},
		{
			name: "stopnotifynewtransactions",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("stopnotifynewtransactions")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewStopNotifyNewTransactionsCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"stopnotifynewtransactions","params":[],"id":1}`,
			unmarshalled: &cdrjson.StopNotifyNewTransactionsCmd{},
		},
		{
			name: "rescan",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("rescan", "0000000000000000000000000000000000000000000000000000000000000123")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewRescanCmd("0000000000000000000000000000000000000000000000000000000000000123")
			},
			marshalled: `{"jsonrpc":"1.0","method":"rescan","params":["0000000000000000000000000000000000000000000000000000000000000123"],"id":1}`,
			unmarshalled: &cdrjson.RescanCmd{
				BlockHashes: "0000000000000000000000000000000000000000000000000000000000000123",
			},
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Marshal the command as created by the new static command
		// creation function.
		marshalled, err := cdrjson.MarshalCmd("1.0", testID, test.staticCmd())
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

		// Ensure the command is created without error via the generic
		// new command creation function.
		cmd, err := test.newCmd()
		if err != nil {
			t.Errorf("Test #%d (%s) unexpected NewCmd error: %v ",
				i, test.name, err)
		}

		// Marshal the command as created by the generic new command
		// creation function.
		marshalled, err = cdrjson.MarshalCmd("1.0", testID, cmd)
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
