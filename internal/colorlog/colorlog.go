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
	ANSIColourCyan   = "\033[36m" // Info (Cyan)
	ANSIReset        = "\033[0m"  // Reset color

	CLIPREFIX = ANSIColourCyan + "syncsnipe →" + ANSIReset + " "
	NEWLINE   = "\n"
)

// Generic log formatter to avoid repetition
func logWithColor(color, format string, v ...any) {
	log.Printf(color+" "+format+ANSIReset, v...)
}

func Info(format string, v ...any)    { logWithColor(ANSIColourCyan, format, v...) }
func Success(format string, v ...any) { logWithColor(ANSIColourGreen, format, v...) }
func Warn(format string, v ...any)    { logWithColor(ANSIColourYellow, format, v...) }
func Error(format string, v ...any)   { logWithColor(ANSIColourRed, format, v...) }

func CLISuccess(format string, v ...any) {
	fmt.Fprintf(os.Stdout, CLIPREFIX+ANSIColourGreen+format+ANSIReset+NEWLINE, v...)
}

func CLIError(format string, v ...any) {
	fmt.Fprintf(os.Stdout, CLIPREFIX+ANSIColourRed+format+ANSIReset+NEWLINE, v...)
}
