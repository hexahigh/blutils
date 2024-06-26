package verbprint

import (
	"io"
	"log"

	color "github.com/hexahigh/go-lib/ansicolor"
)

var defaultVerbosityMap = map[int]string{0: "ERROR", 1: "WARN", 2: "INFO", 3: "DEBUG"}

type VerboseLogger struct {
	verbosityLevel int
	logger         *log.Logger
	color          int // -1 = no color, 0 = auto, 1 = force, 2 = force true color
	verbosityMap   map[int]string
}

// New returns a new VerboseLogger. You will need to call InitColor() in order to use colors.
func New(verbosityLevel int, logger *log.Logger, color int) *VerboseLogger {
	if logger == nil {
		panic("No logger provided to verbose logger")
	}
	return &VerboseLogger{
		verbosityLevel: verbosityLevel,
		logger:         logger,
		color:          color,
		verbosityMap:   defaultVerbosityMap,
	}
}

func (v *VerboseLogger) InitColor() {
	if v.color >= 0 && color.Supports256Colors() || v.color >= 1 {
		red := color.Red
		yellow := color.Yellow
		green := color.Green
		purple := color.Purple

		if color.SupportsTrueColor() || v.color >= 2 {
			v.Println(3, "Terminal supports full color")
			red = color.Red24bit
			yellow = color.Yellow24bit
			green = color.Green24bit
			purple = color.Purple24bit
		}

		v.verbosityMap = map[int]string{0: red + "ERROR" + color.Reset, 1: yellow + "WARN" + color.Reset, 2: green + "INFO" + color.Reset, 3: purple + "DEBUG" + color.Reset}
	}
}

// If verbosityLevel is >= minLevel, the message will be printed. Appends the level to the message with optional coloring. All other arguments are handled in the manner of fmt.Println
func (v *VerboseLogger) Println(minLevel int, msg ...any) {
	if v.verbosityLevel >= minLevel {
		msg = append([]any{"[" + v.verbosityMap[minLevel] + "]"}, msg...)
		v.logger.Println(msg...)
	}
}

// If verbosityLevel is >= minLevel, the message will be printed. Appends the level to the message with optional coloring. All other arguments are handled in the manner of fmt.Printf
func (v *VerboseLogger) Printf(minLevel int, format string, msg ...any) {
	if v.verbosityLevel >= minLevel {
		msg = append([]any{"[" + v.verbosityMap[minLevel] + "]"}, msg...)
		v.logger.Printf(format, msg...)
	}
}

// If verbosityLevel is >= minLevel, the message will be printed. All other arguments are handled in the manner of fmt.Println
func (v *VerboseLogger) PrintlnC(minLevel int, msg ...any) {
	if v.verbosityLevel >= minLevel {
		v.Println(minLevel, msg...)
	}
}

func (v *VerboseLogger) PrintW(minLevel int) io.Writer {
	if v.verbosityLevel >= minLevel {
		return v.logger.Writer()
	}

	return io.Discard
}
