package file

import "strings"

func TrimExt(name string) string {
	if isWindows() && hasExt(name) {
		name = strings.TrimSuffix(name, extWin)
	}

	return name
}
