package reporter_test

import (
	"aoc-cli/executor"
	"aoc-cli/reporter"
	"errors"
	"testing"
	"time"
)

func TestGetReport(t *testing.T) {
	cases := []struct {
		name           string
		inResult       *executor.Result
		want           string
		expectedStatus reporter.Status
	}{
		{
			"report no exp if expectation is empty and result is present",
			createSuccessResult("the result"),
			"",
			reporter.StatusNoExp,
		},
		{
			"report no exp if expectation is empty and result is empty",
			createSuccessResult(""),
			"",
			reporter.StatusNoExp,
		},
		{
			"report success if expectation matches the result",
			createSuccessResult("the result"),
			"the result",
			reporter.StatusPassed,
		},
		{
			"report failure if expectation does not match the result",
			createSuccessResult("the wrong result"),
			"the result",
			reporter.StatusFailed,
		},
		{
			"report failure if command fails but output matches",
			createFailureResult("the failed result", errors.New("error output")),
			"the failed result",
			reporter.StatusFailed,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			report := reporter.GetReport(c.inResult, c.want)

			if report.Status != c.expectedStatus {
				t.Errorf("Expected %s but got %s", c.expectedStatus, report.Status)
			}
		})
	}
}

func createSuccessResult(out string) *executor.Result {
	return &executor.Result{Out: out, Err: nil, Duration: time.Millisecond}
}

func createFailureResult(out string, err error) *executor.Result {
	return &executor.Result{Out: out, Err: err, Duration: time.Millisecond}
}
