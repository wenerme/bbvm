package bbvm

import (
	"context"
	"github.com/wenerme/bbvm/bbasm"
	"log"
	"strconv"
	"strings"
)

type StringHdr struct {
	H    int
	V    string
	Free bool
}

func (v *StringHdr) Handler() int {
	return v.H
}

func StdStringRes(rt bbasm.Runtime, std *Std) *Std {
	id := 0
	hdrs := map[int]*StringHdr{}
	return &Std{
		AllocString: func(ctx context.Context) StringHandler {
			id--
			hdr := &StringHdr{
				H: id,
			}
			hdrs[hdr.H] = hdr
			return hdr
		},
		FreeString: func(ctx context.Context, hdr StringHandler) {
			h := hdr.(*StringHdr)
			h.Free = true
			hdrs[h.H] = nil
		},
		StringGet: func(ctx context.Context, hdr StringHandler) string {
			return hdr.(*StringHdr).V
		},
		StringSet: func(ctx context.Context, hdr StringHandler, v string) {
			hdr.(*StringHdr).V = v
		},
		StringOf: func(ctx context.Context, hdr int) StringHandler {
			if hdr >= 0 {
				return &StringHdr{
					V: rt.GetString(hdr),
				}
			}
			h := hdrs[hdr]
			// handle address
			if h == nil {
				log.Println("invalid string hdr", hdr)
			}
			return h
		},
		//StringToHandler: func(ctx context.Context, hdr StringHandler) int {
		//	return hdr.(*StringHdr).H
		//},
	}
}
func StdStringFunc(rt bbasm.Runtime, std *Std) *Std {
	return &Std{
		StringToInt: func(ctx context.Context, hdr StringHandler) (int, error) {
			return strconv.Atoi(std.StringGet(ctx, hdr))
		},
		StringToFloat: func(ctx context.Context, hdr StringHandler) (float32, error) {
			v, err := strconv.ParseFloat(std.StringGet(ctx, hdr), 32)
			return float32(v), err
		},
		StringCopy: func(ctx context.Context, dst StringHandler, src StringHandler) {
			std.StringSet(ctx, dst, std.StringGet(ctx, src))
		},
		StringConcat: func(ctx context.Context, a StringHandler, b StringHandler) {
			std.StringSet(ctx, a, std.StringGet(ctx, a)+std.StringGet(ctx, b))
		},
		StringLength: func(ctx context.Context, hdr StringHandler) int {
			return len(std.StringGet(ctx, hdr))
		},
		StringGetAscii: func(ctx context.Context, hdr StringHandler, idx int) int {
			return int(std.StringGet(ctx, hdr)[idx])
		},
		StringSetAscii: func(ctx context.Context, hdr StringHandler, idx, v int) {
			s := []byte(std.StringGet(ctx, hdr))
			s[idx] = byte(v % 0xFF)
			std.StringSet(ctx, hdr, string(s))
		},
		StringCompare: func(ctx context.Context, a StringHandler, b StringHandler) int {
			return strings.Compare(std.StringGet(ctx, a), std.StringGet(ctx, b))
		},
		StringFind: func(ctx context.Context, hdr StringHandler, f StringHandler, offset int) int {
			s := std.StringGet(ctx, hdr)
			return strings.Index(s[offset:], std.StringGet(ctx, f))
		},
		StringLeft: func(ctx context.Context, dst StringHandler, hdr StringHandler, len int) {
			std.StringSet(ctx, dst, std.StringGet(ctx, hdr)[0:len])
		},
		StringRight: func(ctx context.Context, dst StringHandler, hdr StringHandler, l int) {
			s := std.StringGet(ctx, hdr)
			std.StringSet(ctx, dst, s[len(s)-l:])
		},
		StringMid: func(ctx context.Context, dst StringHandler, hdr StringHandler, idx int, len int) {
			s := std.StringGet(ctx, hdr)
			std.StringSet(ctx, dst, s[idx:len])
		},
	}
}
