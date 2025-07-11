package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples and
usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications. This application is a tool to 
generate the needed files to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("config", "c", "$HOME/pSched", "Specify folder for project instance.")
	initCmd.Flags().StringP("directory", "d", "./", "Specify folder for kernel source.")
	initCmd.Flags().StringP("source", "s", "$HOME/pSched", `Specify custom folder for kernel
mods cloned repo.`)
	initCmd.Flags().Bool("clone", false, `When specifying a custom kernel mods path, specify
you want the tool to clone the repo. (default false)`)
}
