/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new feature branch and set up workspace based on the configured base repo",
	Long: `This command initializes a new feature by creating a branch in the features 
repository and preparing a working directory. By default, the feature is based 
on the configured base repository and the specified base branch.

The new feature branch is named "feature-<name>" and is set up to track future 
changes independently of other features.
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlag(
			fmt.Sprintf("features.%s.base-branch", args[0]),
			cmd.Flags().Lookup("base-branch"),
		)
		cobra.CheckErr(err)
		viper.Set(fmt.Sprintf("features.%s.status", args[0]), "open")
		viper.Set("features.last", args[0])
		cobra.CheckErr(viper.WriteConfig())

		fmt.Println("new feature:", args[0])
	},
}

func init() {
	featureCmd.AddCommand(newCmd)

	newCmd.Flags().StringVarP(&basePath, "base-branch", "b", "", "Branch from base repo to start working from.")
	cobra.CheckErr(newCmd.MarkFlagRequired("base-branch"))
}
