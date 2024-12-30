package cli

import (
	"aoc-cli/executor"
	"aoc-cli/reporter"
	"aoc-cli/runner"
	"io"
	"sort"
	"strings"
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

type CLI struct {
	Out io.Writer
}

func (c *CLI) PrintReports(reports runner.ReportMap, hideLevel HideLevel) {
	output := ""

	for _, path := range sortedPaths(reports) {
		if reports[path].Status == reporter.StatusPassed {
			output += GreenBG + " PASSED " + Green
		} else if reports[path].Status == reporter.StatusFailed {
			output += RedBG + " FAILED " + Red
		} else if reports[path].Status == reporter.StatusNoExp {
			output += BlueBG + " NO EXP " + Blue
		} else if reports[path].Status == reporter.StatusExec {
			output += WhiteBG + "  EXEC  " + White
		}
		output += " " + path + "\n" + ResetColor
		output += printDetails(hideLevel, reports[path].Result, reports[path].Expected)
	}
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
