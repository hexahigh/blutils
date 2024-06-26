package cmd

import (
	_ "embed"
	"strings"
)

//go:embed version
var versionFile string

func VersionParser(key string) string {
	for _, l := range strings.Split(versionFile, "\n") {
		parts := strings.SplitN(l, "=", 2)
		if len(parts) == 2 && parts[0] == key {
			return strings.TrimSpace(parts[1])
		}
	}
	return ""
}
