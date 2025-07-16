package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jerebenitez/sched-cli/lib"
	godiffpatch "github.com/sourcegraph/go-diff-patch"
	"github.com/spf13/cobra"
)

// generateDiffCmd represents the generateDiff command
var generateDiffCmd = &cobra.Command{
	Use:   "generate-diff [path to modified file]",
	Short: "Generate a patch for a modified file in the source tree:",
	Long: `Generate a patch for a modified file in the source tree:
  - If the source tree is a git repo, the original version is taken from HEAD.
  - If not, the path to a backup must be provided through --original (-o).`,
  	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runGenerateDiff(args[0])
	},
}

func runGenerateDiff(patchPath string) {
	// TODO: Refactor this function
	// TODO: Get path from --src
	// TODO: Test function
	fmt.Printf("Generating diff for %s...\n", patchPath)

	path, _ := filepath.Split(patchPath)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatalf("Error creating folders!. Quitting.\n")
	}

	var originalContent []byte

	if lib.IsGitRepo(path) {
		cmd := exec.Command("git", "show", fmt.Sprintf("HEAD:%s", path))
		cmd.Dir = path
		out, err := cmd.Output()
		if err != nil {
			log.Fatalf("Error reading original from git: %v", err)
			return
		}
		originalContent = out
	} else {
		originalPath := rootCmd.PersistentFlags().Lookup("original").Value.String()
		if originalPath == "" {
			log.Fatalf("Folder is not a git repository, need to provide path to original with --original.")
			return
		}

		originalFullPath := filepath.Join(originalPath, path)
		content, err := os.ReadFile(originalFullPath)
		if err != nil {
			log.Fatalf("Cannot read original file from %s: %v", originalFullPath, err)
			return
		}
		originalContent = content
	}

	newContent, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Cannot read current file: %v", err)
		return
	}

	patch := godiffpatch.GeneratePatch(path, string(originalContent), string(newContent))
	if err := os.WriteFile(patchPath, []byte(patch), 0644); err != nil {
		log.Fatalf("Cannot write patch to %s: %v", patchPath, err)
		return
	}

	fmt.Println("Diff generated successfully.")
}

func init() {
	rootCmd.AddCommand(generateDiffCmd)
	generateDiffCmd.Flags().StringP("original", "o", "", "Path to backup of file. Mandatory if the source tree is not a git repo.")
}
