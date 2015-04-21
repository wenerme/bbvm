package bbvm
import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strconv"
	log "github.com/golang/glog"
)

type InstHandler func(*Inst)
const HANDLE_ALL = math.MaxInt32

type Register interface {
	Get() int
	Set(int)
}
type register struct {
	val int
}
func (r *register)Get() int {
	return r.val
}
func (r *register)Set(v int) {
	r.val = v
}

type monitorRegister struct {
	val int
	Changed bool
}

func (r *monitorRegister)Set(v int) {
	r.val=v
	r.Changed = true
}

func (r *monitorRegister)Get() int {
	return r.val
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
}
func (v *vm)Load(rom []byte) {
	v.mem = make([]byte, len(rom))
	copy(v.mem, rom)
	v.rb.Set(len(v.mem))
	v.rs.Set(v.rb.Get())
	//	v.stack = make([]byte, 1024)
	v.mem = append(v.mem, make([]byte, 1024)...)
	v.stack = v.mem[len(rom):]
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

func multiPortHandler(a, b int, m map[string]InstHandler) InstHandler {
	h, exists := m[multiPortKey(a, b)]
	if !exists {
		h, exists = m[multiPortKey(HANDLE_ALL, b)]
	}
	if !exists {
		h, exists = m[multiPortKey(a, HANDLE_ALL)]
	}
	return h
}
func multiPortKey(a, b int) string {
	var sa, sb string
	if a == HANDLE_ALL {
		sa = "*"
	}else {
		sa = strconv.Itoa(a)
	}
	if b == HANDLE_ALL {
		sb = "*"
	}else {
		sb = strconv.Itoa(b)
	}
	return sa+","+sb
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
	err := v.inst.UnmarshalBinary(v.mem[v.rp.Get():])
	if err != nil { panic(err)}
	v.Proc()
	if !v.rp.Changed {
		v.rp.Set(v.rp.Get() + v.inst.Opcode.Len())
	}
}
func (v *vm)Pop() int {
	v.rs.Set(v.rs.Get() - 4)
	i := binary.LittleEndian.Uint32(v.mem[v.rs.Get():])
	return int(i)
}
func (v *vm)Push(i int) {
	binary.LittleEndian.PutUint32(v.mem[v.rs.Get():], uint32(i))
	v.rs.Set(v.rs.Get() + 4)
}
func (v *vm)Call(addr int) {
	v.Push(v.rp.Get() + v.inst.Opcode.Len())
	v.Jump(addr)
}
func (v *vm)Jump(addr int) {
	v.rp.Set(addr)
}

func (v *vm)GetInt(addr int) int {
	return int(binary.LittleEndian.Uint32(v.mem[addr:]))
}

func (v *vm)SetInt(addr int, i int) {
	binary.LittleEndian.PutUint32(v.mem[addr:], uint32(i))
}

func (v *vm)GetStr(addr int) (string, bool) {
	if addr < 0 {
		if s, ok := v.strPool.Get(addr).Get().(string); ok {
			return s, true
		}
		return "", false
	}
	end := addr
	for ; v.mem[end] != 0; end +=1 {}
	if end == addr { return "", true}
	return string(v.mem[addr:end]), true
}

func (v *vm)Report() string {
	s := fmt.Sprintf("# rp:%d rf:%d rs:%d rb:%d r0:%d r1:%d r2:%d r3:%d\n%s",
	v.rp.Get(), v.rf.Get(), v.rs.Get(), v.rb.Get(),
	v.r0.Get(), v.r1.Get(), v.r2.Get(), v.r3.Get(), v.inst)
	return s
}
func (v *vm)Exit() {
	v.exited = true
}
func (v *vm)IsExited() bool {
	return v.exited
}

func (v *vm)Proc() {
	i := &v.inst
	switch v.inst.Opcode{
		case OP_EXIT:
		v.Exit()
		case OP_NOP:
		// NOP
		case OP_CAL:
		i.B.Set(calculate(i.A.Get(), i.B.Get(), i.CalculateType, i.DataType))
		case OP_CALL:
		v.Call(i.A.Get())
		case OP_JMP:
		v.Jump(i.A.Get())
		case OP_IN:
		h := v.In(i.A.Get(), i.B.Get())
		if h == nil {log.Error("Can not handle IN ", i)}else {h(i)}
		case OP_OUT:
		h := v.Out(i.A.Get(), i.B.Get())
		if h == nil {log.Error("Can not handle OUT ", i)}else {h(i)}
		case OP_RET:
		v.Jump(v.Pop())
		case OP_CMP:
		case OP_JPC:
		case OP_PUSH:
		v.Push(i.A.Get())
		case OP_POP:
		i.A.Set(v.Pop())
		case OP_LD:
		// 没考虑类型
		i.A.Set(i.B.Get())
	}
}

func calculate(a int, b int, t CalculateType, d DataType) int {
	var oa, ob, oc float64
	if d == T_FLOAT {
		oa = float64(math.Float32frombits(uint32(a)))
		ob = float64(math.Float32frombits(uint32(b)))
	}else {
		oa = float64(a)
		ob = float64(b)
	}
	switch t{
		case CAL_ADD: oc = oa + ob
		case CAL_SUB: oc = oa - ob
		case CAL_MUL: oc = oa * ob
		case CAL_DIV: oc = oa / ob
		case CAL_MOD: oc = float64(int(oa) % int(ob))
		default:
		log.Error("Unhandled calculate")
	}
	var ret int
	switch d{
		case T_DWORD: ret = int(oc)
		case T_WORD: ret = int(int16(oc))
		case T_BYTE: ret = int(byte(oc))
		case T_FLOAT: ret = int(math.Float32bits(float32(oc)))
		case T_INT: ret = int(oc)
	}
	return ret
}

func NewVM() *vm {
	v := &vm{
		in:make(map[string]InstHandler),
		out:make(map[string]InstHandler),

	}
	v.inst.VM = v
	v.inst.A.Vm = v
	v.inst.B.Vm = v
	return v
}

type VM interface {

}



type Inst struct {
	DataType DataType
	CompareType CompareType
	CalculateType CalculateType
	Opcode Opcode
	A Operand
	B Operand

	VM *vm
}
type Operand struct {
	Val uint32
	AddrMode AddressMode
	Vm *vm
}
func (o *Operand)Get() int {
	switch (o.AddrMode) {
		case AM_REGISTER:
		return o.Vm.Register(RegisterType(o.Val)).Get();
		case AM_REGISTER_DEFERRED:
		return o.Vm.GetInt(o.Vm.Register(RegisterType(o.Val)).Get());
		case AM_IMMEDIATE:
		return int(o.Val);
		case AM_DIRECT:
		return o.Vm.GetInt(int(o.Val));
	}
	panic("Unexcepted")
}
func (o *Operand)Uint() uint {
	return uint(o.Get())
}
func (o *Operand)Float() float32 {
	return math.Float32frombits(uint32(o.Get()))
}
func (o *Operand)SetFloat(v float32) {
	o.Set(int(math.Float32bits(v)))
}
func (o *Operand)Str() string {
	return "Operand STR"
}
func (o *Operand)SetStr(v string) {

}
func (o *Operand)Set(i int) {
	switch (o.AddrMode) {
		case AM_REGISTER:
		o.Vm.Register(RegisterType(o.Val)).Set(i);
		case AM_REGISTER_DEFERRED:
		o.Vm.SetInt(o.Vm.Register(RegisterType(o.Val)).Get(), i);
		case AM_IMMEDIATE:
		log.Error("Set a IMMEDIATE operand")
		case AM_DIRECT:
		o.Vm.SetInt(int(o.Val), i);
		default:
		panic("Unexcepted")
	}
}
func (o Operand)String() string {
	switch o.AddrMode{
		case AM_REGISTER:
		return fmt.Sprintf("%s", RegisterType(o.Val))
		case AM_REGISTER_DEFERRED:
		return fmt.Sprintf("[ %s ]", RegisterType(o.Val))
		case AM_IMMEDIATE:
		return fmt.Sprintf("%d", o.Val)
		case AM_DIRECT:
		return fmt.Sprintf("[ %d ]", o.Val)
	}
	return "UNKNOWN Operand"
}
var ErrDataToShort = errors.New("Data to short to UnmarshalBinary")

func (i *Inst)MarshalBinary() ([]byte, error) {
	data := make([]byte, i.Opcode.Len())
	data[0] = uint8(i.Opcode) << 4

	switch i.Opcode.Len() {
		case 1: // 无操作数
		case 5: // 一个操作数
		data[0] |= uint8(i.A.AddrMode)
		i.A.AddrMode = AddressMode(data[0] & 0xF)
		binary.LittleEndian.PutUint32(data[1:], i.A.Val)
		case 10: // 两个操作数
		data[0] |= uint8(i.DataType)

		data[1] |= uint8(i.A.AddrMode << 2 | i.B.AddrMode)
		binary.LittleEndian.PutUint32(data[2:], i.A.Val)
		binary.LittleEndian.PutUint32(data[6:], i.B.Val)

		if i.Opcode == OP_CAL {
			data[1] |= uint8(i.CalculateType) << 4
		}

		case 6: // JPC
		data[0] |= uint8(i.CompareType)
		data[1] |= uint8(i.A.AddrMode)
		binary.LittleEndian.PutUint32(data[2:], i.A.Val)
	}
	return data, nil
}
func (i *Inst)UnmarshalBinary(data []byte) error {
	/*
指令码 + 数据类型 + 特殊用途字节 + 寻址方式 + 第一个操作数 + 第二个操作数
0x 0       0         0           0        00000000     00000000

无操作数 1byte
指令码 + 无用
0x 0       0
一个操作数 5byte
指令码 + 寻址方式 + 第一个操作数
0x 0       0        00000000
两个操作数 10byte
指令码 + 数据类型 + 保留字节 + 寻址方式 + 第一个操作数 + 第二个操作数
0x 0       0         0        0        00000000     00000000
JPC指令 6byte
指令码 + 比较操作 + 保留字节 + 寻址方式 + 第一个操作数
0x 0       0         0        0        00000000
	*/
	i.Opcode = Opcode(data[0] >> 4)
	if i.Opcode.Len() > len(data) {
		return ErrDataToShort
	}
	switch i.Opcode.Len() {
		case 1: // 无操作数
		case 5: // 一个操作数
		i.A.AddrMode = AddressMode(data[0] & 0xF)
		i.A.Val = binary.LittleEndian.Uint32(data[1:])
		case 10: // 两个操作数
		i.DataType = DataType(data[0] & 0xF)

		addrMode := data[1] & 0xF
		i.A.AddrMode = AddressMode(addrMode/4)
		i.B.AddrMode = AddressMode(addrMode%4)
		i.A.Val = binary.LittleEndian.Uint32(data[2:])
		i.B.Val = binary.LittleEndian.Uint32(data[6:])

		if i.Opcode == OP_CAL {
			i.CalculateType = CalculateType(data[1] >> 4)
		}

		case 6: // JPC
		i.CompareType = CompareType(data[0] & 0xF)
		i.A.AddrMode = AddressMode(data[1] & 0xF)
		i.A.Val = binary.LittleEndian.Uint32(data[2:])
	}
	return nil
}

func (i Inst)String() string {
	switch i.Opcode.Len(){
		case 1:
		return fmt.Sprint(i.Opcode)
		case 5:
		return fmt.Sprintf("%s %s", i.Opcode, i.A)
		case 10:
		switch i.Opcode{
			case OP_CAL:
			return fmt.Sprintf("%s %s %s %s, %s", i.Opcode, i.DataType, i.CalculateType, i.A, i.B)
			case OP_LD:
			return fmt.Sprintf("%s %s %s, %s", i.Opcode, i.DataType, i.A, i.B)
			case OP_CMP:
			return fmt.Sprintf("%s %s %s, %s", i.Opcode, i.CompareType, i.A, i.B)
		}
		return fmt.Sprintf("%s %s %s", i.Opcode, i.A, i.B)
		case 6:
		return fmt.Sprintf("%s %s %s", i.Opcode, i.CompareType, i.A)
	}
	return "Unknown opcodec len"
}