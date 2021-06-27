package bbasm

import "context"

type Runtime interface {
	Memory

	Push(val int) // push stack
	Pop() int     // pop stack
	Register(typ RegisterType) Register
	Jump(addr int)
	Exit()
	In(ctx context.Context, a int, b int)
	Out(ctx context.Context, a int, b int)
}

type Register interface {
	Get() int
	Set(int)
	Float() float32
	SetFloat(v float32)
}

type Memory interface {
	GetInt(int) int
	SetInt(int, int)
	GetFloat(int) float32
	SetFloat(int, float32)
	GetString(int) string
}
