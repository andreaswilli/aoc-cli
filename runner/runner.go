package runner

import (
	"aoc-cli/engine"
	"aoc-cli/executor"
	"aoc-cli/expectation"
	fswatcher "aoc-cli/fs_watcher"
	"aoc-cli/reporter"
	"aoc-cli/trigger"
	"io/fs"
	"strings"
	"sync"
	"time"
)

type ReportMap map[string]reporter.Report

type Runner struct {
	FS            fs.FS
	EngineManager *engine.EngineManager
}

func NewRunner(fsys fs.FS, engineManager *engine.EngineManager) Runner {
	return Runner{FS: fsys, EngineManager: engineManager}
}

func (r Runner) Run(path string) (reportChan chan *reporter.Report, err error) {
	return r.run(path, r.createOneShotTrigger)
}

func (r Runner) Watch(path string) (reportChan chan *reporter.Report, err error) {
	return r.run(path, r.createFsWatchTrigger)
}

func (r Runner) run(
	path string,
	createTrigger func(paths []string) trigger.Trigger,
) (
	reportChan chan *reporter.Report,
	err error,
) {
	reportChan = make(chan *reporter.Report)
	dirs, err := r.getDirs(path)

	if err != nil {
		close(reportChan)
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(dirs))

	for _, dir := range dirs {
		subFS, subErr := fs.Sub(r.FS, dir)

		if subErr != nil {
			wg.Done()
			continue
		}

		engine := r.EngineManager.FindAppropriateEngine(subFS)

		if engine == nil {
			wg.Done()
			continue
		}

		cmd, cmdErr := engine.GetCmd(dir)

		if cmdErr != nil {
			close(reportChan)
			return reportChan, cmdErr
		}

		go func() {
			defer wg.Done()
			paths := append(engine.ExtraFiles, dir)
			for result := range executor.Execute(cmd, createTrigger(paths)) {
				expected := expectation.GetExpectation(dir, r.FS)
				report := reporter.GetReport(dir, engine.Name, result, expected)
				reportChan <- report
			}
		}()
	}

	go func() {
		wg.Wait()
		close(reportChan)
	}()
	return
}

func (r Runner) getDirs(path string) ([]string, error) {
	if strings.Contains(path, "/") {
		return []string{path}, nil
	}

	dirEntries, err := fs.ReadDir(r.FS, path)

	if err != nil {
		return nil, err
	}

	dirs := []string{}
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			dirs = append(dirs, path+"/"+dirEntry.Name())
		}
	}

	return dirs, nil
}

func (r Runner) createOneShotTrigger(paths []string) trigger.Trigger {
	return &trigger.OneShotTrigger{}
}

func (r Runner) createFsWatchTrigger(paths []string) trigger.Trigger {
	fsWatcher := &fswatcher.FsWatcher{
		FS:            r.FS,
		WatchPaths:    paths,
		CheckInterval: time.Second,
	}
	return &trigger.FsWatchTrigger{FsWatcher: fsWatcher}
}
