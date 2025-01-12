package engine

import (
	"errors"
	"os/exec"
	"strings"
)

type Engine struct {
	Name       string   `json:"name"`
	Cmd        string   `json:"cmd"`
	EntryFile  string   `json:"entryFile"`
	ExtraFiles []string `json:"extraFiles"`
}

func (e *Engine) GetCmd(path string) (*exec.Cmd, error) {
	cmdStr := strings.ReplaceAll(e.Cmd, "{{entryFile}}", path+"/"+e.EntryFile)
	binary := strings.Split(cmdStr, " ")[0]
	args := strings.Split(cmdStr, " ")[1:]

	if binary == "" {
		return nil, errors.New("invalid command")
	}

	return exec.Command(binary, args...), nil
}
