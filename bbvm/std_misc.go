package bbvm

import (
	"context"
	"github.com/wenerme/bbvm/bbasm"
	"math"
	"math/rand"
	"time"
)

func StdMisc(rt bbasm.Runtime, std *Std) *Std {
	random := rand.New(rand.NewSource(0))

	return &Std{
		FloatToInt: func(ctx context.Context, v float32) int {
			return int(v)
		},
		IntToFloat: func(ctx context.Context, v int) float32 {
			return float32(v)
		},
		Sin: func(ctx context.Context, a float32) float32 {
			return float32(math.Sin(float64(a)))
		},
		Cos: func(ctx context.Context, a float32) float32 {
			return float32(math.Cos(float64(a)))
		},
		Tan: func(ctx context.Context, a float32) float32 {
			return float32(math.Tan(float64(a)))
		},
		Sqrt: func(ctx context.Context, a float32) float32 {
			return float32(math.Sqrt(float64(a)))
		},
		IntAbs: func(ctx context.Context, a int) int {
			if a >= 0 {
				return a
			}
			return -a
		},
		FloatAbs: func(ctx context.Context, a float32) float32 {
			if a >= 0 {
				return a
			}
			return -a
		},
		Delay: func(ctx context.Context, msec int) {
			time.Sleep(time.Duration(msec) * time.Millisecond)
		},
		RandSeed: func(ctx context.Context, seed int) {
			random.Seed(int64(seed))
		},
		Rand: func(ctx context.Context, n int) int {
			return random.Intn(n)
		},
		VmTest: func(ctx context.Context) {
		},
	}
}
