package cmd

import (
	"crypto/rand"
	"log"
	"math/big"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

type BitflipParams struct {
	BitsToFlip *int
	Percentage *int
	MinOffset  *int
	ChunkSize  *int
}

var bitflipParams BitflipParams

func init() {
	rootCmd.AddCommand(bitflipCmd)

	bitflipParams.BitsToFlip = bitflipCmd.Flags().IntP("bits", "b", 0, "Number of bits to flip")
	bitflipParams.Percentage = bitflipCmd.Flags().IntP("percentage", "p", 0, "Percentage of bits to flip (Will be ignored if set to 0 or --bits is set)")
	bitflipParams.MinOffset = bitflipCmd.Flags().IntP("min-offset", "m", 0, "Minimum offset")
	bitflipParams.ChunkSize = bitflipCmd.Flags().IntP("chunk", "c", 1, "If >1, flips bits in chunks of this size")
}

var bitflipCmd = &cobra.Command{
	Use:   "bitflip [filename]",
	Short: "Simulates a bitflip",
	Long:  `Simulates a bitflip`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		var filename string
		var fileContent []byte

		var err error

		saveFile := func() error {
			// Write back to file
			return os.WriteFile(filename, fileContent, 0644)
		}

		mainFunc := func() {
			filename = args[0]
			fileContent, err = os.ReadFile(filename)
			if err != nil {
				log.Fatalf("Failed to read file: %v", err)
			}

			maxPos := big.NewInt(int64(len(fileContent)))

			verbosePrintln(3, "maxPos:", maxPos)

			bitsToFlip := *bitflipParams.BitsToFlip
			for i := 0; i < bitsToFlip; i++ {
				// Choose a random position
				pos, err := rand.Int(rand.Reader, maxPos)
				if err != nil {
					log.Fatalf("Failed to generate random number: %v", err)
				}

				if *bitflipParams.MinOffset != 0 && pos.Int64() < int64(*bitflipParams.MinOffset) {
					i--
					continue
				}

				// Perform bit flip
				fileContent[pos.Int64()] ^= 1

				if *bitflipParams.ChunkSize > 1 && pos.Int64()+int64(*bitflipParams.ChunkSize) < maxPos.Int64() {
					for j := 0; j < *bitflipParams.ChunkSize; j++ {
						fileContent[pos.Int64()+int64(j)] ^= 1
					}
				}

				verbosePrintln(3, "Flipped at offset", pos.Int64())
			}

			err = saveFile()
			if err != nil {
				verbosePrintln(0, "Failed to save file:", err)
			}

			verbosePrintln(2, bitsToFlip**bitflipParams.ChunkSize, "bits flipped")

			os.Exit(0)
		}

		go mainFunc()

		verbosePrintln(2, "Running. Press CTRL-C to exit.")
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sc

		verbosePrintln(2, "Saving file before exiting...")

		_ = saveFile()

		// Run os.Exit just in case
		os.Exit(0)
	},
}
