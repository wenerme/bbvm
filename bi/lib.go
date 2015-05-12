package bi
import (
	"image"
	"bytes"
	"os"
	"encoding/binary"
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

func DecodeLibGray2(r *bytes.Reader, bo binary.ByteOrder) (*Gray2, error) {
	buf := make([]byte, 4)
	_, err := r.Read(buf)
	if err != nil { return nil, err }
	// Length
	// _ = int(bo.Uint32(buf))

	_, err = r.Read(buf[:2])
	if err != nil { return nil, err }
	w := int(bo.Uint16(buf))
	_, err = r.Read(buf[:2])
	if err != nil { return nil, err }
	h := int(bo.Uint16(buf))

	_, err=r.Read(buf)
	if err != nil { return nil, err }
	_, err=r.Read(buf)
	if err != nil { return nil, err }

	i := NewGray2(image.Rect(0, 0, w, h))
	_, err=r.Read(i.Pix)
	if err != nil { return i, err }

	return i, nil
}