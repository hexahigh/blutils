package cmd

import (
	_ "embed"
	"io"
	"strings"

	color "github.com/hexahigh/go-lib/ansicolor"
)

//go:embed version
var versionFile string

var verbosityMap = map[int]string{0: "ERROR", 1: "WARN", 2: "INFO", 3: "DEBUG"}

func initColor() {
	if !*Params.NoColor {
		red := color.Red
		yellow := color.Yellow
		green := color.Green
		blue := color.Purple

		if color.SupportsTrueColor() || *Params.TrueColor {
			VerbosePrintln(3, "Terminal supports full color")
			red = color.Red24bit
			yellow = color.Yellow24bit
			green = color.Green24bit
			blue = color.Purple24bit
		}

		verbosityMap = map[int]string{0: red + "ERROR" + color.Reset, 1: yellow + "WARN" + color.Reset, 2: green + "INFO" + color.Reset, 3: blue + "DEBUG" + color.Reset}
	}
}

func VerbosePrintln(minLevel int, msg ...any) {
	if *Params.Verbosity >= minLevel {
		msg = append([]any{"[" + verbosityMap[minLevel] + "]"}, msg...)
		logger.Println(msg...)
	}
}

func VerbosePrintf(minLevel int, format string, msg ...any) {
	if *Params.Verbosity >= minLevel {
		msg = append([]any{"[" + verbosityMap[minLevel] + "]"}, msg...)
		logger.Printf(format, msg...)
	}
}

// Prints your message without any log level stuff
func VerbosePrintlnC(minLevel int, msg ...any) {
	if *Params.Verbosity >= minLevel {
		logger.Println(msg...)
	}
}

// Uses io.Writer to print your message
func VerbosePrintlnW(minLevel int) io.Writer {
	if *Params.Verbosity >= minLevel {
		return logger.Writer()
	}

	return io.Discard
}

func VersionParser(key string) string {
	for _, l := range strings.Split(versionFile, "\n") {
		parts := strings.SplitN(l, "=", 2)
		if len(parts) == 2 && parts[0] == key {
			return strings.TrimSpace(parts[1])
		}
	}
	return ""
}
