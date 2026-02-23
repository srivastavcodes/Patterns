package table_driven_test

import (
	"bytes"
	"slices"
	"testing"
)

var mapTests = map[string]struct {
	input  []byte
	output []byte
}{
	"juhi name": {
		input:  []byte("juhi"),
		output: []byte("ihuj"),
	},
	"parth name": {
		input:  []byte("parth"),
		output: []byte("htrap"),
	},
}

func TestFlagParser(t *testing.T) {
	for name, test := range mapTests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if slices.Reverse(test.input); !bytes.Equal(test.input, test.output) {
				t.Fatalf("got %v, want %v", test.input, test.output)
			}
		})
	}
}
