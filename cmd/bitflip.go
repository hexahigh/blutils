package cmd

import (
	"crypto/rand"
	"log"
	"math/big"
	"os"

	"github.com/spf13/cobra"
)

type BitflipParams struct {
	BitsToFlip *int
	Percentage *int
	MinOffset  *int
}

var bitflipParams BitflipParams

func init() {
	rootCmd.AddCommand(bitflipCmd)

	bitflipParams.BitsToFlip = bitflipCmd.Flags().IntP("bits", "b", 0, "Number of bits to flip")
	bitflipParams.Percentage = bitflipCmd.Flags().IntP("percentage", "p", 0, "Percentage of bits to flip (Will be ignored if set to 0 or --bits is set)")
	bitflipParams.MinOffset = bitflipCmd.Flags().IntP("min-offset", "m", 0, "Minimum offset")
}

var bitflipCmd = &cobra.Command{
	Use:   "bitflip [filename]",
	Short: "Simulates a bitflip",
	Long:  `Simulates a bitflip`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		fileContent, err := os.ReadFile(filename)
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}

		bitsToFlip := *bitflipParams.BitsToFlip
		for i := 0; i < bitsToFlip; i++ {
			// Choose a random position
			maxPos := big.NewInt(int64(len(fileContent)))
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

			verbosePrintln(3, "Flipped at offset", pos.Int64())
		}

		// Write back to file
		err = os.WriteFile(filename, fileContent, 0644)
		if err != nil {
			log.Fatalf("Failed to write file: %v", err)
		}

		verbosePrintln(2, bitsToFlip, "bits flipped")
	},
}
