package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export package paths to stdout",
	Long: `Export package paths to stdout. (alias of "gup list --import-path")
To save as a file, run "gup export > path/to/file"`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := list(true); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
