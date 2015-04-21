package bbvm

func handleInStr(i *Inst) {
	v, p, o := i.VM, i.B.Get(), i.A // port and param

	switch p{
		case 2:
		// 2 | 申请字符串句柄 | 申请到的句柄 |  |  IN():SHDL<br>从-1开始查询
		if r, err := v.StrPool().Acquire(); err == nil {
			rlog.Info("Acquire StrRes %d", r.Id())
			o.Set(r.Id())
		}else {
			rlog.Error("Acquire StrRes faield: %s", err)
		}
		case 5:
		// 5 | 复制字符串 | r3的值 | r2:源字符串句柄<br>r3:目标字符串句柄 | IN(r2:SHDL,r3:SHDL):SHDL<br>r3所代表字符串的内容被修改
		rlog.Info("StrRes copy '%s' to %d", o.Str(), v.Register(REG_R3).Get())
		v.StrPool().Get(v.Register(REG_R3).Get()).Set(o.Str())
		// TODO 确认是否返回r3的值
		case 8:
		// 8 | 释放字符串句柄 | r3的值 | r3:字符串句柄 | IN(r3:SHDL):SHDL
		// TODO 确认是否返回r3的值
		hdl := v.Register(REG_R3).Get()
		r := v.StrPool().Get(hdl)
		if r != nil {
			rlog.Info("Release StrRes %d '%v'", r.Id(), r.Get())
			v.StrPool().Release(r)
		}else {
			rlog.Info("Release StrRes %d failed: not exists", hdl)
		}
	}
}

func HandInStr(v VM) {
	v.SetIn(HANDLE_ALL, 2, handleInStr)
	v.SetIn(HANDLE_ALL, 5, handleInStr)
	v.SetIn(HANDLE_ALL, 8, handleInStr)
}