package bi
import (
	"image"
	"os"
	"encoding/binary"
	"./bc"
	"io"
	"errors"
	"fmt"
	_ "golang.org/x/image/bmp"
	"image/draw"
	"image/color"
)

type Config struct {
	image.Config
	Name string
	Format string
	Offset int
}

type imageFormat string
const (
	libRGB565Format imageFormat = "libRGB565"
	libGray2BEFormat imageFormat = "libGray2BE"
	libGray2LEFormat imageFormat = "libGray2LE"
	rlbFormat imageFormat = "rlb"
	dlxFormat imageFormat = "dlx"
)
func DecodeConfig(r io.ReadSeeker) (configs []Config, err error) {
	f, err := detectImageFormat(r)
	if err != nil {return }
	_, configs, err= decodeImage(f, r, false)
	return
}
func Decode(r io.ReadSeeker) (images []image.Image, configs []Config, err error) {
	f, err := detectImageFormat(r)
	if err != nil {return }
	images, configs, err= decodeImage(f, r, true)
	return
}
func DecodeAt(r io.ReadSeeker, n int) (i image.Image, err error) {
	f, err := detectImageFormat(r)
	if err != nil {return }

	_, configs, err := decodeImage(f, r, false)
	if err != nil {return }

	if n >= len(configs) || n < 0 {
		err = errors.New(fmt.Sprintf("Only got %d image here, but you want %d", len(configs), n))
		return
	}
	i, err = decodeOneImage(f, r, configs[n])
	return
}
// will not affect the position
func decodeImage(f imageFormat, r io.ReadSeeker, decodeImages bool) (images []image.Image, configs []Config, err error) {
	pos, err := r.Seek(0, os.SEEK_CUR)
	if err != nil {return }
	defer func() {
		_, e2 := r.Seek(pos, os.SEEK_SET)
		if e2 != nil {
			if err == nil {
				err = e2
			}else {
				err = errors.New(fmt.Sprint(err, e2))
			}
		}
	}()

	var bo binary.ByteOrder
	if f == libGray2BEFormat {
		bo = binary.BigEndian
	}else {
		bo = binary.LittleEndian
	}

	var n uint32
	var offset, length uint32
	var w, h uint16

	err = binary.Read(r, bo, &n)
	if err != nil {return }

	configs = make([]Config, n)
	if decodeImages {
		images = make([]image.Image, n)
	}
	gap := int64(0)
	nextConfigOffset, err := r.Seek(0, os.SEEK_CUR)
	if err != nil {return }

	name := make([]byte, 32)
	if f == rlbFormat {gap = 32}
	for i := uint32(0); i < n; i += 1 {
		r.Seek(nextConfigOffset, os.SEEK_SET)
		if err != nil {return }
		nextConfigOffset += gap + 4

		c := Config{}

		err = binary.Read(r, bo, &offset)
		if err != nil {return }

		// rlb got 32 byte for name
		if f  == rlbFormat {
			_, err = r.Read(name)
			if err != nil { return }
			c.Name = bytesToName(name)
		}

		_, err = r.Seek(int64(offset), os.SEEK_SET)
		if err != nil {return }

		err = binary.Read(r, bo, &length)
		if err != nil {return }

		c.Offset = int(offset)
		c.Format = string(f)

		switch f{
		case rlbFormat:
			c.Config, c.Format, err = image.DecodeConfig(r)
			if err != nil {return }
		case libGray2BEFormat, libGray2LEFormat, libRGB565Format:
			err = binary.Read(r, bo, &w)
			if err != nil {return }
			err = binary.Read(r, bo, &h)
			if err != nil {return }

			c.Width = int(w)
			c.Height = int(h)

			if f == libRGB565Format {
				c.ColorModel = bc.RGB565Model
			}else {
				c.ColorModel = bc.Gray2Model
			}

		default:
			err=errors.New("Unsuported image format "+string(f))
			return
		}

		configs[i] = c

		if decodeImages {
			images[i], err = decodeOneImage(f, r, c)
			if err != nil {return }
		}
	}

	return
}
// will affect the position
func decodeOneImage(f imageFormat, r io.ReadSeeker, c Config) (i image.Image, err error) {
	// Seek to this config offset
	switch f{
	case libGray2BEFormat, libGray2LEFormat, libRGB565Format:
		// + 4 ignore length + 2 ignore size
		_, err = r.Seek(int64(c.Offset + 6), os.SEEK_SET)
		if err != nil {return }
	default:
		// + 4 ignore length
		_, err = r.Seek(int64(c.Offset + 4), os.SEEK_SET)
		if err != nil {return }
	}

	switch f{
	case rlbFormat:
		i, _, err = image.Decode(r)
		if err != nil {return }
		// 所有 alpha 为 0,重置为 0xff
		dy, dx := i.Bounds().Dy(), i.Bounds().Dx()
		for y := 0; y < dy; y +=1 {
			for x := 0; x < dx; x +=1 {
				i.(draw.Image).Set(x, y, setAlpha(i.At(x, y), 0xff))
			}
		}
	case libGray2BEFormat, libGray2LEFormat:
		img := NewGray2(image.Rect(0, 0, c.Width, c.Height))
		_, err = r.Read(img.Pix)
		if err != nil {return }
		// 保证所有数据为 BG
		if f == libGray2LEFormat {
			reverseBit2(img.Pix)
		}
		i = img
	case libRGB565Format:
		_, err = r.Seek(int64(c.Offset + 4), os.SEEK_SET)
		if err != nil {return }

		img := NewRGB565(image.Rect(0, 0, c.Width, c.Height))
		_, err = r.Read(img.Pix)
		if err != nil {return }
		reverseBit16(img.Pix)
		i = img
	default:
		err = errors.New("Unknow image format: "+string(f))
	}

	return
}

// will not affect the position
func detectImageFormat(r io.ReadSeeker) (f imageFormat, err error) {
	pos, err := r.Seek(0, os.SEEK_CUR)
	if err != nil {return }
	defer func() {
		_, e2 := r.Seek(pos, os.SEEK_SET)
		if e2 != nil {
			if err == nil {
				err = e2
			}else {
				err = errors.New(fmt.Sprint(err, e2))
			}
		}
	}()

	buf := make([]byte, 4)
	var bo binary.ByteOrder = binary.LittleEndian

	_, err = r.Read(buf)
	if err != nil { return  }
	n := int(bo.Uint32(buf))
	if n > 0xffff {
		bo = binary.BigEndian
		n = int(bo.Uint32(buf))
		if n > 0xffff {
			err = errors.New("Wrong format,image nummber exceed 0xFFFF: "+fmt.Sprint(n))
		}
	}


	var offset, length uint32
	var w, h uint16
	err = binary.Read(r, bo, &offset)
	if err != nil { return  }

	_, err = r.Seek(int64(offset), os.SEEK_SET)
	if err != nil { return  }
	err = binary.Read(r, bo, &length)
	if err != nil { return  }

	_, err = r.Read(buf)
	if err != nil { return  }
	if buf[0] == 'B' && buf[1] == 'M' {
		f = rlbFormat
		return
	}

	w = bo.Uint16(buf)
	h = bo.Uint16(buf[2:])

	if length == uint32(w)*uint32(h)/4+12 {
		if bo == binary.BigEndian {
			f = libGray2BEFormat
		}else {
			f = libGray2LEFormat
		}
	}else if length == uint32(w)*uint32(h)*2+12 {
		f = libRGB565Format
	}
	if f == "" {
		err = errors.New(fmt.Sprintf("Unknow image format n: %d length: %d size: %d*%d", n, length, w, h))
	}
	return
}

func reverseBit2(b []byte) {
	// 00 10 11 01 ->
	// 01 11 10 00
	l := len(b)
	for i := 0; i < l; i ++ {
		v := b[i]
		b[i] = v&3 << 6 | (v >> 2)& 3 << 4 | (v >> 4)& 3 << 2 | (v >> 6)
	}
}

func reverseBit16(b []byte) {
	l := len(b)
	if l % 2 != 0 { l -= 1 }
	for i := 0; i < l; i += 2 {
		b[i], b[i+1] = b[i+1], b[i]
	}
}

func bytesToName(b []byte) string {
	for i := 0; i < len(b); i ++ {
		if b[i] ==0 {
			return string(b[0:i])
		}
	}
	return string(b)
}

func setAlpha(c color.Color, alpha uint8) (result color.Color) {
	switch c.(type){
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