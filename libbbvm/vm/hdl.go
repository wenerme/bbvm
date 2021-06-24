package vm

import (
	. "github.com/wenerme/bbvm/libbbvm/asm"
)

// 非单独的 IN 和 OUT 端口处理
type miscHandler struct{}

var Misc miscHandler

func (o miscHandler) All(v VM) {
	o.Data(v)
}
func (miscHandler) Data(v VM) {
	v.Attr()["data-ptr"] = 0
	v.SetOut(13, 0, hdlDataFunc)
	v.SetOut(14, 0, hdlDataFunc)
	v.SetOut(15, 0, hdlDataFunc)
	v.SetIn(HANDLE_ALL, 22, hdlDataFunc)
}

// IN
//22 | 重定位数据指针 | r3的值 | r2:数据位置 | r3中为任意值
// OUT
//13 | 从数据区读取整数 | 0 |  | r3的值变为读取的整数
//14 | 从数据区读取字符串 | 0 | r3:目标字符串句柄 | r3所指字符串的内容变为读取的字符串
//15 | 从数据区读取浮点数 | 0 |  | r3的值变为读取的浮点数
func hdlDataFunc(i *Inst) {
	if i.Opcode == OP_IN && i.B.Get() == 22 {
		i.VM.Attr()["data-ptr"] = i.VM.r2.Get()
		log.Debug("Change data pointer to %d", i.VM.r2.Get())
		i.A.Set(i.VM.r3.Get()) // 返回 r3 的值
		return
	}
	{
		v, p, _ := i.VM, i.A.Get(), i.B // port and param
		pos, r3 := v.attr["data-ptr"].(int), &v.r3
		switch p {
		case 13:
			r := v.GetInt(pos)
			log.Debug("Read int at %d get %d", pos, r)
			r3.Set(r)
			pos += 4
		case 14:
			s := v.MustGetStr(pos)
			log.Debug("Read str at %d get '%s'", pos, s)
			pos += len(s) + 1 // +1 for \0
			r3.SetStr(s)
		case 15:
			r3.Set(v.GetInt(pos))
			log.Debug("Read float at %d get %.6f", pos, r3.Float32())
			pos += 4
		}
		v.attr["data-ptr"] = pos
	}
}
