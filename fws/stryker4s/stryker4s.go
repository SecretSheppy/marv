package stryker4s

import (
	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/internal/languages"
	"github.com/SecretSheppy/marv/internal/mtelib"
	"github.com/SecretSheppy/marv/internal/mutations"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var meta = fwlib.Meta{
	Name:     "stryker4s",
	Language: languages.Scala,
	URL:      "https://github.com/stryker-mutator/stryker4s",
}

type YamlConfig struct {
	MTEJson string `yaml:"mte-json"`
}

type YamlWrapper struct {
	Cfg *YamlConfig `yaml:"stryker4s"`
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

type Stryker4s struct {
	yml *YamlWrapper
	mte *mtelib.MTE
}

func NewStryker4s() *Stryker4s {
	return &Stryker4s{yml: &YamlWrapper{}}
}

func (s *Stryker4s) Meta() *fwlib.Meta {
	return &meta
}

func (s *Stryker4s) Yaml() fwlib.FWConfig {
	return s.yml
}

func (s *Stryker4s) LoadResults() error {
	log.Info().Msgf("%s - loading results", s.Meta().Name)
	var err error
	s.mte, err = mtelib.NewMTE(s.yml.Cfg.MTEJson)
	return err
}

func (s *Stryker4s) TransformResults() error {
	log.Info().Msgf("%s - transforming results", s.Meta().Name)

	bar := fwlib.NewProgressbar(s.mte.RawMutationsCount(), "transforming")
	s.mte.Transform(bar)
	fwlib.FinishProgressbar(bar)

	return nil
}

func (s *Stryker4s) Mutations() mutations.Mutations {
	return s.mte.Mutations()
}

func (s *Stryker4s) ReadLines(file string) ([]string, error) {
	return s.mte.ReadLines(file), nil
}
