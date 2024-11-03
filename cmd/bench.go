package cmd

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(benchCmd)

	benchCmd.Flags().IntP("cpu", "c", 0, "Number of CPU workers")
	benchCmd.Flags().IntP("timeout", "t", 10, "Maximum time in seconds")

	configBindFlags(*benchCmd)
}

var benchCmd = &cobra.Command{
	Use:   "bench",
	Short: "Simple benchmarking tool",
	Long:  `Simple benchmarking tool`,
	Run: func(cmd *cobra.Command, args []string) {
		cn := commandToConfigString(*cmd)
		startTime := time.Now()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var wg sync.WaitGroup
		opsCount := make(chan int, viper.GetInt(cn+".cpu"))
		done := make(chan bool)

		go func() {
			log.Infoln("Starting timer")
			if viper.GetInt(cn+".timeout") > 0 {
				time.Sleep(time.Duration(viper.GetInt(cn+".timeout")) * time.Second)
				cancel()
			}
		}()

		if viper.GetInt(cn+".cpu") > 0 {
			log.Infof("Starting %d CPU workers\n", viper.GetInt(cn+".cpu"))
		}

		for i := 0; i < viper.GetInt(cn+".cpu"); i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				count := 0
				for {
					select {
					case <-ctx.Done():
						opsCount <- count
						return
					default:
						a, b := 0, 1
						for j := 0; j < 10000; j++ {
							a, b = b, a+b
						}
						count++
					}
				}
			}(i)
		}

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		go func() {
			<-c
			cancel()
			done <- true
		}()

		wg.Wait()
		close(opsCount)

		endTime := time.Now()
		duration := endTime.Sub(startTime).Seconds()

		totalOps := 0
		for ops := range opsCount {
			totalOps += ops
		}

		log.Infof("Ran for %.2f seconds\n", duration)
		log.Infof("Total operations: %d\n", totalOps)
		log.Infof("Operations per second: %.2f\n", float64(totalOps)/float64(duration))
	},
}
