package cmd

import (
	"fmt"

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

func runGenerateDiff(path string) {
	fmt.Printf("Generating diff for %s...\n", path)

	// Crear carpetas de ser necesario

	// Checkear si es un repo de git
	//		si no lo es, necesito tener orig para copiar ese archivo
	//		si lo es, copiar el archivo de git

	// Generar el diff con la versión que haya guardado en orig/
}

func init() {
	rootCmd.AddCommand(generateDiffCmd)
	generateDiffCmd.Flags().StringP("original", "o", "", "Path to backup of file. Mandatory if the source tree is not a git repo.")
}
