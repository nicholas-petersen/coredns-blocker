package blocker

import (
	"strings"
	"testing"

	"github.com/coredns/caddy"
)

func TestSetup(t *testing.T) {
	tests := []struct {
		input string
		shouldErr   bool
		expectedErr string
	}{
		// positive
		{"blocker hosts", false, ""},
		//negative
		{"blocker", true, "Wrong argument count or unexpected line ending after 'blocker'"},
		{"blocker hosts hosts", true, "Wrong argument count or unexpected line ending after 'hosts'"},
	}

	for i, test := range tests {
		c := caddy.NewTestController("dns", test.input)
		_, err := fileParse(c)

		if test.shouldErr && err == nil {
			t.Errorf("Test %d: expected error but found %s for input %s", i, err, test.input)
		}

		if err != nil {
			if !test.shouldErr {
				t.Fatalf("Test %d: expected no error but found one for input %s, got: %v", i, test.input, err)
			}

			if !strings.Contains(err.Error(), test.expectedErr) {
				t.Error(err)
				t.Errorf("Test %d: expected error to contain: %v, found error: %v, input: %s", i, test.expectedErr, err, test.input)
			}
		}
 	}
}
