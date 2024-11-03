package cmd

import (
	_ "embed"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	Verbosity  *int8
	ConfigDir  *string
	ConfigFile *string
}

var Viper *viper.Viper

func init() {
	rootCmd.PersistentFlags().Int8P("verbosity", "v", 4, "verbosity level. 0=panic, 1=fatal, 2=error, 3=warn, 4=info, 5=debug, 6=trace")
	rootCmd.PersistentFlags().StringP("config-dir", "d", getDefaultConfigDir(), "Directory containing data and config files")
	rootCmd.PersistentFlags().StringP("config-file", "c", "config.toml", "Name of the config file, with extension")
	rootCmd.ParseFlags(os.Args[1:])

	err := os.MkdirAll(viper.GetString("config-dir"), 0700)
	if err != nil {
		log.Warnf("Could not create config dir. Cause: %v", err)
	}
	viper.SetConfigName(viper.GetString("config-file"))
	viper.SetConfigType("toml")
	viper.AddConfigPath(viper.GetString("config-dir"))
	viper.BindPFlags(rootCmd.PersistentFlags())
	viper.SetEnvPrefix("BLUTILS")
	viper.AutomaticEnv()

	log.SetLevel(log.Level(*Params.Verbosity))
	log.SetFormatter(&log.TextFormatter{})

	for key, value := range viper.GetViper().AllSettings() {
		log.WithFields(log.Fields{
			key: value,
		}).Debug("Command Flag")
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

func getDefaultConfigDir() string {
	var dir string
	dir, err := os.UserConfigDir()
	dir = path.Join(dir, "go-vt")
	if err != nil {
		log.Warnf("Could not get user config dir, using PWD. Cause: %v", err)
		dir, _ = os.Getwd()
	}
	return dir
}
