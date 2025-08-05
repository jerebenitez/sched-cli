package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: `Update patches after their corresponding kernel source
		files have changed.`,
	Long: 
`This command rebuilds your modified file from "orig/ + patch" and runs
a 3-way merge with the updated source tree version. The merged result is
then used to regenerate the patch.

How it works
- Reconstructs your modified file using the original file and existing
patch
- Merges it with the updated file from the source tree
- Produces an updated patch based on the merged result

You can use any merge tool by setting "tool" and "toolFormat" in your
config file.
The tool receives three files in the following order:

1. Upstream file (latest from source tree)
2. Original file (from orig/)
3. Reconstructed version (orig/ + patch)

You can change the order using placeholders in "toolFormat". Example:

tool: meld
toolFormat: "%[2]s %[3]s %[1]s"

This opens "meld" with:
- Original file on the left
- Reconstructed version in the middle
- Upstream file on the right`,
	Run: func(cmd *cobra.Command, args []string) {
		tool := viper.GetString("tool")
		toolFormat := viper.GetString("toolFormat")
		formated := fmt.Sprintf(toolFormat, "file1", "file2", "file3")

		fmt.Printf("%s " + formated, tool)
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolP("all", "a", false, "Update all files.")
	updateCmd.Flags().StringP("file", "f", "", "Update specific file.")
}
