package bbvm

import (
	"context"
	"encoding/binary"
	"github.com/juju/errors"
	"github.com/wenerme/bbvm/bbasm"
	"go.uber.org/zap"
	"math"
)

var _ bbasm.Runtime = &Instance{}

type Instance struct {
	Memory []byte
	Order  binary.ByteOrder

	RP bbasm.Register
	RF bbasm.Register
	RS bbasm.Register
	RB bbasm.Register
	R0 bbasm.Register
	R1 bbasm.Register
	R2 bbasm.Register
	R3 bbasm.Register

	Inst *bbasm.Inst
	Std  *Std

	done bool
}

func (rt *Instance) GetInt(addr int) int {
	return int(rt.Order.Uint32(rt.Memory[addr:]))
}

func (rt *Instance) SetInt(addr int, val int) {
	rt.Order.PutUint32(rt.Memory[addr:], uint32(val))
}
func (rt *Instance) GetFloat(addr int) float32 {
	return math.Float32frombits(uint32(rt.GetInt(addr)))
}

func (rt *Instance) SetFloat(addr int, val float32) {
	rt.SetInt(addr, int(math.Float32bits(val)))
}

// GetString return \0 terminated string in memory
func (rt *Instance) GetString(addr int) string {
	l := 0
	slice := rt.Memory[addr:]
	for i, v := range slice {
		if v == 0 {
			l = i
			break
		}
	}
	s, err := rt.Std.BytesToString(slice[:l])
	if err != nil {
		panic(err)
	}
	zap.S().Debugf("GetString %v+%v -> %v", addr, l, s)
	return s
}
func (rt *Instance) Push(val int) {
	rt.SetInt(rt.RS.Get(), val)
	rt.RS.Set(rt.RS.Get() + 4)
}

func (rt *Instance) Pop() int {
	rt.RS.Set(rt.RS.Get() - 4)
	return rt.GetInt(rt.RS.Get())
}

func (rt *Instance) Register(registerType bbasm.RegisterType) bbasm.Register {
	switch registerType {
	case bbasm.RP:
		return rt.RP
	case bbasm.RF:
		return rt.RF
	case bbasm.RS:
		return rt.RS
	case bbasm.RB:
		return rt.RB
	case bbasm.R0:
		return rt.R0
	case bbasm.R1:
		return rt.R1
	case bbasm.R2:
		return rt.R2
	case bbasm.R3:
		return rt.R3
	}
	panic(errors.New("invalid " + registerType.String()))
}

func (rt *Instance) Jump(addr int) {
	rt.RP.Set(addr)
}

func (rt *Instance) Exit() {
	rt.done = true
}

func (rt *Instance) In(ctx context.Context, a int, b int) {
	rt.Std.Execute(ctx, rt, rt.Inst)
}

func (rt *Instance) Out(ctx context.Context, a int, b int) {
	rt.Std.Execute(ctx, rt, rt.Inst)
}

func NewInstance() *Instance {
	sr := &Instance{
		Order: binary.LittleEndian,
		RP:    &Register{Label: "RP"},
		RF:    &Register{Label: "RF"},
		RS:    &Register{Label: "RS"},
		RB:    &Register{Label: "RB"},
		R0:    &Register{Label: "R0"},
		R1:    &Register{Label: "R1"},
		R2:    &Register{Label: "R2"},
		R3:    &Register{Label: "R3"},
		Inst:  &bbasm.Inst{},
	}
	sr.Inst.A.RT = sr
	sr.Inst.B.RT = sr
	return sr
}

func (rt *Instance) Load(rom []byte) {
	rt.done = false
	rt.Memory = make([]byte, len(rom))
	copy(rt.Memory, rom)
	// stack
	rt.Memory = append(rt.Memory, make([]byte, 1024)...)
	rt.Reset()
}

func (rt *Instance) Reset() {
	rt.RP.Set(0)
	rt.RF.Set(0)
	rt.RB.Set(len(rt.Memory) - 1024)
	rt.RS.Set(rt.RB.Get())
	rt.R0.Set(0)
	rt.R1.Set(0)
	rt.R2.Set(0)
	rt.R3.Set(0)
}
func (rt *Instance) Run(ctx context.Context) error {
	for !rt.done {
		if err := rt.Step(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (rt *Instance) Step(ctx context.Context) (err error) {
	rp := rt.RP
	a := rp.Get()
	mem := rt.Memory
	if rp.Get() >= len(mem)-1024 {
		rt.done = true
		return nil
	}

	inst := rt.Inst
	err = inst.UnmarshalBinary(mem[rp.Get():])
	if err != nil {
		return err
	}
	// todo return error
	bbasm.Execute(ctx, rt, inst)
	if rp.Get() == a {
		rp.Set(a + inst.Opcode.Len())
	}
	return
}

type Register struct {
	V     int
	Label string
}

func (sr *Register) Float() float32 {
	return math.Float32frombits(uint32(sr.V))
}
func (sr *Register) SetFloat(v float32) {
	sr.V = int(math.Float32bits(v))
}
func (sr *Register) Get() int {
	return sr.V
}
func (sr *Register) Set(v int) {
	sr.V = v
}
