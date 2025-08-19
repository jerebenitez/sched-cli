package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/jerebenitez/sched-cli/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	kernel, sched, branch string
	gitInit bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "config",
	Aliases: []string{"cfg"},
	Short: "Configure global flags for the CLI.",
	Run: func(cmd *cobra.Command, args []string) {
		// If any flag was provided, write config
		cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
			if f.Changed {
				if f.Name == "kernel" || f.Name == "scheduler" {
					path, err := filepath.Abs(f.Value.String())
					cobra.CheckErr(err)
					err = f.Value.Set(path)
					cobra.CheckErr(err)
				}
				err := viper.WriteConfig()
				cobra.CheckErr(err)
			}
		})

		// Show current config
		if !lib.HasProvidedFlags(cmd) {
			for k, v := range viper.AllSettings() {
				fmt.Printf("\033[1m%s\033[0m: %v\n", k, v)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.PersistentFlags().StringVarP(&kernel, "kernel", "k", "/usr/src", "Path to kernel source tree.")
	err := viper.BindPFlag("kernel", initCmd.PersistentFlags().Lookup("kernel"))
	cobra.CheckErr(err)

	initCmd.PersistentFlags().StringVarP(&sched, "scheduler", "s", ".", "Path to project scheduler source tree.")
	err = viper.BindPFlag("scheduler", initCmd.PersistentFlags().Lookup("scheduler"))
	cobra.CheckErr(err)

	// initCmd.PersistentFlags().StringVarP(&branch, "branch", "b", "main", "Branch of scheduler to work under.")
	// err = viper.BindPFlag("branch", initCmd.PersistentFlags().Lookup("branch"))
	// cobra.CheckErr(err)

	initCmd.Flags().BoolVarP(&gitInit, "git-init", "g", true, "Init a git repo under kernel source tree (recommended).")
}
