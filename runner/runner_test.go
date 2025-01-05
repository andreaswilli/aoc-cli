package runner_test

import (
	"aoc-cli/engine"
	"aoc-cli/reporter"
	"aoc-cli/runner"
	"testing"
	"testing/fstest"
)

var mapFS = fstest.MapFS{
	"2024/day_01/solution.nix": &fstest.MapFile{},
	"2024/day_01/expected.txt": &fstest.MapFile{Data: []byte("HelloWorld!\n")},
	"2024/day_02/solution.nix": &fstest.MapFile{},
	"2024/day_02/expected.txt": &fstest.MapFile{Data: []byte("Something else")},
	"2024/day_03/other.nix":    &fstest.MapFile{},
	"2024/unrelated-file.txt":  &fstest.MapFile{},
}

var mapFSEngines = fstest.MapFS{
	"aoc-cli.json": &fstest.MapFile{
		Data: []byte(`{"engines":[
      {
        "name": "echo",
        "cmd": "echo HelloWorld!",
        "entryFile": "solution.nix"
      }
    ]}`),
	},
}

func TestRun(t *testing.T) {
	cases := []struct {
		name         string
		path         string
		numReports   int
		wantStatuses map[string]reporter.Status
	}{
		{
			"single day is runnable",
			"2024/day_01",
			2,
			map[string]reporter.Status{
				"2024/day_01": reporter.StatusPassed,
			},
		},
		{
			"single day is not runnable",
			"2024/day_03",
			0,
			nil,
		},
		{
			"two days are runnable",
			"2024",
			4,
			map[string]reporter.Status{
				"2024/day_01": reporter.StatusPassed,
				"2024/day_02": reporter.StatusFailed,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			em, err := engine.NewEngineManager(mapFSEngines)

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			r := runner.NewRunner(mapFS, em)
			reportChan, err := r.Run(c.path)

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			var got []runner.ReportMap
			for reportMap := range reportChan {
				got = append(got, reportMap)
			}

			if len(got) != c.numReports {
				t.Errorf("Want %d report maps, got %d", c.numReports, len(got))
			}

			if c.numReports > 0 {
				finalReport := got[len(got)-1]
				for path, wantStatus := range c.wantStatuses {
					if finalReport[path].Status != wantStatus {
						t.Errorf("Want status %v, got %v", wantStatus, finalReport[path].Status)
					}
				}
			}
		})
	}
}
