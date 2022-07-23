package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/Akimon658/gup/file"
	"github.com/Akimon658/gup/internal/goutil"
	"github.com/Akimon658/gup/internal/print"
	"github.com/Akimon658/gup/internal/slice"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update binaries installed by 'go install'",
	Long: `Update binaries installed by 'go install'

If you execute '$ gup update', gup gets the package path of all commands
under $GOPATH/bin and automatically updates commands to the latest version.`,
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(gup(cmd, args))
	},
}

func init() {
	updateCmd.Flags().BoolP("dry-run", "n", false, "perform the trial update with no changes")
	rootCmd.AddCommand(updateCmd)
}

// gup is main sequence.
// All errors are handled in this function.
func gup(cmd *cobra.Command, args []string) int {
	dryRun, err := cmd.Flags().GetBool("dry-run")
	if err != nil {
		print.Fatal(fmt.Errorf("%s: %w", "can not parse command line argument (--dry-run)", err))
	}

	pkgs, err := getPackageInfo()
	if err != nil {
		print.Fatal(err)
	}
	pkgs = extractUserSpecifyPkg(pkgs, args)

	if len(pkgs) == 0 {
		print.Fatal("unable to update package: no package information")
	}
	return update(pkgs, dryRun)
}

type updateResult struct {
	pkg goutil.Package
	err error
}

func update(pkgs []goutil.Package, dryRun bool) int {
	result := 0
	countFmt := "[%" + pkgDigit(pkgs) + "d/%" + pkgDigit(pkgs) + "d]"
	dryRunManager := goutil.NewGoPaths()

	print.Info("update binary under $GOPATH/bin or $GOBIN")
	signals := make(chan os.Signal, 1)
	if dryRun {
		if err := dryRunManager.StartDryRunMode(); err != nil {
			print.Err(fmt.Errorf("can not change to dry run mode: %w", err))
			return 1
		}
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP,
			syscall.SIGQUIT, syscall.SIGABRT)
		go catchSignal(signals, dryRunManager)
	}

	ch := make(chan updateResult)
	updater := func(p goutil.Package, result chan updateResult) {
		var err error
		if p.ImportPath == "" {
			err = fmt.Errorf(" %s is not installed by 'go install' (or permission incorrect)", p.Name)
		} else {
			if err = goutil.Install(p.ImportPath, p.BuildFlags); err != nil {
				err = fmt.Errorf(" %s %w", p.Name, err)
			}
		}

		p.SetLatestVer()
		r := updateResult{
			pkg: p,
			err: err,
		}
		result <- r
	}

	// update all package
	for _, v := range pkgs {
		go updater(v, ch)
	}

	// print result
	for i := 0; i < len(pkgs); i++ {
		v := <-ch
		if v.err == nil {
			print.Info(fmt.Sprintf(countFmt+" %s (%s)",
				i+1, len(pkgs), v.pkg.ImportPath, v.pkg.CurrentToLatestStr()))
		} else {
			result = 1
			print.Err(fmt.Errorf(countFmt+"%s", i+1, len(pkgs), v.err.Error()))
		}
	}

	if dryRun {
		if err := dryRunManager.EndDryRunMode(); err != nil {
			print.Err(fmt.Errorf("can not change dry run mode to normal mode: %w", err))
			return 1
		}
		close(signals)
	}

	return result
}

func catchSignal(c chan os.Signal, dryRunManager *goutil.GoPaths) {
	for {
		select {
		case <-c:
			if err := dryRunManager.EndDryRunMode(); err != nil {
				print.Err(fmt.Errorf("can not change dry run mode to normal mode: %w", err))
			}
			os.Exit(0)
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func pkgDigit(pkgs []goutil.Package) string {
	return strconv.Itoa(len(strconv.Itoa(len(pkgs))))
}

func getPackageInfo() ([]goutil.Package, error) {
	goBin, err := goutil.GoBin()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "can't find installed binaries", err)
	}

	binList, err := goutil.BinaryPathList(goBin)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "can't get binary-paths installed by 'go install'", err)
	}

	pkgs, err := goutil.GetPackageInformation(binList)
	if err != nil {
		return nil, err
	}

	return pkgs, nil
}

func extractUserSpecifyPkg(pkgs []goutil.Package, targets []string) []goutil.Package {
	result := []goutil.Package{}
	tmp := []string{}
	if len(targets) == 0 {
		return pkgs
	}

	for _, v := range targets {
		v = file.AddExt(v)
		for _, w := range pkgs {
			if v == w.Name {
				result = append(result, w)
				tmp = append(tmp, w.Name)
			}
		}
	}

	if len(tmp) != len(targets) {
		for _, target := range targets {
			if !slice.Contains(tmp, target) {
				print.Warn("not found '" + target + "' package in $GOPATH/bin or $GOBIN")
			}
		}
	}

	return result
}
