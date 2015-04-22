package bbvm
import (
	"bytes"
	"math"
	"strconv"
	"fmt"
	"time"
)
type in struct { }
var IN in

// 2 | 申请字符串句柄 | 申请到的句柄 |  |  IN():SHDL<br>从-1开始查询
// 5 | 复制字符串 | r3的值 | r2:源字符串句柄<br>r3:目标字符串句柄 | IN(r2:SHDL,r3:SHDL):SHDL<br>r3所代表字符串的内容被修改
// 8 | 释放字符串句柄 | r3的值 | r3:字符串句柄 | IN(r3:SHDL):SHDL
// 9 | 比较字符串 | 两字符串的差值 相同为0，大于为1,小于为-1 | r2:基准字符串<br>r3:比较字符串 | IN(r2:SHDL,r3:SHDL):int
//12 | 获取字符的ASCII码 | ASCII码 | r2:字符位置<br>r3:字符串 | GBK编码,返回的结果范围为有符号的 8bit值,因此对中文操作时返回负数
//13 | 将给定字符串中指定索引的字符替换为给定的ASCII代表的字符 | r3的值 | r1:ASCII码<br>r2:字符位置<br>r3:目标字符串 | r3所代表字符串的内容被修改, 要求r3是句柄才能修改r3的值,给出的ASCII会进行模256的处理
func inStrFunc(i *Inst) {
	v, p, o := i.VM, i.B.Get(), i.A // port and param

	switch p{
		case 2:
		if r, err := v.StrPool().Acquire(); err == nil {
			rlog.Info("Acquire StrRes %d", r.Id())
			o.Set(r.Id())
		}else {
			rlog.Error("Acquire StrRes faield: %s", err)
		}
		case 5:
		rlog.Info("StrRes copy '%s' to %d", v.r2.Str(), v.Register(REG_R3).Get())
		v.r3.SetStr(v.r2.Str())
		o.Set(v.r3.Get())
		case 8:
		hdl := v.Register(REG_R3).Get()
		r := v.StrPool().Get(hdl)
		if r != nil {
			rlog.Info("Release StrRes %d '%v'", r.Id(), r.Get())
			v.StrPool().Release(r)
		}else {
			rlog.Info("Release StrRes %d failed: not exists", hdl)
		}
		// TODO 确认是否返回r3的值
		o.Set(v.r3.Get())
		case 9:

		r := bytes.Compare([]byte(v.MustGetStr(v.r3.Get())), []byte(v.MustGetStr(v.r2.Get())))
		//		rlog.Info("Str compare %s %s = %d",v.MustGetStr(v.r2.Get()),v.MustGetStr(v.r3.Get()),r)
		o.Set(r)

		case 12:
		s, i := v.r3.Str(), v.r2.Get()
		if len(s) <i {
			log.Error("Get char of '%s' at %d out of rang", s, i)
			o.Set(0)
		}else {
			o.Set(int(int8(s[i])))
		}
		case 13:

		r, c, i := v.r3.StrRes(), v.r1.Get(), v.r2.Get()
		if r != nil {
			s := r.Get().(string)
			b := []byte(s)
			if len(b)<i {
				log.Error("Set char of '%s'@%d at %d to '%c' failed:out of range", s,v.r3.Get(), i, c)
			}else {
				b[i] = byte(c%256)
				r.Set(string(b))
				break
			}
		}else {
			log.Error("Set char of %d at %d to %c failed:not a str res", v.r3.Get(), i, c)
		}
		o.Set(v.r3.Get())
	}
}


//0 | 浮点数转换为整数 | 整数 | r3:浮点数 | int(r3.float)
//1 | 整数转换为浮点数 | 浮点数 | r3:整数 | float(r3.int)
//2 | 申请字符串句柄 | 申请到的句柄 |  |  strPool.acquire
//3 | 字符串转换为整数 | 整数 | r3:字符串句柄,__地址__ | float(r3.str);若r3的值不是合法的字符串句柄则返回r3的值
//4 | 整数转换为字符串 | 返回的值为r3:整数 | r2:目标字符串_句柄_<br>r3:整数 | r2.str=str(r3.int);return r3.int;r2所代表字符串的内容被修改
//10 | 整数转换为浮点数再转换为字符串 | r3的值 | r2:目标字符串<br>r3:整数 | r2所代表字符串的内容被修改
//11 | 字符串转换为浮点数 | 浮点数 | r3:字符串 |
func inConvFunc(i *Inst) {
	v, p, o := i.VM, i.B.Get(), i.A // port and param

	switch p {
		case 0:
		o.Set(int(math.Float32frombits(uint32(v.r3.Get()))))
		case 1:
		o.SetFloat32(float32(v.r3.Get()))
		case 3:
		if s, fine := v.GetStr(v.Register(REG_R3).Get()); fine {
			if r, ok := strconv.Atoi(s); ok == nil {
				o.Set(r)
			}else {
				log.Error("Convert atoi faield:"+s)
			}
		}else {
			log.Error("GetStr faield")
		}
		case 4:
		if s, fine := v.GetStr(v.Register(REG_R3).Get()); fine {
			if r, ok := strconv.Atoi(s); ok == nil {
				o.Set(r)
			}else {
				log.Error("Convert atoi faield:"+s)
			}
		}else {
			log.Error("GetStr faield")
		}
		case 10:
		v.r2.SetStr(fmt.Sprintf(FORMAT_FLOAT, float32(v.r3.Get())))
		o.Set(v.r3.Get())
		case 11:
		f, err := strconv.ParseFloat(v.r3.Str(), 32)
		if err!= nil {
			log.Error(err.Error())
			o.Set(0)
		}else {
			o.SetFloat32(float32(f))
		}


	}
}

func (in)StrFunc(v VM) {
	v.SetIn(HANDLE_ALL, 2, inStrFunc)
	v.SetIn(HANDLE_ALL, 5, inStrFunc)
	v.SetIn(HANDLE_ALL, 8, inStrFunc)
	v.SetIn(HANDLE_ALL, 9, inStrFunc)
	v.SetIn(HANDLE_ALL, 12, inStrFunc)
	v.SetIn(HANDLE_ALL, 13, inStrFunc)
}
func (in)ConvFunc(v VM) {
	v.SetIn(HANDLE_ALL, 0, inConvFunc)
	v.SetIn(HANDLE_ALL, 1, inConvFunc)
	v.SetIn(HANDLE_ALL, 3, inConvFunc)
	v.SetIn(HANDLE_ALL, 4, inConvFunc)
	v.SetIn(HANDLE_ALL, 10, inConvFunc)
	v.SetIn(HANDLE_ALL, 11, inConvFunc)
}
func (in)Misc(v VM) {
	v.SetIn(HANDLE_ALL, 14, inMisc)
	v.SetIn(HANDLE_ALL, 15, inMisc)
}
//14 | （功用不明） | 65535 |  |
//15 | 获取嘀嗒计数 | 嘀嗒计数 |  | 这里不知道他是怎么算的这个数字,但是会随着时间增长就是了
func inMisc(i *Inst) {
	_, p, o := i.VM, i.B.Get(), i.A // port and param

	switch p{
		case 14:
		o.Set(65535)
		case 15:
		o.Set(int(time.Now().Unix()))
	}
}