package cmd

import (
	_ "embed"
	"os"

	"github.com/sirupsen/logrus"
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

var Params struct {
	Verbosity *int8
}

var Log *logrus.Logger

func init() {
	Params.Verbosity = rootCmd.PersistentFlags().Int8P("verbosity", "v", 4, "verbosity level. 0=panic, 1=fatal, 2=error, 3=warn, 4=info, 5=debug, 6=trace")
	rootCmd.ParseFlags(os.Args[1:])

	Log = logrus.New()
	Log.SetLevel(logrus.Level(*Params.Verbosity))
	Log.SetFormatter(&logrus.TextFormatter{})
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
