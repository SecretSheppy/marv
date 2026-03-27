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

	"github.com/SecretSheppy/marv/decompilers/dcomplib"
)

// VFServer is a decompiler that utilizes the Vineflower decompiler but calls it through http requests to the
// vineflower-server process that I wrote for Marv. This process is slower than Garlic, but not by that much. Due to
// its great compatibility, this is the default Java decompiler that Marv will use.
type VFServer struct {
	cmd *exec.Cmd
}

func (v *VFServer) ExePath() string {
	return path.Join(dcomplib.ExeBasePath(), "vineflower-server.jar")
}

func (v *VFServer) Setup() error {
	var wg sync.WaitGroup
	wg.Add(1)

	v.cmd = exec.Command("java", "-jar", v.ExePath())
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

type VFServerErrorResponse struct {
	Message string `json:"message"`
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
		return nil, errors.New(fmt.Sprintf("%s returned status %d", v.ExePath(), resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	errRes := &VFServerErrorResponse{}
	if err := json.Unmarshal(body, errRes); err != nil {
		return nil, err
	}
	if errRes.Message != "" {
		return nil, errors.New(errRes.Message)
	}

	res := &VFServerResponse{}
	if err := json.Unmarshal(body, res); err != nil {
		return nil, err
	}

	var val string
	for _, val = range res.Output {
		break
	}

	return []byte(val), nil
}

func (v *VFServer) String() string {
	return v.ExePath()
}
