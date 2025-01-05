package expectation

import (
	"io/fs"
)

func GetExpectation(path string, fsys fs.FS) string {
	bytes, err := fs.ReadFile(fsys, path+"/expected.txt")

	if err != nil {
		return ""
	}

	return string(bytes)
}
