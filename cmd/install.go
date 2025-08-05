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
	Long: `Copies new files over and applies patches to existing files`,
	Run: func(cmd *cobra.Command, args []string) {
		dir := viper.GetString("sched")
		src := viper.GetString("kernel")

		err := runInstall(installConfig{
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

	source := filepath.Join(cfg.Dir, "src")
	files, err := lib.ReadRecursiveDir(os.DirFS(source))
	if err != nil {
		return err
	}

	if err := lib.InstallFiles(cfg.Src, files); err != nil {
		return err
	}

	fmt.Fprintf(cfg.Out, "Files installed!\n")

	fmt.Fprintf(cfg.Out, "Applying patches...\n")

	patchesPath := filepath.Join(cfg.Dir, "patches")
	patches, err := lib.ReadRecursiveDir(os.DirFS(patchesPath))
	if err != nil {
		return err
	}

	if err := lib.ApplyPatches(cfg.Src, cfg.Dir, patches); err != nil {
		return err
	}

	fmt.Fprintf(cfg.Out, "Installation completed. You may now compile and install the kernel.")

	fmt.Fprintf(cfg.Out, "Scheduler installed.")
	return nil
}


func init() {
	rootCmd.AddCommand(installCmd)
}
