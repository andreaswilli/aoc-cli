package engine_test

import (
	"aoc-cli/engine"
	"errors"
	"io/fs"
	"reflect"
	"testing"
	"testing/fstest"
)

var mapFSWithEngine = fstest.MapFS{
	"aoc-cli.json": &fstest.MapFile{
		Data: []byte(`{"engines":[
      {"name": "node", "cmd": "node {{entryFile}}", "entryFile": "index.js"},
      {"name": "python", "cmd": "python {{entryFile}}", "entryFile": "main.py"}
    ]}`),
	},
}

var mapFSNoEngine = fstest.MapFS{
	"aoc-cli.json": &fstest.MapFile{
		Data: []byte(`{"engines":[]}`),
	},
}

var mapFSEmpty = fstest.MapFS{}

func TestNewEngineManager(t *testing.T) {
	cases := []struct {
		name    string
		fsys    fs.FS
		wantErr error
	}{
		{
			"get error for empty fs",
			mapFSEmpty,
			errors.New("'aoc-cli.json' not found, it must exist in the root of your project"),
		},
		{
			"get error if no engines are defined",
			mapFSNoEngine,
			errors.New("no engines are defined in 'aoc-cli.json'"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, gotErr := engine.NewEngineManager(c.fsys)
			assertEqual(t, gotErr, c.wantErr)
			assertEqual(t, got, nil)
		})
	}
}

func TestEngineManager_Get(t *testing.T) {
	cases := []struct {
		testName   string
		engineName string
		want       *engine.Engine
		wantErr    error
	}{
		{
			"get existing engine",
			"node",
			&engine.Engine{
				Name:      "node",
				Cmd:       "node {{entryFile}}",
				EntryFile: "index.js",
			},
			nil,
		},
		{
			"get error for non-existing engine",
			"non-existing",
			nil,
			errors.New("engine \"non-existing\" not found"),
		},
	}
	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			em, err := engine.NewEngineManager(mapFSWithEngine)

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			got, gotErr := em.Get(c.engineName)

			if c.wantErr == nil {
				assertEqual(t, gotErr, nil)
			} else {
				assertEqual(t, gotErr.Error(), c.wantErr.Error())
			}

			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("Get() = %+v, want %+v", got, c.want)
			}
		})
	}
}

func TestEngineManager_FindAppropriateEngine(t *testing.T) {
	cases := []struct {
		name string
		fsys fs.FS
		want *engine.Engine
	}{
		{
			"find appropriate engine",
			fstest.MapFS{"index.js": &fstest.MapFile{}},
			&engine.Engine{
				Name:      "node",
				Cmd:       "node {{entryFile}}",
				EntryFile: "index.js",
			},
		},
		{
			"find no appropriate engine",
			fstest.MapFS{},
			nil,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			em, err := engine.NewEngineManager(mapFSWithEngine)

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			got := em.FindAppropriateEngine(c.fsys)

			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("FindAppropriateEngine() = %v, want %v", got, c.want)
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
