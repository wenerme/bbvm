package parser

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestParseInt(t *testing.T) {
	ERR := int32(math.MaxInt32 - 2)
	tests := []struct {
		Val string
		I   int32
	}{
		{"10", 10},
		{"0x10", 16},
		{"0X10", 16},
		{"0b10", 2},
		{"0B10", 2},
		{"010", 8},

		{"1A", ERR},
		{"0x1G", ERR},
		{"0y10", ERR},
		{"0b102", ERR},
		{"0B1a0", ERR},
		{"019", ERR},
	}

	for _, v := range tests {
		i, e := parseInt(v.Val)
		if v.I == ERR {
			assert.Error(t, e)
		} else {
			assert.NoError(t, e)
			assert.EqualValues(t, v.I, i)
		}
	}
}
