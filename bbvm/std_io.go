package bbvm

import (
	"context"
	"fmt"
	"github.com/wenerme/bbvm/bbasm"
	"io"
)

func NewPrintToWriter(out io.Writer) StdBuilder {
	return func(rt bbasm.Runtime, std *Std, ) *Std {
		return &Std{
			PrintLnInt: func(ctx context.Context, v int) {
				_, _ = fmt.Fprintln(out, v)
			},
			PrintLnString: func(ctx context.Context, v StringHandler) {
				_, _ = fmt.Fprintln(out, std.StringGet(ctx, v))
			},
			PrintInt: func(ctx context.Context, v int) {
				_, _ = fmt.Fprint(out, v)
			},
			PrintString: func(ctx context.Context, v StringHandler) {
				_, _ = fmt.Fprint(out, std.StringGet(ctx, v))
			},
			PrintFloat: func(ctx context.Context, v float32) {
				_, _ = fmt.Fprintf(out, "%.6f", v)
			},
			PrintChar: func(ctx context.Context, v int) {
				_, _ = fmt.Fprintf(out, "%c", v)
			},
		}
	}
}

func NewInputFromReader(in io.Reader) StdBuilder {
	return func(rt bbasm.Runtime, std *Std) *Std {
		return &Std{
			InputInt: func(ctx context.Context) (v int) {
				_, _ = fmt.Scanf("%d", &v)
				return v
			},
			InputFloat: func(ctx context.Context) (v float32) {
				_, _ = fmt.Scanf("%f", &v)
				return
			},
			InputString: func(ctx context.Context, dst StringHandler) {
				v := ""
				_, _ = fmt.Scanf("%s", &v)
				std.StringSet(ctx, dst, v)
			},
		}
	}
}
