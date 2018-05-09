// Copyright (c) 2014 The btcsuite developers
// Copyright (c) 2015-2018 The commanderu developers
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

// TestWalletSvrCmds tests all of the wallet server commands marshal and
// unmarshal into valid results include handling of optional fields being
// omitted in the marshalled command, while optional fields with defaults have
// the default assigned on unmarshalled commands.
func TestWalletSvrCmds(t *testing.T) {
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
			name: "addmultisigaddress",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("addmultisigaddress", 2, []string{"031234", "035678"})
			},
			staticCmd: func() interface{} {
				keys := []string{"031234", "035678"}
				return cdrjson.NewAddMultisigAddressCmd(2, keys, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"addmultisigaddress","params":[2,["031234","035678"]],"id":1}`,
			unmarshalled: &cdrjson.AddMultisigAddressCmd{
				NRequired: 2,
				Keys:      []string{"031234", "035678"},
				Account:   nil,
			},
		},
		{
			name: "addmultisigaddress optional",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("addmultisigaddress", 2, []string{"031234", "035678"}, "test")
			},
			staticCmd: func() interface{} {
				keys := []string{"031234", "035678"}
				return cdrjson.NewAddMultisigAddressCmd(2, keys, cdrjson.String("test"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"addmultisigaddress","params":[2,["031234","035678"],"test"],"id":1}`,
			unmarshalled: &cdrjson.AddMultisigAddressCmd{
				NRequired: 2,
				Keys:      []string{"031234", "035678"},
				Account:   cdrjson.String("test"),
			},
		},
		{
			name: "createmultisig",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("createmultisig", 2, []string{"031234", "035678"})
			},
			staticCmd: func() interface{} {
				keys := []string{"031234", "035678"}
				return cdrjson.NewCreateMultisigCmd(2, keys)
			},
			marshalled: `{"jsonrpc":"1.0","method":"createmultisig","params":[2,["031234","035678"]],"id":1}`,
			unmarshalled: &cdrjson.CreateMultisigCmd{
				NRequired: 2,
				Keys:      []string{"031234", "035678"},
			},
		},
		{
			name: "dumpprivkey",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("dumpprivkey", "1Address")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewDumpPrivKeyCmd("1Address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"dumpprivkey","params":["1Address"],"id":1}`,
			unmarshalled: &cdrjson.DumpPrivKeyCmd{
				Address: "1Address",
			},
		},
		{
			name: "estimatefee",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("estimatefee", 6)
			},
			staticCmd: func() interface{} {
				return cdrjson.NewEstimateFeeCmd(6)
			},
			marshalled: `{"jsonrpc":"1.0","method":"estimatefee","params":[6],"id":1}`,
			unmarshalled: &cdrjson.EstimateFeeCmd{
				NumBlocks: 6,
			},
		},
		{
			name: "estimatepriority",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("estimatepriority", 6)
			},
			staticCmd: func() interface{} {
				return cdrjson.NewEstimatePriorityCmd(6)
			},
			marshalled: `{"jsonrpc":"1.0","method":"estimatepriority","params":[6],"id":1}`,
			unmarshalled: &cdrjson.EstimatePriorityCmd{
				NumBlocks: 6,
			},
		},
		{
			name: "getaccount",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("getaccount", "1Address")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewGetAccountCmd("1Address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaccount","params":["1Address"],"id":1}`,
			unmarshalled: &cdrjson.GetAccountCmd{
				Address: "1Address",
			},
		},
		{
			name: "getaccountaddress",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("getaccountaddress", "acct")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewGetAccountAddressCmd("acct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaccountaddress","params":["acct"],"id":1}`,
			unmarshalled: &cdrjson.GetAccountAddressCmd{
				Account: "acct",
			},
		},
		{
			name: "getaddressesbyaccount",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("getaddressesbyaccount", "acct")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewGetAddressesByAccountCmd("acct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaddressesbyaccount","params":["acct"],"id":1}`,
			unmarshalled: &cdrjson.GetAddressesByAccountCmd{
				Account: "acct",
			},
		},
		{
			name: "getbalance",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("getbalance")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewGetBalanceCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getbalance","params":[],"id":1}`,
			unmarshalled: &cdrjson.GetBalanceCmd{
				Account: nil,
				MinConf: cdrjson.Int(1),
			},
		},
		{
			name: "getbalance optional1",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("getbalance", "acct")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewGetBalanceCmd(cdrjson.String("acct"), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getbalance","params":["acct"],"id":1}`,
			unmarshalled: &cdrjson.GetBalanceCmd{
				Account: cdrjson.String("acct"),
				MinConf: cdrjson.Int(1),
			},
		},
		{
			name: "getbalance optional2",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("getbalance", "acct", 6)
			},
			staticCmd: func() interface{} {
				return cdrjson.NewGetBalanceCmd(cdrjson.String("acct"), cdrjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getbalance","params":["acct",6],"id":1}`,
			unmarshalled: &cdrjson.GetBalanceCmd{
				Account: cdrjson.String("acct"),
				MinConf: cdrjson.Int(6),
			},
		},
		{
			name: "getnewaddress",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("getnewaddress")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewGetNewAddressCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnewaddress","params":[],"id":1}`,
			unmarshalled: &cdrjson.GetNewAddressCmd{
				Account:   nil,
				GapPolicy: nil,
			},
		},
		{
			name: "getnewaddress optional",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("getnewaddress", "acct", "ignore")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewGetNewAddressCmd(cdrjson.String("acct"), cdrjson.String("ignore"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnewaddress","params":["acct","ignore"],"id":1}`,
			unmarshalled: &cdrjson.GetNewAddressCmd{
				Account:   cdrjson.String("acct"),
				GapPolicy: cdrjson.String("ignore"),
			},
		},
		{
			name: "getrawchangeaddress",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("getrawchangeaddress")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewGetRawChangeAddressCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawchangeaddress","params":[],"id":1}`,
			unmarshalled: &cdrjson.GetRawChangeAddressCmd{
				Account: nil,
			},
		},
		{
			name: "getrawchangeaddress optional",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("getrawchangeaddress", "acct")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewGetRawChangeAddressCmd(cdrjson.String("acct"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawchangeaddress","params":["acct"],"id":1}`,
			unmarshalled: &cdrjson.GetRawChangeAddressCmd{
				Account: cdrjson.String("acct"),
			},
		},
		{
			name: "getreceivedbyaccount",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("getreceivedbyaccount", "acct")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewGetReceivedByAccountCmd("acct", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaccount","params":["acct"],"id":1}`,
			unmarshalled: &cdrjson.GetReceivedByAccountCmd{
				Account: "acct",
				MinConf: cdrjson.Int(1),
			},
		},
		{
			name: "getreceivedbyaccount optional",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("getreceivedbyaccount", "acct", 6)
			},
			staticCmd: func() interface{} {
				return cdrjson.NewGetReceivedByAccountCmd("acct", cdrjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaccount","params":["acct",6],"id":1}`,
			unmarshalled: &cdrjson.GetReceivedByAccountCmd{
				Account: "acct",
				MinConf: cdrjson.Int(6),
			},
		},
		{
			name: "getreceivedbyaddress",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("getreceivedbyaddress", "1Address")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewGetReceivedByAddressCmd("1Address", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaddress","params":["1Address"],"id":1}`,
			unmarshalled: &cdrjson.GetReceivedByAddressCmd{
				Address: "1Address",
				MinConf: cdrjson.Int(1),
			},
		},
		{
			name: "getreceivedbyaddress optional",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("getreceivedbyaddress", "1Address", 6)
			},
			staticCmd: func() interface{} {
				return cdrjson.NewGetReceivedByAddressCmd("1Address", cdrjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaddress","params":["1Address",6],"id":1}`,
			unmarshalled: &cdrjson.GetReceivedByAddressCmd{
				Address: "1Address",
				MinConf: cdrjson.Int(6),
			},
		},
		{
			name: "gettransaction",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("gettransaction", "123")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewGetTransactionCmd("123", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettransaction","params":["123"],"id":1}`,
			unmarshalled: &cdrjson.GetTransactionCmd{
				Txid:             "123",
				IncludeWatchOnly: cdrjson.Bool(false),
			},
		},
		{
			name: "gettransaction optional",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("gettransaction", "123", true)
			},
			staticCmd: func() interface{} {
				return cdrjson.NewGetTransactionCmd("123", cdrjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettransaction","params":["123",true],"id":1}`,
			unmarshalled: &cdrjson.GetTransactionCmd{
				Txid:             "123",
				IncludeWatchOnly: cdrjson.Bool(true),
			},
		},
		{
			name: "importprivkey",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("importprivkey", "abc")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewImportPrivKeyCmd("abc", nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc"],"id":1}`,
			unmarshalled: &cdrjson.ImportPrivKeyCmd{
				PrivKey: "abc",
				Label:   nil,
				Rescan:  cdrjson.Bool(true),
			},
		},
		{
			name: "importprivkey optional1",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("importprivkey", "abc", "label")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewImportPrivKeyCmd("abc", cdrjson.String("label"), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc","label"],"id":1}`,
			unmarshalled: &cdrjson.ImportPrivKeyCmd{
				PrivKey: "abc",
				Label:   cdrjson.String("label"),
				Rescan:  cdrjson.Bool(true),
			},
		},
		{
			name: "importprivkey optional2",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("importprivkey", "abc", "label", false)
			},
			staticCmd: func() interface{} {
				return cdrjson.NewImportPrivKeyCmd("abc", cdrjson.String("label"), cdrjson.Bool(false), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc","label",false],"id":1}`,
			unmarshalled: &cdrjson.ImportPrivKeyCmd{
				PrivKey: "abc",
				Label:   cdrjson.String("label"),
				Rescan:  cdrjson.Bool(false),
			},
		},
		{
			name: "importprivkey optional3",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("importprivkey", "abc", "label", false, 12345)
			},
			staticCmd: func() interface{} {
				return cdrjson.NewImportPrivKeyCmd("abc", cdrjson.String("label"), cdrjson.Bool(false), cdrjson.Int(12345))
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc","label",false,12345],"id":1}`,
			unmarshalled: &cdrjson.ImportPrivKeyCmd{
				PrivKey:  "abc",
				Label:    cdrjson.String("label"),
				Rescan:   cdrjson.Bool(false),
				ScanFrom: cdrjson.Int(12345),
			},
		},
		{
			name: "keypoolrefill",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("keypoolrefill")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewKeyPoolRefillCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"keypoolrefill","params":[],"id":1}`,
			unmarshalled: &cdrjson.KeyPoolRefillCmd{
				NewSize: cdrjson.Uint(100),
			},
		},
		{
			name: "keypoolrefill optional",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("keypoolrefill", 200)
			},
			staticCmd: func() interface{} {
				return cdrjson.NewKeyPoolRefillCmd(cdrjson.Uint(200))
			},
			marshalled: `{"jsonrpc":"1.0","method":"keypoolrefill","params":[200],"id":1}`,
			unmarshalled: &cdrjson.KeyPoolRefillCmd{
				NewSize: cdrjson.Uint(200),
			},
		},
		{
			name: "listaccounts",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("listaccounts")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewListAccountsCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listaccounts","params":[],"id":1}`,
			unmarshalled: &cdrjson.ListAccountsCmd{
				MinConf: cdrjson.Int(1),
			},
		},
		{
			name: "listaccounts optional",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("listaccounts", 6)
			},
			staticCmd: func() interface{} {
				return cdrjson.NewListAccountsCmd(cdrjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listaccounts","params":[6],"id":1}`,
			unmarshalled: &cdrjson.ListAccountsCmd{
				MinConf: cdrjson.Int(6),
			},
		},
		{
			name: "listlockunspent",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("listlockunspent")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewListLockUnspentCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"listlockunspent","params":[],"id":1}`,
			unmarshalled: &cdrjson.ListLockUnspentCmd{},
		},
		{
			name: "listreceivedbyaccount",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("listreceivedbyaccount")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewListReceivedByAccountCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[],"id":1}`,
			unmarshalled: &cdrjson.ListReceivedByAccountCmd{
				MinConf:          cdrjson.Int(1),
				IncludeEmpty:     cdrjson.Bool(false),
				IncludeWatchOnly: cdrjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaccount optional1",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("listreceivedbyaccount", 6)
			},
			staticCmd: func() interface{} {
				return cdrjson.NewListReceivedByAccountCmd(cdrjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[6],"id":1}`,
			unmarshalled: &cdrjson.ListReceivedByAccountCmd{
				MinConf:          cdrjson.Int(6),
				IncludeEmpty:     cdrjson.Bool(false),
				IncludeWatchOnly: cdrjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaccount optional2",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("listreceivedbyaccount", 6, true)
			},
			staticCmd: func() interface{} {
				return cdrjson.NewListReceivedByAccountCmd(cdrjson.Int(6), cdrjson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[6,true],"id":1}`,
			unmarshalled: &cdrjson.ListReceivedByAccountCmd{
				MinConf:          cdrjson.Int(6),
				IncludeEmpty:     cdrjson.Bool(true),
				IncludeWatchOnly: cdrjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaccount optional3",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("listreceivedbyaccount", 6, true, false)
			},
			staticCmd: func() interface{} {
				return cdrjson.NewListReceivedByAccountCmd(cdrjson.Int(6), cdrjson.Bool(true), cdrjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[6,true,false],"id":1}`,
			unmarshalled: &cdrjson.ListReceivedByAccountCmd{
				MinConf:          cdrjson.Int(6),
				IncludeEmpty:     cdrjson.Bool(true),
				IncludeWatchOnly: cdrjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("listreceivedbyaddress")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewListReceivedByAddressCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[],"id":1}`,
			unmarshalled: &cdrjson.ListReceivedByAddressCmd{
				MinConf:          cdrjson.Int(1),
				IncludeEmpty:     cdrjson.Bool(false),
				IncludeWatchOnly: cdrjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress optional1",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("listreceivedbyaddress", 6)
			},
			staticCmd: func() interface{} {
				return cdrjson.NewListReceivedByAddressCmd(cdrjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[6],"id":1}`,
			unmarshalled: &cdrjson.ListReceivedByAddressCmd{
				MinConf:          cdrjson.Int(6),
				IncludeEmpty:     cdrjson.Bool(false),
				IncludeWatchOnly: cdrjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress optional2",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("listreceivedbyaddress", 6, true)
			},
			staticCmd: func() interface{} {
				return cdrjson.NewListReceivedByAddressCmd(cdrjson.Int(6), cdrjson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[6,true],"id":1}`,
			unmarshalled: &cdrjson.ListReceivedByAddressCmd{
				MinConf:          cdrjson.Int(6),
				IncludeEmpty:     cdrjson.Bool(true),
				IncludeWatchOnly: cdrjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress optional3",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("listreceivedbyaddress", 6, true, false)
			},
			staticCmd: func() interface{} {
				return cdrjson.NewListReceivedByAddressCmd(cdrjson.Int(6), cdrjson.Bool(true), cdrjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[6,true,false],"id":1}`,
			unmarshalled: &cdrjson.ListReceivedByAddressCmd{
				MinConf:          cdrjson.Int(6),
				IncludeEmpty:     cdrjson.Bool(true),
				IncludeWatchOnly: cdrjson.Bool(false),
			},
		},
		{
			name: "listsinceblock",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("listsinceblock")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewListSinceBlockCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":[],"id":1}`,
			unmarshalled: &cdrjson.ListSinceBlockCmd{
				BlockHash:           nil,
				TargetConfirmations: cdrjson.Int(1),
				IncludeWatchOnly:    cdrjson.Bool(false),
			},
		},
		{
			name: "listsinceblock optional1",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("listsinceblock", "123")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewListSinceBlockCmd(cdrjson.String("123"), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":["123"],"id":1}`,
			unmarshalled: &cdrjson.ListSinceBlockCmd{
				BlockHash:           cdrjson.String("123"),
				TargetConfirmations: cdrjson.Int(1),
				IncludeWatchOnly:    cdrjson.Bool(false),
			},
		},
		{
			name: "listsinceblock optional2",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("listsinceblock", "123", 6)
			},
			staticCmd: func() interface{} {
				return cdrjson.NewListSinceBlockCmd(cdrjson.String("123"), cdrjson.Int(6), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":["123",6],"id":1}`,
			unmarshalled: &cdrjson.ListSinceBlockCmd{
				BlockHash:           cdrjson.String("123"),
				TargetConfirmations: cdrjson.Int(6),
				IncludeWatchOnly:    cdrjson.Bool(false),
			},
		},
		{
			name: "listsinceblock optional3",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("listsinceblock", "123", 6, true)
			},
			staticCmd: func() interface{} {
				return cdrjson.NewListSinceBlockCmd(cdrjson.String("123"), cdrjson.Int(6), cdrjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":["123",6,true],"id":1}`,
			unmarshalled: &cdrjson.ListSinceBlockCmd{
				BlockHash:           cdrjson.String("123"),
				TargetConfirmations: cdrjson.Int(6),
				IncludeWatchOnly:    cdrjson.Bool(true),
			},
		},
		{
			name: "listtransactions",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("listtransactions")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewListTransactionsCmd(nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":[],"id":1}`,
			unmarshalled: &cdrjson.ListTransactionsCmd{
				Account:          nil,
				Count:            cdrjson.Int(10),
				From:             cdrjson.Int(0),
				IncludeWatchOnly: cdrjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional1",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("listtransactions", "acct")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewListTransactionsCmd(cdrjson.String("acct"), nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct"],"id":1}`,
			unmarshalled: &cdrjson.ListTransactionsCmd{
				Account:          cdrjson.String("acct"),
				Count:            cdrjson.Int(10),
				From:             cdrjson.Int(0),
				IncludeWatchOnly: cdrjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional2",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("listtransactions", "acct", 20)
			},
			staticCmd: func() interface{} {
				return cdrjson.NewListTransactionsCmd(cdrjson.String("acct"), cdrjson.Int(20), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct",20],"id":1}`,
			unmarshalled: &cdrjson.ListTransactionsCmd{
				Account:          cdrjson.String("acct"),
				Count:            cdrjson.Int(20),
				From:             cdrjson.Int(0),
				IncludeWatchOnly: cdrjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional3",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("listtransactions", "acct", 20, 1)
			},
			staticCmd: func() interface{} {
				return cdrjson.NewListTransactionsCmd(cdrjson.String("acct"), cdrjson.Int(20),
					cdrjson.Int(1), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct",20,1],"id":1}`,
			unmarshalled: &cdrjson.ListTransactionsCmd{
				Account:          cdrjson.String("acct"),
				Count:            cdrjson.Int(20),
				From:             cdrjson.Int(1),
				IncludeWatchOnly: cdrjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional4",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("listtransactions", "acct", 20, 1, true)
			},
			staticCmd: func() interface{} {
				return cdrjson.NewListTransactionsCmd(cdrjson.String("acct"), cdrjson.Int(20),
					cdrjson.Int(1), cdrjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct",20,1,true],"id":1}`,
			unmarshalled: &cdrjson.ListTransactionsCmd{
				Account:          cdrjson.String("acct"),
				Count:            cdrjson.Int(20),
				From:             cdrjson.Int(1),
				IncludeWatchOnly: cdrjson.Bool(true),
			},
		},
		{
			name: "listunspent",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("listunspent")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewListUnspentCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[],"id":1}`,
			unmarshalled: &cdrjson.ListUnspentCmd{
				MinConf:   cdrjson.Int(1),
				MaxConf:   cdrjson.Int(9999999),
				Addresses: nil,
			},
		},
		{
			name: "listunspent optional1",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("listunspent", 6)
			},
			staticCmd: func() interface{} {
				return cdrjson.NewListUnspentCmd(cdrjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[6],"id":1}`,
			unmarshalled: &cdrjson.ListUnspentCmd{
				MinConf:   cdrjson.Int(6),
				MaxConf:   cdrjson.Int(9999999),
				Addresses: nil,
			},
		},
		{
			name: "listunspent optional2",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("listunspent", 6, 100)
			},
			staticCmd: func() interface{} {
				return cdrjson.NewListUnspentCmd(cdrjson.Int(6), cdrjson.Int(100), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[6,100],"id":1}`,
			unmarshalled: &cdrjson.ListUnspentCmd{
				MinConf:   cdrjson.Int(6),
				MaxConf:   cdrjson.Int(100),
				Addresses: nil,
			},
		},
		{
			name: "listunspent optional3",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("listunspent", 6, 100, []string{"1Address", "1Address2"})
			},
			staticCmd: func() interface{} {
				return cdrjson.NewListUnspentCmd(cdrjson.Int(6), cdrjson.Int(100),
					&[]string{"1Address", "1Address2"})
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[6,100,["1Address","1Address2"]],"id":1}`,
			unmarshalled: &cdrjson.ListUnspentCmd{
				MinConf:   cdrjson.Int(6),
				MaxConf:   cdrjson.Int(100),
				Addresses: &[]string{"1Address", "1Address2"},
			},
		},
		{
			name: "lockunspent",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("lockunspent", true, `[{"txid":"123","vout":1}]`)
			},
			staticCmd: func() interface{} {
				txInputs := []cdrjson.TransactionInput{
					{Txid: "123", Vout: 1},
				}
				return cdrjson.NewLockUnspentCmd(true, txInputs)
			},
			marshalled: `{"jsonrpc":"1.0","method":"lockunspent","params":[true,[{"txid":"123","vout":1,"tree":0}]],"id":1}`,
			unmarshalled: &cdrjson.LockUnspentCmd{
				Unlock: true,
				Transactions: []cdrjson.TransactionInput{
					{Txid: "123", Vout: 1},
				},
			},
		},
		{
			name: "sendfrom",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("sendfrom", "from", "1Address", 0.5)
			},
			staticCmd: func() interface{} {
				return cdrjson.NewSendFromCmd("from", "1Address", 0.5, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5],"id":1}`,
			unmarshalled: &cdrjson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     cdrjson.Int(1),
				Comment:     nil,
				CommentTo:   nil,
			},
		},
		{
			name: "sendfrom optional1",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("sendfrom", "from", "1Address", 0.5, 6)
			},
			staticCmd: func() interface{} {
				return cdrjson.NewSendFromCmd("from", "1Address", 0.5, cdrjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5,6],"id":1}`,
			unmarshalled: &cdrjson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     cdrjson.Int(6),
				Comment:     nil,
				CommentTo:   nil,
			},
		},
		{
			name: "sendfrom optional2",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("sendfrom", "from", "1Address", 0.5, 6, "comment")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewSendFromCmd("from", "1Address", 0.5, cdrjson.Int(6),
					cdrjson.String("comment"), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5,6,"comment"],"id":1}`,
			unmarshalled: &cdrjson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     cdrjson.Int(6),
				Comment:     cdrjson.String("comment"),
				CommentTo:   nil,
			},
		},
		{
			name: "sendfrom optional3",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("sendfrom", "from", "1Address", 0.5, 6, "comment", "commentto")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewSendFromCmd("from", "1Address", 0.5, cdrjson.Int(6),
					cdrjson.String("comment"), cdrjson.String("commentto"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5,6,"comment","commentto"],"id":1}`,
			unmarshalled: &cdrjson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     cdrjson.Int(6),
				Comment:     cdrjson.String("comment"),
				CommentTo:   cdrjson.String("commentto"),
			},
		},
		{
			name: "sendmany",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("sendmany", "from", `{"1Address":0.5}`)
			},
			staticCmd: func() interface{} {
				amounts := map[string]float64{"1Address": 0.5}
				return cdrjson.NewSendManyCmd("from", amounts, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendmany","params":["from",{"1Address":0.5}],"id":1}`,
			unmarshalled: &cdrjson.SendManyCmd{
				FromAccount: "from",
				Amounts:     map[string]float64{"1Address": 0.5},
				MinConf:     cdrjson.Int(1),
				Comment:     nil,
			},
		},
		{
			name: "sendmany optional1",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("sendmany", "from", `{"1Address":0.5}`, 6)
			},
			staticCmd: func() interface{} {
				amounts := map[string]float64{"1Address": 0.5}
				return cdrjson.NewSendManyCmd("from", amounts, cdrjson.Int(6), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendmany","params":["from",{"1Address":0.5},6],"id":1}`,
			unmarshalled: &cdrjson.SendManyCmd{
				FromAccount: "from",
				Amounts:     map[string]float64{"1Address": 0.5},
				MinConf:     cdrjson.Int(6),
				Comment:     nil,
			},
		},
		{
			name: "sendmany optional2",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("sendmany", "from", `{"1Address":0.5}`, 6, "comment")
			},
			staticCmd: func() interface{} {
				amounts := map[string]float64{"1Address": 0.5}
				return cdrjson.NewSendManyCmd("from", amounts, cdrjson.Int(6), cdrjson.String("comment"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendmany","params":["from",{"1Address":0.5},6,"comment"],"id":1}`,
			unmarshalled: &cdrjson.SendManyCmd{
				FromAccount: "from",
				Amounts:     map[string]float64{"1Address": 0.5},
				MinConf:     cdrjson.Int(6),
				Comment:     cdrjson.String("comment"),
			},
		},
		{
			name: "sendtoaddress",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("sendtoaddress", "1Address", 0.5)
			},
			staticCmd: func() interface{} {
				return cdrjson.NewSendToAddressCmd("1Address", 0.5, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendtoaddress","params":["1Address",0.5],"id":1}`,
			unmarshalled: &cdrjson.SendToAddressCmd{
				Address:   "1Address",
				Amount:    0.5,
				Comment:   nil,
				CommentTo: nil,
			},
		},
		{
			name: "sendtoaddress optional1",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("sendtoaddress", "1Address", 0.5, "comment", "commentto")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewSendToAddressCmd("1Address", 0.5, cdrjson.String("comment"),
					cdrjson.String("commentto"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendtoaddress","params":["1Address",0.5,"comment","commentto"],"id":1}`,
			unmarshalled: &cdrjson.SendToAddressCmd{
				Address:   "1Address",
				Amount:    0.5,
				Comment:   cdrjson.String("comment"),
				CommentTo: cdrjson.String("commentto"),
			},
		},
		{
			name: "settxfee",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("settxfee", 0.0001)
			},
			staticCmd: func() interface{} {
				return cdrjson.NewSetTxFeeCmd(0.0001)
			},
			marshalled: `{"jsonrpc":"1.0","method":"settxfee","params":[0.0001],"id":1}`,
			unmarshalled: &cdrjson.SetTxFeeCmd{
				Amount: 0.0001,
			},
		},
		{
			name: "signmessage",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("signmessage", "1Address", "message")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewSignMessageCmd("1Address", "message")
			},
			marshalled: `{"jsonrpc":"1.0","method":"signmessage","params":["1Address","message"],"id":1}`,
			unmarshalled: &cdrjson.SignMessageCmd{
				Address: "1Address",
				Message: "message",
			},
		},
		{
			name: "signrawtransaction",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("signrawtransaction", "001122")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewSignRawTransactionCmd("001122", nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122"],"id":1}`,
			unmarshalled: &cdrjson.SignRawTransactionCmd{
				RawTx:    "001122",
				Inputs:   nil,
				PrivKeys: nil,
				Flags:    cdrjson.String("ALL"),
			},
		},
		{
			name: "signrawtransaction optional1",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("signrawtransaction", "001122", `[{"txid":"123","vout":1,"tree":0,"scriptPubKey":"00","redeemScript":"01"}]`)
			},
			staticCmd: func() interface{} {
				txInputs := []cdrjson.RawTxInput{
					{
						Txid:         "123",
						Vout:         1,
						ScriptPubKey: "00",
						RedeemScript: "01",
					},
				}

				return cdrjson.NewSignRawTransactionCmd("001122", &txInputs, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122",[{"txid":"123","vout":1,"tree":0,"scriptPubKey":"00","redeemScript":"01"}]],"id":1}`,
			unmarshalled: &cdrjson.SignRawTransactionCmd{
				RawTx: "001122",
				Inputs: &[]cdrjson.RawTxInput{
					{
						Txid:         "123",
						Vout:         1,
						ScriptPubKey: "00",
						RedeemScript: "01",
					},
				},
				PrivKeys: nil,
				Flags:    cdrjson.String("ALL"),
			},
		},
		{
			name: "signrawtransaction optional2",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("signrawtransaction", "001122", `[]`, `["abc"]`)
			},
			staticCmd: func() interface{} {
				txInputs := []cdrjson.RawTxInput{}
				privKeys := []string{"abc"}
				return cdrjson.NewSignRawTransactionCmd("001122", &txInputs, &privKeys, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122",[],["abc"]],"id":1}`,
			unmarshalled: &cdrjson.SignRawTransactionCmd{
				RawTx:    "001122",
				Inputs:   &[]cdrjson.RawTxInput{},
				PrivKeys: &[]string{"abc"},
				Flags:    cdrjson.String("ALL"),
			},
		},
		{
			name: "signrawtransaction optional3",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("signrawtransaction", "001122", `[]`, `[]`, "ALL")
			},
			staticCmd: func() interface{} {
				txInputs := []cdrjson.RawTxInput{}
				privKeys := []string{}
				return cdrjson.NewSignRawTransactionCmd("001122", &txInputs, &privKeys,
					cdrjson.String("ALL"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122",[],[],"ALL"],"id":1}`,
			unmarshalled: &cdrjson.SignRawTransactionCmd{
				RawTx:    "001122",
				Inputs:   &[]cdrjson.RawTxInput{},
				PrivKeys: &[]string{},
				Flags:    cdrjson.String("ALL"),
			},
		},
		{
			name: "verifyseed",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("verifyseed", "abc")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewVerifySeedCmd("abc", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifyseed","params":["abc"],"id":1}`,
			unmarshalled: &cdrjson.VerifySeedCmd{
				Seed:    "abc",
				Account: nil,
			},
		},
		{
			name: "verifyseed optional",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("verifyseed", "abc", 5)
			},
			staticCmd: func() interface{} {
				account := cdrjson.Uint32(5)
				return cdrjson.NewVerifySeedCmd("abc", account)
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifyseed","params":["abc",5],"id":1}`,
			unmarshalled: &cdrjson.VerifySeedCmd{
				Seed:    "abc",
				Account: cdrjson.Uint32(5),
			},
		},
		{
			name: "walletlock",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("walletlock")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewWalletLockCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"walletlock","params":[],"id":1}`,
			unmarshalled: &cdrjson.WalletLockCmd{},
		},
		{
			name: "walletpassphrase",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("walletpassphrase", "pass", 60)
			},
			staticCmd: func() interface{} {
				return cdrjson.NewWalletPassphraseCmd("pass", 60)
			},
			marshalled: `{"jsonrpc":"1.0","method":"walletpassphrase","params":["pass",60],"id":1}`,
			unmarshalled: &cdrjson.WalletPassphraseCmd{
				Passphrase: "pass",
				Timeout:    60,
			},
		},
		{
			name: "walletpassphrasechange",
			newCmd: func() (interface{}, error) {
				return cdrjson.NewCmd("walletpassphrasechange", "old", "new")
			},
			staticCmd: func() interface{} {
				return cdrjson.NewWalletPassphraseChangeCmd("old", "new")
			},
			marshalled: `{"jsonrpc":"1.0","method":"walletpassphrasechange","params":["old","new"],"id":1}`,
			unmarshalled: &cdrjson.WalletPassphraseChangeCmd{
				OldPassphrase: "old",
				NewPassphrase: "new",
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
