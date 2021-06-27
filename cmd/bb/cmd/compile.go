package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wenerme/bbvm/bbasm/parser"
	"go.uber.org/zap"
	"io/ioutil"
	"path"
)

var compileOpt = struct {
	output string
	input  string
}{}

var compileCmd = &cobra.Command{
	Use:     "compile",
	Aliases: []string{"c"},
	Run: func(cmd *cobra.Command, args []string) {
		if compileOpt.input == "" {
			zap.S().Fatal("missing input")
		}

		if compileOpt.output == "" {
			if len(args) > 1 {
				compileOpt.output = args[1]
			}
		}
		if compileOpt.output == "" {
			ext := path.Ext(compileOpt.input)
			compileOpt.output = compileOpt.input[:len(compileOpt.input)-len(ext)] + ".bbin"
		}
		zap.S().Infof("compile %v to %v", compileOpt.input, compileOpt.output)
		file, err := ioutil.ReadFile(compileOpt.input)
		if err != nil {
			zap.S().With("err", err).Fatal("failed to read input %v", compileOpt.input)
		}
		bytes, err := parser.Compile(string(file))
		if err != nil {
			zap.S().With("err", err).Fatal("compile failed %v", compileOpt.input)
		}
		err = ioutil.WriteFile(compileOpt.output, bytes, 0644)
		if err != nil {
			zap.S().With("err", err).Fatal("failed to write output %v", compileOpt.output)
		}
	},
}

func init() {
	rootCmd.AddCommand(compileCmd)
	compileCmd.Flags().StringVarP(&compileOpt.output, "output", "o", "", "output bbin")
	compileCmd.Flags().StringVarP(&compileOpt.input, "input", "i", "", "input bbasm")
}
