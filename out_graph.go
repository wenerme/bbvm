package bbvm

func (out)Graphic(v VM) {

}

func outGraphicFunc(i *Inst) {
	v, p, _ := i.VM, i.A.Get(), i.B // port and param
	r2, r3 := &v.r2, &v.r3
	pagePool := v.Attr()["page-pool"].(ResPool)
	/*
16 | 设定模拟器屏幕 | 0 | r2:宽, r3:高 |  SETLCD(WIDTH,HEIGHT)
17 | 申请画布句柄 | 0 ,r3:PAGE句柄 | - | CREATEPAGE()
18 | 释放画布句柄 | 0 | r3:PAGE句柄 |  DELETEPAGE(PAGE)
19 | 申请图片句柄并从文件载入像素资源 | r3:资源句柄 | r3:文件名, r2:资源索引 |  LOADRES(FILE$,ID)
20 | 复制图片到画布上 | 0 | r3:地址,其他参数在该地址后 |  SHOWPIC(PAGE,PIC,DX,DY,W,H,X,Y,MODE)
21 | 显示画布 | 0 | r3:PAGE句柄 |  FLIPPAGE(PAGE)
22 | 复制画布 | 0 | r2:目标PAGE句柄,r3:源PAGE句柄 |  BITBLTPAGE(DEST,SRC)
23 | 填充画布 | 0 | r3:参数地址 |  FILLPAGE(PAGE,X,Y,WID,HGT,COLOR)
24 | 写入画布某点颜色 | 0 | r3:参数地址 |  PIXEL(PAGE,X,Y,COLOR)
25 | 读取画布某点颜色 | 0 | r3:参数地址 |  READPIXEL(PAGE,X,Y)
	 */
	switch p{
	case 16:
		log.Info("SETLCD(%d, %d)", r2.Get(), r3.Get())
	case 17:
		r, err := pagePool.Acquire()
		if err != nil {log.Warning("CREATEPAGE() faield: %s", err.Error()); break}
		log.Info("CREATEPAGE() -> %d", r.Id())
		r3.Set(r.Id())
	case 18:
		r := pagePool.Get(r3.Get())
		if r == nil {log.Warning("DELETEPAGE(%s) faield: Page not exists", r3.Get()); break}
		pagePool.Release(r)
		log.Info("DELETEPAGE(%s)", r3.Get())
	case 19:
		fn, idx := r3.Str(), r2.Get()
		r3.Set(0)
		log.Info("LOADRES(%s,%d) -> %d", fn, idx, r3.Get())
	case 20:
		log.Info("SHOWPIC(PAGE,PIC,DX,DY,W,H,X,Y,MODE)")
	case 21:
		log.Info("FLIPPAGE(%d)", r3.Get())
	case 22:
		log.Info("BITBLTPAGE(%d,%d)", r2.Get(), r3.Get())
	case 23:
	case 24:
	case 25:

		/*
26 | 释放图片句柄 | 0 | r3:资源句柄 |  FREERES(ID)

34 | 判定某键是否按下 | 0;r3 | r3:KEY |  KEYPRESS(KEY)
35 | 清屏 | 0 |  |
36 | 按行列定位光标 | 0 | r2:行,r3:列 |  LOCATE(LINE,ROW)
37 | 设定文字颜色 | 0 | r3:参数地址 |  COLOR(FRONT,BACK,FRAME)
38 | 设定文字字体大小 | 0 | r3:FONT |  FONT(F)
39 | 等待按键 | r3:按键 | - |  WAITKEY()
40 | 获取图片宽度 | r3 | r3 |  GETPICWID(PIC)
41 | 获取图片高度 | r3 | r3 |  GETPICHGT(PIC)
42 | 按坐标定位光标 | - | r2:行,r3:列 |  PIXLOCATE(LINE,ROW)
43 | 复制部分画布 | - | r3:参数地址 |  STRETCHBLTPAGE(X,Y,DEST,SRC)
44 | 设定背景模式 | r3:MODE | - |  SETBKMODE(mode)
45 | 获取按键的字符串 | 0 | r3:字符串句柄,用于存储结果 |  InKey$
46 | 获取按键的ASCII码 | 0 | r3:KEYPRESS |  INKEY()
		 */
	}
}