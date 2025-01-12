package reporter_test

import (
	"aoc-cli/executor"
	"aoc-cli/reporter"
	"errors"
	"reflect"
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
    {
      "report exec if result is pending",
      createPendingResult(),
      "the result",
      reporter.StatusExec,
    },
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			report := reporter.GetReport("2024/day_01", "node", c.inResult, c.want)
      assertEqual(t, report.Path, "2024/day_01")
      assertEqual(t, report.EngineName, "node")
      assertEqual(t, report.Status, c.expectedStatus)
		})
	}
}

func createSuccessResult(out string) *executor.Result {
  return &executor.Result{
  	Pending: false,
  	Out: out,
  	Err: nil,
  	Duration: time.Millisecond,
  }
}

func createFailureResult(out string, err error) *executor.Result {
	return &executor.Result{
  	Pending: false,
		Out: out,
		Err: err,
		Duration: time.Millisecond,
	}
}

func createPendingResult() *executor.Result {
  return &executor.Result{
  	Pending: true,
  	Out: "",
  	Err: nil,
  	Duration: 0,
  }
}

func assertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Expected '%v' but got '%v'", want, got)
	}
}
