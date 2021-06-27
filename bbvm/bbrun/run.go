package bbrun

import (
	"context"
	"github.com/wenerme/bbvm/bbvm"
	"go.uber.org/zap"
	"os"
	"reflect"
)

func Run(name string, rom []byte) error {
	vm := bbvm.NewInstance()
	vm.Std = &bbvm.Std{}
	vm.Std.Use(bbvm.StdBase(vm, vm.Std))
	vm.Std.Use(bbvm.StdGBK(vm, vm.Std))
	vm.Std.Use(bbvm.StdStringRes(vm, vm.Std))
	vm.Std.Use(bbvm.StdStringFunc(vm, vm.Std))
	vm.Std.Use(bbvm.NewPrintToWriter(os.Stdout)(vm, vm.Std))
	// report missing std
	rv := reflect.ValueOf(vm.Std).Elem()
	rt := rv.Type()
	n := rv.NumField()
	for i := 0; i < n; i++ {
		f := rv.Field(i)
		if f.IsNil() {
			zap.S().With("std", rt.Field(i).Name).Debug("missing std")
		}
	}
	vm.Load(rom)
	return vm.Run(context.Background())
}
