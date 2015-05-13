package bbvm
import (
	"image/color"
	"image"
	"os"
	"./bi"
	"fmt"
)


type page struct {
	Graphic
	ps PenStyle
	wid int
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
	return g
}

func (*graphDev)LoadPicture(fn string, n int) (p Picture, err error) {
	fp, err := os.Open(fn)
	if err != nil {return }
	i, err := bi.DecodeAt(fp, n)
	if err != nil {return }
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