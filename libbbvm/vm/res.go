package vm
import (
	"io"
	"github.com/spacemonkeygo/errors"
)

var ErrPoolLimitReached = errors.NewClass("No more resource can be acquired")

// 资源对象
type Res interface {
	Id() int
	Get() interface{}
	Set(interface{})
}
type ResValCreator func(ResPool) interface{}

// 资源池
type ResPool interface {
	Acquire() (Res, error)
	Release(Res)
	Get(int) Res
	SetCreator(ResValCreator)
	Creator() ResValCreator
}
type res struct {
	id int
	val interface{}
}
func (r *res)Get() interface{} {
	return r.val
}
func (r *res)Set(v interface{}) {
	r.val = v
}
func (r *res)Id() int {
	return r.id
}
type resPool struct {
	pool  map[int]Res
	start int
	step  int
	current int
	reuse bool
	limit int
	creator ResValCreator
}


func (p *resPool)Acquire() (Res, error) {
	if len(p.pool) >= p.limit {
		return nil, ErrPoolLimitReached.New("Can not acquire %v", p)
	}
	c := p.current
	if p.reuse {
		c = p.start
	}
	for ;; c += p.step {
		if p.pool[c] == nil {
			p.current = c
			return p.create(), nil
		}
	}
	panic("Unreachable")
}
func (p *resPool)create() Res {
	r := &res{p.current, nil}
	if p.creator != nil {
		r.val = p.creator(p)
	}
	p.pool[r.id] = r
	return r
}

func (p *resPool)Release(r Res) {
	if c, ok := r.(io.Closer); ok {
		if err := c.Close(); err != nil {
			log.Error("Release resource failed: %s", err.Error())
		}
	}

	delete(p.pool, r.Id())
}
func (p *resPool)Get(i int) Res {
	return p.pool[i]
}
func (p *resPool)SetCreator(c ResValCreator) {
	p.creator = c
}
func (p *resPool)Creator() ResValCreator {
	return p.creator
}