package garlic

import (
	"os"
	"os/exec"
	"path"
)

type Garlic struct{}

func (g *Garlic) ExePath() string {
	dir := os.Getenv("LIB_PATH")
	if dir == "" {
		wd, _ := os.Getwd()
		dir = path.Join(wd, "lib")
	}
	return path.Join(dir, "garlic")
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
