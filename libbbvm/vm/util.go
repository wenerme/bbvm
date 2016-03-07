package vm
import (
	"golang.org/x/image/bmp"
	"os"
	"image"
	"github.com/op/go-logging"
)


var log = logging.MustGetLogger("vm")

// 初始化 Log
func init() {
	//	format := logging.MustStringFormatter("%{color}%{time:15:04:05} %{level:.4s} %{shortfunc} %{color:reset} %{message}", )
	format := logging.MustStringFormatter("%{color}%{time:15:04:05.000} %{level:.4s} %{longfile} %{shortfunc} %{color:reset} %{message}")
	backend1 := logging.NewLogBackend(os.Stdout, "", 0)
	backend1Formatter := logging.NewBackendFormatter(backend1, format)
	logging.SetBackend(backend1Formatter)
	logging.SetLevel(logging.DEBUG, "vm")
}

func saveImage(i image.Image, fn string) {
	p, err := os.Create(fn)
	if err != nil {panic(err)}
	err = bmp.Encode(p, i)
	if err != nil {panic(err)}
}
