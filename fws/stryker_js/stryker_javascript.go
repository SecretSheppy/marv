package stryker_js

import (
	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/internal/mtelib"
	"github.com/SecretSheppy/marv/internal/mutations"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var meta = fwlib.Meta{
	Name: "stryker-js",
	URL:  "https://github.com/stryker-mutator/stryker-net",
}

type YamlConfig struct {
	MTEJson string `yaml:"mte-json"`
}

type YamlWrapper struct {
	Cfg *YamlConfig `yaml:"stryker-js"`
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

type StrykerJS struct {
	yml *YamlWrapper
	mte *mtelib.MTE
}

func NewStrykerJS() *StrykerJS {
	return &StrykerJS{yml: &YamlWrapper{}}
}

func (s *StrykerJS) Meta() *fwlib.Meta {
	return &meta
}

func (s *StrykerJS) Yaml() fwlib.FWConfig {
	return s.yml
}

func (s *StrykerJS) LoadResults() error {
	log.Info().Msgf("%s - loading results", s.Meta().Name)
	var err error
	s.mte, err = mtelib.NewMTE(s.yml.Cfg.MTEJson)
	return err
}

func (s *StrykerJS) TransformResults() error {
	log.Info().Msgf("%s - transforming results", s.Meta().Name)

	bar := fwlib.NewProgressbar(s.mte.RawMutationsCount(), "transforming")
	s.mte.Transform(bar)
	fwlib.FinishProgressbar(bar)

	return nil
}

func (s *StrykerJS) Mutations() mutations.Mutations {
	return s.mte.Mutations()
}

func (s *StrykerJS) ReadLines(file string) ([]string, error) {
	return s.mte.ReadLines(file), nil
}
