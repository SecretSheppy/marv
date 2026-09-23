package dcomplib

import (
	"os"

	"github.com/rs/zerolog/log"
)

// ExeBasePath returns either the MARV_LIB_PATH environment variable or the cwd joined with "lib".
func ExeBasePath() string {
	dir := os.Getenv("MARV_LIB_PATH")
	if dir == "" {
		dir, _ = os.Getwd()
		log.Warn().Str("pwd", dir).Msgf("MARV_LIB_PATH environment variable not set, using working directory")
	}
	return dir
}
