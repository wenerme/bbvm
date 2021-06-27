package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wenerme/bbvm/bbvm/fyneui"
	"go.uber.org/zap"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bb",
	Short: "bb is a tool for BeBasic Source Code and Virtual Machine",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fyneui.RunApp(nil)
	},
}

var rootOpts = struct {
	debug bool
}{}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is bb.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&rootOpts.debug, "debug", "D", false, "enable debug")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".bb" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".bb")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	if rootOpts.debug {
		logger, _ := zap.NewDevelopmentConfig().Build()
		zap.ReplaceGlobals(logger)
	}
}
