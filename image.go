package bbvm
import (
	"image/draw"
	"image/color"
	"image"
)
// Use pen to draw a image
type Pen struct {
	draw.Image
	Color color.Color
	Pos image.Point
}
func (p *Pen)Rect(r image.Rectangle) {
	rect := r.Intersect(p.Bounds())
	//	old := p.Pos
	min, max := rect.Min, rect.Max
	p.Pos = min
	p.LineTo(image.Point{min.X, max.Y})
	p.LineTo(max)
	p.LineTo(image.Point{max.X, min.Y})
	p.LineTo(min)
}
func (p *Pen)Line(a image.Point, b image.Point) {
	p.DrawLine(a.X, a.Y, b.X, b.Y)
}
func (p *Pen)DrawLine(x0, y0, x1, y1 int) {
	// http://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm
	dx := x1-x0
	dy := y1-y0

	if dy < 0 {dy = -dy}
	if dx < 0 {dx = -dx}

	sx := -1
	if x0 < x1 {sx = 1}
	sy := -1
	if y0 < y1 { sy = 1}

	err := dx-dy
	var e2 int
	x := x0
	y := y0

	for {
		p.Set(x, y, p.Color)

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
func (p *Pen)LineTo(b image.Point) {
	p.Line(p.Pos, b)
	p.Pos = b
}
func (p *Pen)Circle(x0 int, y0 int, r int) {
	// http://en.wikipedia.org/wiki/Midpoint_circle_algorithm
	x, y := r, 0
	re := 1-x
	for x >= y {
		p.Set(+x + x0, +y + y0, p.Color)
		p.Set(+y + x0, +x + y0, p.Color)
		p.Set(-x + x0, +y + y0, p.Color)
		p.Set(-y + x0, +x + y0, p.Color)
		p.Set(-x + x0, -y + y0, p.Color)
		p.Set(-y + x0, -x + y0, p.Color)
		p.Set(+x + x0, -y + y0, p.Color)
		p.Set(+y + x0, -x + y0, p.Color)
		y +=1
		if (re<0) {
			re += 2 * y + 1;
		} else {
			x--;
			re += 2 * (y - x) + 1;
		}
	}
}
