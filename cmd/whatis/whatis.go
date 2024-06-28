package whatis

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	root "github.com/hexahigh/blutils/cmd"
	"github.com/hexahigh/go-lib/ctinfo"
	"github.com/hexahigh/go-lib/sniff"
	"github.com/spf13/cobra"
)

type WhatIsParams struct {
}

var whatIsParams WhatIsParams

func init() {
	root.RootCmd.AddCommand(whatIsCmd)

}

var whatIsCmd = &cobra.Command{
	Use:   "whatis [files...]",
	Short: "Gets information about files",
	Long:  `Gets information about files`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		type Output struct {
			ContentType string   `json:"content_type"`
			Extensions  []string `json:"extensions"`
			Filetype    string   `json:"filetype"`
		}

		var output Output

		for _, filePath := range args {
			// Skip if file is a directory
			if info, err := os.Stat(filePath); err == nil && info.IsDir() {
				root.Logger.Println(3, "Skipping directory:", filePath)
				continue
			}

			// Read the first 1024 bytes of the file
			file, err := os.Open(filePath)
			if err != nil {
				root.Logger.Println(0, "Failed to read file:", err)
			}

			defer file.Close()

			buffer := make([]byte, 1024)
			n, err := file.Read(buffer)
			if err != nil && err != io.EOF {
				log.Printf("Failed to read file %s: %v", filePath, err)
				continue
			}

			output.ContentType = sniff.DetectContentType(buffer[:n])
			output.Extensions = ctinfo.GetCTInfo(output.ContentType).Extensions
			output.Filetype = ctinfo.GetCTInfo(output.ContentType).Filetype

			fmt.Println("Filename:", filepath.Base(filePath))
			fmt.Println("Content type:", output.ContentType)
			fmt.Println("Extensions:", output.Extensions)
			fmt.Println("Filetype:", output.Filetype)
			fmt.Println()

		}
	},
}
