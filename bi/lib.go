package bi
import (
	"image"
	"bytes"
	"os"
	"encoding/binary"
	"image/color"
	"./bc"
)

type LibConfig struct {
	image.Config
	Format string
	Offset int
	Name string
}


func DecodeLibConfig(r *bytes.Reader) ([]LibConfig, error) {
	buf := make([]byte, 4)
	_, err := r.Read(buf)
	if err != nil { return nil, err }
	n := int(binary.LittleEndian.Uint32(buf))
	configs := make([]LibConfig, n)
	for i := 0; i < n; i ++ {
		buf = make([]byte, 4)
		_, err = r.Read(buf)
		if err != nil { return configs, err }
		c := LibConfig{}
		c.Offset = int(binary.LittleEndian.Uint32(buf))

		r.Seek(int64(c.Offset), os.SEEK_SET)

		cfg, err := DecodeLibRGB565OneConfig(r)
		f := "lib"
		if err != nil { return configs, err }
		c.Format = f
		c.Config = cfg
		configs[i] = c
	}

	return configs, nil
}

func DecodeLibRGB565OneConfig(r *bytes.Reader) (image.Config, error) {
	buf := make([]byte, 4)
	cfg := image.Config{}
	cfg.ColorModel = bc.BGR565Model
	_, err := r.Read(buf)
	if err != nil { return cfg, err }
	_ = int(binary.LittleEndian.Uint32(buf))
	_, err = r.Read(buf[:2])
	if err != nil { return cfg, err }
	w := int(binary.LittleEndian.Uint16(buf))
	_, err = r.Read(buf[:2])
	if err != nil { return cfg, err }
	h := int(binary.LittleEndian.Uint16(buf))

	cfg.Height = h
	cfg.Width = w
	return cfg, nil
}
type RGB565 struct {
	// Tow byte pre pixel
	Pix []uint8

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
		return bc.RGB565(0)
	}
	i := p.PixOffset(x, y)
	return bc.RGB565(binary.LittleEndian.Uint16(p.Pix[i:]))
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

// Gray16 is an in-memory image whose At method returns color.Gray2 values.
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
	l = l >> uint((3-(x % 4))*2)
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

func DecodeLibRGB565One(r *bytes.Reader) (image.Image, error) {
	buf := make([]byte, 4)
	_, err := r.Read(buf)
	if err != nil { return nil, err }
	_ = int(binary.LittleEndian.Uint32(buf))

	_, err = r.Read(buf[:2])
	if err != nil { return nil, err }
	w := int(binary.LittleEndian.Uint16(buf))
	_, err = r.Read(buf[:2])
	if err != nil { return nil, err }
	h := int(binary.LittleEndian.Uint16(buf))

	// Skip 8 byte
	_, err=r.Read(buf)
	if err != nil { return nil, err }
	_, err=r.Read(buf)
	if err != nil { return nil, err }

	i := NewRGB565(image.Rect(0, 0, w, h))
	_, err=r.Read(i.Pix)
	if err != nil { return i, err }

	return i, nil
}
func DecodeLibGray2One(r *bytes.Reader) (image.Image, error) {
	buf := make([]byte, 4)
	_, err := r.Read(buf)
	if err != nil { return nil, err }
	_ = int(binary.LittleEndian.Uint32(buf))

	_, err = r.Read(buf[:2])
	if err != nil { return nil, err }
	w := int(binary.LittleEndian.Uint16(buf))
	_, err = r.Read(buf[:2])
	if err != nil { return nil, err }
	h := int(binary.LittleEndian.Uint16(buf))

	_, err=r.Read(buf)
	if err != nil { return nil, err }
	_, err=r.Read(buf)
	if err != nil { return nil, err }

	i := NewGray2(image.Rect(0, 0, w, h))
	_, err=r.Read(i.Pix)
	if err != nil { return i, err }

	return i, nil
}