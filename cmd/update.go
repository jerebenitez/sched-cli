package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "A brief description of your command",
	Long: 
`Updates one or more patches when the corresponding source files in the kernel
have changed.
Reconstructs your modified version using orig/ + patch, then performs a 3-way
merge with the updated file in the source tree. The result is used to 
regenerate the patch.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("update called")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolP("all", "a", false, "Update all files.")
	updateCmd.Flags().StringP("file", "f", "", "Update specific file.")
	updateCmd.Flags().StringP("tool", "t", "diff3", `Specify tool to perform 3-way merge. You can set the
default one, and add new ones in the config file.`)
}
