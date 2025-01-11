package fswatcher_test

import (
	"aoc-cli/fs_watcher"
	"testing"
	"testing/fstest"
	"time"
)

var mapFS = fstest.MapFS{
	"2024/day_01/solution.js": &fstest.MapFile{},
	"2024/day_02/solution.js": &fstest.MapFile{},
}

func TestFsWatcher_WaitForAnyChange(t *testing.T) {
	cases := []struct {
		name                 string
		watchPaths           []string
		change               string
		shouldRegisterChange bool
	}{
		{
			"register change if watched file changes",
			[]string{"2024/day_01/solution.js"},
			"2024/day_01/solution.js",
			true,
		},
		{
			"register change if file in watched folder changes",
			[]string{"2024/day_01"},
			"2024/day_01/solution.js",
			true,
		},
		{
			"register change if new file in watched folder is added",
			[]string{"2024/day_01"},
			"2024/day_01/new_file.txt",
			true,
		},
		{
			"do not register change if unwatched file changes",
			[]string{"2024/day_01"},
			"2024/day_02/solution.js",
			false,
		},
		{
			"do not register change if no files are changed",
			[]string{"2024/day_01"},
			"",
			false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			watcher := fswatcher.FsWatcher{
				FS:         mapFS,
				WatchPaths: c.watchPaths,
				CheckInterval:   time.Millisecond,
			}

			doneChan := make(chan bool)
			var err error

			go func() {
				err = watcher.WaitForAnyChange()
				doneChan <- true
				close(doneChan)
			}()

			if c.change != "" {
				time.Sleep(1 * time.Millisecond)
				mapFS[c.change] = &fstest.MapFile{
					ModTime: time.Now(),
					Data:    []byte("new content"),
				}
			}

			select {
			case <-doneChan:
				if !c.shouldRegisterChange {
					t.Error("should not register change but did")
				}
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			case <-time.After(100 * time.Millisecond):
				if c.shouldRegisterChange {
					t.Error("should register change but timed out")
				}
			}
		})
	}
}
