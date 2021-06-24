package vm

import (
	"bytes"
	"fmt"
	"github.com/op/go-logging"
	"github.com/wenerme/bbvm/libbbvm/asm"
	. "github.com/wenerme/bbvm/libbbvm/asm"
	"math"
	"strconv"
	"strings"
	"time"
)

type in struct{}

var IN in

func (in) Str(v VM) {
	v.SetIn(HANDLE_ALL, 2, inStrFunc)
	v.SetIn(HANDLE_ALL, 5, inStrFunc)
	v.SetIn(HANDLE_ALL, 6, inStrFunc)
	v.SetIn(HANDLE_ALL, 7, inStrFunc)
	v.SetIn(HANDLE_ALL, 8, inStrFunc)
	v.SetIn(HANDLE_ALL, 9, inStrFunc)
	v.SetIn(HANDLE_ALL, 12, inStrFunc)
	v.SetIn(HANDLE_ALL, 13, inStrFunc)
	v.SetIn(HANDLE_ALL, 34, inStrFunc)
	v.SetIn(HANDLE_ALL, 35, inStrFunc)
	v.SetIn(HANDLE_ALL, 36, inStrFunc)
	v.SetIn(HANDLE_ALL, 37, inStrFunc)
	v.SetIn(HANDLE_ALL, 38, inStrFunc)
	v.SetIn(HANDLE_ALL, 39, inStrFunc)
}
func (in) Conv(v VM) {
	v.SetIn(HANDLE_ALL, 0, inConvFunc)
	v.SetIn(HANDLE_ALL, 1, inConvFunc)
	v.SetIn(HANDLE_ALL, 3, inConvFunc)
	v.SetIn(HANDLE_ALL, 4, inConvFunc)
	v.SetIn(HANDLE_ALL, 10, inConvFunc)
	v.SetIn(HANDLE_ALL, 11, inConvFunc)
	v.SetIn(HANDLE_ALL, 32, inConvFunc)
	v.SetIn(HANDLE_ALL, 33, inConvFunc)
}
func (in) Misc(v VM) {
	v.SetIn(HANDLE_ALL, 14, inMiscFunc)
	v.SetIn(HANDLE_ALL, 15, inMiscFunc)
	v.SetIn(HANDLE_ALL, 23, inMiscFunc)
	v.SetIn(HANDLE_ALL, 24, inMiscFunc)
	v.SetIn(HANDLE_ALL, 25, inMiscFunc)
}
func (in) Math(v VM) {
	v.SetIn(HANDLE_ALL, 16, inMathFunc)
	v.SetIn(HANDLE_ALL, 17, inMathFunc)
	v.SetIn(HANDLE_ALL, 18, inMathFunc)
	v.SetIn(HANDLE_ALL, 19, inMathFunc)
	v.SetIn(HANDLE_ALL, 20, inMathFunc)
	v.SetIn(HANDLE_ALL, 21, inMathFunc)
}
func (i in) All(v VM) {
	i.Conv(v)
	i.Str(v)
	i.Misc(v)
	i.Math(v)
}

// 2 | 申请字符串句柄 | 申请到的句柄 |  |  IN():SHDL<br>从-1开始查询
// 5 | 复制字符串 | r3的值 | r2:源字符串句柄<br>r3:目标字符串句柄 | IN(r2:SHDL,r3:SHDL):SHDL<br>r3所代表字符串的内容被修改
// 6 | 连接字符串 | r3的值 | r2:源字符串<br>r3:目标字符串 | r3.str=r3.str+r2.str
// 7 | 获取字符串长度 | 字符串长度 | r3:字符串 | strlen(r3.str)
// 8 | 释放字符串句柄 | r3的值 | r3:字符串句柄 | IN(r3:SHDL):SHDL
// 9 | 比较字符串 | 两字符串的差值 相同为0，大于为1,小于为-1 | r2:基准字符串<br>r3:比较字符串 | IN(r2:SHDL,r3:SHDL):int
//12 | 获取字符的ASCII码 | ASCII码 | r2:字符位置<br>r3:字符串 | GBK编码,返回的结果范围为有符号的 8bit值,因此对中文操作时返回负数
//13 | 将给定字符串中指定索引的字符替换为给定的ASCII代表的字符 | r3的值 | r1:ASCII码<br>r2:字符位置<br>r3:目标字符串 | r3所代表字符串的内容被修改, 要求r3是句柄才能修改r3的值,给出的ASCII会进行模256的处理
//34 | 获取字符的ASCII码 | ASCII码 | r3:字符串 |
//35 | 左取字符串 | r3的值 | r1:截取长度<br>r2:源字符串<br>r3:目标字符串 | r3所代表字符串的内容被修改 （此端口似乎不正常）
//36 | 右取字符串 | r3的值 | r1:截取长度<br>r2:源字符串<br>r3:目标字符串 | r3所代表字符串的内容被修改
//37 | 中间取字符串 | r0截取长度 | r0:截取长度<br>r1:截取位置<br>r2:源字符串<br>r3:目标字符串 | r3所代表字符串的内容被修改
//38 | 查找字符串 | 位置 | r1:起始位置<br>r2:子字符串<br>r3:父字符串 |
//39 | 获取字符串长度 | 字符串长度 | r3:字符串 |
func inStrFunc(i *asm.Inst) {
	v, p, o := i.VM, i.B.Get(), i.A // port and param
	r2, r3 := v.Register(asm.REG_R2), v.Register(asm.REG_R3)
	switch p {
	case 2:
		if r, err := v.StrPool().Acquire(); err == nil {
			log.Info("Acquire StrRes %d", r.Id())
			o.Set(r.Id())
		} else {
			log.Error("Acquire StrRes faield: %s", err)
			o.Set(r3.Get())
		}
	case 5:
		log.Debug("strcpy '%s' to %d", r2, r3.Get())
		r3.SetStr(r2.Str())
		o.Set(r3.Get())
	case 6:
		if log.IsEnabledFor(logging.DEBUG) {
			log.Debug("strcat(%s,%s)", r2.Str(), r3.Str())
		}
		v.r3.SetStr(r3.Str() + r2.Str())
		o.Set(v.r3.Get())
	case 7:
		log.Debug("strlen(%s)", r3.Str())
		if s, err := v.GetStr(r3.Get()); err == nil {
			o.Set(len(s))
		} else {
			log.Error("strlen failed: %s", err.Error())
			o.Set(r3.Get())
		}
	case 8:
		hdl := r3.Get()
		r := v.StrPool().Get(hdl)
		if r != nil {
			log.Info("Release StrRes %d '%v'", r.Id(), r.Get())
			v.StrPool().Release(r)
		} else {
			log.Info("Release StrRes %d failed: not exists", hdl)
		}
		o.Set(r3.Get())
	case 9:
		r := bytes.Compare([]byte(v.MustGetStr(v.r3.Get())), []byte(v.MustGetStr(v.r2.Get())))
		log.Debug(`strcmp("%s", '%s") = %d`, v.MustGetStr(v.r2.Get()), v.MustGetStr(v.r3.Get()), r)
		o.Set(r)

	case 12:
		s, i := v.r3.Str(), v.r2.Get()
		if len(s) < i {
			log.Error("Get char of '%s' at %d out of rang", s, i)
			o.Set(0)
		} else {
			o.Set(int(int8(s[i])))
		}
	case 13:

		r, c, i := v.r3.StrRes(), v.r1.Get(), v.r2.Get()
		if r != nil {
			s := r.Get().(string)
			b := []byte(s)
			if len(b) < i {
				log.Error("Set char of '%s'@%d at %d to '%c' failed:out of range", s, v.r3.Get(), i, c)
			} else {
				b[i] = byte(c % 256)
				r.Set(string(b))
				break
			}
		} else {
			log.Error("Set char of %d at %d to %c failed:not a str res", v.r3.Get(), i, c)
		}
		o.Set(v.r3.Get())

	case 34:
		s, i := v.r3.Str(), 0
		if len(s) < i {
			log.Error("Get char of '%s' at %d out of rang", s, i)
			o.Set(0)
		} else {
			o.Set(int(int8(s[i])))
		}
	case 35:
		s, l, r := v.r2.Str(), v.r1.Get(), v.r3.StrRes()
		if len(s) < l {
			log.Error("Left substring('%s',%d) failed:", s, l)
		} else {
			r.Set(string(s[:l]))
		}
		o.Set(v.r3.Get())
	case 36:
		s, l, r := v.r2.Str(), v.r1.Get(), v.r3.StrRes()
		if len(s) < l {
			log.Error("Right substring('%s',%d) failed:", s, l)
		} else {
			r.Set(string(s[len(s)-l:]))
		}
		o.Set(v.r3.Get())
	case 37:
		l, p, s, r := v.r0.Get(), v.r1.Get(), v.r2.Str(), v.r3.StrRes()
		r.Set(string(s[p : p+l]))
		o.Set(l)
	case 38:
		s, sub, i := v.r3.Str(), v.r2.Str(), v.r1.Get()
		if len(s) > i {
			o.Set(strings.Index(s[i:], sub) + i)
		} else {
			log.Error("Index '%s' of '%s' at %d: out of range", sub, s, i)
			o.Set(-1)
		}
	case 39:
		o.Set(len(v.r3.Str()))
	}
}

//0 | 浮点数转换为整数 | 整数 | r3:浮点数 | int(r3.float)
//1 | 整数转换为浮点数 | 浮点数 | r3:整数 | float(r3.int)
//2 | 申请字符串句柄 | 申请到的句柄 |  |  strPool.acquire
//3 | 字符串转换为整数 | 整数 | r3:字符串句柄,__地址__ | float(r3.str);若r3的值不是合法的字符串句柄则返回r3的值
//4 | 整数转换为字符串 | 返回的值为r3:整数 | r2:目标字符串_句柄_<br>r3:整数 | r2.str=str(r3.int);return r3.int;r2所代表字符串的内容被修改
//10 | 整数转换为浮点数再转换为字符串 | r3的值 | r2:目标字符串<br>r3:整数 | r2所代表字符串的内容被修改
//11 | 字符串转换为浮点数 | 浮点数 | r3:字符串 |
//32 | 整数转换为字符串 | r3的值 | r1:整数<br>r3:目标字符串 | r3所代表字符串的内容被修改
//33 | 字符串转换为整数 | 整数 | r3:字符串 |
func inConvFunc(i *Inst) {
	v, p, o := i.VM, i.B.Get(), i.A // port and param

	switch p {
	case 0:
		o.Set(int(math.Float32frombits(uint32(v.r3.Get()))))
	case 1:
		o.SetFloat32(float32(v.r3.Get()))
	case 3:
		s, err := v.GetStr(v.r3.Get())
		if err != nil {
			o.Set(v.r3.Get())
			break
		}
		if r, err := strconv.Atoi(s); err == nil {
			o.Set(r)
		} else {
			log.Error("Atoi(%s) faield: %s", s, err)
		}
	case 4:
		v.r2.SetStr(strconv.Itoa(v.r3.Get()))
		o.Set(v.r3.Get())
	case 10:
		v.r2.SetStr(float32ToStr(float32(v.r3.Get())))
		o.Set(v.r3.Get())
	case 11:
		f, err := strconv.ParseFloat(v.r3.Str(), 32)
		if err != nil {
			log.Error(err.Error())
			o.Set(0)
		} else {
			o.SetFloat32(float32(f))
		}

	case 32:
		v.r3.SetStr(fmt.Sprint(v.r1.Get()))
		o.Set(v.r3.Get())
	case 33:
		if f, ok := strconv.ParseFloat(v.r3.Str(), 32); ok == nil {
			o.Set(int(int32(f)))
		} else {
			log.Error(ok.Error())
			o.Set(0)
		}
	}
}

//14 | （功用不明） | 65535 |  |
//15 | 获取嘀嗒计数 | 嘀嗒计数 |  | 这里不知道他是怎么算的这个数字,但是会随着时间增长就是了
//23 | 读内存数据 | 地址内容 | r3:地址 |
//24 | 写内存数据 | r3的值 | r2:待写入数据<br>r3:待写入地址 |
//25 | 获取环境值 | 环境值 |  |
func inMiscFunc(i *Inst) {
	v, p, o := i.VM, i.B.Get(), i.A // port and param

	switch p {
	case 14:
		o.Set(65535)
	case 15:
		o.Set(int(time.Now().Unix()))
	case 23:
		o.Set(v.GetInt(v.r3.Get()))
	case 24:
		v.SetInt(v.r3.Get(), v.r2.Get())
		o.Set(v.r3.Get())
	case 25:
		o.Set(0) // FIXME 环境值为0
	}
}

//16 | 求正弦值 | X!的正弦值 | r3:X! |
//17 | 求余弦值 | X!的余弦值 | r3:X! |
//18 | 求正切值 | X!的正切值 | r3:X! |
//19 | 求平方根值 | X!的平方根值 | r3:X! |
//20 | 求绝对值 | X%的绝对值 | r3:X% |
//21 | 求绝对值 | X!的绝对值 | r3:X! |
func inMathFunc(i *Inst) {
	v, p, o := i.VM, i.B.Get(), i.A // port and param
	a, b := float64(v.r3.Float32()), float64(0)
	switch p {
	case 16:
		b = math.Sin(a)
	case 17:
		b = math.Cos(a)
	case 18:
		b = math.Tan(a)
	case 19:
		b = math.Sqrt(a)
	case 20:
		i := v.r3.Get()
		if i < 0 {
			i = -i
		}
		o.Set(i)
		log.Error("DO IN 20 %d", v.r3.Get())
		return
	case 21:
		b = math.Abs(a)
	}

	o.SetFloat32(float32(b))
}
