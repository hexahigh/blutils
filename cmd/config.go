package cmd

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var configDefaults = map[string]map[string]interface{}{
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
	"ascii85": {
		"decode": false,
		"input":  "",
		"output": "",
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
	"hello": {
		"world": map[string]interface{}{
			"foo": map[string]interface{}{
				"bar": 42,
				"baz": true,
			},
			"qux": map[string]interface{}{
				"quux": "hello",
			},
		},
		"corge": map[string]interface{}{
			"grault": 24,
		},
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
	command.Flags().VisitAll(func(flag *pflag.Flag) {
		err := viper.BindPFlag(commandToConfigString(command)+"."+flag.Name, flag)
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

//* I refuse to remove this comment as i will probably need it in the future, you have no idea how much pain this has caused me.
/* func commandToConfigString(c cobra.Command) string {
	if isRootCommand(c) {
		return c.Name()
	}
	if c.Parent() == nil {
		return ""
	}
	return commandToConfigString(*c.Parent()) + "." + c.Name()
} */

func commandToConfigString(c cobra.Command) string {
	log.Infoln("Command:", c.Name())
	configString := c.Name()
	for parent := c.Parent(); parent != nil; parent = parent.Parent() {
		if parent.Name() != "blutils" {
			configString = parent.Name() + "." + configString
		} else {
			break
		}
	}
	return configString
}

func isRootCommand(c cobra.Command) bool {
	return c.Name() == "blutils"
}
