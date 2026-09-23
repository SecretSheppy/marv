package vineflower_server

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"regexp"
	"strings"
	"sync"
	"syscall"

	"github.com/SecretSheppy/marv/decompilers/dcomplib"
	"github.com/rs/zerolog/log"
)

// NOTE: use port 8081 so that the main server can start on 8080 whilst this process shuts down.
const port = 8081

var (
	re = regexp.MustCompile("vineflower-server(-\\d+\\.\\d+\\.\\d+)?\\.jar")

	errVineflowerServerNotFound = errors.New("vineflower server not found")
)

// VFServer is a decompiler that utilizes the Vineflower decompiler but calls it through http requests to the
// vineflower-server process that I wrote for Marv. This process is slower than Garlic, but not by that much. Due to
// its great compatibility, this is the default Java decompiler that Marv will use.
type VFServer struct {
	cmd    *exec.Cmd
	ctx    context.Context
	cancel context.CancelFunc
	exe    string
}

func (v *VFServer) exePath() error {
	libPath := dcomplib.ExeBasePath()
	entries, err := os.ReadDir(libPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if re.MatchString(entry.Name()) {
			v.exe = path.Join(libPath, entry.Name())
			return nil
		}
	}

	return errVineflowerServerNotFound
}

func (v *VFServer) ExePath() string {
	return v.exe
}

func (v *VFServer) Setup() error {
	var wg sync.WaitGroup
	wg.Add(1)
	if err := v.exePath(); err != nil {
		if errors.Is(err, errVineflowerServerNotFound) {
			log.Fatal().Err(err).Msg("decompiler 'vineflower server' not found in MARV_LIB_PATH or working directory")
		}
		return err
	}

	v.cmd = exec.Command("java", "-jar", v.exe, fmt.Sprintf("--server.port=%d", port))
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

	v.ctx, v.cancel = context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			return
		case <-sigs:
			if err := v.Teardown(); err != nil {
				log.Error().Err(err).Msgf("Failed to kill subprocess %s", v.exe)
			}
			return
		}
	}(v.ctx)

	return nil
}

func (v *VFServer) Teardown() error {
	if v.cmd != nil {
		return v.cmd.Process.Kill()
	}
	v.cancel()
	return nil
}

type VFServerErrorResponse struct {
	Message string `json:"message"`
}

type VFServerResponse struct {
	Output map[string]string `json:"output"`
}

func (v *VFServer) Decompile(p string) ([]byte, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/vineflower?source=%s", port, p))
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
	return v.exe
}
