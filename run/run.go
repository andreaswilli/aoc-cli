package run

import (
	"fmt"
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

func Watch(command string, filePath string) <-chan Result {
	outChan := make(chan Result)
	go func() {
		for {
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()

			outChan <- Run(command)

			doneChan := make(chan bool)

			go func(doneChan chan bool) {
				defer func() {
					doneChan <- true
				}()

				err := watchFile(filePath)
				if err != nil {
					fmt.Println(err)
				}

				fmt.Println("File has been changed")
			}(doneChan)

			<-doneChan
		}
	}()

	return outChan
}

func watchFile(filePath string) error {
	initialStat, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	for {
		stat, err := os.Stat(filePath)
		if err != nil {
			return err
		}

		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			break
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}
