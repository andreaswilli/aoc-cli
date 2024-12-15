package main

import (
	"aoc-cli/run"
	"fmt"
	"os"
	"os/exec"
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
	White   = "\033[97;27m"
)

func main() {
	subcommand := os.Args[1]
	path := os.Args[2]
	fileToRun := fmt.Sprintf("%s/solution.nix", path)

	command := fmt.Sprintf("nix eval --quiet --experimental-features pipe-operator --extra-experimental-features nix-command --extra-experimental-features flakes --file %s", fileToRun)

	if subcommand == "run" {
		printResult(run.Run(command), path)
	} else if subcommand == "watch" {
		for result := range run.Watch(command, path) {
			clearScreen()
			printResult(result, path)
		}
	} else {
		fmt.Printf("Unknown subcommand '%s'\n", subcommand)
	}
}

func printResult(result run.Result, path string) {
	if result.Err != nil {
		fmt.Println(Red + result.Out + result.Err.Error() + Reset)
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
		fmt.Println("\n\n" + result.Out)
	} else if result.Out == expectedOutput {
		printBadge("PASSED", path)
		fmt.Println("\n\n" + result.Out)
	} else {
		printBadge("FAILED", path)
		fmt.Println("\n\n" + "Got:\n" + result.Out + "\nExpected:\n" + expectedOutput)
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
	} else {
		color = Blue
		colorBG = BlueBG
	}

	fmt.Printf("%s %s %s %s%s", colorBG, text, color, path, Reset)
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

}
