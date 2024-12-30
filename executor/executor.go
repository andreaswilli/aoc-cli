package executor

import (
	"aoc-cli/trigger"
	"os/exec"
	"time"
)

type Result struct {
	Pending  bool
	Out      string
	Err      error
	Duration time.Duration
}

func Execute(cmd *exec.Cmd, trigger trigger.Trigger) chan *Result {
	outChan := make(chan *Result)

	go func() {
		for range trigger.Listen() {
			outChan <- &Result{
				Pending:  true,
				Out:      "",
				Err:      nil,
				Duration: 0,
			}

			// commands can only run once
			clonedCmd := *cmd

			start := time.Now()
			outputByteArray, err := clonedCmd.CombinedOutput()
			end := time.Now()

			outChan <- &Result{
				Pending:  false,
				Out:      string(outputByteArray),
				Err:      err,
				Duration: end.Sub(start),
			}
		}
		close(outChan)
	}()

	return outChan
}
