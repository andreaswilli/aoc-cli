package main

import (
	"aoc-cli/cli"
	"aoc-cli/engine"
	"aoc-cli/reporter"
	"aoc-cli/runner"
	"fmt"
	"os"
	"sync"
)

func main() {
	CLI := cli.CLI{Args: os.Args[1:], Out: os.Stdout}
	filesystem := os.DirFS(".")
	engineManager, err := engine.NewEngineManager(filesystem)

	if err != nil {
		fmt.Printf("Engine error: %v", err)
	}

	r := runner.NewRunner(filesystem, engineManager)

	userCmd := CLI.GetUserCmd()

	if userCmd == nil {
		os.Exit(1)
	}

	var reportChan <-chan *reporter.Report
	var runErr error

	switch userCmd.SubCmd {
	case cli.SubCmdRun:
		reportChan, runErr = r.Run(userCmd.Path)
	case cli.SubCmdWatch:
		reportChan, runErr = r.Watch(userCmd.Path)
	}

	if runErr != nil {
		fmt.Printf("Execution error: %v", err)
	}

	reportMap := runner.ReportMap{}
	mutex := &sync.Mutex{}

	for report := range reportChan {
		mutex.Lock()
		reportMap[report.Path] = *report
		CLI.PrintReports(reportMap, cli.HidePassed)
		mutex.Unlock()
	}
}
