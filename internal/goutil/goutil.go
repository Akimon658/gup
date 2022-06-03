package goutil

import (
	"errors"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/adrg/xdg"
	"github.com/fatih/color"

	"github.com/Akimon658/gup/config"
	"github.com/Akimon658/gup/internal/print"
)

type GoPaths struct {
	GOBIN  string
	GOPATH string
	// TmpPath is tmporary path for dry run
	TmpPath string
}

type Package struct {
	Name       string
	ImportPath string
	// ModulePath is path where go.mod is stored
	ModulePath string
	Version    *Version
	BuildFlags *config.BuildFlags
}

type Version struct {
	Current string
	Latest  string
}

func NewVersion() *Version {
	return &Version{
		Current: "",
		Latest:  "",
	}
}

func (p *Package) SetCurrentVer() {
	p.Version.Current = GetPackageVersion(p.Name)
}

func (p *Package) SetLatestVer() {
	p.Version.Latest = GetPackageVersion(p.Name)
}

// CurrentToLatestStr returns string about the current version and the latest version
func (p *Package) CurrentToLatestStr() string {
	if IsAlreadyUpToDate(*p.Version) {
		return "Already up-to-date: " + color.GreenString(p.Version.Latest)
	}
	return color.GreenString(p.Version.Current) + " to " + color.YellowString(p.Version.Latest)
}

// VersionCheckResultStr returns string about command version check.
func (p *Package) VersionCheckResultStr() string {
	if IsAlreadyUpToDate(*p.Version) {
		return "Already up-to-date: " + color.GreenString(p.Version.Latest)
	}
	return "current: " + color.GreenString(p.Version.Current) + ", latest: " + color.YellowString(p.Version.Latest)
}

func IsAlreadyUpToDate(ver Version) bool {
	return ver.Current == ver.Latest
}

func NewGoPaths() *GoPaths {
	return &GoPaths{
		GOBIN:  goBin(),
		GOPATH: goPath(),
	}
}

// StartDryRunMode changes the GOBIN or GOPATH settings to install the binaries in the temporary directory.
func (gp *GoPaths) StartDryRunMode() error {
	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}

	if gp.GOBIN != "" {
		if err := os.Setenv("GOBIN", tmpDir); err != nil {
			return err
		}
	} else if gp.GOPATH != "" {
		if err := os.Setenv("GOPATH", tmpDir); err != nil {
			return err
		}
	} else {
		return errors.New("$GOPATH and $GOBIN is not set")
	}
	return nil
}

func (gp *GoPaths) EndDryRunMode() error {
	if gp.GOBIN != "" {
		if err := os.Setenv("GOBIN", gp.GOBIN); err != nil {
			return err
		}
	} else if gp.GOPATH != "" {
		if err := os.Setenv("GOPATH", gp.GOPATH); err != nil {
			return err
		}
	} else {
		return errors.New("$GOPATH and $GOBIN is not set")
	}

	if err := gp.removeTmpDir(); err != nil {
		return fmt.Errorf("%s: %w", "temporary directory for dry run remains", err)
	}
	return nil
}

func (gp *GoPaths) removeTmpDir() error {
	if gp.TmpPath != "" {
		if err := os.RemoveAll(gp.TmpPath); err != nil {
			return err
		}
	}
	return nil
}

func Install(importPath string, buildFlags *config.BuildFlags) error {
	if importPath == "command-line-arguments" {
		return errors.New("is devel-binary copied from local environment")
	}

	if err := exec.Command("go", "install", "-ldflags", buildFlags.Ldflags, "-tags", buildFlags.Tags, importPath+"@latest").Run(); err != nil {
		return fmt.Errorf("can't install %s: %w", importPath, err)
	}

	return nil
}

// GetLatestVer executes "go list -m -f {{.Version}} <importPath>@latest"
func GetLatestVer(modulePath string) (string, error) {
	out, err := exec.Command("go", "list", "-m", "-f", "{{.Version}}", modulePath+"@latest").Output()
	if err != nil {
		return "", errors.New("can't check " + modulePath)
	}
	return strings.TrimRight(string(out), "\n"), nil
}

func goPath() string {
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		return gopath
	}
	return build.Default.GOPATH
}

func goBin() string {
	return os.Getenv("GOBIN")
}

func GoBin() (string, error) {
	goBin := goBin()
	if goBin != "" {
		return goBin, nil
	}

	goPath := goPath()
	if goPath == "" {
		return "", errors.New("$GOPATH is not set")
	}
	return filepath.Join(goPath, "bin"), nil
}

// GoVersionWithOptionM returns result of "go version -m"
func GoVersionWithOptionM(bin string) ([]string, error) {
	out, err := exec.Command("go", "version", "-m", bin).Output()
	if err != nil {
		return nil, err
	}
	return strings.Split(string(out), "\n"), nil
}

func BinaryPathList(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	list := []string{}
	for _, e := range entries {
		if !e.IsDir() {
			list = append(list, filepath.Join(path, e.Name()))
		}
	}
	return list, nil
}

func GetPackageInformation(binList []string) ([]Package, error) {
	var pkgs []Package
	conf, err := config.Read(filepath.Join(xdg.ConfigHome, "gup", "package.yml"))
	if err != nil {
		return nil, err
	}

	for _, v := range binList {
		out, err := GoVersionWithOptionM(v)
		if err != nil {
			print.Warn(fmt.Errorf("%s: %w", "can not get package path", err))
			continue
		}

		name := filepath.Base(v)

		pkg := Package{
			Name:       name,
			ImportPath: extractImportPath(out),
			ModulePath: extractModulePath(out),
			Version:    NewVersion(),
			BuildFlags: conf.GetFlags(name),
		}
		pkg.SetCurrentVer()
		pkgs = append(pkgs, pkg)
	}

	return pkgs, nil
}

func GetPackageVersion(cmdName string) string {
	goBin, err := GoBin()
	if err != nil {
		return "unknown"
	}

	out, err := GoVersionWithOptionM(filepath.Join(goBin, cmdName))
	if err != nil {
		return "unknown"
	}

	for _, v := range out {
		vv := strings.TrimSpace(v)
		if len(v) != len(vv) && strings.HasPrefix(vv, "mod") {
			//         mod     github.com/nao1215/subaru       v1.0.2  h1:LU9/1bFyqef3re6FVSFgTFMSXCZvrmDpmX3KQtlHzXA=
			v = strings.TrimLeft(vv, "mod")
			v = strings.TrimSpace(v)

			//github.com/nao1215/subaru       v1.0.2  h1:LU9/1bFyqef3re6FVSFgTFMSXCZvrmDpmX3KQtlHzXA=
			r := regexp.MustCompile(`^[^\s]+(\s)`)
			v = r.ReplaceAllString(v, "")

			// v1.0.2  h1:LU9/1bFyqef3re6FVSFgTFMSXCZvrmDpmX3KQtlHzXA=
			r = regexp.MustCompile(`(\s)[^\s]+$`)

			// v1.0.2
			return r.ReplaceAllString(v, "")
		}
	}
	return "unknown"
}

func extractImportPath(lines []string) string {
	for _, v := range lines {
		vv := strings.TrimSpace(v)
		if len(v) != len(vv) && strings.HasPrefix(vv, "path") {
			vv = strings.TrimLeft(vv, "path")
			vv = strings.TrimSpace(vv)
			return strings.TrimRight(vv, "\n")
		}
	}
	return ""
}

func extractModulePath(lines []string) string {
	for _, v := range lines {
		vv := strings.TrimSpace(v)
		if len(v) != len(vv) && strings.HasPrefix(vv, "mod") {
			//         mod     github.com/nao1215/subaru       v1.0.2  h1:LU9/1bFyqef3re6FVSFgTFMSXCZvrmDpmX3KQtlHzXA=
			v = strings.TrimLeft(vv, "mod")
			v = strings.TrimSpace(v)

			//github.com/nao1215/subaru       v1.0.2  h1:LU9/1bFyqef3re6FVSFgTFMSXCZvrmDpmX3KQtlHzXA=
			r := regexp.MustCompile(`(\s).*$`)
			return r.ReplaceAllString(v, "")
		}
	}
	return ""
}
