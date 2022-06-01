package cmd

import (
	"bufio"
	"errors"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/Akimon658/gup/internal/goutil"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import commands from a file",
	Long: `Import commands from a file.
Use "gup export" to generate a file`,
	Run: func(cmd *cobra.Command, args []string) {
		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			log.Fatal(err)
		}
		if err := runImport(args[0], dryRun); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	importCmd.Flags().BoolP("dry-run", "n", false, "perform the trial update with no changes")
	rootCmd.AddCommand(importCmd)
}

func runImport(path string, dryRun bool) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var pkgs []goutil.Package
	for scanner.Scan() {
		pkg := goutil.Package{
			ImportPath: scanner.Text(),
			Version:    &goutil.Version{Current: "<from import>"},
		}
		pkgs = append(pkgs, pkg)
	}

	if len(pkgs) == 0 {
		return errors.New("given file is empty")
	}

	if code := update(pkgs, dryRun); code != 0 {
		return errors.New("failed to import packages")
	}

	return nil
}
