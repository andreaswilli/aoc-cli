package cli_test

import (
	"aoc-cli/cli"
	"aoc-cli/executor"
	"aoc-cli/reporter"
	"aoc-cli/runner"
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/charmbracelet/x/ansi"
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
	pendingReport = reporter.Report{
		Result:   createPendingResult(),
		Expected: "the result",
		Status:   reporter.StatusExec,
	}
)

var reportStart = "\n" + ansi.CursorUp(0) + ansi.EraseDisplay(0)
var duration = " (50.13ms)"

func TestCLI_GetUserCmd(t *testing.T) {
	cases := []struct {
		name    string
		args    []string
		wantCmd *cli.UserCmd
		wantOut string
	}{
    {
      "return valid command",
      []string{"run", "2024/day_01"},
      &cli.UserCmd{"run", "2024/day_01"},
      "",
    },
		{
			"print error if no arguments are given",
			[]string{},
			nil,
			"please provide a subcommand and a path to run\n",
		},
    {
      "print error for unknown subcommand",
      []string{"unknowncmd"},
      nil,
      "unknown subcommand 'unknowncmd'\n",
    },
    {
      "print error if only one argument is given for valid subcommand",
      []string{"run"},
      nil,
      "please provide a path to run\n",
    },
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out := new(bytes.Buffer)
			cli := cli.CLI{Args: c.args, Out: out}

			got := cli.GetUserCmd()

			if !reflect.DeepEqual(got, c.wantCmd) {
				t.Errorf("GetUserCmd() = %v, want %v", got, c.wantCmd)
			}

			if out.String() != c.wantOut {
				t.Errorf("Got output: %q, want %q", out.String(), c.wantOut)
			}
		})
	}
}

func TestCLI_PrintReports(t *testing.T) {
	cases := []struct {
		name      string
		reports   runner.ReportMap
		hideLevel cli.HideLevel
		output    string
	}{
		{
			"print one successful report",
			runner.ReportMap{"2024/day_01": successfulReport},
			cli.HidePassed,
			reportStart +
				cli.GreenBG + " PASSED " + cli.Green + " 2024/day_01" + cli.Gray +
				duration + cli.ResetColor + "\n",
		},
		{
			"print no additional newline if the result already ends in one",
			runner.ReportMap{"2024/day_02": failedReportWithNewlines},
			cli.HidePassed,
			reportStart +
				cli.RedBG + " FAILED " + cli.Red + " 2024/day_02" + cli.Gray +
				duration + cli.ResetColor + "\n" +
				"\nExpected:\nthe result\n\nGot:\nthe wrong result\n\n",
		},
		{
			"print three different reports in alphabetical order with all details hidden",
			runner.ReportMap{
				"2024/day_03": noExpReport,
				"2024/day_01": successfulReport,
				"2024/day_04": pendingReport,
				"2024/day_02": failedReport,
			},
			cli.HideAll,
			reportStart +
				cli.GreenBG + " PASSED " + cli.Green + " 2024/day_01" + cli.Gray +
				duration + cli.ResetColor + "\n" +
				cli.RedBG + " FAILED " + cli.Red + " 2024/day_02" + cli.Gray +
				duration + cli.ResetColor + "\n" +
				cli.BlueBG + " NO EXP " + cli.Blue + " 2024/day_03" + cli.Gray +
				duration + cli.ResetColor + "\n" +
				cli.WhiteBG + "  EXEC  " + cli.White + " 2024/day_04" + cli.ResetColor + "\n",
		},
		{
			"print three different reports in alphabetical order with passed details hidden",
			runner.ReportMap{
				"2024/day_03": noExpReport,
				"2024/day_01": successfulReport,
				"2024/day_04": pendingReport,
				"2024/day_02": failedReport,
			},
			cli.HidePassed,
			reportStart +
				cli.GreenBG + " PASSED " + cli.Green + " 2024/day_01" + cli.Gray +
				duration + cli.ResetColor + "\n" +
				cli.RedBG + " FAILED " + cli.Red + " 2024/day_02" + cli.Gray +
				duration + cli.ResetColor + "\n" +
				"\nExpected:\nthe result\n\nGot:\nthe wrong result\n\n" +
				cli.BlueBG + " NO EXP " + cli.Blue + " 2024/day_03" + cli.Gray +
				duration + cli.ResetColor + "\n" +
				"\nthe result\n\n" +
				cli.WhiteBG + "  EXEC  " + cli.White + " 2024/day_04" + cli.ResetColor + "\n",
		},
		{
			"print three different reports in alphabetical order with no details hidden",
			runner.ReportMap{
				"2024/day_03": noExpReport,
				"2024/day_01": successfulReport,
				"2024/day_04": pendingReport,
				"2024/day_02": failedReport,
			},
			cli.HideNone,
			reportStart +
				cli.GreenBG + " PASSED " + cli.Green + " 2024/day_01" + cli.Gray +
				duration + cli.ResetColor + "\n" +
				"\nthe result\n\n" +
				cli.RedBG + " FAILED " + cli.Red + " 2024/day_02" + cli.Gray +
				duration + cli.ResetColor + "\n" +
				"\nExpected:\nthe result\n\nGot:\nthe wrong result\n\n" +
				cli.BlueBG + " NO EXP " + cli.Blue + " 2024/day_03" + cli.Gray +
				duration + cli.ResetColor + "\n" +
				"\nthe result\n\n" +
				cli.WhiteBG + "  EXEC  " + cli.White + " 2024/day_04" + cli.ResetColor + "\n",
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
	return &executor.Result{Out: out, Err: nil, Duration: 50125 * time.Microsecond}
}

func createPendingResult() *executor.Result {
	return &executor.Result{Pending: true, Out: "", Err: nil, Duration: 0}
}
