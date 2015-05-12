package bi
import "image/color"

var (
	BGR888Model color.Model = color.ModelFunc(bgr888Model)
	BGR565Model color.Model = color.ModelFunc(bgr565Model)
	RGB565Model color.Model = color.ModelFunc(rgb565Model)
)
type BGR888Color int
func (i BGR888Color)Int() int {
	return int(i)
}
func (i BGR888Color)RGBA() (r, g, b, a uint32) {
	r = uint32(i&0xff)
	r |= r << 8
	g = uint32(i>>8&0xff)
	g |= g << 8
	b = uint32(i>>16&0xff)
	b |= b << 8
	a = 0xff
	a |= a << 8
	return
}
func bgr888Model(c color.Color) color.Color {
	if _, ok := c.(BGR888Color); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	r >>= 8
	g >>= 8
	b >>= 8
	return BGR888Color(b<<16|g<<8|r)
}

type BGR565Color uint16
func (i BGR565Color)RGBA() (r, g, b, a uint32) {
	r = uint32(i&0x1f)
	r |= r << 11
	g = uint32(i>>5&0x3f)
	g |= g << 10
	b = uint32(i>>11&0x1f)
	b |= b << 11
	a = 0xff
	a |= a << 8

	return
}
func bgr565Model(c color.Color) color.Color {
	if _, ok := c.(BGR565Color); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	r >>= 8
	g >>= 8
	b >>= 8

	r >>= 3 & 0x1f
	g >>= 2 & 0x3f
	b >>= 3 & 0x1f

	return BGR565Color(b<<11|g<<5|r)
}

type RGB565Color uint16
func (i RGB565Color)RGBA() (r, g, b, a uint32) {
	b = uint32(i&0x1f)
	b |= b << 11
	g = uint32(i>>5&0x3f)
	g |= g << 10
	r = uint32(i>>11&0x1f)
	r |= r << 11

	a = 0xff
	a |= a << 8

	return
}
func rgb565Model(c color.Color) color.Color {
	if _, ok := c.(BGR565Color); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	r >>= 11 & 0x1f
	g >>= 10 & 0x3f
	b >>= 11 & 0x1f

	return RGB565Color(r<<11|g<<5|b)
}