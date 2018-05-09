// Copyright (c) 2013-2014 The btcsuite developers
// Copyright (c) 2015-2016 The commanderu developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package cdrutil

const (
	// AtomsPerCent is the number of atomic units in one coin cent.
	AtomsPerCent = 1e6

	// AtomsPerCoin is the number of atomic units in one coin.
	AtomsPerCoin = 1e8

	// MaxAmount is the maximum transaction amount allowed in atoms.
	// commanderu - Changeme for release
	MaxAmount = 21e6 * AtomsPerCoin
)