package reporter

import "aoc-cli/executor"

type Status string

const (
	StatusPassed Status = "PASSED"
	StatusFailed Status = "FAILED"
	StatusNoExp  Status = "NO EXP"
	StatusExec   Status = "EXEC"
)

type Report struct {
	Result   *executor.Result
	Expected string
	Status   Status
}

func GetReport(result *executor.Result, expected string) (report Report) {
	report = Report{result, expected, StatusFailed}

	if result.Err != nil {
		return
	}

	if result.Pending {
		report.Status = StatusExec
		return
	}

	if len(expected) == 0 {
		report.Status = StatusNoExp
	} else if result.Out == expected {
		report.Status = StatusPassed
	}
	return
}
