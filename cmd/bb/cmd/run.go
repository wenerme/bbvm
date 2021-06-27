package cmd

import (
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"github.com/wenerme/bbvm/bbasm/parser"
	"github.com/wenerme/bbvm/bbvm/bbrun"
	"go.uber.org/zap"
	"io/ioutil"
	"path"
	"strings"
)

var runOpts = struct {
	Input    string
	Terminal bool
}{

}
var runCmd = &cobra.Command{
	Use:     "run",
	Aliases: []string{"r"},
	Run: func(cmd *cobra.Command, args []string) {
		var bin []byte
		var err error
		if !strings.HasSuffix(runOpts.Input, ".bbin") {
			zap.S().Info("try compile")
			code, err := ioutil.ReadFile(runOpts.Input)
			if err != nil {
				zap.S().With("err", err).Fatal("failed to read file")
			}

			bin, err = parser.Compile(string(code))
			if err != nil {
				zap.S().With("err", err).Fatal("compile failed")
			}
		} else {
			bin, err = ioutil.ReadFile(runOpts.Input)
			if err != nil {
				zap.S().With("err", err).Fatal("failed to read file")
			}
		}

		err = bbrun.Run(path.Base(runOpts.Input), bin)
		if err != nil {
			zap.S().With("err", err).Fatal("run failed")
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("expect a .bbasm or .bbin to run")
		}
		runOpts.Input = args[0]
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVarP(&runOpts.Terminal, "terminal", "t", false, "run in terminal")
}
