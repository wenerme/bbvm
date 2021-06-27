package main

import (
	_ "embed"
	"flag"
	"github.com/wenerme/bbvm/bbvm/bbrun"
	"go.uber.org/zap"
)

//go:embed assets/run.bbin
var bin []byte

//go:embed assets/run.json
var info string

func main() {
	debug := flag.Bool("D", false, "Debug")

	flag.Parse()

	if *debug {
		logger, _ := zap.NewDevelopmentConfig().Build()
		zap.ReplaceGlobals(logger)
	}

	err := bbrun.Run(info, bin)
	if err != nil {
		zap.S().With("err", err).Fatal("run failed")
	}
}
