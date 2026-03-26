package vineflower

import (
	"os"
	"os/exec"
	"path"

	"github.com/rs/zerolog/log"
)

// Deprecated: VFServer is now the intended Vineflower decompiler for Marv. This still works but is very, very slow.
type Vineflower struct{}

func (v *Vineflower) ExePath() string {
	dir := os.Getenv("LIB_PATH")
	if dir == "" {
		wd, _ := os.Getwd()
		dir = path.Join(wd, "lib")
	}
	return path.Join(dir, "vineflower.jar")
}

func (v *Vineflower) Setup() error {
	log.Warn().Msg("vineflower (standalone) decompiler support deprecated, decompiling will be slow.")
	return nil
}

func (v *Vineflower) Teardown() error { return nil }

func (v *Vineflower) Decompile(p string) ([]byte, error) {
	cmd := exec.Command("java", "-jar", v.ExePath(), "--log-level=error", p)
	cmd.Env = os.Environ()
	return cmd.Output()
}

func (v *Vineflower) String() string {
	return v.ExePath()
}
