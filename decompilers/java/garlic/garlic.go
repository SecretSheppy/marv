package garlic

import (
	"os/exec"
	"path"

	"github.com/SecretSheppy/marv/decompilers/dcomplib"
	"github.com/rs/zerolog/log"
)

// Garlic is a decompiler that directly calls the garlic executable and reads the decompiled code from stdout. This is
// the fastest decompiler that marv can currently use, and it produces very high quality results.
//
// Compatibility: The Garlic Java decompiler supports windows, but I have not managed to get the garlic.exe
// binary to work correctly.
type Garlic struct{}

func (g *Garlic) ExePath() string {
	return path.Join(dcomplib.ExeBasePath(), "garlic")
}

func (g *Garlic) Setup() error {
	log.Warn().Msgf("garlic decompiler is unstable, using it could cause mutants to be skipped")
	return nil
}

func (g *Garlic) Teardown() error { return nil }

func (g *Garlic) Decompile(p string) ([]byte, error) {
	return exec.Command(g.ExePath(), p).Output()
}

func (g *Garlic) String() string {
	return g.ExePath()
}
