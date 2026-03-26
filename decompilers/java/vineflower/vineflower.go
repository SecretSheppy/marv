package vineflower

import (
	"os/exec"
	"path"

	"github.com/SecretSheppy/marv/decompilers/dcomplib"
	"github.com/rs/zerolog/log"
)

// Vineflower is a decompiler that directly calls the vineflower.jar and reads the decompiled code from stdout. This
// process is very slow, but it does work well and was the first decompiler method implemented.
//
// Deprecated: VFServer is now the intended Vineflower decompiler for Marv. This still works but is very, very slow.
type Vineflower struct{}

func (v *Vineflower) ExePath() string {
	return path.Join(dcomplib.ExeBasePath(), "vineflower.jar")
}

func (v *Vineflower) Setup() error {
	log.Warn().Msg("vineflower (standalone) decompiler support deprecated, decompiling will be slow.")
	return nil
}

func (v *Vineflower) Teardown() error { return nil }

func (v *Vineflower) Decompile(p string) ([]byte, error) {
	return exec.Command("java", "-jar", v.ExePath(), "--log-level=error", p).Output()
}

func (v *Vineflower) String() string {
	return v.ExePath()
}
