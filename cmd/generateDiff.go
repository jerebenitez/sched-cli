package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// generateDiffCmd represents the generateDiff command
var generateDiffCmd = &cobra.Command{
	Use:   "generate-diff [path to modified file]",
	Short: "A brief description of your command",
	Long: `Generate a patch for a modified file in the source tree:
  - If the source tree is a git repo, the original version is taken from HEAD.
  - If not, the path to a backup must be provided through --original (-o).`,
  	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("generateDiff called")
	},
}

func init() {
	rootCmd.AddCommand(generateDiffCmd)
	generateDiffCmd.Flags().StringP("original", "o", "", "Path to backup of file. Mandatory if the source tree is not a git repo.")
}
