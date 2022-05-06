package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Akimon658/gup/internal/print"
)

var rootCmd = &cobra.Command{
	Use: "gup",
	Short: `gup command update binaries installed by 'go install'.
If you update all binaries, just run '$ gup update'`,
}

// Execute run gup process.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		print.Err(err)
	}
}
