package vm

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"image"
	"image/color"
	"image/draw"
)

type sdlFont struct {
	*ttf.Font
}

func NewSDLFont(f *ttf.Font) Font {
	sf := &sdlFont{f}
	return sf
}

func (f *sdlFont) Width() int {
	return f.Height() / 2
}

func (f *sdlFont) Render(s string, fg, bg color.Color) (i image.Image, err error) {
	sf, err := f.RenderUTF8Solid(s, ColorToSdlColor(fg))
	i, err = SurfaceConvertToImage(sf)
	fmt.Println(i.Bounds())
	return
}
func (f *sdlFont) RenderRune(c rune, fg, bg color.Color) (i image.Image, err error) {
	sf, err := f.RenderUTF8Solid(string(c), ColorToSdlColor(fg))
	i, err = SurfaceConvertToImage(sf)
	return
}

func ColorToSdlColor(c color.Color) sdl.Color {
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	return sdl.Color{rgba.R, rgba.G, rgba.B, rgba.A}
}

func SurfaceConvertToImage(s *sdl.Surface) (img image.Image, err error) {
	switch s.Format.Format {
	case sdl.PIXELFORMAT_RGBA8888, sdl.PIXELFORMAT_RGBX8888:
		i := &image.RGBA{Rect: image.Rect(0, 0, int(s.W), int(s.H))}
		i.Pix = s.Pixels()
		img = i
	case sdl.PIXELFORMAT_INDEX8:
		i := image.NewRGBA(image.Rect(0, 0, int(s.W), int(s.H)))
		key, err := s.GetColorKey()
		if err != nil {
			return nil, err
		}
		//		fmt.Println(s.PixelNum(), len(s.Pixels()), len(i.Pix), s.W, s.H, s.Pitch, s.Format.BitsPerPixel, s.Format.BytesPerPixel, key)
		l := len(s.Pixels())
		var r, g, b, a uint8
		for n := 0; n < l; n += 1 {
			pixel := s.Pixels()[n]
			if pixel == uint8(key) {
				r, g, b, a = 0, 0, 0, 0
			} else {
				r, g, b, a = sdl.GetRGBA(uint32(pixel), s.Format)
			}

			p := i.PixOffset(n%int(s.Pitch), n/int(s.Pitch))
			i.Pix[p], i.Pix[p+1], i.Pix[p+2], i.Pix[p+3] = r, g, b, a
		}
		img = i
	default:
		sx, err := s.ConvertFormat(sdl.PIXELFORMAT_RGBX8888, 0)
		if err != nil {
			return nil, err
		}
		i := &image.RGBA{Rect: image.Rect(0, 0, int(s.W), int(s.H))}
		i.Pix = sx.Pixels()
		img = i
	}
	return
}

func setAllAlpha(i image.Image, alpha uint8) {
	dy, dx := i.Bounds().Dy(), i.Bounds().Dx()
	for y := 0; y < dy; y += 1 {
		for x := 0; x < dx; x += 1 {
			i.(draw.Image).Set(x, y, setAlpha(i.At(x, y), 0xff))
		}
	}
}

func setAlpha(c color.Color, alpha uint8) (result color.Color) {
	switch c.(type) {
	case color.NRGBA:
		xc := c.(color.NRGBA)
		xc.A = alpha
		result = xc
	default:
		r, g, b, _ := c.RGBA()
		result = color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), alpha}
	}
	return
}
