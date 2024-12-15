package run

import (
	"os"
	"os/exec"
	"strings"
	"time"
)

type Result struct {
	Out string
	Err error
}

func Run(command string) Result {
	name := strings.Split(command, " ")
	cmd := exec.Command(name[0], name[1:]...)

	outputByteArray, err := cmd.CombinedOutput()

	return Result{Out: string(outputByteArray), Err: err}
}

func Watch(command string, dirPath string) <-chan Result {
	outChan := make(chan Result)
	go func() {
		for {
			outChan <- Run(command)

			doneChan := make(chan bool)

			go func(doneChan chan bool) {
				defer func() {
					doneChan <- true
				}()

				err := watchFiles([]string{dirPath, "lib/nix"})
				if err != nil {
          panic(err)
				}
			}(doneChan)

			<-doneChan
		}
	}()

	return outChan
}

func watchFiles(filePaths []string) error {
	initialStatsMap, err := getStats(filePaths)
	if err != nil {
		return err
	}

	for {
		statsMap, err := getStats(filePaths)
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

		time.Sleep(1 * time.Second)
	}

	return nil
}

func getStats(filePaths []string) (map[string]os.FileInfo, error) {
	statsMap := make(map[string]os.FileInfo)
	for _, path := range filePaths {
		initialStat, err := os.Stat(path)
		if err != nil {
			return nil, err
		}

		if initialStat.IsDir() {
			entries, err := os.ReadDir(path)
			if err != nil {
				return nil, err
			}

			for _, entry := range entries {
				initialStat, err := os.Stat(path + "/" + entry.Name())
				if err != nil {
					return nil, err
				}
				if !initialStat.IsDir() {
					statsMap[path+"/"+entry.Name()] = initialStat
				}
			}

		} else {
			statsMap[path] = initialStat
		}
	}
	return statsMap, nil
}
