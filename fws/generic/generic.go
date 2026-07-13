package generic

import (
	"encoding/json"
	"os"
	"path"
	"strings"

	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/internal/mutations"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var meta = &fwlib.Meta{
	Name: "generic",
	URL:  "https://github.com/SecretSheppy/marv",
}

type YamlConfig struct {
	FWName   string `yaml:"framework"`
	MarvJson string `yaml:"marv-json"`
	SrcDir   string `yaml:"src-dir"`
}

type YamlWrapper struct {
	Cfg *YamlConfig `yaml:"generic"`
}

func (y *YamlWrapper) Init() interface{} {
	return &YamlWrapper{Cfg: &YamlConfig{}}
}

func (y *YamlWrapper) Load(yml []byte) (bool, error) {
	if err := yaml.Unmarshal(yml, y); err != nil {
		return false, err
	}
	if y.Cfg == nil {
		return false, nil
	}
	if y.Cfg.FWName != "" {
		meta.Name = y.Cfg.FWName
	}
	return y.Cfg.MarvJson != "" && y.Cfg.SrcDir != "", nil
}

type Generic struct {
	yml   *YamlWrapper
	ms    mutations.Mutations
	files map[string][]string
}

func NewGeneric() *Generic {
	return &Generic{yml: &YamlWrapper{}}
}

func (g *Generic) Meta() *fwlib.Meta {
	return meta
}

func (g *Generic) Yaml() fwlib.FWConfig {
	return g.yml
}

func (g *Generic) LoadResults() error {
	log.Info().Msgf("%s - loading results", g.Meta().Name)

	var (
		data []byte
		err  error
	)
	data, err = os.ReadFile(g.yml.Cfg.MarvJson)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, &g.ms); err != nil {
		return err
	}
	g.files = make(map[string][]string)
	for file, _ := range g.ms {
		data, err = os.ReadFile(path.Join(g.yml.Cfg.SrcDir, file))
		if err != nil {
			return err
		}
		lines := make([]string, 0)
		for line := range strings.Lines(string(data)) {
			lines = append(lines, strings.ReplaceAll(line, "\n", ""))
		}
		g.files[file] = lines
	}
	return nil
}

func (g *Generic) TransformResults() error {
	return nil
}

func (g *Generic) Mutations() mutations.Mutations {
	return g.ms
}

func (g *Generic) ReadLines(file string) ([]string, error) {
	return g.files[file], nil
}
