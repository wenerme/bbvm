package bbvm
import (
	"fmt"
	"math"
)

type InstHandler func(*Inst)
const HANDLE_ALL = math.MaxInt32


type VM interface {
	SetOut(int, int, InstHandler)
	SetIn(int, int, InstHandler)
	StrPool() ResPool
	MustGetStr(int) string
	Attr() map[string]interface{}
}

type Register interface {
	Get() int
	Set(int)
}
type register struct {
	Val int
	VM *vm
}

func (r *register)Get() int {
	return r.Val
}
func (r *register)Set(v int) {
	r.Val = v
}
func (o *register)Float32() float32 {
	return math.Float32frombits(uint32(o.Get()))
}
func (o *register)SetFloat32(v float32) {
	o.Set(int(math.Float32bits(v)))
}
func (o *register)StrRes() Res {
	return o.VM.StrPool().Get(o.Get())
}
func (o *register)Str() string {
	if s, ok := o.VM.GetStr(o.Get()); ok {
		return s
	}else {
		log.Error("register string res %d not exists", o.Get())
		return ""
	}
}
func (o *register)SetStr(v string) {
	if r := o.StrRes(); r!= nil {
		o.StrRes().Set(v)
	}else {
		log.Error("register string res %d not exists", o.Get())
	}
}

type monitorRegister struct {
	register
	Changed bool
}

func (r *monitorRegister)Set(v int) {
	r.Val=v
	r.Changed = true
}

func (r *monitorRegister)Get() int {
	return r.Val
}


type vm struct {
	mem []byte
	in map[string]InstHandler
	out map[string]InstHandler

	rp monitorRegister
	rf register
	rs register
	rb register
	r0 register
	r1 register
	r2 register
	r3 register

	inst Inst

	exited bool

	stack []byte

	strPool ResPool
	// 存储各个模块的上下文信息
	attr map[string]interface{}
}
func (v *vm)Load(rom []byte) {
	log.Info("Load ROM. size: %d", len(rom))
	v.Reset()
	v.mem = make([]byte, len(rom))
	copy(v.mem, rom)
	v.rb.Set(len(v.mem))
	v.rs.Set(v.rb.Get())

	log.Info("Init stack. size: %d", 1024)
	v.mem = append(v.mem, make([]byte, 1024)...)
	v.stack = v.mem[len(rom):]
}
func (v *vm)Reset() {
	v.r0.Set(0)
	v.r1.Set(0)
	v.r2.Set(0)
	v.r3.Set(0)
	v.rs.Set(0)
	v.rb.Set(0)
	v.rp.Set(0)
	v.rf.Set(0)
}
func (v *vm)SetIn(a, b int, h InstHandler) {
	v.in[multiPortKey(a, b)] = h
}
func (v *vm)SetOut(a, b int, h InstHandler) {
	v.out[multiPortKey(a, b)] = h
}
func (v *vm)In(a, b int) InstHandler {
	return multiPortHandler(a, b, v.in)
}
func (v *vm)Out(a, b int) InstHandler {
	return multiPortHandler(a, b, v.out)
}


func (v *vm)Register(t RegisterType) Register {
	switch t{
		case REG_RP: return &v.rp
		case REG_RF: return &v.rf
		case REG_RS: return &v.rs
		case REG_RB: return &v.rb
		case REG_R0: return &v.r0
		case REG_R1: return &v.r1
		case REG_R2: return &v.r2
		case REG_R3: return &v.r3
	}
	return nil
}
func (v *vm)Loop() {
	v.rp.Changed = false

	if v.rp.Get() >= len(v.mem) - len(v.stack) {
		log.Info("Run over, exit")
		v.Exit()
		return
	}

	err := v.inst.UnmarshalBinary(v.mem[v.rp.Get():])
	if err != nil { panic(err)}
	v.Proc()
	if !v.rp.Changed {
		v.rp.Set(v.rp.Get() + v.inst.Opcode.Len())
	}
}

func (v *vm)Report() string {
	s := fmt.Sprintf("%s #rp:%d rf:%d rs:%d rb:%d r0:%d r1:%d r2:%d r3:%d",
	v.inst, v.rp.Get(), v.rf.Get(), v.rs.Get(), v.rb.Get(),
	v.r0.Get(), v.r1.Get(), v.r2.Get(), v.r3.Get(), )
	return s
}
func (v *vm)Exit() {
	v.exited = true
}
func (v *vm)IsExited() bool {
	return v.exited
}

func (v *vm)StrPool() ResPool {
	return v.strPool
}

func (v *vm)Attr() map[string]interface{} {
	return v.attr
}



func NewVM() *vm {
	v := &vm{
		attr: make(map[string]interface{}),
		in:make(map[string]InstHandler),
		out:make(map[string]InstHandler),
	}
	v.strPool = newStrPool()
	v.inst.VM = v
	v.inst.A.VM = v
	v.inst.B.VM = v
	v.r0.VM = v
	v.r1.VM = v
	v.r2.VM = v
	v.r3.VM = v
	v.rs.VM = v
	v.rb.VM = v
	v.rp.VM = v

	v.Reset()

	return v
}



