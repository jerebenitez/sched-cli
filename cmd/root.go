package cmd

import (
	"os"
	"path/filepath"

	"github.com/jerebenitez/sched-cli/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sched-cli",
	Short: "CLI app to manage the development of Project Scheduler.",
	Long: ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "~/.config/sched/config.yaml", "Config file")
	// Create config folder if it doesn't exists
	configPath, err := rootCmd.PersistentFlags().GetString("config")
	cobra.CheckErr(err)

	path, file := filepath.Split(configPath)
	path, err = lib.ResolvePath(path)
	cobra.CheckErr(err)

	err = os.MkdirAll(path, os.ModePerm)
	cobra.CheckErr(err)

	configPath = filepath.Join(path, file)
	err = lib.Touch(configPath)
	cobra.CheckErr(err)

	initConfigViper()
}

func initConfigViper() {
	if cfgFile != "" {
		path, err := lib.ResolvePath(cfgFile)
		cobra.CheckErr(err)
		viper.SetConfigFile(path)
	} else {
		configPath, err := rootCmd.PersistentFlags().GetString("config")
		cobra.CheckErr(err)

		configPath, err = lib.ResolvePath(configPath)
		cobra.CheckErr(err)
		path, _ := filepath.Split(configPath)

		viper.AddConfigPath(path)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	err := viper.ReadInConfig()
	cobra.CheckErr(err)
}


