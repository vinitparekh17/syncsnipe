package colorlog

import (
	"log"
	"os"
)

const (
	ANSIColourGreen  = "\033[32m" // Success (Green)
	ANSIColourYellow = "\033[33m" // Warning (Yellow)
	ANSIColourRed    = "\033[31m" // Error (Red)
	ANSIColourCyan   = "\x1b[36m" // Info (Cyan)
	ANSIReset        = "\033[0m"  // Reset color
)

func Info(format string, v ...any) {
	log.Printf(ANSIColourCyan+" "+format+ANSIReset, v...)
}

func Success(format string, v ...any) {
	log.Printf(ANSIColourGreen+" "+format+ANSIReset, v...)
}

func Warn(format string, v ...any) {
	log.Printf(ANSIColourYellow+" "+format+ANSIReset, v...)
}

func Error(format string, v ...any) {
	log.Printf(ANSIColourRed+" "+format+ANSIReset, v...)
}

// Fatal logs the message and exits the program with status code 1
// Intended use of this function is limited to the main package and CLI commands to show user feedback on error
// avoid using this function in internal packages; just return the error as much as possible
func Fatal(format string, v ...any) {
	log.Fatalf(ANSIColourRed+" "+format+ANSIReset, v...)
}

// Complete logs the message and exits the program with status code 0
// Intended use of this function is limited to cli commands to show user feedback on successful completion
func Complete(format string, v ...any) {
	Success(format, v...)
	os.Exit(0)
}
