package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jerebenitez/sched-cli/lib"
	godiffpatch "github.com/sourcegraph/go-diff-patch"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type generateConfig struct {
	FilePath string
	Original string
	Dir		 string
	Src		 string
	Out		 io.Writer
}

//TODO: Fix helpers to get file from repo/folder
// generateDiffCmd represents the generateDiff command
var generateDiffCmd = &cobra.Command{
	Use:   "generate-diff [path to modified file]",
	Aliases: []string{"gd"},
	Short: "Generate a patch for a modified file in the source tree:",
	Long: `Generate a patch for a modified file in the source tree:
  - If the source tree is a git repo, the original version is taken from HEAD.
  - If not, the path to a backup must be provided through --original (-o).`,
  	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		original, err := cmd.Flags().GetString("original")
		if err != nil {
			log.Fatalf("could not read --original: %v", err)
		}

		dir := viper.GetString("sched")
		src := viper.GetString("kernel")

		filePath := args[0]

		err = runGenerateDiff(generateConfig{
			FilePath: filePath,
			Original: original,
			Dir:	  dir,
			Src:	  src,
			Out:	  os.Stdout,
		})

		if err != nil {
			log.Fatalf("generate-diff failed: %v", err)
		}
	},
}

func runGenerateDiff(cfg generateConfig) error {
	fmt.Fprintf(cfg.Out, "Generating diff for %s...\n", cfg.FilePath)
	path, _ := filepath.Split(cfg.FilePath)
	patchPath := filepath.Join(cfg.Dir, path)

	err := os.MkdirAll(patchPath, os.ModePerm)
	if err != nil {
		return err
	}

	var originalContent []byte
	if lib.IsGitRepo(path) {
		originalContent, err = getFromRepo(cfg.FilePath)	
	} else {
		originalContent, err = getFromOriginal(patchPath, cfg.Original)
	}

	if err != nil {
		return err
	}

	newContent, err := os.ReadFile(cfg.FilePath)
	if err != nil {
		return err
	}

	patch := godiffpatch.GeneratePatch(path, string(originalContent), string(newContent))
	if err := os.WriteFile(patchPath, []byte(patch), 0644); err != nil {
		return err
	}

	fmt.Fprintf(cfg.Out, "Diff generated successfully.\n")
	return nil
}

// TODO: Change args and move to lib
func getFromRepo(path string) ([]byte, error) {
	cmd := exec.Command("git", "show", fmt.Sprintf("HEAD:%s", path))
	cmd.Dir = path

	out, err := cmd.Output()
	if err != nil {
		return []byte{}, err
	}

	return out, nil
}

// TODO: Change args and move to lib
func getFromOriginal(path, original string) ([]byte, error) {
	if original == "" {
		return []byte{}, fmt.Errorf("folder is not a git repository, need to provide path to original with --original")
	}

	originalPath := filepath.Join(original, path)
	content, err := os.ReadFile(originalPath)
	if err != nil {
		return []byte{}, err
	}
	return content, nil
}

func init() {
	rootCmd.AddCommand(generateDiffCmd)
	generateDiffCmd.Flags().StringP("original", "o", "", "Path to backup of file. Mandatory if the source tree is not a git repo.")
}
