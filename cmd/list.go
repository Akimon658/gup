package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/Akimon658/gup/internal/goutil"
	"github.com/Akimon658/gup/internal/print"
	ls "github.com/Akimon658/gup/list"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed commands",
	Long:  "List informations of installed commands",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		listPackagePaths, err := cmd.Flags().GetBool("package-path")
		if err != nil {
			log.Fatal(err)
		}

		if err := list(listPackagePaths); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	listCmd.Flags().Bool("package-path", false, "print package paths only")
	rootCmd.AddCommand(listCmd)
}

func list(listPackagePaths bool) error {
	pkgs, err := getPackageInfo()
	if err != nil {
		return err
	}

	if len(pkgs) == 0 {
		print.Fatal("unable to list up package: no package information")
	}

	if listPackagePaths {
		fmt.Println(ls.PackagePaths(pkgs))
	} else {
		printPackageList(pkgs)
	}

	return nil
}

// PackageList list up command package in $GOPATH/bin or $GOBIN
func printPackageList(pkgs []goutil.Package) {
	max := 0
	for _, v := range pkgs {
		if len(v.Name) > max {
			max = len(v.Name)
		}
	}

	for _, v := range pkgs {
		fmt.Fprintf(print.Stdout, "%"+strconv.Itoa(max)+"s: %s%s\n",
			v.Name,
			v.ImportPath,
			color.GreenString("@"+goutil.GetPackageVersion(v.Name)))
	}
}
