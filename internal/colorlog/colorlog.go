package colorlog

import (
	"fmt"
	"log"
)

const (
	ANSIColourGreen  = "\033[32m" // Success (Green)
	ANSIColourYellow = "\033[33m" // Warning (Yellow)
	ANSIColourRed    = "\033[31m" // Error (Red)
	ANSIColourCyan   = "\x1b[36m" // Info (Cyan)
	ANSIReset        = "\033[0m"  // Reset color

	CLIPREFIX = ANSIColourCyan + "syncsnipe →" + ANSIReset + " "
	NEWLINE   = "\n"
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

func CLISuccess(format string, v ...any) {
	fmt.Printf(CLIPREFIX+ANSIColourGreen+format+ANSIReset+NEWLINE, v...)
}

func CLIError(format string, v ...any) {
	fmt.Printf(CLIPREFIX+ANSIColourRed+format+ANSIReset+NEWLINE, v...)
}
