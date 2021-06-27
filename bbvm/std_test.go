package bbvm

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStdUse(t *testing.T) {
	a := &Std{
		PrintInt: func(ctx context.Context, v int) {
			t.Fail()
		},
		PrintChar: func(ctx context.Context, v int) {

		},
	}
	b := &Std{
		PrintInt: func(ctx context.Context, v int) {
		},
		VmTest: func(ctx context.Context) {
		},
	}
	a.Use(b)
	assert.NotNil(t, a.VmTest)
	assert.NotNil(t, a.PrintChar)
	a.PrintInt(nil, 0)
}
