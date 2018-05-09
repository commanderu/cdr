rpcclient
=========

[![Build Status](http://img.shields.io/travis/commanderu/cdrd.svg)](https://travis-ci.org/commanderu/cdrd)
[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/commanderu/cdrd/rpcclient)

rpcclient implements a Websocket-enabled commanderu JSON-RPC client package written
in [Go](http://golang.org/).  It provides a robust and easy to use client for
interfacing with a commanderu RPC server that uses a cdrd compatible commanderu
JSON-RPC API.

## Status

This package is currently under active development.  It is already stable and
the infrastructure is complete.  However, there are still several RPCs left to
implement and the API is not stable yet.

## Documentation

* [API Reference](http://godoc.org/github.com/commanderu/cdrd/rpcclient)
* [cdrd Websockets Example](https://github.com/commanderu/cdrd/tree/master/rpcclient/examples/cdrdwebsockets)
  Connects to a cdrd RPC server using TLS-secured websockets, registers for
  block connected and block disconnected notifications, and gets the current
  block count
* [cdrwallet Websockets Example](https://github.com/commanderu/cdrd/tree/master/rpcclient/examples/cdrwalletwebsockets)  
  Connects to a cdrwallet RPC server using TLS-secured websockets, registers for
  notifications about changes to account balances, and gets a list of unspent
  transaction outputs (utxos) the wallet can sign

## Major Features

* Supports Websockets (cdrd/cdrwallet) and HTTP POST mode (bitcoin core-like)
* Provides callback and registration functions for cdrd/cdrwallet notifications
* Supports cdrd extensions
* Translates to and from higher-level and easier to use Go types
* Offers a synchronous (blocking) and asynchronous API
* When running in Websockets mode (the default):
  * Automatic reconnect handling (can be disabled)
  * Outstanding commands are automatically reissued
  * Registered notifications are automatically reregistered
  * Back-off support on reconnect attempts

## Installation

```bash
$ go get -u github.com/commanderu/cdrd/rpcclient
```

## License

Package rpcclient is licensed under the [copyfree](http://copyfree.org) ISC
License.
