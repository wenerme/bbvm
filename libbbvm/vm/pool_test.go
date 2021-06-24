package vm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPool(t *testing.T) {
	assert := assert.New(t)
	p := newStrPool()
	r, _ := p.Acquire()

	r.Set("Yes")
	r, _ = p.Acquire()
	r.Set("No")

	assert.Equal("Yes", p.Get(-1).Get())
	assert.Equal("No", p.Get(-2).Get())
}
