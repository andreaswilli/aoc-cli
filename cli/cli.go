package cli

import (
	"aoc-cli/reporter"
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
	Gray       = "\033[37;27m"
	GrayBG     = "\033[37;7m"
)

type CLI struct {
	Out io.Writer
}

type ReportMap map[string]reporter.Report

func (c *CLI) PrintReports(reports ReportMap, hideLevel HideLevel) {
	output := ""

	for _, path := range sortedPaths(reports) {
		if reports[path].Status == reporter.StatusPassed {
			output += GreenBG + " PASSED " + Green
		} else if reports[path].Status == reporter.StatusFailed {
			output += RedBG + " FAILED " + Red
		} else if reports[path].Status == reporter.StatusNoExp {
			output += BlueBG + " NO EXP " + Blue
		}
		output += " " + path + "\n" + ResetColor
		output += printDetails(hideLevel, reports[path].Result.Out, reports[path].Expected)
	}
	c.Out.Write([]byte(output))
}

func sortedPaths(reports ReportMap) []string {
	paths := make([]string, 0, len(reports))

	for path := range reports {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths
}

func printDetails(hideLevel HideLevel, actual string, expected string) string {
	output := ""

	if hideLevel >= HidePassed && expected != "" && expected != actual {
		output += "\nExpected:\n" + expected
		output += newlineIfNeeded(output)
		output += "\nGot:\n" + actual
		output += newlineIfNeeded(output)
		output += "\n"
	} else if hideLevel == HideNone || hideLevel == HidePassed && expected == "" {
		output += "\n" + actual
		output += newlineIfNeeded(output)
		output += "\n"
	}
	return output
}

func newlineIfNeeded(s string) string {
	if !strings.HasSuffix(s, "\n") {
		return "\n"
	}
	return ""
}
