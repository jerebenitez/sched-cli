package cmd

import (
	"fmt"
	"strings"

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
	fmt.Println("Running checks on patches and file existence.")
	error := false
	dir := rootCmd.PersistentFlags().Lookup("dir").Value.String()
	dir = strings.Trim(dir, "/")
	src := rootCmd.PersistentFlags().Lookup("src").Value.String()
	src = strings.Trim(src, "/")
	origFiles := lib.ReadRecursiveDir(dir, "orig")

	for _, file := range origFiles {
		srcFile := 	src + "/" + file
		origFile := dir + "/orig/" + file

		if lib.FileExists(srcFile) {
			fmt.Printf("Checking %s...", file)
			if lib.FilesAreDifferent(origFile, srcFile) {
				fmt.Printf(" [CHANGED]\n")
			} else {
				fmt.Printf(" [NOT CHANGED]\n")
			}
		} else {
			fmt.Printf("ERROR: %s missing in source tree.\n", file)
			error = true
		}
	}

	if (error) {
		fmt.Println("Errors detected, exiting...")
		return
	}

	fmt.Println("Testing patches...")
	patches := lib.ReadRecursiveDir(dir, "patches")

	for _, patch := range patches {
		srcFile := src + "/" + lib.TrimExtension(patch)
		patchFile := dir + "/patches/" + patch

		if lib.FileExists(srcFile) {
			if lib.CanApplyPatch(srcFile, patchFile) {
				fmt.Printf("Patch %s can be applied.\n", patch)
			} else {
				fmt.Printf("ERROR: Patch %s CAN'T be applied.\n", patch)
			}
		} else {
			fmt.Printf("ERROR: %s not applicable to source tree.\n", patch)
		}
	}
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
