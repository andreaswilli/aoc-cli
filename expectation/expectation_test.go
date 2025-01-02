package expectation_test

import (
	"aoc-cli/expectation"
	"testing"
	"testing/fstest"
)

var mapFS = fstest.MapFS{
	"2024/day_02/expected.txt": &fstest.MapFile{Data: []byte("expected result")},
}

func TestGetExpectation(t *testing.T) {
	cases := []struct {
		name string
		path string
		want string
	}{
		{
			"return empty string if there is no expectation",
			"2024/day_01/",
			"",
		},
		{
			"return file content if there is an expectation",
			"2024/day_02/",
			"expected result",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := expectation.GetExpectation(c.path, mapFS)
			if got != c.want {
				t.Errorf("GetExpectation() = %q, want %q", got, c.want)
			}
		})
	}
}
