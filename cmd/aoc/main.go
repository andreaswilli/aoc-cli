package main

import (
	"aoc-cli/run"
	"fmt"
	"os"
)

var (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	White  = "\033[97m"
)

func main() {
	subcommand := os.Args[1]
	path := os.Args[2]
	fileToRun := fmt.Sprintf("%s/solution.nix", path)

	command := fmt.Sprintf("nix eval --experimental-features pipe-operator --extra-experimental-features nix-command --extra-experimental-features flakes --file %s", fileToRun)

	if subcommand == "run" {
		printResult(run.Run(command))
	} else if subcommand == "watch" {
		for result := range run.Watch(command, fileToRun) {
			printResult(result)
		}
	} else {
		fmt.Printf("Unknown subcommand '%s'\n", subcommand)
	}
}

func printResult(result run.Result) {
	if result.Err != nil {
		fmt.Println(Red + result.Out + result.Err.Error() + Reset)
	} else {
		fmt.Println(result.Out)
	}
}
