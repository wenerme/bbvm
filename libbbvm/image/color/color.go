package color
import "image/color"

var (
	BGR888Model color.Model = color.ModelFunc(bgr888Model)
	BGR565Model color.Model = color.ModelFunc(bgr565Model)
	RGB565Model color.Model = color.ModelFunc(rgb565Model)
	Gray2Model color.Model = color.ModelFunc(gray2Model)
)
type BGR888 struct {
	V uint32
}
func (i BGR888)RGBA() (r, g, b, a uint32) {
	v := i.V
	r = uint32(v&0xff)
	r |= r << 8
	g = uint32(v>>8&0xff)
	g |= g << 8
	b = uint32(v>>16&0xff)
	b |= b << 8
	a = 0xff
	a |= a << 8
	return
}
func bgr888Model(c color.Color) color.Color {
	if _, ok := c.(BGR888); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	r >>= 8
	g >>= 8
	b >>= 8
	return BGR888{b<<16|g<<8|r}
}

type BGR565 struct {
	V uint16
}
func (i BGR565)RGBA() (r, g, b, a uint32) {
	v := i.V
	r = uint32(v&0x1f)
	r |= r << 11
	g = uint32(v>>5&0x3f)
	g |= g << 10
	b = uint32(v>>11&0x1f)
	b |= b << 11
	a = 0xff
	a |= a << 8

	return
}
func bgr565Model(c color.Color) color.Color {
	if _, ok := c.(BGR565); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	r >>= 8
	g >>= 8
	b >>= 8

	r >>= 3 & 0x1f
	g >>= 2 & 0x3f
	b >>= 3 & 0x1f

	return BGR565{uint16(b<<11|g<<5|r)}
}

type RGB565 struct {
	V uint16
}
func (i RGB565)RGBA() (r, g, b, a uint32) {
	v := i.V
	b = uint32(v&0x1f)
	b |= b << 11
	g = uint32(v>>5&0x3f)
	g |= g << 10
	r = uint32(v>>11&0x1f)
	r |= r << 11

	a = 0xff
	a |= a << 8

	return
}
func rgb565Model(c color.Color) color.Color {
	if _, ok := c.(BGR565); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	r >>= 11 & 0x1f
	g >>= 10 & 0x3f
	b >>= 11 & 0x1f

	return RGB565{uint16(r<<11|g<<5|b)}
}

// Gray16 represents a 2-bit grayscale color.
type Gray2 struct {
	Y uint8
}

func (c Gray2) RGBA() (r, g, b, a uint32) {
	y := uint32(c.Y) << 14
	return y, y, y, 0xffff
}


func gray2Model(c color.Color) color.Color {
	if _, ok := c.(Gray2); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	y := (299*r + 587*g + 114*b + 500) / 1000
	return Gray2{uint8(y >> 14)}
}
