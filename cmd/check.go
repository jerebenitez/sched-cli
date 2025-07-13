package cmd

import (
	"fmt"
	"github.com/jerebenitez/sched-cli/lib"

	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check whether modifications are compatible with current source tree.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: checkImpl,
}

func checkImpl(cmd *cobra.Command, args []string) {
	src := rootCmd.PersistentFlags().Lookup("src")
	origFiles := lib.ReadRecursiveDir(src.Value.String(), "orig")
	patchesFiles := lib.ReadRecursiveDir(src.Value.String(), "patches")

	// Compare that there are patches for every file in orig/
	fmt.Print("Checking patches and original files... ")
	missing := lib.CheckPatches(origFiles, patchesFiles)

	if (len(missing) == 0) {
		fmt.Println("OK.")
	} else {
		fmt.Println()
		for _, p := range missing {
			fmt.Printf("ERROR: %s missing its %s file.\n", p.File, p.Missing)
		}
	}

	// Check that files from orig/ exists in kernel source

	// Check that patches can be applied cleanly
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
