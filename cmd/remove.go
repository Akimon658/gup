package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Akimon658/gup/file"
	"github.com/Akimon658/gup/internal/goutil"
	"github.com/Akimon658/gup/internal/print"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove binaries under $GOPATH/bin or $GOBIN",
	Long: `Remove commands in $GOPATH/bin or $GOBIN.
If you want to specify multiple binaries at once, separate them with space.
[e.g.] gup remove a_cmd b_cmd c_cmd`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			print.Fatal(err)
		}

		os.Exit(remove(args, force))
	},
}

func init() {
	removeCmd.Flags().BoolP("force", "f", false, "Forcibly remove the file")
	rootCmd.AddCommand(removeCmd)
}

func remove(args []string, force bool) int {
	gobin, err := goutil.GoBin()
	if err != nil {
		print.Fatal(err)
	}

	code := 0
	for _, v := range args {
		v = file.AddExt(v)

		target := filepath.Join(gobin, v)
		stat, err := os.Stat(target)
		if err != nil {
			print.Err(err)
			code = 1
			continue
		}

		if stat.IsDir() {
			print.Err(fmt.Errorf("no such file: %s", target))
			code = 1
			continue
		}

		if !force {
			if !print.Question(fmt.Sprintf("Do you want to remove %s?", target)) {
				fmt.Println("Removal cancelled")
				continue
			}
		}

		if err := os.Remove(target); err != nil {
			print.Err(err)
			code = 1
			continue
		}

		fmt.Println("Removed " + target)
	}
	return code
}
