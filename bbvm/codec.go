package bbvm

import "golang.org/x/text/encoding/simplifiedchinese"

func GBKToString(bytes []byte) (string, error) {
	b, err := simplifiedchinese.GBK.NewDecoder().Bytes(bytes)
	return string(b), err
}

func UTF8ToString(bytes []byte) (string, error) {
	return string(bytes), nil
}
