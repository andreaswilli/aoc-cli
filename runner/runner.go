package runner

import (
	"aoc-cli/executor"
	"aoc-cli/expectation"
	"aoc-cli/reporter"
	"aoc-cli/trigger"
	"io/fs"
	"os/exec"
	"strings"
	"sync"
)

type ReportMap map[string]reporter.Report

type Runner struct {
	FS           fs.FS
	SourceSuffix string
}

func NewRunner(fsys fs.FS, sourceSuffix string) Runner {
	return Runner{FS: fsys, SourceSuffix: sourceSuffix}
}

func (r Runner) Run(path string) (reportChan chan ReportMap, err error) {
	reportChan = make(chan ReportMap)
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
		cmd := exec.Command("nix", "eval", "--quiet", "--experimental-features", "pipe-operator", "--extra-experimental-features", "nix-command", "--extra-experimental-features", "flakes", "--file", item)
		go func() {
			defer wg.Done()
			for result := range executor.Execute(cmd, &trigger.OneShotTrigger{}) {
				dir := strings.TrimSuffix(item, r.SourceSuffix)
				expected := expectation.GetExpectation(dir, r.FS)

				report := reporter.GetReport(result, expected)

				mutex.Lock()
				reportMap[item] = report
				mutex.Unlock()

				reportChan <- reportMap
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
	items, err := fs.Glob(r.FS, path+"/"+r.SourceSuffix)

	if err != nil {
		return []string{}, err
	}

	nestedItems, err := fs.Glob(r.FS, path+"/**/"+r.SourceSuffix)

	if err != nil {
		return []string{}, err
	}

	return append(items, nestedItems...), nil
}
