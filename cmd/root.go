/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	_ "embed"
	"log"
	"os"

	color "github.com/hexahigh/go-lib/ansicolor"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "blutils",
	Short: "Utility program",
	Long:  `Utility program`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

//go:embed version
var version string

var rootParams RootParams
var logger *log.Logger

var verbosityMap = map[int]string{0: "ERROR", 1: "WARN", 2: "INFO", 3: "DEBUG"}

type RootParams struct {
	Verbosity *int
	NoColor   *bool
	TrueColor *bool
}

func init() {
	rootParams.Verbosity = rootCmd.PersistentFlags().IntP("verbosity", "v", 2, "Verbosity level (0-3)")
	rootParams.NoColor = rootCmd.PersistentFlags().Bool("no-color", false, "Disable color output in log")
	rootParams.TrueColor = rootCmd.PersistentFlags().Bool("true-color", false, "Force true color output in log")
	rootCmd.ParseFlags(os.Args[1:])

	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)

	if !*rootParams.NoColor {
		red := color.Red
		yellow := color.Yellow
		green := color.Green
		blue := color.Purple

		if color.SupportsTrueColor() || *rootParams.TrueColor {
			verbosePrintln(3, "Terminal supports full color")
			red = color.Red24bit
			yellow = color.Yellow24bit
			green = color.Green24bit
			blue = color.Purple24bit
		}

		verbosityMap = map[int]string{0: red + "ERROR" + color.Reset, 1: yellow + "WARN" + color.Reset, 2: green + "INFO" + color.Reset, 3: blue + "DEBUG" + color.Reset}
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

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
