// token_test.go
// This file is part of Epixel CRM Software
// File started in 2022-03

package epicrm_apiparts

import (
	"testing"
)

func TestExtractBearerTokenFromHeader(t *testing.T) {
	// TODO cases: illegal characters, newlines, escapes, etc.
	expouts := map[string][]byte {
		"": nil,
		"Bearer": nil,
		"Bearer ": nil,
		"Bearer TestToken": []byte("TestToken"),
	}

	for inp, eout := range expouts {
		out := extractBearerTokenFromHeader(inp)

		// Only one of them is nil or they do not match
		if ((eout == nil) != (out == nil)) || (string(out) != string(eout)) {
			t.Errorf("ExtractFromHeader(%q) returned %#v (expected %#v)\n",
				inp, out, eout)
		}
	}
}
