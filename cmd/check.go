package cmd

import (
	"fmt"
	"log"
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

		if exists, err := lib.FileExists(srcFile); exists {
			fmt.Printf("\tChecking %s...", file)
			if diff, err := lib.FilesAreDifferent(origFile, srcFile); diff {
				fmt.Printf(" [CHANGED]\n")
			} else {
				if err != nil {
					fmt.Printf(" [NOT CHANGED]\n")
				} else {
					log.Fatal(err)
				}
			}
		} else if err == nil {
			fmt.Printf("\tERROR: %s does not exists.\n", file)
			error = true
		} else {
			log.Fatal(err)
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

		if exists, err := lib.FileExists(srcFile); exists {
			if canApply, err := lib.CanApplyPatch(srcFile, patchFile); canApply {
				fmt.Printf("\tPatch %s can be applied.\n", patch)
			} else if err == nil {
				fmt.Printf("\tERROR: Patch %s CAN'T be applied.\n", patch)
			} else {
				log.Fatal(err)
			}
		} else if err == nil {
			fmt.Printf("\tERROR: %s not applicable to source tree.\n", patch)
		} else {
			log.Fatal(err)
		}
	}
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
