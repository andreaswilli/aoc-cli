package reporter_test

import (
	"aoc-cli/reporter"
	"reflect"
	"testing"
)

func TestGetReport(t *testing.T) {
	cases := []struct {
		name   string
		got    string
		want   string
		report reporter.Report
	}{
		{
			"report no exp if expectation is empty and result is present",
			"the result",
			"",
			reporter.Report{"the result", "", reporter.StatusNoExp},
		},
		{
			"report no exp if expectation is empty and result is empty",
			"",
			"",
			reporter.Report{"", "", reporter.StatusNoExp},
		},
		{
			"report success if expectation matches the result",
			"the result",
			"the result",
			reporter.Report{"the result", "the result", reporter.StatusPassed},
		},
		{
			"report failure if expectation does not match the result",
			"the wrong result",
			"the result",
			reporter.Report{"the wrong result", "the result", reporter.StatusFailed},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			report := reporter.GetReport(c.got, c.want)

			if !reflect.DeepEqual(report, c.report) {
				t.Errorf("Expected %s but got %s", c.report, report)
			}
		})
	}
}
