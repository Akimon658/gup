package notify

import (
	"github.com/Akimon658/gup/internal/assets"
	"github.com/Akimon658/gup/internal/print"
	"github.com/gen2brain/beeep"
)

// Info notify information message at desktop
func Info(title, message string) {
	err := beeep.Notify(title, message, assets.InfoIconPath())
	if err != nil {
		print.Warn(err)
	}
}

// Warn notify warning message at desktop
func Warn(title, message string) {
	err := beeep.Notify(title, message, assets.WarningIconPath())
	if err != nil {
		print.Warn(err)
	}
}
