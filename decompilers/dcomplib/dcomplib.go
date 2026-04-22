package dcomplib

import (
	"os"
)

// ExeBasePath returns either the MARV_LIB_PATH environment variable or the cwd joined with "lib".
func ExeBasePath() string {
	dir := os.Getenv("MARV_LIB_PATH")
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	}
	return dir
}
