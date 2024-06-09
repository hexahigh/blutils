//go:build !no_bitflip

package cmd

import (
	"crypto/rand"
	"log"
	"math/big"
	"os"
	"os/signal"
	"syscall"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

type BitflipParams struct {
	BitsToFlip *int
	Percentage *int
	MinOffset  *int
	ChunkSize  *int
	NoProgress *bool
	Extreme    *bool
}

var bitflipParams BitflipParams

func init() {
	rootCmd.AddCommand(bitflipCmd)

	bitflipParams.BitsToFlip = bitflipCmd.Flags().IntP("bits", "b", 0, "Number of bits to flip")
	bitflipParams.Percentage = bitflipCmd.Flags().IntP("percentage", "p", 0, "Percentage of bits to flip (Will be ignored if 0 or if --bits is set)")
	bitflipParams.MinOffset = bitflipCmd.Flags().IntP("min-offset", "m", 0, "Minimum offset")
	bitflipParams.ChunkSize = bitflipCmd.Flags().IntP("chunk", "c", 1, "If >1, flips bits in chunks of this size")
	bitflipParams.NoProgress = bitflipCmd.Flags().BoolP("no-progress", "n", false, "Disable progress bar")
	bitflipParams.Extreme = bitflipCmd.Flags().BoolP("extreme", "e", false, "Flips to a random byte instead")
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

		maxPos := big.NewInt(int64(len(fileContent)))

		verbosePrintln(3, "maxPos:", maxPos)

		if *bitflipParams.Percentage > 0 && *bitflipParams.BitsToFlip == 0 {
			bitsToFlip := maxPos.Int64() * int64(*bitflipParams.Percentage) / 100
			*bitflipParams.BitsToFlip = int(bitsToFlip)
		}

		var pb *progressbar.ProgressBar

		if !*bitflipParams.NoProgress {
			pb = progressbar.Default(int64(*bitflipParams.BitsToFlip))
		}

		bitsToFlip := *bitflipParams.BitsToFlip
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		for i := 0; i < bitsToFlip; i++ {
			select {
			case <-sigChan:
				verbosePrintln(2, "Interrupt received, saving file...")
				err = os.WriteFile(filename, fileContent, 0644)
				if err != nil {
					verbosePrintln(0, "Failed to save file:", err)
				} else {
					verbosePrintln(2, "File saved successfully.")
				}
				os.Exit(0)
			default:
				pos, err := rand.Int(rand.Reader, maxPos)
				if err != nil {
					verbosePrintln(0, "Failed to generate random number: %v", err)
				}

				if *bitflipParams.MinOffset != 0 && pos.Int64() < int64(*bitflipParams.MinOffset) {
					i-- // Decrement counter to retry this iteration
					continue
				}

				if *bitflipParams.Extreme {
					randomByte := make([]byte, 1)
					_, err = rand.Read(randomByte)
					if err != nil {
						verbosePrintln(0, "Failed to generate random byte: %v", err)
					}
					fileContent[pos.Int64()] = randomByte[0]
				} else {
					// Perform bit flip
					fileContent[pos.Int64()] ^= 1

					if *bitflipParams.ChunkSize > 1 && pos.Int64()+int64(*bitflipParams.ChunkSize) < maxPos.Int64() {
						for j := 0; j < *bitflipParams.ChunkSize; j++ {
							fileContent[pos.Int64()+int64(j)] ^= 1
						}
					}
				}

				pb.Add(1)
				verbosePrintln(3, "Flipped at offset", pos.Int64())
			}
		}

		// Save file after all bits have been flipped
		err = os.WriteFile(filename, fileContent, 0644)
		if err != nil {
			verbosePrintln(0, "Failed to save file:", err)
		}

		verbosePrintln(2, bitsToFlip**bitflipParams.ChunkSize, "bits flipped")
	},
}
