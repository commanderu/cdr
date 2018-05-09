txscript
========

[![Build Status](http://img.shields.io/travis/commanderu/cdrd.svg)](https://travis-ci.org/commanderu/cdrd)
[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/commanderu/cdrd/txscript)

Package txscript implements the commanderu transaction script language.  There is
a comprehensive test suite.

This package has intentionally been designed so it can be used as a standalone
package for any projects needing to use or validate commanderu transaction scripts.

## commanderu Scripts

commanderu provides a stack-based, FORTH-like language for the scripts in
the commanderu transactions.  This language is not turing complete
although it is still fairly powerful.

## Installation and Updating

```bash
$ go get -u github.com/commanderu/cdrd/txscript
```

## Examples

* [Standard Pay-to-pubkey-hash Script](http://godoc.org/github.com/commanderu/cdrd/txscript#example-PayToAddrScript)  
  Demonstrates creating a script which pays to a commanderu address.  It also
  prints the created script hex and uses the DisasmString function to display
  the disassembled script.

* [Extracting Details from Standard Scripts](http://godoc.org/github.com/commanderu/cdrd/txscript#example-ExtractPkScriptAddrs)  
  Demonstrates extracting information from a standard public key script.

* [Manually Signing a Transaction Output](http://godoc.org/github.com/commanderu/cdrd/txscript#example-SignTxOutput)  
  Demonstrates manually creating and signing a redeem transaction.

## License

Package txscript is licensed under the [copyfree](http://copyfree.org) ISC
License.