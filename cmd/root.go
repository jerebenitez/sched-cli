package cmd

import (
	"os"

	"github.com/spf13/cobra"
)



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
	rootCmd.PersistentFlags().StringP("src", "s", "/usr/src", "Path to kernel source code.")
	rootCmd.PersistentFlags().StringP("dir", "d", ".", "Path to shceduler source code.")
}


