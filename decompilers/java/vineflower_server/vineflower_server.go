package vineflower_server

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
)

type VFServer struct {
	cmd *exec.Cmd
}

func (v *VFServer) JarPath() string {
	dir := os.Getenv("LIB_PATH")
	if dir == "" {
		wd, _ := os.Getwd()
		dir = path.Join(wd, "lib")
	}
	return path.Join(dir, "vineflower-server.jar")
}

func (v *VFServer) Setup() error {
	var wg sync.WaitGroup
	wg.Add(1)

	v.cmd = exec.Command("java", "-jar", v.JarPath())
	v.cmd.Env = os.Environ()
	stdout, err := v.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := v.cmd.Start(); err != nil {
		return err
	}

	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "Started VineflowerServer") {
				break
			}
		}
	}()

	wg.Wait()
	return nil
}

func (v *VFServer) Teardown() error {
	if v.cmd != nil {
		return v.cmd.Process.Kill()
	}
	return nil
}

type VFServerResponse struct {
	Output map[string]string `json:"output"`
}

func (v *VFServer) Decompile(p string) ([]byte, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/vineflower?source=%s", p))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("%s returned status %d", v.JarPath(), resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := &VFServerResponse{}
	if err := json.Unmarshal(body, res); err != nil {
		return nil, err
	}

	return []byte(res.Output[p]), nil
}

func (v *VFServer) String() string {
	return v.JarPath()
}
