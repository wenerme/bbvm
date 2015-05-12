package bi
import (
	"image"
	"image/color"
	"./bc"
)



// RGB565 is an in-memory image whose At method returns bc.RGB565 values.
type RGB565 struct {
	// Pix holds the image's pixels, as gray values in big-endian format. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*2].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
}

func (*RGB565)ColorModel() color.Model {
	return bc.RGB565Model
}

func (i *RGB565)Bounds() image.Rectangle {
	return i.Rect
}

func (i *RGB565)At(x, y int) color.Color {
	return i.RGB565At(x, y)
}

func (p *RGB565) RGB565At(x, y int) bc.RGB565 {
	if !(image.Point{x, y}.In(p.Rect)) {
		return bc.RGB565{0}
	}
	i := p.PixOffset(x, y)
	return bc.RGB565{uint16(p.Pix[i+1]) | uint16(p.Pix[i+0])<<8}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *RGB565) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*2
}

func NewRGB565(r image.Rectangle) *RGB565 {
	w, h := r.Dx(), r.Dy()
	buf := make([]uint8, 2*w*h)
	return &RGB565{buf, 2*w, r}
}






// Gray2 is an in-memory image whose At method returns color.Gray2 values.
type Gray2 struct {
	// Pix holds the image's pixels, as gray values in big-endian format. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)/4].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
}

func (*Gray2)ColorModel() color.Model { return bc.Gray2Model }

func (i *Gray2)Bounds() image.Rectangle {
	return i.Rect
}

func (i *Gray2)At(x, y int) color.Color { return i.Gray2At(x, y) }

func (p *Gray2) Gray2At(x, y int) bc.Gray2 {
	if !(image.Point{x, y}.In(p.Rect)) {
		return bc.Gray2{0}
	}
	i := p.PixOffset(x, y)

	l := p.Pix[i]
	l = l >> uint((x % 4)*2)
	return bc.Gray2{l & 0x3}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *Gray2) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)/4
}

func NewGray2(r image.Rectangle) *Gray2 {
	w, h := r.Dx(), r.Dy()
	w4 := w/4
	if w % 4 != 0 {w4+=1}
	buf := make([]uint8, w4*h)
	return &Gray2{buf, w4, r}
}