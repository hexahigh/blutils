package cmd

import (
	"crypto/rand"
	"math/big"
	"os"
	"os/signal"
	"syscall"

	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(bitflipCmd)

	bitflipCmd.Flags().IntP("bits", "b", 0, "Number of bits to flip")
	bitflipCmd.Flags().IntP("percentage", "p", 0, "Percentage of bits to flip (Will be ignored if 0 or if --bits is set)")
	bitflipCmd.Flags().IntP("min-offset", "m", 0, "Minimum offset")
	bitflipCmd.Flags().IntP("chunk", "c", 1, "If >1, flips bits in chunks of this size")
	bitflipCmd.Flags().BoolP("no-progress", "n", false, "Disable progress bar")
	bitflipCmd.Flags().BoolP("extreme", "e", false, "Flips to a random byte instead")
}

var bitflipCmd = &cobra.Command{
	Use:   "bitflip [filename]",
	Short: "Simulates a bitflip",
	Long: `Simulates a bitflip.
	If extreme mode is used it even works like a disk shredder! But i wouldn't recommend using it for that purpose
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cn := cmd.Name()
		filename := args[0]
		fileContent, err := os.ReadFile(filename)
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}

		maxPos := big.NewInt(int64(len(fileContent)))

		log.Debugln("maxPos:", maxPos)

		if viper.GetInt(cn+".percentage") > 0 && viper.GetInt(cn+".bits") == 0 {
			bitsToFlip := maxPos.Int64() * int64(viper.GetInt(cn+".percentage")) / 100
			viper.Set(cn+".bits", bitsToFlip)
		}

		var pb *progressbar.ProgressBar

		if !viper.GetBool(cn + ".no-progress") {
			pb = progressbar.Default(viper.GetInt64(cn + ".bits"))
		}

		bitsToFlip := viper.GetInt64(cn + ".bits")
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		for i := 0; i < int(bitsToFlip); i++ {
			select {
			case <-sigChan:
				log.Infoln("Interrupt received, saving file...")
				err = os.WriteFile(filename, fileContent, 0644)
				if err != nil {
					log.Errorln("Failed to save file:", err)
				} else {
					log.Infoln("File saved successfully.")
				}
				os.Exit(0)
			default:
				pos, err := rand.Int(rand.Reader, maxPos)
				if err != nil {
					log.Errorf("Failed to generate random number: %v", err)
				}

				if viper.GetInt64(cn+".min-offset") != 0 && pos.Int64() < int64(viper.GetInt(cn+".min-offset")) {
					i-- // Decrement counter to retry this iteration
					continue
				}

				if viper.GetBool(cn + ".extreme") {
					randomByte := make([]byte, 1)
					_, err = rand.Read(randomByte)
					if err != nil {
						log.Errorf("Failed to generate random byte: %v", err)
					}
					fileContent[pos.Int64()] = randomByte[0]
				} else {
					// Perform bit flip
					fileContent[pos.Int64()] ^= 1

					if viper.GetInt64(cn+".chunk") > 1 && pos.Int64()+int64(viper.GetInt64(cn+".chunk")) < maxPos.Int64() {
						for j := 0; j < viper.GetInt(cn+".chunk"); j++ {
							fileContent[pos.Int64()+int64(j)] ^= 1
						}
					}
				}

				pb.Add(1)
				log.Debugf("Flipped at offset %d", pos.Int64())
			}
		}

		// Save file after all bits have been flipped
		err = os.WriteFile(filename, fileContent, 0644)
		if err != nil {
			log.Errorln("Failed to save file:", err)
		}

		log.Infoln(bitsToFlip*viper.GetInt64(cn+".chunk"), "bits flipped")
	},
}
