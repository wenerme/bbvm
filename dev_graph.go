package bbvm
import (
	"image/color"
	"image"
	"os"
	"./bi"
	"fmt"
	"image/draw"
)


type page struct {
	Graphic
	ps PenStyle
	wid int
	bs BrushStyle
}
func NewPage(w, h int) Page {
	p := &page{}
	p.Graphic = NewImageGraphic(image.NewRGBA(image.Rect(0, 0, w, h)))
	return p
}
func (p *page)SetPen(ps PenStyle, w int, c color.Color) {
	p.ps, p.wid = ps, w
	p.SetColor(c)
}
func (p *page)SetBrushStyle(bs BrushStyle) {
	p.bs = bs
}

type picture struct {
	image.Image
	name string
}
func (p *picture)Name() string {
	return p.name
}

/*
LoadPicture(fn string, i int) (Picture, error)
SetLcd(int, int)
PagePool() ResPool
PicPool() ResPool
FlipPage(Page)
SetFont(FontType)
SetFontStyle(color.Color, color.Color)
SetBackgroundMode(BackgroundMode)
Print(string, ... interface{})
Locate(int, int)
LocatePixel(int, int)
 */
type graphDev struct {
	pagePool ResPool
	picPool ResPool
	screen Page
	w, h int
}

func NewGraphDev(w, h int) GraphDev {
	g := &graphDev{}
	g.w, g.h = w, h
	g.screen = g.createPage()
	g.pagePool = &resPool{
		pool: make(map[int]Res),
		start:0, step:1, reuse:true, limit: 999,
		creator:func(ResPool) interface{} {return g.createPage()},
	}
	g.picPool = &resPool{
		pool: make(map[int]Res),
		start:0, step:1, reuse:true, limit: 999,
		//		creator:func(ResPool) interface{} {return g.createPage()},
	}
	return g
}
// 加载图片资源
// 资源索引从 1 开始
func (*graphDev)LoadPicture(fn string, n int) (p Picture, err error) {
	fp, err := os.Open(fn)
	if err != nil {return }
	i, err := bi.DecodeAt(fp, n - 1)// 该索引需要从 0 开始
	if err != nil {return }
	saveImage(i, fmt.Sprintf("load-picture-%d.bmp", n))

	p = &picture{i, fmt.Sprintf("%s#%d", fn, n)}
	return
}
func (*graphDev)SetLcd(w int, h int) {
	log.Info("SetLcd(%d,%d)", w, h)
}

func (g *graphDev)PagePool() ResPool {
	return g.pagePool
}
func (g *graphDev)PicPool() ResPool {
	return g.picPool
}
func (g *graphDev)Screen() Page {
	return g.screen
}

func (g *graphDev)FlipPage(Page) {
}

func (g *graphDev)SetFont(FontType) {

}
func (g *graphDev)SetFontStyle(color.Color, color.Color) {

}
func (g *graphDev)SetBackgroundMode(BackgroundMode) {

}
func (g *graphDev)Print(string, ... interface{}) {

}
func (g *graphDev)Locate(int, int) {

}
func (g *graphDev)LocatePixel(int, int) {

}
func (g *graphDev)createPage() Page {
	p := NewPage(g.w, g.h)
	p.SetColor(color.Black)
	p.FillRect(image.Rect(0, 0, g.w, g.h))
	return p
}


type Printer struct {
	draw.Image
	font Font
	loc image.Point
	FontColor color.Color
	BackgroundColor color.Color
}
func (p *Printer)SetFont(f Font) {
	p.font = f
}
func (p *Printer)Font() Font {
	return p.font
}
func (p *Printer)Print(str string) (err error) {
	for _, c := range str {
		err = p.PrintRune(c)
		if err != nil {return }
	}
	return
}
func (p *Printer)Write(b []byte) (int, error) {
	for i, c := range string(b) {
		err := p.PrintRune(c)
		if err != nil {return i, err }
	}
	return len(b), nil
}
func (p *Printer)PrintRune(c rune) (err error) {
	switch c{
	case '\t':
		p.loc.X += p.font.Width()*4
	case ' ':
		p.loc.X += p.font.Width()
	case '\n':
		p.loc.X = 0
		p.loc.Y += p.font.Height()
	default:
		i, err := p.Font().Render(string(c), p.FontColor, p.BackgroundColor)
		if err != nil {return err}
		area := image.Rectangle{p.loc, image.Pt(p.loc.X + i.Bounds().Dx(), p.loc.Y + i.Bounds().Dy())}
		draw.Draw(p, area, i, image.ZP, draw.Over)
		p.loc.X += i.Bounds().Dx()// FIXME + Dx 还是+ Width(可能会遇到多字节字符)
	}

	if p.loc.Y + p.font.Height() > p.Bounds().Dy() {
		p.NewPage()
		p.loc.Y = 0
	}
	if p.loc.X + p.font.Width() > p.Bounds().Dx() {
		p.loc.Y += p.font.Height()
	}

	return
}
func (p *Printer)Locate(row, col int) {
	p.LocatePixel(row * p.font.Height(), col * p.font.Width())
}
func (p *Printer)LocatePixel(x, y int) {
	p.loc = image.Pt(x, y)
}
func (p *Printer)render() {

}
func (p *Printer)NewPage() {

}