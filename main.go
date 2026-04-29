package main

import (
	"os"

	"github.com/SecretSheppy/marv/internal/cmds"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	cmds.Execute()
}
