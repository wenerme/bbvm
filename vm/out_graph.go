package vm
import (
	"image"
	"golang.org/x/image/draw"
	"./bi/bc"
)

func (out)Graphic(v VM) {
	v.Attr()["graph-dev"] = NewGraphDev(240, 320)
	v.SetOut(16, HANDLE_ALL, outGraphicFunc)
	v.SetOut(17, HANDLE_ALL, outGraphicFunc)
	v.SetOut(18, HANDLE_ALL, outGraphicFunc)
	v.SetOut(19, HANDLE_ALL, outGraphicFunc)
	v.SetOut(20, HANDLE_ALL, outGraphicFunc)
	v.SetOut(21, HANDLE_ALL, outGraphicFunc)
	v.SetOut(22, HANDLE_ALL, outGraphicFunc)
	v.SetOut(23, HANDLE_ALL, outGraphicFunc)
	v.SetOut(24, HANDLE_ALL, outGraphicFunc)
	v.SetOut(25, HANDLE_ALL, outGraphicFunc)
	v.SetOut(26, HANDLE_ALL, outGraphicFunc)
	v.SetOut(40, HANDLE_ALL, outGraphicFunc)
	v.SetOut(41, HANDLE_ALL, outGraphicFunc)
	v.SetOut(64, HANDLE_ALL, outGraphicFunc)
	v.SetOut(65, HANDLE_ALL, outGraphicFunc)
	v.SetOut(66, HANDLE_ALL, outGraphicFunc)
	v.SetOut(67, HANDLE_ALL, outGraphicFunc)
	v.SetOut(68, HANDLE_ALL, outGraphicFunc)
	v.SetOut(69, HANDLE_ALL, outGraphicFunc)
}
/*
64 | 设置画笔 | 0 | r3:参数地址 |  SETPEN(PAGE,STYLE,WID,COLOR)
65 | 设置刷子 | 0 | r2:PAGE r3:STYLE |  SETBRUSH(PAGE,STYLE)
66 | 移动画笔 | 0 | r1,r2,r3:PAGE,X,Y |  MOVETO(PAGE,X,Y)
67 | 画线 | 0 | r1,r2,r3:PAGE,X,Y |  LINETO(PAGE,X,Y)
68 | 画矩形 | 0 | r3:参数地址 |  RECTANGLE(PAGE,LEFT,TOP,RIGHT,BOTTOM)
69 | 画圆 | 0 | r3:参数地址 |  CIRCLE(PAGE,CX,CY,CR)
 */
func outDrawFunc(i *Inst) {

}
/*
19 | 申请图片句柄并从文件载入像素资源 | r3:资源句柄 | r3:文件名, r2:资源索引 |  LOADRES(FILE$,ID)
20 | 复制图片到画布上 | 0 | r3:地址,其他参数在该地址后 |  SHOWPIC(PAGE,PIC,DX,DY,W,H,X,Y,MODE)
26 | 释放图片句柄 | 0 | r3:资源句柄 |  FREERES(ID)
40 | 获取图片宽度 | r3 | r3 |  GETPICWID(PIC)
41 | 获取图片高度 | r3 | r3 |  GETPICHGT(PIC)
 */
func outPictureFunc(i *Inst) {

}
/*
17 | 申请画布句柄 | 0 ,r3:PAGE句柄 | - | CREATEPAGE()
18 | 释放画布句柄 | 0 | r3:PAGE句柄 |  DELETEPAGE(PAGE)
21 | 显示画布 | 0 | r3:PAGE句柄 |  FLIPPAGE(PAGE)
22 | 复制画布 | 0 | r2:目标PAGE句柄,r3:源PAGE句柄 |  BITBLTPAGE(DEST,SRC)
23 | 填充画布 | 0 | r3:参数地址 |  FILLPAGE(PAGE,X,Y,WID,HGT,COLOR)
24 | 写入画布某点颜色 | 0 | r3:参数地址 |  PIXEL(PAGE,X,Y,COLOR)
25 | 读取画布某点颜色 | 0 | r3:参数地址 |  READPIXEL(PAGE,X,Y)

 */
func outPageFunc(i *Inst) {

}
/*
37 | 设定文字颜色 | 0 | r3:参数地址 |  COLOR(FRONT,BACK,FRAME)
 */
func outPrintFunc(i *Inst) {

}
/*
16 | 设定模拟器屏幕 | 0 | r2:宽, r3:高 |  SETLCD(WIDTH,HEIGHT)
35 | 清屏 | 0 |  |
 */
func outGraphicFunc(i *Inst) {
	v, p, _ := i.VM, i.A.Get(), i.B // port and param
	r1, r2, r3 := &v.r1, &v.r2, &v.r3
	gd := v.Attr()["graph-dev"].(GraphDev)
	pagePool := gd.PagePool()
	picPool := gd.PicPool()
	_, _, _ = gd, picPool, pagePool
	switch p{
	case 16:
		log.Debug("SETLCD(%d,%d)", r2.Get(), r3.Get())
		gd.SetLcd(r2.Get(), r3.Get())
	case 17:
		r, err := pagePool.Acquire()
		if err != nil {log.Warning("CREATEPAGE() faield: %s", err.Error()); break}
		log.Debug("CREATEPAGE() -> %d", r.Id())
		r3.Set(r.Id())
	case 18:
		r := pagePool.Get(r3.Get())
		if r == nil {log.Warning("DELETEPAGE(%s) faield: Page not exists", r3.Get()); break}
		pagePool.Release(r)
		log.Debug("DELETEPAGE(%s)", r3.Get())
	case 19:
		fn, idx := r3.Str(), r2.Get()

		if res, err := picPool.Acquire(); err == nil {
			pic, err := gd.LoadPicture(fn, idx)
			if err != nil {
				log.Error("LOADRES(%s,%d) faield: %s", fn, idx, err.Error())
			}
			res.Set(pic)
			r3.Set(res.Id())
			log.Debug("LOADRES(%s,%d) -> %d", fn, idx, r3.Get())
		}else {
			log.Error("LOADRES(%s,%d) faield: %s", fn, idx, err.Error())
		}
	case 20:
		args := newArgs(r3.Get(), v, 9)
		pi, picId, dx, dy, w, h, x, y, mode := args.Next9Int()
		log.Debug("SHOWPIC(%d,%d,%d,%d,%d,%d,%d,%d,%d)", pi, picId, dx, dy, w, h, x, y, mode)
		if pg := getPage(gd, pi); pg != nil {
			if pic := getPic(gd, picId); pic != nil {
				dest := image.Rect(dx, dy, w, h)
				// TODO 透明模式,以 0xFF00FF, 0xF800F8(RGB565) 作为透明色
				draw.Draw(pg, dest, pic, image.Pt(x, y), draw.Over)
			}
		}
	case 21:
		log.Debug("FLIPPAGE(%d)", r3.Get())
		if pg := getPage(gd, r3.Get()); pg != nil {
			draw.Copy(gd.Screen(), image.Pt(0, 0), pg, pg.Bounds(), nil)
		}
	case 22:
		log.Debug("BITBLTPAGE(%d,%d)", r2.Get(), r3.Get())
		if pgDest := getPage(gd, r2.Get()); pgDest != nil {
			if pgSrc := getPage(gd, r3.Get()); pgSrc != nil {
				draw.Copy(pgDest, image.Pt(0, 0), pgSrc, pgSrc.Bounds(), nil)
			}
		}
	case 23:
		args := newArgs(r3.Get(), v, 6)
		pi, x, y, wid, hgt, c := args.Next6Int()
		log.Debug("FILLPAGE(%d,%d,%d,%d,%d,%d)", pi, x, y, wid, hgt, c)
		if pg := getPage(gd, pi); pg != nil {
			c := pg.Color()
			pg.FillRect(image.Rect(x, y, wid, hgt))
			pg.SetColor(c)
		}
	case 24:
		args := newArgs(r3.Get(), v, 4)
		pi, x, y, c := args.Next4Int()
		log.Debug("PIXEL(%d,%d,%d,%d)", pi, x, y, c)
		if pg := getPage(gd, pi); pg != nil {
			// TODO 边界检查
			pg.Set(x, y, bc.BGR888{uint32(c)})
		}
	case 25:
		args := newArgs(r3.Get(), v, 4)
		pageId, x, y := args.Next3Int()
		log.Debug("READPIXEL(%d,%d,%d)", pageId, x, y)
		if pg := getPage(gd, pageId); pg != nil {
			// TODO 边界检查
			r3.Set(int(bc.BGR888Model.Convert(pg.At(x, y)).(bc.BGR888).V))
		}
	case 26:
		idx := r3.Get()
		if res := picPool.Get(idx); res != nil {
			log.Error("FREERES(%d)", idx)
			picPool.Release(res)
		}else {
			log.Error("FREERES(%d) faield: not exists", idx)
		}
	/*
36 | 按行列定位光标 | 0 | r2:行,r3:列 |  LOCATE(LINE,ROW)

38 | 设定文字字体大小 | 0 | r3:FONT |  FONT(F)
40 | 获取图片宽度 | r3 | r3 |  GETPICWID(PIC)
41 | 获取图片高度 | r3 | r3 |  GETPICHGT(PIC)
42 | 按坐标定位光标 | - | r2:行,r3:列 |  PIXLOCATE(LINE,ROW)
43 | 复制部分画布 | - | r3:参数地址 |  STRETCHBLTPAGE(X,Y,DEST,SRC)
44 | 设定背景模式 | r3:MODE | - |  SETBKMODE(mode)
	 */
	case 40:
		pi := r3.Get()
		if pic := getPic(gd, pi); pic != nil {
			r3.Set(pic.Bounds().Dx())
		}else {
			r3.Set(0)
		}
		log.Debug("GETPICWID(%d) -> %d", pi, r3.Get())
	case 41:
		pi := r3.Get()
		if pic := getPic(gd, pi); pic != nil {
			r3.Set(pic.Bounds().Dy())
		}else {
			r3.Set(0)
		}
		log.Debug("GETPICHGT(%d) -> %d", pi, r3.Get())
	/*
64 | 设置画笔 | 0 | r3:参数地址 |  SETPEN(PAGE,STYLE,WID,COLOR)
65 | 设置刷子 | 0 | r2:PAGE r3:STYLE |  SETBRUSH(PAGE,STYLE)
66 | 移动画笔 | 0 | r1,r2,r3:PAGE,X,Y |  MOVETO(PAGE,X,Y)
67 | 画线 | 0 | r1,r2,r3:PAGE,X,Y |  LINETO(PAGE,X,Y)
68 | 画矩形 | 0 | r3:参数地址 |  RECTANGLE(PAGE,LEFT,TOP,RIGHT,BOTTOM)
69 | 画圆 | 0 | r3:参数地址 |  CIRCLE(PAGE,CX,CY,CR)
*/
	case 64:
		args := newArgs(r3.Get(), v, 4)
		pi, style, w, c := args.Next4Int()
		log.Debug("SETPEN(%d,%d,%d,%d)", pi, style, w, c)
		if pg := getPage(gd, pi); pg != nil {
			pg.SetPen(PenStyle(style), w, bc.BGR888{uint32(c)})
		}
	case 65:
		pi, style := r2.Get(), BrushStyle(r3.Get())
		log.Debug("SETBRUSH(%d,%d)", pi, style)
		if pg := getPage(gd, pi); pg != nil {
			pg.SetBrushStyle(style)
		}
	case 66:
		pi, x, y := r1.Get(), r2.Get(), r3.Get()
		log.Debug("MOVETO(%d,%d,%d)", pi, x, y)
		if pg := getPage(gd, pi); pg != nil {
			pg.MoveTo(image.Pt(x, y))
		}
	case 67:
		pi, x, y := r1.Get(), r2.Get(), r3.Get()
		log.Debug("LINETO(%d,%d,%d)", pi, x, y)
		if pg := getPage(gd, pi); pg != nil {
			pg.LineTo(image.Pt(x, y))
		}
	case 68:
		// RECTANGLE(PAGE,LEFT,TOP,RIGHT,BOTTOM)
		args := newArgs(r3.Get(), v, 5)
		pi, left, top, right, bottom := args.Next5Int()
		log.Debug("RECTANGLE(%d,%d,%d,%d,%d)", pi, left, top, right, bottom)
		if pg := getPage(gd, pi); pg != nil {
			pg.Rect(image.Rect(left, top, right, bottom))
		}
	case 69:
		args := newArgs(r3.Get(), v, 4)
		pi, cx, cy, cr := args.Next4Int()
		log.Debug("CIRCLE(%d,%d,%d,%d)", pi, cx, cy, cr)
		if pg := getPage(gd, pi); pg != nil {
			pg.Circle(cx, cy, cr)
		}
	}
}

func getPage(g GraphDev, id int) Page {
	if id == -1 {
		return g.Screen()
	}
	if id < 0 {
		log.Error("Page handle is invalid: %d", id)
		return nil
	}
	p := g.PagePool()
	r := p.Get(id)
	if r == nil || r.Get() == nil {
		log.Error("GetPage #%d not exists", id)
		return nil
	}
	return r.Get().(Page)
}
func getPic(g GraphDev, id int) Picture {
	p := g.PicPool()
	r := p.Get(id)
	if r == nil || r.Get() == nil {
		log.Error("GetPic #%d not exists", id)
		return nil
	}
	return r.Get().(Picture)
}

func newArgs(addr int, vm VM, n int) args {
	a := args(make([]int, n))
	for i := 0; i < n; i ++ {
		a[i] = vm.GetInt(addr)
		addr += 4
	}
	return a
}
type args []int
func (a args)Next2Int() (int, int) {
	return a[1], a[0]
}
func (a args)Next3Int() (int, int, int) {
	return a[2], a[1], a[0]
}
func (a args)Next4Int() (int, int, int, int) {
	return a[3], a[2], a[1], a[0]
}
func (a args)Next5Int() (int, int, int, int, int) {
	return a[4], a[3], a[2], a[1], a[0]
}
func (a args)Next6Int() (int, int, int, int, int, int) {
	return a[5], a[4], a[3], a[2], a[1], a[0]
}

func (a args)Next9Int() (int, int, int, int, int, int, int, int, int) {
	return a[8], a[7], a[6], a[5], a[4], a[3], a[2], a[1], a[0]
}
