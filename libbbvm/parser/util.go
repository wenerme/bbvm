package parser

import "strconv"

func parseInt(val string) (ret int32, e error) {
	var i int64
	if len(val) > 2 && val[0] == '0' {
		switch val[0:2] {
		case "0x", "0X":
			i, e = strconv.ParseInt(val[2:], 16, 32)
		case "0b", "0B":
			i, e = strconv.ParseInt(val[2:], 2, 32)
		default:
			i, e = strconv.ParseInt(val[1:], 8, 32)
		}
	} else {
		i, e = strconv.ParseInt(val, 10, 32)
	}

	ret = int32(i)
	return
}
