package cmd

import (
	_ "embed"
	"log"
	"os"

	"github.com/hexahigh/blutils/lib/verbprint"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "blutils",
	Short: "Utility program",
	Long:  `Utility program`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

var Params params
var logLogger *log.Logger
var Logger *verbprint.VerboseLogger

type params struct {
	Verbosity *int
	NoColor   *bool
	TrueColor *bool
}

func init() {
	Params.Verbosity = RootCmd.PersistentFlags().IntP("verbosity", "v", 2, "Verbosity level (0-3)")
	Params.NoColor = RootCmd.PersistentFlags().Bool("no-color", false, "Disable color output in log")
	Params.TrueColor = RootCmd.PersistentFlags().Bool("true-color", false, "Force true color output in log")
	RootCmd.ParseFlags(os.Args[1:])

	logLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime)

	colorNum := 0

	if *Params.NoColor {
		colorNum = -1
	} else if *Params.TrueColor {
		colorNum = 2
	}

	Logger = verbprint.New(*Params.Verbosity, logLogger, colorNum)
	Logger.InitColor()

}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
