package cmd

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var configDefaults = map[string]map[string]any{
	"blutils": {
		"verbosity":              4,
		"config-dir":             getDefaultConfigDir(),
		"config-file":            "config.toml",
		"output":                 "text",
		"forceColors":            false,
		"disableColors":          false,
		"disableLevelTruncation": false,
		"padLevelText":           true,
		"fullTimestamp":          false,
		"disableTimestamp":       false,
		"prettyPrint":            false,
	},
	"bitflip": {
		"bits":        0,
		"percentage":  0,
		"min-offset":  0,
		"chunk":       1,
		"no-progress": false,
		"extreme":     false,
	},
	"bench": {
		"cpu":     0,
		"timeout": 10,
	},
	"update": {
		"repo": "hexahigh/blutils",
		"tag":  "",
		"temp": filepath.Join(os.TempDir(), "blutils-build"),
	},
}

func configLoadDefaults() {
	for k, v := range configDefaults {
		for k2, v2 := range v {
			viper.SetDefault(k+"."+k2, v2)
		}
	}
}

func configBindFlags(command cobra.Command) {
	command.LocalFlags().VisitAll(func(flag *pflag.Flag) {
		err := viper.BindPFlag(command.Name()+"."+flag.Name, flag)
		if err != nil {
			log.Fatalf("Error initializing viper: %v", err)
		}
	})
	command.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		err := viper.BindPFlag(command.Name()+"."+flag.Name, flag)
		if err != nil {
			log.Fatalf("Error initializing viper: %v", err)
		}
	})
}

func writeDefaults() error {
	return writeDefaultsAs(viper.ConfigFileUsed())
}

func writeDefaultsAs(path string) error {
	newViper := viper.New()
	for k, v := range configDefaults {
		for k2, v2 := range v {
			newViper.SetDefault(k+"."+k2, v2)
		}
	}
	err := newViper.WriteConfigAs(path)
	return err
}

func getDefault(key string) any {
	return viper.Get(key)
}
