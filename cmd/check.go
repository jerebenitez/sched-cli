package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

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
	// Read files from src/ and orig/
	readRecursiveDir(".", 0)
}

func readRecursiveDir(name string, level int) {
	files, err := os.ReadDir(name)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fileName := name + "/" + file.Name()
		fmt.Println(strings.Repeat(" ", level * 2) + file.Name())
		fi, err := os.Lstat(fileName)

		if err != nil {
			log.Fatal(err)
		}

		switch mode := fi.Mode(); {
		case mode.IsDir(): 
			readRecursiveDir(fileName, level + 1)
		}
	}

}

func init() {
	rootCmd.AddCommand(checkCmd)
}
