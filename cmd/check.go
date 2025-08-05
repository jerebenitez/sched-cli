package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/jerebenitez/sched-cli/lib"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type checkConfig struct {
	Dir string
	Src string
	Out io.Writer
}

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Aliases: []string{"c"},
	Short: "Check whether modifications are compatible with current source tree.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running checks on patches and file existence.")

		dir := viper.GetString("sched")
		src := viper.GetString("kernel")

		err := runCheck(checkConfig{
			Dir: dir,
			Src: src,
			Out: os.Stdout,
		})

		if err != nil {
			log.Fatal(err)
		}
	},
}

func runCheck(cfg checkConfig) error {
	if err := checkFiles(cfg); err != nil {
		return err
	}

	if err := checkPatches(cfg); err != nil {
		return fmt.Errorf("patches check failed: %w", err)
	}

	fmt.Fprintln(cfg.Out, "All checks passed.")

	return nil
}

// TODO: Change args and move to lib
func checkFiles(cfg checkConfig) error {
	origFiles, err := lib.ReadRecursiveDir(os.DirFS(filepath.Join(cfg.Dir, "orig")))
	if err != nil {
		return err
	}

	error := false
	for _, file := range origFiles {
		srcFile := filepath.Join(cfg.Src, file)
		origFile := filepath.Join(cfg.Dir, "orig", file)

		exists, err := lib.FileExists(srcFile)
		if err != nil {
			return err
		}

		if exists {
			fmt.Fprintf(cfg.Out, "\tChecking %s...", file)

			diff, err := lib.FilesAreDifferent(origFile, srcFile)
			if err != nil {
				return err
			}

			if diff {
				fmt.Fprintf(cfg.Out, " [CHANGED]\n")
			} else {
				fmt.Fprintf(cfg.Out, " [NOT CHANGED]\n")
			}
		} else {
			fmt.Fprintf(cfg.Out, "\tERROR: %s does not exist.\n", file)
			error = true
		}
	}

	if (error) {
		return fmt.Errorf("errors detected")
	}

	return nil
}

// TODO: Change args and move to lib
func checkPatches(cfg checkConfig) error {
	fmt.Fprintf(cfg.Out, "Testing patches.\n")

	patches, err := lib.ReadRecursiveDir(os.DirFS(filepath.Join(cfg.Dir, "patches")))
	if err != nil {
		return err
	}

	for _, patch := range patches {
		srcFile := filepath.Join(cfg.Src, lib.TrimExtension(patch))
		patchFile := filepath.Join(cfg.Dir, "patches", patch)

		exists, err := lib.FileExists(srcFile)
		if err != nil {
			return err
		}

		if exists {
			canApply, err := lib.ApplyPatch(srcFile, patchFile, true)
			if err != nil {
				return err
			}

			if canApply {
				fmt.Fprintf(cfg.Out, "\tPatch %s can be applied.\n", patch)
			} else {
				fmt.Fprintf(cfg.Out, "\tERROR: Patch %s CAN'T be applied.\n", patch)
			}
		} else {
			fmt.Fprintf(cfg.Out, "\tERROR: %s not applicable to source tree.\n", patch)
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
