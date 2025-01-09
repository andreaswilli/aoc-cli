package runner

import (
	"aoc-cli/engine"
	"aoc-cli/executor"
	"aoc-cli/expectation"
	"aoc-cli/reporter"
	"aoc-cli/trigger"
	"io/fs"
	"strings"
	"sync"
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
			for result := range executor.Execute(cmd, &trigger.OneShotTrigger{}) {
				expected := expectation.GetExpectation(dir, r.FS)
				report := reporter.GetReport(dir, result, expected)
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
