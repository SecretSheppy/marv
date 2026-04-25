package go_mutesting

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/internal/languages"
	"github.com/SecretSheppy/marv/internal/mutations"
	"github.com/aymanbagabas/go-udiff"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var meta = fwlib.Meta{
	Name:     "go-mutesting",
	Language: languages.Go,
	URL:      "https://github.com/zimmski/go-mutesting", // NOTE: for actively maintained fork see https://github.com/avito-tech/go-mutesting
}

type YamlConfig struct {
	JsonReport string `yaml:"json-report"`
}

type YamlWrapper struct {
	Cfg *YamlConfig `yaml:"go-mutesting"`
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
	return y.Cfg.JsonReport != "", nil
}

type Report struct {
	Escaped   []Mutation `json:"escaped"`
	Timeouted []Mutation `json:"timeouted"`
	Killed    []Mutation `json:"killed"`
	Errored   []Mutation `json:"errored"`
}

type Mutation struct {
	Mutator Mutator `json:"mutator"`
}

type Mutator struct {
	MutatorName        string `json:"mutatorName"`
	OriginalSourceCode string `json:"originalSourceCode"`
	MutatedSourceCode  string `json:"mutatedSourceCode"`
	OriginalFilePath   string `json:"originalFilePath"`
	OriginalStartLine  int    `json:"originalStartLine"`
}

type GoMutesting struct {
	yml    *YamlWrapper
	report Report
	ms     mutations.Mutations
	files  map[string][]string
}

func NewGoMutesting() *GoMutesting {
	return &GoMutesting{yml: &YamlWrapper{}}
}

func (g *GoMutesting) Meta() *fwlib.Meta {
	return &meta
}

func (g *GoMutesting) Yaml() fwlib.FWConfig {
	return g.yml
}

func (g *GoMutesting) LoadResults() error {
	log.Info().Msgf("%s - loading results", g.Meta().Name)

	data, err := os.ReadFile(g.yml.Cfg.JsonReport)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &g.report)
}

func (g *GoMutesting) TransformResults() error {
	g.ms = make(mutations.Mutations)
	g.files = make(map[string][]string)
	if err := g.transformResults(g.report.Escaped, mutations.Survived); err != nil {
		return err
	}
	if err := g.transformResults(g.report.Timeouted, mutations.Timeout); err != nil {
		return err
	}
	if err := g.transformResults(g.report.Killed, mutations.Killed); err != nil {
		return err
	}
	return g.transformResults(g.report.Errored, mutations.Crashed)
}

func (g *GoMutesting) transformResults(ms []Mutation, status mutations.Status) error {
	for _, mutation := range ms {
		mutator := mutation.Mutator
		lines := g.addOrGetFile(mutator)

		edits := udiff.Strings(mutator.OriginalSourceCode, mutator.MutatedSourceCode)
		d, err := udiff.ToUnifiedDiff("old", "new", mutator.OriginalSourceCode, edits, 0)
		if err != nil {
			return err
		}

		var (
			removedLineCount, hunkStartLine int
			replacement                     strings.Builder
		)
		for _, h := range d.Hunks {
			hunkStartLine = h.FromLine

			for _, l := range h.Lines {
				switch l.Kind {
				case udiff.Delete:
					removedLineCount++
				case udiff.Insert, udiff.Equal:
					replacement.WriteString(l.Content)
				}
			}
		}

		startLine := mutator.OriginalStartLine - 1
		if startLine <= 0 || mutator.MutatorName == "loop/range_break" {
			startLine = hunkStartLine - 1
		}
		endLine := startLine + removedLineCount - 1

		m := &mutations.Mutation{
			Operation: mutator.MutatorName,
			Start: &mutations.Range{
				Line: startLine,
				Char: 0,
			},
			End: &mutations.Range{
				Line: endLine,
				Char: len(lines[endLine]),
			},
			Status:      status,
			Replacement: replacement.String(),
		}

		g.ms.Append(mutator.OriginalFilePath, m)
	}
	return nil
}

func (g *GoMutesting) addOrGetFile(mutator Mutator) []string {
	path := mutator.OriginalFilePath
	if g.files[path] == nil {
		lines := make([]string, 0)
		for line := range strings.Lines(mutator.OriginalSourceCode) {
			lines = append(lines, strings.ReplaceAll(line, "\n", ""))
		}
		g.files[path] = lines
	}
	return g.files[path]
}

func (g *GoMutesting) Mutations() mutations.Mutations {
	return g.ms
}

func (g *GoMutesting) ReadLines(file string) ([]string, error) {
	return g.files[file], nil
}
