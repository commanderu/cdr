// Copyright (c) 2014 The btcsuite developers
// Copyright (c) 2015-2016 The commanderu developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package cdrjson_test

import (
	"testing"

	"github.com/commanderu/cdrd/cdrjson"
)

// TestErrorCodeStringer tests the stringized output for the ErrorCode type.
func TestErrorCodeStringer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   cdrjson.ErrorCode
		want string
	}{
		{cdrjson.ErrDuplicateMethod, "ErrDuplicateMethod"},
		{cdrjson.ErrInvalidUsageFlags, "ErrInvalidUsageFlags"},
		{cdrjson.ErrInvalidType, "ErrInvalidType"},
		{cdrjson.ErrEmbeddedType, "ErrEmbeddedType"},
		{cdrjson.ErrUnexportedField, "ErrUnexportedField"},
		{cdrjson.ErrUnsupportedFieldType, "ErrUnsupportedFieldType"},
		{cdrjson.ErrNonOptionalField, "ErrNonOptionalField"},
		{cdrjson.ErrNonOptionalDefault, "ErrNonOptionalDefault"},
		{cdrjson.ErrMismatchedDefault, "ErrMismatchedDefault"},
		{cdrjson.ErrUnregisteredMethod, "ErrUnregisteredMethod"},
		{cdrjson.ErrNumParams, "ErrNumParams"},
		{cdrjson.ErrMissingDescription, "ErrMissingDescription"},
		{0xffff, "Unknown ErrorCode (65535)"},
	}

	// Detect additional error codes that don't have the stringer added.
	if len(tests)-1 != int(cdrjson.TstNumErrorCodes) {
		t.Errorf("It appears an error code was added without adding an " +
			"associated stringer test")
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		result := test.in.String()
		if result != test.want {
			t.Errorf("String #%d\n got: %s want: %s", i, result,
				test.want)
			continue
		}
	}
}

// TestError tests the error output for the Error type.
func TestError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   cdrjson.Error
		want string
	}{
		{
			cdrjson.Error{Message: "some error"},
			"some error",
		},
		{
			cdrjson.Error{Message: "human-readable error"},
			"human-readable error",
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		result := test.in.Error()
		if result != test.want {
			t.Errorf("Error #%d\n got: %s want: %s", i, result,
				test.want)
			continue
		}
	}
}
