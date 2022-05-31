package list

import (
	"strings"

	"github.com/Akimon658/gup/internal/goutil"
)

func ImportPaths(pkgs []goutil.Package) string {
	var sb strings.Builder

	for i := range pkgs {
		sb.WriteString(pkgs[i].ImportPath)
		sb.WriteString("\n")
	}

	return sb.String()
}
