package infection

import (
	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/internal/languages"
	"github.com/SecretSheppy/marv/internal/mtelib"
	"github.com/SecretSheppy/marv/internal/mutations"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var meta = fwlib.Meta{
	Name:     "infection",
	Language: languages.Php,
	URL:      "https://https://infection.github.io/",
}

type YamlConfig struct {
	MTEJson string `yaml:"mte-json"`
}

type YamlWrapper struct {
	Cfg *YamlConfig `yaml:"infection"`
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

type Infection struct {
	yml *YamlWrapper
	mte *mtelib.MTE
}

func NewInfection() *Infection {
	return &Infection{yml: &YamlWrapper{}}
}

func (i *Infection) Meta() *fwlib.Meta {
	return &meta
}

func (i *Infection) Yaml() fwlib.FWConfig {
	return i.yml
}

func (i *Infection) LoadResults() error {
	log.Info().Msgf("%s - loading results", i.Meta().Name)
	var err error
	i.mte, err = mtelib.NewMTE(i.yml.Cfg.MTEJson)
	return err
}

func (i *Infection) TransformResults() error {
	log.Info().Msgf("%s - transforming results", i.Meta().Name)

	bar := fwlib.NewProgressbar(i.mte.RawMutationsCount(), "transforming")
	i.mte.Transform(bar)
	fwlib.FinishProgressbar(bar)

	i.correctLineLengthOverhangs()
	return nil
}

// see infection/README.md as to why this method is necessary and for visual examples of it in practise.
func (i *Infection) correctLineLengthOverhangs() {
	for file, conflicts := range i.mte.Mutations() {
		lines := i.mte.ReadLines(file)
		for _, conflict := range conflicts {
			for _, mutation := range conflict.Mutations {
				endLineLength := len(lines[mutation.End.Line])
				if mutation.End.Char > endLineLength {
					mutation.End.Char = endLineLength
				}
			}
		}
	}
}

func (i *Infection) Mutations() mutations.Mutations {
	return i.mte.Mutations()
}

func (i *Infection) ReadLines(file string) ([]string, error) {
	return i.mte.ReadLines(file), nil
}
