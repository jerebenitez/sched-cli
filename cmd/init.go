package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	initCmd.Flags().StringP("config", "c", "$HOME/pSched", "Specify folder for project instance. Default: $HOME/pSched")
	initCmd.Flags().StringP("directory", "d", "./", "Specify folder for kernel source. Default: /usr/src")
	initCmd.Flags().StringP("source", "s", "$HOME/pSched", "Specify custom folder for kernel mods cloned repo. Default: $HOME/pSched")
	initCmd.Flags().Bool("clone", false, "When specifying a custom kernel mods path, specify whether you want the tool to clone the repo. Default: false")
}
