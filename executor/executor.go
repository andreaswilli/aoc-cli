package executor

import (
	"aoc-cli/trigger"
	"os/exec"
	"time"
)

type Result struct {
	Out      string
	Err      error
	Duration time.Duration
}

func Execute(cmd *exec.Cmd, trigger trigger.Trigger) chan *Result {
	outChan := make(chan *Result)

	go func() {
		for range trigger.Listen() {
     // commands can only run once
      clonedCmd := *cmd

			start := time.Now()
			outputByteArray, err := clonedCmd.CombinedOutput()
			end := time.Now()

			result := &Result{
				Out: string(outputByteArray),
				Err: err,
				Duration: end.Sub(start),
			}
			outChan <- result
		}
		close(outChan)
	}()

	return outChan
}
