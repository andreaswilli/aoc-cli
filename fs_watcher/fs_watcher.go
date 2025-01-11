package fswatcher

import (
	"io/fs"
	"time"
)

type FsWatcher struct {
	FS         fs.FS
	WatchPaths []string
	CheckInterval   time.Duration
}

func (f *FsWatcher) WaitForAnyChange() error {
	initialStatsMap, err := f.getStats()
	if err != nil {
		return err
	}

	for {
		statsMap, err := f.getStats()
		if err != nil {
			return err
		}

		if len(statsMap) != len(initialStatsMap) {
			break
		}

		var modified bool
		for path, stat := range statsMap {
			if initialStatsMap[path] == nil ||
				initialStatsMap[path].Size() != stat.Size() ||
				initialStatsMap[path].ModTime() != stat.ModTime() {
				modified = true
				break
			}
		}

		if modified {
			break
		}

		time.Sleep(f.CheckInterval)
	}

	return nil
}

func (f *FsWatcher) getStats() (map[string]fs.FileInfo, error) {
	statsMap := make(map[string]fs.FileInfo)
	for _, path := range f.WatchPaths {
		stat, err := fs.Stat(f.FS, path)
		if err != nil {
			return nil, err
		}

		if stat.IsDir() {
			entries, err := fs.ReadDir(f.FS, path)
			if err != nil {
				return nil, err
			}

			for _, entry := range entries {
				stat, err := fs.Stat(f.FS, path+"/"+entry.Name())
				if err != nil {
					return nil, err
				}
				if !stat.IsDir() {
					statsMap[path+"/"+entry.Name()] = stat
				}
			}

		} else {
			statsMap[path] = stat
		}
	}
	return statsMap, nil
}
