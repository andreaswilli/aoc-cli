package runner

import (
	"aoc-cli/executor"
	"aoc-cli/reporter"
	"aoc-cli/trigger"
	"io/fs"
	"os/exec"
	"sync"
)

type ReportMap map[string]reporter.Report

type Runner struct {
	FS fs.GlobFS
}

func (r Runner) Run(path string) (reportChan chan *ReportMap, err error) {
	reportChan = make(chan *ReportMap)
	items, err := r.getMatches(path)

	if err != nil {
		close(reportChan)
		return
	}

	reportMap := ReportMap{}
  mutex := &sync.Mutex{}

	wg := &sync.WaitGroup{}
	wg.Add(len(items))

	for _, item := range items {
		cmd := exec.Command("echo", item)
		go func() {
			defer wg.Done()
			for result := range executor.Execute(cmd, &trigger.OneShotTrigger{}) {
				report := reporter.GetReport(result, item + "\n")

        mutex.Lock()
				reportMap[item] = report
        mutex.Unlock()

				reportChan <- &reportMap
			}
		}()
	}

	go func() {
		wg.Wait()
		close(reportChan)
	}()
	return
}

func (r Runner) getMatches(path string) ([]string, error) {
	items, err := r.FS.Glob(path + "/solution.nix")

	if err != nil {
		return []string{}, err
	}

	nestedItems, err := r.FS.Glob(path + "/**/solution.nix")

	if err != nil {
		return []string{}, err
	}

	return append(items, nestedItems...), nil
}
