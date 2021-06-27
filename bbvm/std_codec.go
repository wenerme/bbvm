package bbvm

import (
	"github.com/wenerme/bbvm/bbasm"
	"golang.org/x/text/encoding/simplifiedchinese"
)

func StdUTF8(rt bbasm.Runtime, std *Std) *Std {
	return &Std{
		BytesToString: func(b []byte) (string, error) {
			return string(b), nil
		},
		StringToBytes: func(s string) ([]byte, error) {
			return []byte(s), nil
		},
	}
}
func StdGBK(rt bbasm.Runtime, std *Std) *Std {
	return &Std{
		BytesToString: func(b []byte) (string, error) {
			bytes, err := simplifiedchinese.GBK.NewDecoder().Bytes(b)
			return string(bytes), err
		},
		StringToBytes: func(s string) ([]byte, error) {
			return simplifiedchinese.GBK.NewEncoder().Bytes([]byte(s))
		},
	}
}
