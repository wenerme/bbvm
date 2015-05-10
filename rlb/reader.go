package rlb
import (
	"image"
	"bytes"
	"encoding/binary"
	"github.com/op/go-logging"
	"os"
	_ "golang.org/x/image/bmp"
)

var log = logging.MustGetLogger("rlb")

// 初始化 Log
func init() {
	format := logging.MustStringFormatter("%{color}%{time:15:04:05} %{level:.4s} %{shortfunc} %{color:reset} %{message}", )
	//	format := logging.MustStringFormatter("%{color}%{time:15:04:05.000} %{longfile} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}", )
	backend1 := logging.NewLogBackend(os.Stdout, "", 0)
	backend1Formatter := logging.NewBackendFormatter(backend1, format)
	logging.SetBackend(backend1Formatter)
}

type RlbConfig struct {
	image.Config
	Name string
	Format string
	Offset int
}
func Decode(r *bytes.Reader) ([]image.Image, []RlbConfig, error) {
	configs, err := DecodeConfig(r)
	if err != nil { return nil, nil, err }

	n := len(configs)
	images := make([]image.Image, n)
	for i := 0; i < n; i ++ {
		r.Seek(int64(configs[i].Offset + 4), os.SEEK_SET)
		images[i], _, err = image.Decode(r)
		if err != nil { return images, configs, err }
	}
	return images, configs, nil
}

func DecodeConfig(r *bytes.Reader) ([]RlbConfig, error) {
	buf := make([]byte, 4)
	name := make([]byte, 32)
	_, err := r.Read(buf)
	if err != nil { return nil, err }
	n := int(binary.LittleEndian.Uint32(buf))
	configs := make([]RlbConfig, n)
	for i := 0; i < n; i ++ {
		buf = make([]byte, 4)
		_, err = r.Read(buf)
		if err != nil { return configs, err }
		c := RlbConfig{}
		c.Offset = int(binary.LittleEndian.Uint32(buf))
		_, err = r.Read(name)
		if err != nil { return configs, err }
		c.Name = bytesToName(name)

		// +4 ignore length
		r.Seek(int64(c.Offset + 4), os.SEEK_SET)

		cfg, f, err := image.DecodeConfig(r)
		if err != nil { return configs, err }
		c.Format = f
		c.Config = cfg
		configs[i] = c
	}

	return configs, nil
}

func bytesToName(b []byte) string {
	for i := 0; i < len(b); i ++ {
		if b[i] ==0 {
			return string(b[0:i])
		}
	}
	return string(b)
}