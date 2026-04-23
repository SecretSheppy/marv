package mull

import (
	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/internal/languages"
	"github.com/SecretSheppy/marv/internal/mutations"
	"github.com/SecretSheppy/marv/pkg/mtelib"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var meta = fwlib.Meta{
	Name:     "Mull",
	Language: languages.Cpp,
	URL:      "https://mull-project.com/",
}

type YamlConfig struct {
	MTEJson string `yaml:"mte-json"`
}

type YamlWrapper struct {
	Cfg *YamlConfig `yaml:"mull"`
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
	return y.Cfg.MTEJson != "", nil
}

type Mull struct {
	yml *YamlWrapper
	mte *mtelib.MTE
}

func NewMull() *Mull {
	return &Mull{yml: &YamlWrapper{}}
}

func (m *Mull) Meta() *fwlib.Meta {
	return &meta
}

func (m *Mull) Yaml() fwlib.FWConfig {
	return m.yml
}

func (m *Mull) LoadResults() error {
	log.Info().Msgf("%s - loading results", m.Meta().Name)
	var err error
	m.mte, err = mtelib.NewMTE(m.yml.Cfg.MTEJson)
	return err
}

func (m *Mull) TransformResults() error {
	log.Info().Msgf("%s - transforming results", m.Meta().Name)
	m.mte.Transform()
	return nil
}

func (m *Mull) Mutations() mutations.Mutations {
	return m.mte.Mutations()
}

func (m *Mull) ReadLines(file string) ([]string, error) {
	return m.mte.ReadLines(file), nil
}
