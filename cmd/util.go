package cmd

import (
	_ "embed"
	"strings"
)

//go:embed version
var versionFile string

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

func versionParser(key string) string {
	for _, l := range strings.Split(versionFile, "\n") {
		parts := strings.SplitN(l, ":", 2)
		if len(parts) == 2 && parts[0] == key {
			return strings.TrimSpace(parts[1])
		}
	}
	return ""
}
