package vineflower

import (
	"os"
	"os/exec"
	"path"
)

func vineflowerJarPath() string {
	dir := os.Getenv("LIB_PATH")
	if dir == "" {
		wd, _ := os.Getwd()
		dir = path.Join(wd, "lib")
	}
	return path.Join(dir, "vineflower.jar")
}

func Help() ([]byte, error) {
	cmd := exec.Command("java", "-jar", vineflowerJarPath(), "--help")
	cmd.Env = os.Environ()
	return cmd.Output()
}

func Decompile(path string) ([]byte, error) {
	cmd := exec.Command("java", "-jar", vineflowerJarPath(), "--log-level=error", path)
	cmd.Env = os.Environ()
	return cmd.Output()
}
