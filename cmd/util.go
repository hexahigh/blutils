package cmd

import (
	_ "embed"
)

//go:embed version
var version string

func verbosePrintln(minLevel int, msg ...any) {
	if *rootParams.Verbosity >= minLevel {
		msg = append([]any{"[" + verbosityMap[minLevel] + "]"}, msg...)
		logger.Println(msg...)
	}
}

func verbosePrintf(minLevel int, format string, msg ...any) {
	if *rootParams.Verbosity >= minLevel {
		msg = append([]any{"[" + verbosityMap[minLevel] + "]"}, msg...)
		logger.Printf(format, msg...)
	}
}

// Prints your message without any log level stuff
func verbosePrintlnC(minLevel int, msg ...any) {
	if *rootParams.Verbosity >= minLevel {
		logger.Println(msg...)
	}
}
