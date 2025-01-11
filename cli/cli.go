package cli

import (
	"aoc-cli/executor"
	"aoc-cli/reporter"
	"aoc-cli/runner"
	"io"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/x/ansi"
)

type HideLevel int

const (
	HideAll    HideLevel = 0
	HidePassed HideLevel = 1
	HideNone   HideLevel = 2
)

const (
	ResetColor = "\033[0m"
	Red        = "\033[31;27m"
	RedBG      = "\033[31;7m"
	Green      = "\033[32;27m"
	GreenBG    = "\033[32;7m"
	Blue       = "\033[34;27m"
	BlueBG     = "\033[34;7m"
	White      = "\033[97;27m"
	WhiteBG    = "\033[97;7m"
	Gray       = "\033[37;27m"
)

type UserCmd struct {
	SubCmd string
	Path   string
}

type CLI struct {
	Args         []string
	Out          io.Writer
	writtenLines int
}

const (
	SubCmdRun   = "run"
	SubCmdWatch = "watch"
)

var subCmds = []string{SubCmdRun, SubCmdWatch}

func (c *CLI) GetUserCmd() *UserCmd {
	if len(c.Args) == 0 {
		c.Out.Write([]byte("please provide a subcommand and a path to run\n"))
		return nil
	}

	subCmd := c.Args[0]
	if !slices.Contains(subCmds, subCmd) {
		c.Out.Write([]byte("unknown subcommand '" + subCmd + "'\n"))
		return nil
	}

	if len(c.Args) == 1 {
		c.Out.Write([]byte("please provide a path to run\n"))
		return nil
	}

	path := c.Args[1]
	return &UserCmd{SubCmd: subCmd, Path: path}
}

func (c *CLI) PrintReports(reports runner.ReportMap, hideLevel HideLevel) {
	output := "\n"
	output += ansi.CursorUp(c.writtenLines)
	output += ansi.EraseDisplay(0)

	for _, path := range sortedPaths(reports) {
		report := reports[path]
		if report.Status == reporter.StatusPassed {
			output += GreenBG + " PASSED " + Green
		} else if report.Status == reporter.StatusFailed {
			output += RedBG + " FAILED " + Red
		} else if report.Status == reporter.StatusNoExp {
			output += BlueBG + " NO EXP " + Blue
		} else if report.Status == reporter.StatusExec {
			output += WhiteBG + "  EXEC  " + White
		}
		output += " " + path
		if report.Status != reporter.StatusExec {
			output += Gray + " (" + report.Result.Duration.Round(10*time.Microsecond).String() + ")"
		}
		output += ResetColor + "\n"
		output += printDetails(hideLevel, report.Result, report.Expected)
	}

	c.writtenLines = strings.Count(output, "\n")
	c.Out.Write([]byte(output))
}

func sortedPaths(reports runner.ReportMap) []string {
	paths := make([]string, 0, len(reports))

	for path := range reports {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths
}

func printDetails(
	hideLevel HideLevel,
	result *executor.Result,
	expected string,
) (output string) {
	if result.Pending {
		return
	}
	if hideLevel >= HidePassed && expected != "" && expected != result.Out {
		output += "\nExpected:\n" + expected
		output += newlineIfNeeded(output)
		output += "\nGot:\n" + result.Out
		output += newlineIfNeeded(output)
		output += "\n"
	} else if hideLevel == HideNone || hideLevel == HidePassed && expected == "" {
		output += "\n" + result.Out
		output += newlineIfNeeded(output)
		output += "\n"
	}
	return
}

func newlineIfNeeded(s string) string {
	if !strings.HasSuffix(s, "\n") {
		return "\n"
	}
	return ""
}
