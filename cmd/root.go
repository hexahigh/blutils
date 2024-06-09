package cmd

import (
	_ "embed"
	"log"
	"os"

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

var rootParams RootParams
var logger *log.Logger

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

	initColor()

}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
