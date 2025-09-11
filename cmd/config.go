/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/jerebenitez/sched-cli/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var basePath, featuresPath string

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure the local path to the repositories used by this tool.",
	Long: `This command lets you set or update the location of both the repo you're working
on, and the repo where you're storing the patches. The configured path is stored
locally so that other commands know where to find and operate on your repos
without needing to specify the path each time.`,
	Run: func(cmd *cobra.Command, args []string) {
		noFlags := true
		cmd.Flags().VisitAll(func (f *pflag.Flag) {
			if f.Changed {
				path, err := lib.ResolvePath(f.Value.String())
				cobra.CheckErr(err)
				cobra.CheckErr(f.Value.Set(path))
				fmt.Printf("%s repo set to %s\n", f.Name, path)
				noFlags = false
			}
		})

		if noFlags {
			// Show current config
			for k, v := range viper.AllSettings() {
				fmt.Printf("\033[1m%s\033[0m: %v\n", k, v)
			}
		} else {
			cobra.CheckErr(viper.WriteConfig())
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().StringVarP(&basePath, "base", "b", "", "Path to base repo, where you'll make modifications.")
	err := viper.BindPFlag("base", configCmd.Flags().Lookup("base"))
	cobra.CheckErr(err)

	configCmd.Flags().StringVarP(&featuresPath, "features", "f", "", "Path to repo where you'll store the patches.")
	err = viper.BindPFlag("features", configCmd.Flags().Lookup("features"))
	cobra.CheckErr(err)
}
