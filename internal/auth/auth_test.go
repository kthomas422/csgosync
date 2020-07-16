// Copyright 2020 Kyle Thomas. All rights reserved.

/*
	File:		csgosync/internal/auth/auth_test.go
	Language:	Go 1.14
	Dev Env:	Linux 5.7

	This file contains the functions for testing the authention module with the csgo sync application.
*/

package auth

import "testing"

func TestSanitizeInput(t *testing.T) {
	var got string
	var tests = []struct {
		input, wanted string
	}{
		{"pass\r\n", "pass"},
		{"pass\n\r", "pass"},
		{"pass\r", "pass"},
		{"pass\n", "pass"},
		{"pass\n\n", "pass"},
		{"pass\r\r", "pass"},
		{"p\ra\nss\r", "pass"},
		{"pass", "pass"},
	}
	for _, test := range tests {
		got = sanitizeInput(test.input)
		if got != test.wanted {
			t.Error("got:", got, "wanted:", test.wanted)
		}
	}
}
