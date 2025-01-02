package runner_test

import (
	"aoc-cli/runner"
	"testing"
	"testing/fstest"
)

var mapFS = fstest.MapFS{
	"2024/day_01/solution.nix": &fstest.MapFile{},
	"2024/day_02/solution.nix": &fstest.MapFile{},
	"2024/day_03/other.nix":    &fstest.MapFile{},
}

func TestRun(t *testing.T) {
	cases := []struct {
		name string
		path string
		want []*runner.ReportMap
	}{
		{
			"single test case exists",
			"2024/day_01",
			[]*runner.ReportMap{{}, {}},
		},
		{
			"single test case does not exist",
			"2024/day_03",
			[]*runner.ReportMap{},
		},
		{
			"two test cases exist",
			"2024",
			[]*runner.ReportMap{{}, {}, {}, {}},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := runner.NewRunner(mapFS, "solution.nix")

			reportChan, err := r.Run(c.path)

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			var got []runner.ReportMap
			for reportMap := range reportChan {
				got = append(got, reportMap)
			}

			if len(got) != len(c.want) {
				t.Errorf("Want %d report maps, got %d", len(c.want), len(got))
			}
		})
	}
}
