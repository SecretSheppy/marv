package stryker_net

import (
	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/internal/languages"
	"github.com/SecretSheppy/marv/internal/mtelib"
	"github.com/SecretSheppy/marv/internal/mutations"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var meta = fwlib.Meta{
	Name:     "stryker-net",
	Language: languages.CSharp,
	URL:      "https://github.com/stryker-mutator/stryker-net",
}

type YamlConfig struct {
	MTEJson string `yaml:"mte-json"`
}

type YamlWrapper struct {
	Cfg *YamlConfig `yaml:"stryker-net"`
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

type StrykerNet struct {
	yml *YamlWrapper
	mte *mtelib.MTE
}

func NewStrykerNet() *StrykerNet {
	return &StrykerNet{yml: &YamlWrapper{}}
}

func (s *StrykerNet) Meta() *fwlib.Meta {
	return &meta
}

func (s *StrykerNet) Yaml() fwlib.FWConfig {
	return s.yml
}

func (s *StrykerNet) LoadResults() error {
	log.Info().Msgf("%s - loading results", s.Meta().Name)
	var err error
	s.mte, err = mtelib.NewMTE(s.yml.Cfg.MTEJson)
	return err
}

func (s *StrykerNet) TransformResults() error {
	log.Info().Msgf("%s - transforming results", s.Meta().Name)

	bar := fwlib.NewProgressbar(s.mte.RawMutationsCount(), "transforming")
	s.mte.Transform(bar)
	fwlib.FinishProgressbar(bar)

	return nil
}

func (s *StrykerNet) Mutations() mutations.Mutations {
	return s.mte.Mutations()
}

func (s *StrykerNet) ReadLines(file string) ([]string, error) {
	return s.mte.ReadLines(file), nil
}
