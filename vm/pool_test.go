package vm
import (
	"testing"
	"fmt"
)

func TestPool(t *testing.T) {
	p := newStrPool()
	r, _ := p.Acquire()

	r.Set("Yes")
	r, _ = p.Acquire()
	r.Set("No")

	fmt.Print(p.Get(-1))
	fmt.Print(p.Get(-2))
}