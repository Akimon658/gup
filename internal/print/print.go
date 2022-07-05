package print

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"

	"github.com/Akimon658/gup/internal/cmdinfo"
)

var (
	// Stdout is new instance of Writer which handles escape sequence for stdout.
	Stdout = colorable.NewColorableStdout()
	// Stderr is new instance of Writer which handles escape sequence for stderr.
	Stderr = colorable.NewColorableStderr()
)

// Info print information message at STDOUT.
func Info(msg string) {
	fmt.Fprintf(Stdout, "%s:%s: %s\n",
		cmdinfo.Name, color.GreenString("INFO "), msg)
}

// Warn print warning message at STDERR.
func Warn(err interface{}) {
	fmt.Fprintf(Stderr, "%s:%s: %v\n",
		cmdinfo.Name, color.YellowString("WARN "), err)
}

// Err print error message at STDERR.
func Err(err interface{}) {
	fmt.Fprintf(Stderr, "%s:%s: %v\n",
		cmdinfo.Name, color.HiYellowString("ERROR"), err)
}

// Fatal print dying message at STDERR.
func Fatal(err interface{}) {
	fmt.Fprintf(Stderr, "%s:%s: %v\n",
		cmdinfo.Name, color.RedString("FATAL"), err)
	os.Exit(1)
}

// Question displays the question in the terminal and receives an answer from the user.
func Question(ask string) bool {
	fmt.Print(ask + " [Y/n] ")

	var input string
	_, err := fmt.Scanln(&input)
	if err != nil && !strings.Contains(err.Error(), "expected newline") {
		fmt.Fprint(os.Stderr, err.Error())
		return false
	}

	switch strings.ToLower(input) {
	case "y", "yes", "":
		return true
	case "n", "no":
		return false
	default:
		fmt.Println("Prease enter y or n")
		return Question(ask)
	}
}
