package cmd

import (
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/sideband"

	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

type UpdateParams struct {
	Repo    *string
	Tag     *string
	TempDir *string
}

var updateParams UpdateParams

func init() {
	rootCmd.AddCommand(updateCmd)

	updateParams.Repo = updateCmd.Flags().StringP("repo", "r", "hexahigh/blutils", "Repository to update from, must be hosted on github.com")
	updateParams.Tag = updateCmd.Flags().StringP("tag", "t", "", "Force update to a specific tag")
	updateParams.TempDir = updateCmd.Flags().StringP("temp", "T", filepath.Join(os.TempDir(), "blutils-build"), "Temporary directory")
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update blutils",
	Long:  `Update blutils`,
	Run: func(cmd *cobra.Command, args []string) {
		var tag string

		// Check if Go is installed
		_, err := exec.LookPath("go")
		if err != nil {
			verbosePrintln(0, "Go is not installed")
			return
		}

		if *updateParams.Tag == "" {
			verbosePrintln(2, "Getting latest tag from github...")

			// Get list of tags from github
			response, err := http.Get("https://api.github.com/repos/" + *updateParams.Repo + "/git/refs/tags")
			if err != nil {
				verbosePrintln(0, "Error fetching tags from github:", err)
				return
			}
			defer response.Body.Close()

			type TagReference struct {
				Ref string `json:"ref"`
			}

			var tagReferences []TagReference
			var refs []string

			if err := json.NewDecoder(response.Body).Decode(&tagReferences); err != nil {
				verbosePrintln(0, "Error decoding JSON response:", err)
				return
			}

			for _, ref := range tagReferences {
				refs = append(refs, strings.TrimLeft(ref.Ref, "refs/tags/"))
			}
			// Sort tags
			semver.Sort(refs)

			verbosePrintln(3, "Tags:", refs)

			// Get latest tag
			tag = refs[len(refs)-1]
		}

		verbosePrintln(2, "Using tag:", tag)

		// Clean temp dir
		verbosePrintln(2, "Cleaning temp dir...")
		if err := os.RemoveAll(*updateParams.TempDir); err != nil {
			verbosePrintln(0, "Error cleaning temp dir:", err)
			return
		}

		// Clone
		verbosePrintln(2, "Cloning repo...")
		git.PlainClone(*updateParams.TempDir, false, &git.CloneOptions{
			URL:           "https://github.com/" + *updateParams.Repo,
			Progress:      sideband.NewMuxer(sideband.Sideband64k, verbosePrintlnW(3)),
			ReferenceName: plumbing.NewTagReferenceName(tag),
			Depth:         1,
		})

		// Build
		verbosePrintln(2, "Building...")
		buildPath := filepath.Join(*updateParams.TempDir)
		command := exec.Command("go", "build", "-ldflags", "-s -w", "-o", "blutils", ".")
		command.Dir = buildPath
		output, err := command.CombinedOutput()
		if err != nil {
			verbosePrintln(0, "Error during build:", err)
			verbosePrintln(0, "Build output:", string(output))
			return
		}
		verbosePrintln(2, "Build successful")

		exePath, _ := os.Executable()
		oldExePath := exePath + ".bak"

		verbosePrintln(2, "Moving old executable to:", oldExePath)
		if err := os.Rename(exePath, oldExePath); err != nil {
			verbosePrintln(0, "Error moving old executable:", err)
			return
		}

		verbosePrintln(2, "Moving new executable to:", exePath)
		if err := os.Rename(filepath.Join(buildPath, "blutils"), exePath); err != nil {
			verbosePrintln(0, "Error moving new executable:", err)
			return
		}

		verbosePrintln(2, "Cleaning up...")
		if err := os.RemoveAll(*updateParams.TempDir); err != nil {
			verbosePrintln(0, "Error removing temp dir:", err)
		}
		if err := os.Remove(oldExePath); err != nil {
			verbosePrintln(0, "Error removing old executable:", err)
		}
		verbosePrintln(2, "Update successful")
	},
}
