package cli_test

import (
	"aoc-cli/cli"
	"aoc-cli/executor"
	"aoc-cli/reporter"
	"bytes"
	"testing"
	"time"
)

var (
	successfulReport = reporter.Report{
		Result:   createSuccessResult("the result"),
		Expected: "the result",
		Status:   reporter.StatusPassed,
	}
	failedReport = reporter.Report{
		Result:   createSuccessResult("the wrong result"),
		Expected: "the result",
		Status:   reporter.StatusFailed,
	}
	failedReportWithNewlines = reporter.Report{
		Result:   createSuccessResult(failedReport.Result.Out + "\n"),
		Expected: failedReport.Expected + "\n",
		Status:   failedReport.Status,
	}
	noExpReport = reporter.Report{
		Result:   createSuccessResult("the result"),
		Expected: "",
		Status:   reporter.StatusNoExp,
	}
)

func TestCLI_PrintReports(t *testing.T) {
	cases := []struct {
		name      string
		reports   cli.ReportMap
		hideLevel cli.HideLevel
		output    string
	}{
		{
			"print one successful report",
			cli.ReportMap{"2024/day_01": successfulReport},
			cli.HidePassed,
			cli.GreenBG + " PASSED " + cli.Green + " 2024/day_01\n" + cli.ResetColor,
		},
		{
			"print no additional newline if the result already ends in one",
			cli.ReportMap{"2024/day_02": failedReportWithNewlines},
			cli.HidePassed,
			cli.RedBG + " FAILED " + cli.Red + " 2024/day_02\n" + cli.ResetColor +
				"\nExpected:\nthe result\n\nGot:\nthe wrong result\n\n",
		},
		{
			"print three different reports in alphabetical order with all details hidden",
			cli.ReportMap{
				"2024/day_03": noExpReport,
				"2024/day_01": successfulReport,
				"2024/day_02": failedReport,
			},
			cli.HideAll,
			cli.GreenBG + " PASSED " + cli.Green + " 2024/day_01\n" + cli.ResetColor +
				cli.RedBG + " FAILED " + cli.Red + " 2024/day_02\n" + cli.ResetColor +
				cli.BlueBG + " NO EXP " + cli.Blue + " 2024/day_03\n" + cli.ResetColor,
		},
		{
			"print three different reports in alphabetical order with passed details hidden",
			cli.ReportMap{
				"2024/day_03": noExpReport,
				"2024/day_01": successfulReport,
				"2024/day_02": failedReport,
			},
			cli.HidePassed,
			cli.GreenBG + " PASSED " + cli.Green + " 2024/day_01\n" + cli.ResetColor +
				cli.RedBG + " FAILED " + cli.Red + " 2024/day_02\n" + cli.ResetColor +
				"\nExpected:\nthe result\n\nGot:\nthe wrong result\n\n" +
				cli.BlueBG + " NO EXP " + cli.Blue + " 2024/day_03\n" + cli.ResetColor +
				"\nthe result\n\n",
		},
		{
			"print three different reports in alphabetical order with no details hidden",
			cli.ReportMap{
				"2024/day_03": noExpReport,
				"2024/day_01": successfulReport,
				"2024/day_02": failedReport,
			},
			cli.HideNone,
			cli.GreenBG + " PASSED " + cli.Green + " 2024/day_01\n" + cli.ResetColor +
				"\nthe result\n\n" +
				cli.RedBG + " FAILED " + cli.Red + " 2024/day_02\n" + cli.ResetColor +
				"\nExpected:\nthe result\n\nGot:\nthe wrong result\n\n" +
				cli.BlueBG + " NO EXP " + cli.Blue + " 2024/day_03\n" + cli.ResetColor +
				"\nthe result\n\n",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out := new(bytes.Buffer)
			cli := cli.CLI{Out: out}

			cli.PrintReports(c.reports, c.hideLevel)

			if out.String() != c.output {
				t.Error(
					"\n=== Expected: ===\n" + c.output +
						"\n=== but got: ===\n" + out.String(),
				)
			}
		})
	}
}

func createSuccessResult(out string) *executor.Result {
	return &executor.Result{Out: out, Err: nil, Duration: time.Millisecond}
}
