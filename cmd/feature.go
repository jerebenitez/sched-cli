package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// featureCmd represents the feature command
var featureCmd = &cobra.Command{
	Use:   "feature",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		base := viper.GetString("base")
		features := viper.GetString("features")

		if base == "" || features == "" {
			log.Fatal("Improperly configured project. Run `oot config` for more details.")
		}
	},
}

func init() {
	rootCmd.AddCommand(featureCmd)
}
