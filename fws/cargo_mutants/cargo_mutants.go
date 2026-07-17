package cargo_mutants

import (
	"encoding/json"
	"os"
	"path"

	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/internal/mutations"
	"github.com/SecretSheppy/marv/pkg/fio"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var meta = fwlib.Meta{
	Name: "cargo-mutants",
	URL:  "https://github.com/sourcefrog/cargo-mutants",
}

type YamlConfig struct {
	TestWorkDir   string `yaml:"test-work-dir"`
	MutantsOutDir string `yaml:"mutants-out-dir"`
}

type YamlWrapper struct {
	Cfg *YamlConfig `yaml:"cargo-mutants"`
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
	return y.Cfg.TestWorkDir != "" && y.Cfg.MutantsOutDir != "", nil
}

type Pos struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type Span struct {
	Start *Pos `json:"start"`
	End   *Pos `json:"end"`
}

type Mutant struct {
	Name        string `json:"name"`
	File        string `json:"file"`
	Span        *Span  `json:"span"`
	Replacement string `json:"replacement"`
	Genre       string `json:"genre"`
}

type Scenario struct {
	Mutant *Mutant `json:"Mutant"`
}

type Outcome struct {
	Scenario *Scenario `json:"scenario"`
	Summary  string    `json:"summary"`
}

func (o *Outcome) GetMarvStatus() mutations.Status {
	switch o.Summary {
	case "CaughtMutant":
		return mutations.Killed
	case "Unviable":
		return mutations.Crashed
	case "Timeout":
		return mutations.Timeout
	default:
		return mutations.Survived
	}
}

func (o *Outcome) ToMarvMutation() *mutations.Mutation {
	return &mutations.Mutation{
		ID:          uuid.New(),
		Description: o.Scenario.Mutant.Name, // TODO: trim first section
		Operation:   o.Scenario.Mutant.Genre,
		Start: &mutations.Range{
			Line: o.Scenario.Mutant.Span.Start.Line - 1,
			Char: o.Scenario.Mutant.Span.Start.Column - 1,
		},
		End: &mutations.Range{
			Line: o.Scenario.Mutant.Span.End.Line - 1,
			Char: o.Scenario.Mutant.Span.End.Column - 1,
		},
		Status:      o.GetMarvStatus(),
		Replacement: o.Scenario.Mutant.Replacement,
	}
}

type Outcomes struct {
	Outcomes []*Outcome `json:"outcomes"`
}

type CargoMutants struct {
	yml      *YamlWrapper
	outcomes Outcomes
	ms       mutations.Mutations
	files    map[string][]string
}

func NewCargoMutants() *CargoMutants {
	return &CargoMutants{yml: &YamlWrapper{}}
}

func (c *CargoMutants) Meta() *fwlib.Meta {
	return &meta
}

func (c *CargoMutants) Yaml() fwlib.FWConfig {
	return c.yml
}

func (c *CargoMutants) LoadResults() error {
	log.Info().Msgf("%s - loading results", c.Meta().Name)

	outcomesJSON := path.Join(c.yml.Cfg.MutantsOutDir, "outcomes.json")
	data, err := os.ReadFile(outcomesJSON)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &c.outcomes)
}

func (c *CargoMutants) TransformResults() error {
	log.Info().Msgf("%s - transforming results", c.Meta().Name)
	bar := fwlib.NewProgressbar(len(c.outcomes.Outcomes), "transforming")
	c.ms = make(mutations.Mutations)
	c.files = make(map[string][]string)
	for _, outcome := range c.outcomes.Outcomes {
		file := outcome.Scenario.Mutant.File
		if c.files[file] == nil {
			relFile := path.Join(c.yml.Cfg.TestWorkDir, file)
			lines, err := fio.ReadLines(relFile)
			if err != nil {
				return err
			}
			c.files[file] = lines
		}
		c.ms.Append(file, outcome.ToMarvMutation())
		bar.Add(1)
	}
	fwlib.FinishProgressbar(bar)
	return nil
}

func (c *CargoMutants) Mutations() mutations.Mutations {
	return c.ms
}

func (c *CargoMutants) ReadLines(file string) ([]string, error) {
	return c.files[file], nil
}
