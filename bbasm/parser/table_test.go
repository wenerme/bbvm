package parser

import (
	"github.com/stretchr/testify/assert"
	"github.com/wenerme/bbvm/bbasm"
	"testing"
)

func TestLookup(t *testing.T) {
	tests := []struct {
		t   interface{}
		v   string
		exp interface{}
	}{
		{bbasm.INT, "Int", bbasm.INT},
		{bbasm.INT, "int", bbasm.INT},
		{bbasm.INT, "word", bbasm.WORD},
		{bbasm.R0, "rP", bbasm.RP},
		{bbasm.R0, "rxz", nil},
	}
	for _, test := range tests {
		assert.Equal(t, test.exp, Lookup(test.t, test.v))
	}
}
