package garlic

import (
	"os"
	"os/exec"
	"path"

	"github.com/SecretSheppy/marv/decompilers/dcomplib"
)

type Garlic struct{}

func (g *Garlic) ExePath() string {
	return path.Join(dcomplib.ExeBasePath(), "garlic")
}

func (g *Garlic) Setup() error    { return nil }
func (g *Garlic) Teardown() error { return nil }

func (g *Garlic) Decompile(p string) ([]byte, error) {
	cmd := exec.Command(g.ExePath(), p)
	cmd.Env = os.Environ()
	return cmd.Output()
}

func (g *Garlic) String() string {
	return g.ExePath()
}
