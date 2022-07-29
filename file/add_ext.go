package file

import (
	"runtime"
	"strings"
)

const extWin = ".exe"

func AddExt(name string) string {
	if isWindows() && !hasExt(name) {
		name += extWin
	}

	return name
}

func hasExt(name string) bool {
	return strings.HasSuffix(name, extWin)
}

func isWindows() bool {
	return runtime.GOOS == "windows"
}
