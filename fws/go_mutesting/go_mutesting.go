package go_mutesting

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/internal/languages"
	"github.com/SecretSheppy/marv/internal/mutations"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
	"gopkg.in/yaml.v3"
)

var meta = fwlib.Meta{
	Name:     "go-mutesting",
	Language: languages.Go,
	URL:      "https://github.com/zimmski/go-mutesting",
}

type YamlConfig struct {
	Json string `yaml:"json"`
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
	return y.Cfg.Json != "", nil
}

func (y *YamlWrapper) SourceCodeDir() string {
	return ""
}

type Mutator struct {
	MutatorName        string `json:"mutatorName"`
	OriginalSourceCode string `json:"originalSourceCode"`
	MutatedSourceCode  string `json:"mutatedSourceCode"`
	OriginalFilePath   string `json:"originalFilePath"`
	OriginalStartLine  int    `json:"originalStartLine"`
}

type Mutation struct {
	Mutator Mutator `json:"mutator"`
	Diff    string  `json:"diff"`
}

type Statuses struct {
	Escaped   []Mutation `json:"escaped"`
	Timeouted []Mutation `json:"timeouted"`
	Killed    []Mutation `json:"killed"`
	Errored   []Mutation `json:"errored"`
}

type GoMutesting struct {
	yml    *YamlWrapper
	report Statuses
	files  map[string]string
	ms     mutations.Mutations
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

	report, err := os.ReadFile(g.yml.Cfg.Json)
	if err != nil {
		return err
	}
	return json.Unmarshal(report, &g.report)
}

func (g *GoMutesting) TransformResults() error {
	log.Info().Msgf("%s - transforming results", g.Meta().Name)

	bar := progressbar.NewOptions(
		len(g.report.Escaped)+len(g.report.Timeouted)+len(g.report.Killed)+len(g.report.Escaped),
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionSetDescription("transforming"),
		progressbar.OptionSetRenderBlankState(true))

	g.ms = make(mutations.Mutations)
	g.files = make(map[string]string)
	if err := g.transformResults(g.report.Escaped, mutations.Survived, bar); err != nil {
		return err
	}
	if err := g.transformResults(g.report.Timeouted, mutations.Timeout, bar); err != nil {
		return err
	}
	if err := g.transformResults(g.report.Killed, mutations.Killed, bar); err != nil {
		return err
	}
	if err := g.transformResults(g.report.Errored, mutations.Crashed, bar); err != nil {
		return err
	}

	bar.Finish()
	fmt.Println()
	return nil
}

func (g *GoMutesting) transformResults(ms []Mutation, status mutations.Status, bar *progressbar.ProgressBar) error {
	for _, mutation := range ms {
		mutator := mutation.Mutator

		if g.files[mutator.OriginalFilePath] == "" {
			g.files[mutator.OriginalFilePath] = mutator.OriginalSourceCode
		}

		var (
			diff              = strings.TrimPrefix(mutation.Diff, "--- Original\n+++ New\n")
			replacement       strings.Builder
			deletedLinesCount int
			lenLastDeleted    int
		)
		for line := range strings.Lines(diff) {
			switch true {
			case strings.HasPrefix(line, "-"):
				deletedLinesCount++
				lenLastDeleted = len(line) - 2 // NOTE: discounts the "-" prefix and "\n" suffix
			case strings.HasPrefix(line, "+"):
				replacement.WriteString(strings.TrimPrefix(line, "+"))
			}
		}

		sl := &mutations.Mutation{
			Description: mutator.MutatorName,
			Operation:   mutator.MutatorName,
			Start: &mutations.Range{
				Line: mutator.OriginalStartLine - 1,
				Char: 0,
			},
			End: &mutations.Range{
				Line: (mutator.OriginalStartLine - 1) + (deletedLinesCount - 1),
				Char: lenLastDeleted,
			},
			Status:      status,
			Replacement: replacement.String(),
		}

		g.ms.Append(mutator.OriginalFilePath, sl)
		bar.Add(1)
	}
	return nil
}

func (g *GoMutesting) Mutations() mutations.Mutations {
	return g.ms
}

func (g *GoMutesting) ReadLines(file string) ([]string, error) {
	lines := make([]string, 0)
	for line := range strings.Lines(g.files[file]) {
		lines = append(lines, strings.ReplaceAll(line, "\n", ""))
	}
	return lines, nil
}
