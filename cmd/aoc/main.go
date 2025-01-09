package main

import (
	"aoc-cli/cli"
	"aoc-cli/engine"
	"aoc-cli/run"
	"aoc-cli/runner"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	Reset   = "\033[0m"
	Red     = "\033[31;27m"
	RedBG   = "\033[31;7m"
	Green   = "\033[32;27m"
	GreenBG = "\033[32;7m"
	Yellow  = "\033[33;27m"
	Blue    = "\033[34;27m"
	BlueBG  = "\033[34;7m"
	Purple  = "\033[35;27m"
	Cyan    = "\033[36;27m"
	Gray    = "\033[37;27m"
	GrayBG  = "\033[37;7m"
	White   = "\033[97;27m"
)

func main() {
	CLI := cli.CLI{Out: os.Stdout}
	filesystem := os.DirFS(".")
	engineManager, err := engine.NewEngineManager(filesystem)

	if err != nil {
		fmt.Printf("Engine error: %v", err)
	}

	r := runner.NewRunner(filesystem, engineManager)

	reportChan, err := r.Run("2024")

	if err != nil {
		fmt.Printf("Unexpected error: %v", err)
	}

	reportMap := runner.ReportMap{}
	mutex := &sync.Mutex{}

	for report := range reportChan {
		mutex.Lock()
		reportMap[report.Path] = *report
		CLI.PrintReports(reportMap, cli.HidePassed)
		mutex.Unlock()
		fmt.Println("=========")
	}
}

func main_2() {
	subcommand := os.Args[1]
	path := os.Args[2]

	if strings.Contains(path, "/") {
		fileToRun := fmt.Sprintf("%s/solution.nix", path)
		command := fmt.Sprintf("nix eval --quiet --experimental-features pipe-operator --extra-experimental-features nix-command --extra-experimental-features flakes --file %s", fileToRun)

		if subcommand == "run" {
			printBadge(" EXEC ", path)
			result := run.Run(command)
			clearLine()
			printResult(result, path, "all")
		} else if subcommand == "watch" {
			for result := range run.Watch(command, path) {
				clearScreen()
				printResult(result, path, "all")
			}
		} else {
			fmt.Printf("Unknown subcommand '%s'\n", subcommand)
		}
	} else {
		// run all days in this directory
		days, err := os.ReadDir(path)
		if err != nil {
			panic(err)
		}

		dayNames := make([]string, 0)
		for _, day := range days {
			dayNames = append(dayNames, path+"/"+day.Name())
		}
		if subcommand == "run" {
			for _, day := range dayNames {
				printBadge(" EXEC ", day)
				command := getCommand(day)
				result := run.Run(command)
				clearLine()
				printResult(result, day, "hide_successful")
			}
		} else if subcommand == "watch" {
			resultsChan := make(chan map[string]run.Result)

			resultsMap := make(map[string]run.Result)
			for _, day := range dayNames {
				resultsMap[day] = run.Result{}
			}

			for _, day := range dayNames {
				command := getCommand(day)
				go func() {
					for result := range run.Watch(command, day) {
						resultsMap[day] = result
						resultsChan <- resultsMap
					}
				}()
			}
			for results := range resultsChan {
				clearScreen()

				keys := make([]string, 0)
				for k := range results {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, day := range keys {
					printResult(results[day], day, "hide_successful")
				}
			}
		}
	}
}

func getCommand(day string) string {
	fileToRun := fmt.Sprintf("%s/solution.nix", day)
	return fmt.Sprintf("nix eval --quiet --experimental-features pipe-operator --extra-experimental-features nix-command --extra-experimental-features flakes --file %s", fileToRun)
}

func printResult(result run.Result, path string, details string) {
	if result.Err != nil {
		fmt.Println(Red + result.Out + result.Err.Error() + Reset)
		return
	}

	if result.Out == "" {
		printBadge(" EXEC ", path)
		fmt.Println()
		return
	}

	expectedOutputFilePath := fmt.Sprintf("%s/expected.txt", path)

	expectedOutput := ""
	content, err := os.ReadFile(expectedOutputFilePath)
	if err == nil {
		expectedOutput = string(content)
	}

	if expectedOutput == "" {
		printBadge("NO EXP", path)
		fmt.Printf(" (%s)\n", result.Duration.Round(10*time.Microsecond))
		if details == "all" || details == "hide_successful" {
			fmt.Println("\n" + result.Out)
		}
	} else if result.Out == expectedOutput {
		printBadge("PASSED", path)
		fmt.Printf(" (%s)\n", result.Duration.Round(10*time.Microsecond))
		if details == "all" {
			fmt.Println("\n" + result.Out)
		}
	} else {
		printBadge("FAILED", path)
		fmt.Printf(" (%s)\n", result.Duration.Round(10*time.Microsecond))
		if details == "all" || details == "hide_successful" {
			fmt.Println("\n" + "Got:\n" + result.Out + "\nExpected:\n" + expectedOutput)
		}
	}
}

func printBadge(text string, path string) {
	var color string
	var colorBG string

	if text == "PASSED" {
		color = Green
		colorBG = GreenBG
	} else if text == "FAILED" {
		color = Red
		colorBG = RedBG
	} else if text == "NO EXP" {
		color = Blue
		colorBG = BlueBG
	} else {
		color = Gray
		colorBG = GrayBG
	}

	fmt.Printf("%s %s %s %s%s", colorBG, text, color, path, Reset)
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func clearLine() {
	fmt.Print("\r\033[K")
}
