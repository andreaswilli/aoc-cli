package engine_test

import (
	"aoc-cli/engine"
	"errors"
	"os/exec"
	"reflect"
	"testing"
)

func TestEngine_GetCmd(t *testing.T) {
	cases := []struct {
		testName   string
		name       string
		cmd        string
		entryFile  string
		extraFiles []string
		path       string
		want       *exec.Cmd
		wantErr    error
	}{
		{
			"parse command and use entry file",
			"node",
			"node {{entryFile}}",
			"index.js",
      []string{},
			"2024/day_01",
			exec.Command("node", "2024/day_01/index.js"),
			nil,
		},
		{
			"get error for empty command",
			"node",
			"",
			"index.js",
      []string{},
			"2024/day_01",
			nil,
			errors.New("invalid command"),
		},
	}
	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			e := engine.Engine{c.name, c.cmd, c.entryFile, c.extraFiles}
			got, err := e.GetCmd(c.path)

			if err != nil && c.wantErr == nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("GetCmd() = %v, want %v", got, c.want)
			}
		})
	}
}
