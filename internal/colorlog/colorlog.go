package colorlog

import (
	"fmt"
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

// Complete logs a success message in green and exits with 0 in CLI mode; no exit in test mode.
func Complete(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	fullMsg := ANSIColourGreen + " " + msg + ANSIReset
	if os.Getenv("TESTING") == "" {
		log.Print(fullMsg)
		os.Exit(0)
	}
	log.Print(msg) // Just log in test mode
}
