package vineflower

import (
	"os"
	"os/exec"
	"path"
)

type Vineflower struct{}

func (v *Vineflower) JarPath() string {
	dir := os.Getenv("LIB_PATH")
	if dir == "" {
		wd, _ := os.Getwd()
		dir = path.Join(wd, "lib")
	}
	return path.Join(dir, "vineflower.jar")
}

func (v *Vineflower) Help() ([]byte, error) {
	cmd := exec.Command("java", "-jar", v.JarPath(), "--help")
	cmd.Env = os.Environ()
	return cmd.Output()
}

func (v *Vineflower) Decompile(path string) ([]byte, error) {
	cmd := exec.Command("java", "-jar", v.JarPath(), "--log-level=error", path)
	cmd.Env = os.Environ()
	return cmd.Output()
}
