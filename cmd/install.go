package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jerebenitez/sched-cli/lib"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Apply kernel modifications to source tree.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: runInstall,
}

func runInstall(cmd *cobra.Command, args []string) {
	src := rootCmd.PersistentFlags().Lookup("src").Value.String()
	src = strings.Trim(src, "/")
	dir := rootCmd.PersistentFlags().Lookup("dir").Value.String()
	dir = strings.Trim(dir, "/")

	fmt.Printf("Installing files to %s...\n", src)

	files := lib.ReadRecursiveDir(dir, "src")

	for _, file := range files {
		path := lib.GetPath(file)
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			log.Fatalf("Error creating folder %s!. Quitting.\n", path)
			return
		}

		// TODO: Test that this behaves nicely when compiling the kernel
		err := os.Link(filepath.Join(dir, file), filepath.Join(src, file))
		if err != nil {
			log.Fatalf("Error copying file %s!. Quitting.\n", file)
		}

	}
	
	fmt.Println("Applying patches...")

	patches := lib.ReadRecursiveDir(dir, "patches")

	for _, patch := range patches {
		result, err := lib.ApplyPatch(
			filepath.Join(src, lib.TrimExtension(patch)),
			filepath.Join(dir, patch),
			true,
		)
		if err != nil {
			log.Fatalf("Unable to patch %s: %s\n", lib.TrimExtension(patch), err)
		} else if !result {
			log.Fatalf("Unable to patch %s: file not compatible with patch. Run update first!\n", lib.TrimExtension(patch))
		}
	}

	fmt.Println("Installation completed. You may now compile and install the kernel.")
}

func init() {
	rootCmd.AddCommand(installCmd)
}
