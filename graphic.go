package bbvm
import (
	"image/draw"
	"image/color"
	"image"
)
type Graphic interface {
	draw.Image
	Rect(r image.Rectangle)
	FillRect(r image.Rectangle)
	Line(a image.Point, b image.Point)
	LineTo(b image.Point)
	Circle(x0 int, y0 int, r int)
	//	FillCircle(x0 int, y0 int, r int)
	Color() color.Color
	SetColor(color.Color)
}
func NewImageGraphic(i draw.Image) Graphic {
	return &graphic{i, nil, image.Point{0, 0}}
}
// A graphic impl
type graphic struct {
	draw.Image
	c color.Color
	p image.Point
}
func (p *graphic)Color() color.Color {
	return p.c
}
func (p *graphic)SetColor(c color.Color) {
	p.c = c
}
func (p *graphic)Rect(r image.Rectangle) {
	rect := r.Intersect(p.Bounds())
	//	old := p.Pos
	min, max := rect.Min, rect.Max
	p.p = min
	p.LineTo(image.Point{min.X, max.Y})
	p.LineTo(max)
	p.LineTo(image.Point{max.X, min.Y})
	p.LineTo(min)
}
func (p *graphic)FillRect(r image.Rectangle) {
	rect := r.Intersect(p.Bounds())
	//	old := p.Pos
	//	dx := rect.Dx()
	min, max := rect.Min, rect.Max
	for x := min.X; x <=max.X; x+=1 {
		p.DrawLine(x, min.Y, x, max.Y)
	}
}
func (p *graphic)Line(a image.Point, b image.Point) {
	p.DrawLine(a.X, a.Y, b.X, b.Y)
}
func (p *graphic)DrawLine(x0, y0, x1, y1 int) {
	// http://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm
	dx := x1-x0
	dy := y1-y0

	if dy < 0 {dy = -dy}
	if dx < 0 {dx = -dx}

	sx := -1
	if x0 < x1 {sx = 1}
	sy := -1
	if y0 < y1 {sy = 1}

	err := dx-dy
	var e2 int
	x := x0
	y := y0

	for {
		p.Set(x, y, p.c)

		if (x == x1 && y == y1) {
			break;
		}

		e2 = 2*err;
		if (e2 > -1 * dy) {
			err = err - dy
			x = x + sx
		}

		if (e2 < dx) {
			err = err + dx
			y = y + sy
		}
	}
}
func (p *graphic)LineTo(b image.Point) {
	p.Line(p.p, b)
	p.p = b
}
func (p *graphic)Circle(x0 int, y0 int, r int) {
	// http://en.wikipedia.org/wiki/Midpoint_circle_algorithm
	x, y := r, 0
	re := 1-x
	for x >= y {
		p.Set(+x + x0, +y + y0, p.c)
		p.Set(+y + x0, +x + y0, p.c)
		p.Set(-x + x0, +y + y0, p.c)
		p.Set(-y + x0, +x + y0, p.c)
		p.Set(-x + x0, -y + y0, p.c)
		p.Set(-y + x0, -x + y0, p.c)
		p.Set(+x + x0, -y + y0, p.c)
		p.Set(+y + x0, -x + y0, p.c)
		y +=1
		if (re<0) {
			re += 2 * y + 1;
		} else {
			x--;
			re += 2 * (y - x) + 1;
		}
	}
}
