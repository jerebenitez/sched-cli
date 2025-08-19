package lib

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func HasProvidedFlags(cmd *cobra.Command) bool {
	flags := false
	cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		if f.Changed {
			flags = true
		}
	})
	cmd.Flags().VisitAll(func (f *pflag.Flag) {
		if f.Changed {
			flags = true
		}
	})
	
	return flags
}
