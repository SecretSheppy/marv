package dcomplib

import (
	"os"
	"path"
)

// ExeBasePath returns either the LIB_PATH environment variable or the cwd joined with "lib".
func ExeBasePath() string {
	dir := os.Getenv("LIB_PATH")
	if dir == "" {
		wd, _ := os.Getwd()
		dir = path.Join(wd, "lib")
	}
	return dir
}
