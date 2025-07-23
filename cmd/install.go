package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jerebenitez/sched-cli/lib"
	"github.com/spf13/cobra"
)

type installConfig struct {
	Src string
	Dir string
	Out io.Writer
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Aliases: []string{"i"},
	Short: "Apply kernel modifications to source tree.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := cmd.Root().PersistentFlags().GetString("dir")
		if err != nil {
			log.Fatalf("could not read --dir: %v", err)
		}

		src, err := cmd.Root().PersistentFlags().GetString("dir")
		if err != nil {
			log.Fatalf("could not read --src: %v", err)
		}

		fmt.Printf("Installing files to %s...\n", src)
		err = runInstall(installConfig{
			Src: src,
			Dir: dir,
			Out: os.Stdout,
		})

		if err != nil {
			log.Fatalf("install failed: %v", err)
		}
	},
}

func runInstall(cfg installConfig) error {
	fmt.Fprintf(cfg.Out, "Installing files to %s...\n", cfg.Src)
	if err := installFiles(cfg); err != nil {
		return fmt.Errorf("installFiles: %v", err)
	}
	fmt.Fprintf(cfg.Out, "Files installed!\n")

	fmt.Fprintf(cfg.Out, "Applying patches...\n")
	if err := applyPatches(cfg); err != nil {
		return fmt.Errorf("applyPatches: %v", err)
	}
	fmt.Fprintf(cfg.Out, "Installation completed. You may now compile and install the kernel.")

	fmt.Fprintf(cfg.Out, "Scheduler installed.")
	return nil
}

func applyPatches(cfg installConfig) error {
	patchesPath := filepath.Join(cfg.Dir, "patches")

	patches, err := lib.ReadRecursiveDir(os.DirFS(patchesPath))
	if err != nil {
		return fmt.Errorf("readRecursiveDir error: %v", err)
	}

	for _, patch := range patches {
		result, err := lib.ApplyPatch(
			filepath.Join(cfg.Src, lib.TrimExtension(patch)),
			filepath.Join(cfg.Dir, patch),
			true,
		)
		if err != nil {
			return err
		} else if !result {
			return fmt.Errorf(
				"unable to patch %s: file not compatible with patch",
				lib.TrimExtension(patch),
			)
		}
	}

	return nil
}

func installFiles(cfg installConfig) error {
	source := filepath.Join(cfg.Dir, "src")

	files, err := lib.ReadRecursiveDir(os.DirFS(source))
	if err != nil {
		return fmt.Errorf("error while reading dir %s: %v", source, err)
	}

	for _, file := range files {
		path, _ := filepath.Split(file)
		target := filepath.Join(cfg.Src, path)
		if err := os.MkdirAll(target, os.ModePerm); err != nil {
			return fmt.Errorf("error creating folder %s", target)
		}

		// TODO: Test that this behaves nicely when compiling the kernel
		err := os.Link(filepath.Join(source, file), target)
		if err != nil {
			return fmt.Errorf("error copying file %s", file)
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(installCmd)
}
