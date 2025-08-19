package cmd

import (
	"fmt"
	"os/exec"

	"github.com/jerebenitez/sched-cli/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	kernel, sched string
	gitInit bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:`,
	Run: func(cmd *cobra.Command, args []string) {
		// If any flag was provided, write config
		cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
			if f.Changed {
				// TODO: Resolve paths before writing
				err := viper.WriteConfig()
				cobra.CheckErr(err)
			}
		})

		if gitInit && !lib.IsGitRepo(kernel) {
			cmd := exec.Command("git", "init", kernel)
			output, err := cmd.Output()
			cobra.CheckErr(err)
			fmt.Println(string(output))
		}	
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.PersistentFlags().StringVarP(&kernel, "kernel", "k", "/usr/src", "Path to kernel source tree.")
	err := viper.BindPFlag("kernel", initCmd.PersistentFlags().Lookup("kernel"))
	cobra.CheckErr(err)

	initCmd.PersistentFlags().StringVarP(&sched, "sched", "s", ".", "Path to project scheduler source tree.")
	err = viper.BindPFlag("sched", initCmd.PersistentFlags().Lookup("sched"))
	cobra.CheckErr(err)

	initCmd.Flags().BoolVarP(&gitInit, "git-init", "g", true, "Init a git repo under kernel source tree (recommended).")
}
