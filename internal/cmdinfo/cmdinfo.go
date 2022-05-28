package cmdinfo

import (
	"fmt"

	"github.com/harakeishi/curver"
)

const Name = "gup"

// Version return gup command version.
func Version() string {
	return fmt.Sprintf("%s %s (under Apache License version 2.0)", Name, curver.GetVersion())
}
