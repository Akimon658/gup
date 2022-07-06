package file

import (
	"runtime"
	"strings"
)

const extWin = ".exe"

func Extension(name string) string {
	if runtime.GOOS == "windows" && !strings.HasSuffix(name, extWin) {
		name += extWin
	}

	return name
}
