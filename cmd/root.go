package cmd

import (
	_ "embed"
	"os"
	"path"
	"strings"

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

func init() {
	initConfig()
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	defer viper.WriteConfig()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func getDefaultConfigDir() string {
	var dir string
	dir, err := os.UserConfigDir()
	dir = path.Join(dir, "blutils")
	if err != nil {
		log.Warnf("Could not get user config dir, using PWD. Cause: %v", err)
		dir, _ = os.Getwd()
	}
	return dir
}

func initConfig() {
	configLoadDefaults()
	//* Load flags
	rootCmd.PersistentFlags().Int8P("verbosity", "v", 4, "verbosity level. 0=panic, 1=fatal, 2=error, 3=warn, 4=info, 5=debug, 6=trace")
	rootCmd.PersistentFlags().StringP("config-dir", "D", getDefaultConfigDir(), "Directory containing data and config files")
	rootCmd.PersistentFlags().StringP("config-file", "C", "config.toml", "Name of the config file, with extension")
	rootCmd.PersistentFlags().String("output", "text", "Output mode. Supported values are: text, json")
	rootCmd.PersistentFlags().Bool("forceColors", false, "Force colors in log output. Only works if output mode is Text")
	rootCmd.PersistentFlags().Bool("disableColors", false, "Disable colors in log output. Only works if output mode is Text")
	rootCmd.PersistentFlags().Bool("disableLevelTruncation", false, "When colors are enabled, levels are truncated to 4 characters by default, this disables that. Only works if output mode is Text")
	rootCmd.PersistentFlags().Bool("padLevelText", true, "Pads level text for better readability. Only works if output mode is Text")
	rootCmd.PersistentFlags().Bool("fullTimestamp", false, "Show full timestamp in log output. Only works if output mode is Text")
	rootCmd.PersistentFlags().Bool("disableTimestamp", false, "Disable timestamp in log output. Works for both text and json output")
	rootCmd.PersistentFlags().Bool("prettyPrint", false, "Indent json output. Not sure why you would want this, but here you go.")
	rootCmd.ParseFlags(os.Args[1:])

	// viper.BindPFlags(rootCmd.PersistentFlags())
	configBindFlags(*rootCmd)

	//* Load env vars
	viper.SetEnvPrefix("BLUTILS")
	viper.AutomaticEnv()

	//* Set config file options
	viper.SetConfigName(viper.GetString("blutils.config-file"))
	viper.SetConfigType("toml")
	viper.AddConfigPath(viper.GetString("blutils.config-dir"))

	//* Load config file
	err := os.MkdirAll(viper.GetString("blutils.config-dir"), 0700)
	if err != nil {
		log.Warnf("Could not create config dir. Cause: %v", err)
	}
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Warnf("Could not find config file, making one. Cause: %v", err)
			// err = viper.SafeWriteConfigAs(path.Join(viper.GetString("blutils.config-dir"), viper.GetString("blutils.config-file")))
			err = writeDefaultsAs(path.Join(viper.GetString("blutils.config-dir"), viper.GetString("blutils.config-file")))
			if err != nil {
				log.Warnf("Could not write config file. Cause: %v", err)
			}
		} else {
			log.Errorf("Could not read config file. Cause: %v", err)
		}
	}

	//* Unmarshal config
	viper.Unmarshal(&Params)

	log.SetLevel(log.Level(viper.GetInt("blutils.verbosity")))
	switch strings.ToLower(viper.GetString("blutils.output")) {
	case "text":
		log.SetFormatter(&log.TextFormatter{
			ForceColors:            viper.GetBool("blutils.forceColors"),
			DisableColors:          viper.GetBool("blutils.disableColors"),
			DisableLevelTruncation: viper.GetBool("blutils.disableLevelTruncation"),
			PadLevelText:           viper.GetBool("blutils.padLevelText"),
			FullTimestamp:          viper.GetBool("blutils.fullTimestamp"),
			DisableTimestamp:       viper.GetBool("blutils.disableTimestamp"),
		})
	case "json":
		log.SetFormatter(&log.JSONFormatter{
			DisableTimestamp: viper.GetBool("blutils.disableTimestamp"),
			PrettyPrint:      viper.GetBool("blutils.prettyPrint"),
		})
	default:
		log.Fatalf("Unknown output mode: %s", viper.GetString("blutils.output"))
	}

	for key, value := range viper.GetViper().AllSettings() {
		log.WithFields(log.Fields{
			key: value,
		}).Debug("Command Flag")
	}

	defer viper.WriteConfig()
}
