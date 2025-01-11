package trigger

import fswatcher "aoc-cli/fs_watcher"

type FsWatchTrigger struct {
	FsWatcher fswatcher.FsWatcher
}

func (t *FsWatchTrigger) Listen() chan bool {
	nextChan := make(chan bool)

	go func() {
    for {
      nextChan <- true
      t.FsWatcher.WaitForAnyChange()
    }
	}()

	return nextChan
}
