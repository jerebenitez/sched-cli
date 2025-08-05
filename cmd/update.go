package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jerebenitez/sched-cli/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Aliases: []string{"u"},
  	Args: cobra.ExactArgs(1),
	ValidArgsFunction: func (cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		kernelDir := viper.GetString("kernel")
		if kernelDir == "" {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		files, err := filepath.Glob(filepath.Join(kernelDir, toComplete+"*"))
		cobra.CheckErr(err)

		var completions []string
		for _, f := range files {
			info, err := os.Stat(f)
			if err != nil {
				continue
			}

			relPath, _ := filepath.Rel(kernelDir, f)
			if info.IsDir() {
				completions = append(completions, relPath+"/")
			} else {
				completions = append(completions, relPath)
			}
		}

		return completions, cobra.ShellCompDirectiveNoSpace
	},
	Short: `Update patches after their corresponding kernel source
		files have changed.`,
	Long: 
`This command rebuilds your modified file from "orig/ + patch" and runs
a 3-way merge with the updated source tree version. The merged result is
then used to regenerate the patch.

How it works
- Reconstructs your modified file using the original file and existing
patch
- Merges it with the updated file from the source tree
- Produces an updated patch based on the merged result

You can use any merge tool by setting "tool" and "toolFormat" in your
config file.
The tool receives three files in the following order:

1. Upstream file (latest from source tree)
2. Original file (from orig/)
3. Reconstructed version (orig/ + patch)

You can change the order using placeholders in "toolFormat". Example:

tool: meld
toolFormat: "%[2]s %[3]s %[1]s"

This opens "meld" with:
- Original file on the left
- Reconstructed version in the middle
- Upstream file on the right`,
	Run: func(cmd *cobra.Command, args []string) {
		tool := viper.GetString("tool")
		toolFormat := viper.GetString("toolFormat")
		src := viper.GetString("kernel")
		dir := viper.GetString("sched")

		filePath := args[0]
		origFile := filepath.Join(dir, "orig", filePath)
		patchFile := filepath.Join(dir, "patches", filePath+".patch")
		upstreamFile := filepath.Join(src, filePath)

		fmt.Printf("%s %s %s %s\n", filePath, origFile, patchFile, upstreamFile)
		if e, err := filesExist(origFile, patchFile, upstreamFile); err != nil || !e {
			log.Fatalf("Some files didn't exist, unable to update.")
		} 

		// 2. Build temp reconstruction
		reconstructed, err := os.CreateTemp("", "reconstructed-*")
		cobra.CheckErr(err)
		defer os.Remove(reconstructed.Name())
		
		err = lib.CopyFile(origFile, reconstructed.Name())
		if err != nil {
			cobra.CheckErr(err)
		}

		if _, err := lib.ApplyPatch(reconstructed.Name(), patchFile, false); err != nil {
			cobra.CheckErr(err)
		}

		merged, err := os.CreateTemp("", "merged-*")
		cobra.CheckErr(err)
		defer os.Remove(merged.Name())

		// 3. Try automatic merge
		manualMerge := false
		diff3Cmd := exec.Command("diff3", "-m", reconstructed.Name(), origFile, upstreamFile)
		out, err :=diff3Cmd.Output()
		
		if err != nil {
			manualMerge = true
			fmt.Printf("Conflict in %s... launching manual merge\n", filePath)
		} else {
			err = os.WriteFile(merged.Name(), out, 0644)
			cobra.CheckErr(err)
		}

		// 4. If it failed, reconstruct manually
		if manualMerge {
			formatArgs := fmt.Sprintf(toolFormat, upstreamFile, origFile, reconstructed.Name())
			args := strings.Fields(formatArgs)
			mergeCmd := exec.Command(tool, args...)
			mergeCmd.Stdout = os.Stdout
			mergeCmd.Stderr = os.Stderr
			mergeCmd.Stdin = os.Stdin
			if err := mergeCmd.Run(); err != nil {
				fmt.Printf("Manual merge failed: %v\n", err)
				return
			}

			// Assume user saved resolved file in reconstructed
			err := lib.CopyFile(reconstructed.Name(), merged.Name())
			cobra.CheckErr(err)
		}

		// 5. Update original
		err = lib.CopyFile(upstreamFile, origFile)
		cobra.CheckErr(err)

		// 6. Regenerate patch
		if err := os.MkdirAll(filepath.Dir(patchFile), os.ModePerm); err != nil {
			cobra.CheckErr(err)
		}

		diffCmd := exec.Command("diff", "-u", origFile, merged.Name())
		diffOut, err := diffCmd.Output()
		if err == nil || len(diffOut) > 0 {
			err := os.WriteFile(patchFile, diffOut, 0644)
			cobra.CheckErr(err)
		}

		fmt.Printf("Updated patch and orig for %s\n", filePath)
	},
}

func filesExist(files ...string) (bool, error) {
	exist := true
	for _, f := range files {
		exists, err := lib.FileExists(f)
		if err != nil {
			return false, err
		}
		exist = exist && exists
	}

	return exist, nil
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
