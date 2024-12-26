package reporter_test

import (
	"aoc-cli/reporter"
	"fmt"
	"testing"
)

func TestReport(t *testing.T) {
	cases := []struct {
		got    string
		want   string
		status reporter.Status
	}{
		{"the result", "", reporter.StatusNoExp},
		{"the result", "the result", reporter.StatusPassed},
		{"the wrong result", "the result", reporter.StatusFailed},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf(
			"report %q for actual=%q and expected=%q",
			c.status,
			c.got,
			c.want,
		), func(t *testing.T) {
			status := reporter.Report(c.got, c.want)

			if status != c.status {
				t.Errorf("Expected %s but got %s", c.status, status)
			}
		})
	}
}
