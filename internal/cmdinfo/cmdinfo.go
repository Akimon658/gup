package cmdinfo

import (
	"fmt"
	"runtime/debug"
)

const Name = "gup"

// Version return gup command version.
func Version() string {
	version := "unknown"

	info, ok := debug.ReadBuildInfo()
	if ok {
		version = info.Main.Version
	}

	return fmt.Sprintf("%s version %s (under Apache License version 2.0)", Name, version)
}
