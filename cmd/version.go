package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Akimon658/gup/internal/cmdinfo"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show " + cmdinfo.Name + " command version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmdinfo.Version())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
