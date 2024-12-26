package executor_test

import (
	"aoc-cli/executor"
	"aoc-cli/trigger"
	"errors"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

func TestExecute(t *testing.T) {
	cases := []struct {
		desc    string
		cmd     *exec.Cmd
		results []*executor.Result
	}{
		{
			"get output value",
			exec.Command("echo", "hello"),
			[]*executor.Result{{Out: "hello", Err: nil}},
		},
		{
			"get output error",
			exec.Command("ls", "nonexistent"),
			[]*executor.Result{{
				Out: "ls: nonexistent: No such file or directory",
				Err: errors.New("exit status 1"),
			}},
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			trigger := &trigger.OneShotTrigger{}
			results := []*executor.Result{}

			for result := range executor.Execute(c.cmd, trigger) {
				results = append(results, result)

        if len(results) > len(c.results) {
          t.Fatalf("Expected %d results, got %d", len(c.results), len(results))
        }
			}

			assertEqual(t, len(results), len(c.results))
			for i := range c.results {
				assertEqual(t, strings.TrimSpace(results[i].Out), c.results[i].Out)

				expectedErr := c.results[i].Err
				if expectedErr == nil {
					assertEqual(t, results[i].Err, nil)
				} else {
					assertEqual(t, results[i].Err.Error(), c.results[i].Err.Error())
				}

        if results[i].Duration == 0 {
          t.Errorf("Expected duration to be greater than 0")
        }
			}
		})
	}
}

func assertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Expected '%v' but got '%v'", want, got)
	}
}
