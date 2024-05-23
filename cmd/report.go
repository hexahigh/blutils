/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	color "github.com/hexahigh/go-lib/ansicolor"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Params struct {
	Yaml    *bool
	Json    *bool
	Stdout  *bool
	OutFile *string
}

var reportParams Params

func init() {
	rootCmd.AddCommand(reportCmd)

	reportParams.Yaml = reportCmd.Flags().BoolP("yaml", "y", false, "Output report in YAML format")
	reportParams.Json = reportCmd.Flags().BoolP("json", "j", false, "Output report in JSON format")
	reportParams.Stdout = reportCmd.Flags().BoolP("stdout", "s", false, "Output report to stdout")
	reportParams.OutFile = reportCmd.Flags().StringP("out", "o", "", "Output report to file, use - for stdout")

	reportCmd.ParseFlags(os.Args[1:])
}

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Creates a system report",
	Long:  `Creates a system report`,
	Run: func(cmd *cobra.Command, args []string) {
		type BlockDevice struct {
			Name        string   `json:"name"`
			MajMin      string   `json:"maj:min"`
			Rm          bool     `json:"rm"`
			Size        string   `json:"size"`
			Ro          bool     `json:"ro"`
			Type        string   `json:"type"`
			MountPoints []string `json:"mountpoints"`
		}
		type SystemReport struct {
			CPUInfo           map[string]string `json:"cpu_info"`
			OS                map[string]string `json:"os"`
			MEMInfo           map[string]string `json:"mem_info"`
			Env               map[string]string `json:"environment_variables"`
			BlockDevices      []BlockDevice     `json:"block_devices"`
			SwapInfo          map[string]string `json:"swap_info"`
			LscpuInfo         map[string]string `json:"lscpu_info"`
			InstalledPackages []string          `json:"installed_packages"`
			Nproc             int               `json:"nproc"`
		}
		var lsblkOutput struct {
			BlockDevices []BlockDevice `json:"blockdevices"`
		}

		var installedPackages []string

		var memMap, cpuMap, osMap, swapMap, lscpuMap map[string]string

		var out_file string

		//* Functions
		executeCmd := func(cmd *exec.Cmd) ([]byte, error) {
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			return out.Bytes(), err
		}

		removeEmptyLines := func(lines []string) []string {
			var result []string
			for _, line := range lines {
				if line != "" {
					result = append(result, line)
				}
			}
			return result
		}

		printSuccess := func(msg string) {
			verbosePrintlnC(2, color.Green+"[OK]"+color.Reset, msg)
		}

		printFailure := func(msg string) {
			verbosePrintlnC(1, color.Red+"[FAIL]"+color.Reset, msg)
		}

		var fileType string
		if *reportParams.Yaml {
			fileType = "yaml"
		} else if *reportParams.Json {
			fileType = "json"
		} else {
			fileType = "yaml"
		}

		if *reportParams.Yaml && *reportParams.Json {
			verbosePrintln(0, "Cannot use both -Y and -J")
			os.Exit(1)
		}

		if *reportParams.OutFile == "-" {
			*reportParams.Stdout = true
		}

		verbosePrintln(3, "Filetype:", fileType)
		verbosePrintln(3, "Output file:", out_file)

		// Get CPU Info
		verbosePrintln(3, "Getting CPU Info...")
		cpuInfo, err := os.Open("/proc/cpuinfo")
		if err != nil {
			printFailure("cpuinfo")
			verbosePrintln(3, "Error opening /proc/cpuinfo:", err)
		} else {
			defer cpuInfo.Close()

			scanner := bufio.NewScanner(cpuInfo)
			cpuMap = make(map[string]string)
			for scanner.Scan() {
				line := scanner.Text()
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					val := strings.TrimSpace(parts[1])
					cpuMap[key] = val
				}
			}
			if err := scanner.Err(); err != nil {
				printFailure("cpuinfo")
				verbosePrintln(3, "Error reading /proc/cpuinfo:", err)
				return
			}
			printSuccess("cpuinfo")
		}

		// Get OS Info
		verbosePrintln(3, "Getting OS Info...")
		osInfo, err := os.Open("/etc/os-release")
		if err != nil {
			printFailure("os-release")
			verbosePrintln(3, "Error opening /etc/os-release:", err)
		} else {
			defer osInfo.Close()

			scanner := bufio.NewScanner(osInfo)
			osMap = make(map[string]string)
			for scanner.Scan() {
				line := scanner.Text()
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					key := parts[0]
					val := strings.Trim(parts[1], `"`)
					osMap[key] = val
				}
			}

			if err := scanner.Err(); err != nil {
				printFailure("os-release")
				verbosePrintln(3, "Error reading /etc/os-release:", err)
				return
			}
			printSuccess("os-release")
		}

		// Get mem info
		verbosePrintln(3, "Getting Memory Info...")
		memInfo, err := os.Open("/proc/meminfo")
		if err != nil {
			printFailure("meminfo")
			verbosePrintln(3, "Error opening /proc/meminfo:", err)
		} else {
			defer cpuInfo.Close()

			scanner := bufio.NewScanner(memInfo)
			memMap = make(map[string]string)
			for scanner.Scan() {
				line := scanner.Text()
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					val := strings.TrimSpace(parts[1])
					memMap[key] = val
				}
			}

			if err := scanner.Err(); err != nil {
				printFailure("meminfo")
				verbosePrintln(3, "Error reading /proc/meminfo:", err)
				return
			}
			printSuccess("meminfo")
		}

		// Get env
		verbosePrintln(3, "Getting Environment Variables...")
		envMap := make(map[string]string)
		for _, envVar := range os.Environ() {
			parts := strings.SplitN(envVar, "=", 2)
			if len(parts) == 2 {
				key := parts[0]
				val := parts[1]
				envMap[key] = val
			}
		}
		printSuccess("env")

		// Get swap info
		verbosePrintln(3, "Getting Swap Info...")
		swapInfo, err := os.Open("/proc/swaps")
		if err != nil {
			printFailure("swapinfo")
			verbosePrintln(3, "Error opening /proc/swaps:", err)
		} else {
			defer swapInfo.Close()

			scanner := bufio.NewScanner(swapInfo)
			swapMap = make(map[string]string)
			for scanner.Scan() {
				line := scanner.Text()
				parts := strings.Fields(line)
				if len(parts) >= 3 {
					key := parts[0]
					val := strings.Join(parts[1:], " ")
					swapMap[key] = val
				}
			}

			if err := scanner.Err(); err != nil {
				printFailure("swapinfo")
				verbosePrintln(3, "Error reading /proc/swaps:", err)
				return
			}
			printSuccess("swapinfo")
		}

		// Execute lsblk -J
		verbosePrintln(3, "Executing lsblk -J...")
		command := exec.Command("lsblk", "-J")
		out, err := executeCmd(command)
		if err != nil {
			printFailure("lsblk")
			verbosePrintln(3, "Error executing lsblk -J:", err)
		} else {

			// Parse lsblk output
			err = json.Unmarshal(out, &lsblkOutput)
			if err != nil {
				printFailure("lsblk")
				verbosePrintln(3, "Error parsing lsblk output:", err)
				return
			}

			printSuccess("lsblk")
		}

		// Lscpu
		verbosePrintln(3, "Executing lscpu...")
		command = exec.Command("lscpu")
		out, err = executeCmd(command)
		if err != nil {
			printFailure("lscpu")
			verbosePrintln(3, "Error executing lscpu", err)
		} else {
			defer cpuInfo.Close()

			scanner := bufio.NewScanner(bytes.NewReader(out))
			lscpuMap = make(map[string]string)
			for scanner.Scan() {
				line := scanner.Text()
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					val := strings.TrimSpace(parts[1])
					lscpuMap[key] = val
				}
			}

			if err := scanner.Err(); err != nil {
				printFailure("lscpu")
				verbosePrintln(3, "Error executing lscpu", err)
				return
			}

			printSuccess("lscpu")
		}

		// Get installed packages
		verbosePrintln(3, "Getting installed packages...")
		switch osMap["ID"] {
		case "ubuntu", "debian":
			command = exec.Command("dpkg", "--get-selections")
			out, err = executeCmd(command)
			if err != nil {
				printFailure("pkg")
				verbosePrintln(3, "Error executing dpkg --get-selections:", err)
				return
			}
			installedPackages = strings.Split(string(out), "\n")
			for i, pkg := range installedPackages {
				fields := strings.Fields(pkg)
				if len(fields) > 0 {
					installedPackages[i] = fields[0] // Take the first field, which is the package name
				}
			}
			printSuccess("pkg")
		case "rhel":
			command = exec.Command("rpm", "-qa")
			out, err = executeCmd(command)
			if err != nil {
				printFailure("pkg")
				verbosePrintln(3, "Error executing rpm -qa:", err)
				return
			}
			installedPackages = strings.Split(string(out), "\n")
		default:
			printFailure("pkg")
			verbosePrintln(3, "Unsupported/Unknown OS:", osMap["ID"])
			return
		}

		// Remove empty lines from the installed packages list
		installedPackages = removeEmptyLines(installedPackages)

		//Run nproc
		command = exec.Command("nproc")
		out, err = executeCmd(command)
		if err != nil {
			printFailure("nproc")
			verbosePrintln(3, "Error executing nproc:", err)
			return
		} else {
			printSuccess("nproc")
			defer cpuInfo.Close()
		}
		nproc := strings.TrimSpace(string(out))
		nprocInt, err := strconv.Atoi(nproc)
		if err != nil {
			printFailure("nproc")
			verbosePrintln(3, "Error converting nproc to integer:", err)
			return
		} else {
			printSuccess("nproc")
		}

		report := &SystemReport{
			CPUInfo:           cpuMap,
			OS:                osMap,
			MEMInfo:           memMap,
			Env:               envMap,
			BlockDevices:      lsblkOutput.BlockDevices,
			SwapInfo:          swapMap,
			LscpuInfo:         lscpuMap,
			InstalledPackages: installedPackages,
			Nproc:             nprocInt,
		}

		// Marshal the report to either JSON or YAML
		verbosePrintln(3, "Marshalling report...")
		var data []byte
		switch fileType {
		case "yaml":
			data, err = yaml.Marshal(report)
			if err != nil {
				verbosePrintln(0, "Error marshalling report to YAML:", err)
				return
			}
			// Prepend the comment to the data
			var comment string
			comment += "# System report generated by Blutils v" + version + "\n"
			comment += "# Report generated on " + time.Now().Format("2006-01-02 15:04:05") + "\n"
			data = append([]byte(comment), data...)
		case "json":
			data, err = json.Marshal(report)
			if err != nil {
				verbosePrintln(0, "Error marshalling report to JSON:", err)
				return
			}
		}

		if *reportParams.Stdout {
			fmt.Print(string(data))
			return
		} else {
			// Write the data to the output file
			err = os.WriteFile(*reportParams.OutFile, data, fs.ModePerm)
			if err != nil {
				verbosePrintln(0, "Error writing report to file:", err)
				return
			}

			verbosePrintln(2, "System report completed.")
			verbosePrintln(2, "Report written to:", out_file)
			verbosePrintln(1, "The system report contains environment variables, installed packages and more! You may want to review the report if you are planning on sharing it.")
		}
	},
}
