//go:build !no_bench

package bench

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	root "github.com/hexahigh/blutils/cmd"
)

type BenchParams struct {
	CpuWorkers *int
	Timeout    *int
}

var benchParams BenchParams

func init() {
	root.RootCmd.AddCommand(benchCmd)

	benchParams.CpuWorkers = benchCmd.Flags().IntP("cpu", "c", 0, "Number of CPU workers")
	benchParams.Timeout = benchCmd.Flags().IntP("timeout", "t", 10, "Maximum time in seconds")

	benchCmd.ParseFlags(os.Args[1:])
}

var benchCmd = &cobra.Command{
	Use:   "bench",
	Short: "Simple benchmarking tool",
	Long:  `Simple benchmarking tool`,
	Run: func(cmd *cobra.Command, args []string) {
		startTime := time.Now()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var wg sync.WaitGroup
		opsCount := make(chan int, *benchParams.CpuWorkers)
		done := make(chan bool)

		go func() {
			root.VerbosePrintln(3, "Starting timeout")
			if *benchParams.Timeout > 0 {
				time.Sleep(time.Duration(*benchParams.Timeout) * time.Second)
				cancel()
			}
		}()

		if *benchParams.CpuWorkers > 0 {
			fmt.Printf("Starting %d CPU workers\n", *benchParams.CpuWorkers)
		}

		for i := 0; i < *benchParams.CpuWorkers; i++ {
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

		fmt.Printf("Ran for %.2f seconds\n", duration)
		fmt.Printf("Total operations: %d\n", totalOps)
		fmt.Printf("Operations per second: %.2f\n", float64(totalOps)/float64(duration))
	},
}
