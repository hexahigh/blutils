package cmd

import (
	"encoding/ascii85"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(ascii85Cmd)

	ascii85Cmd.Flags().BoolP("decode", "d", false, "Decode")
	ascii85Cmd.Flags().StringP("input", "i", "", "Input file")
	ascii85Cmd.Flags().StringP("output", "o", "", "Output file")

	configBindFlags(*ascii85Cmd)
}

var ascii85Cmd = &cobra.Command{
	Use:   "ascii85",
	Short: "ASCII85",
	Long: `Reads data from stdin or a file and converts it to ASCII85.
ASCII85 is a text encoding that can be used to store binary data.
It is more efficient than base64, base64 increases the size of the data by 33% while ASCII85 increases the size by 25%.`,
	Run: func(cmd *cobra.Command, args []string) {
		cn := commandToConfigString(*cmd)
		inputReader, err := getInputReader(viper.GetString(cn + ".input"))
		if err != nil {
			log.Errorf("Failed to open input: %v", err)
		}

		outputWriter, err := getOutputWriter(viper.GetString(cn + ".output"))
		if err != nil {
			log.Errorf("Failed to open output: %v", err)
		}

		if viper.GetBool(cn + ".decode") {
			decoder := ascii85.NewDecoder(inputReader)
			if _, err := io.Copy(outputWriter, decoder); err != nil {
				log.Errorf("Failed to decode: %v", err)
			}
		} else {
			encoder := ascii85.NewEncoder(outputWriter)
			if _, err := io.Copy(encoder, inputReader); err != nil {
				log.Errorf("Failed to encode: %v", err)
			}
		}

	},
}

func getInputReader(inputFilePath string) (io.Reader, error) {
	if inputFilePath == "" {
		return os.Stdin, nil
	}
	return os.Open(inputFilePath)
}

func getOutputWriter(outputFilePath string) (io.Writer, error) {
	if outputFilePath == "" || outputFilePath == "-" {
		return os.Stdout, nil
	}
	return os.Create(outputFilePath)
}
