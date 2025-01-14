package list

import (
	"strings"

	"github.com/Akimon658/gup/internal/goutil"
)

func PackagePaths(pkgs []goutil.Package) string {
	var sb strings.Builder

	for i := range pkgs {
		sb.WriteString(pkgs[i].ImportPath)
		sb.WriteString("\n")
	}

	return strings.TrimSuffix(sb.String(), "\n")
}
