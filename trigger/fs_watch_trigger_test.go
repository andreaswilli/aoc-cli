package trigger_test

import (
	fswatcher "aoc-cli/fs_watcher"
	"aoc-cli/trigger"
	"testing"
	"testing/fstest"
	"time"
)

var mapFS = fstest.MapFS{
	"2024/day_01/solution.js": &fstest.MapFile{},
	"2024/day_01/input.txt":   &fstest.MapFile{},
}

func TestFsWatchTrigger_Listen(t *testing.T) {
	cases := []struct {
		name           string
		watchPaths     []string
		changes        []string
		triggeredTimes int
	}{
		{
			"should trigger once if no files change",
			[]string{},
			[]string{},
			1,
		},
		{
			"should trigger three times if two watched files change",
			[]string{"2024/day_01"},
			[]string{
				"2024/day_01/solution.js",
				"2024/day_01/input.txt",
			},
			3,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fsWatcher := &fswatcher.FsWatcher{
				FS: mapFS,
				WatchPaths: c.watchPaths,
				CheckInterval: time.Millisecond,
			}
			trigger := trigger.FsWatchTrigger{FsWatcher: fsWatcher}
			triggerChan := trigger.Listen()

			doneChan := make(chan bool)
			triggeredTimes := 0
			go func() {
				for range triggerChan {
					triggeredTimes += 1
					if triggeredTimes == c.triggeredTimes {
						doneChan <- true
						close(doneChan)
						break
					}
				}
			}()

      for _, change := range c.changes {
				time.Sleep(5 * time.Millisecond)
				mapFS[change] = &fstest.MapFile{
					ModTime: time.Now(),
					Data:    []byte("new content"),
				}
			}

			select {
			case <-doneChan:
			case <-time.After(100 * time.Millisecond):
				t.Errorf(
					"should trigger %d times but only triggered %d times before timing out",
					c.triggeredTimes,
					triggeredTimes,
				)
			}
		})
	}
}
