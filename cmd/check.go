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
	Short: `Check whether modifications are compatible with current
		source tree.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running checks on patches and file existence.")

		dir := viper.GetString("scheduler")
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
			fmt.Fprintf(cfg.Out, "Checking %s...", file)

			diff, err := lib.FilesAreDifferent(origFile, srcFile)
			if err != nil {
				return err
			}

			if diff {
				fmt.Fprintf(cfg.Out, " \033[33m[CHANGED]\033[0m\n")
			} else {
				fmt.Fprintf(cfg.Out, " \033[32m[NOT CHANGED]\033[0m\n")
			}
		} else {
			fmt.Fprintf(cfg.Out, "\033[91mERROR: %s does not exist.\033[0m\n", file)
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
				fmt.Fprintf(cfg.Out, "\033[32mPatch %s can be applied.\033[0m\n", patch)
			} else {
				fmt.Fprintf(cfg.Out, "\033[91mERROR: Patch %s CAN'T be applied.\033[0m\n", patch)
			}
		} else {
			fmt.Fprintf(cfg.Out, "\033[91mERROR: %s not applicable to source tree.\033[0m\n", patch)
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
