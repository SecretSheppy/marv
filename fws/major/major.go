package major

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/internal/mutations"
	"github.com/SecretSheppy/marv/pkg/fio"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var meta = fwlib.Meta{
	Name: "Major",
	URL:  "https://mutation-testing.org/",
}

// NOTE: the primary type formatting for the original value i.e., ==(int,int)
var re = regexp.MustCompile(".+\\([a-zA-Z,_]+\\)")

type YamlConfig struct {
	SrcDir    string `yaml:"src-dir"`
	OutputDir string `yaml:"output-dir"`
}

type YamlWrapper struct {
	Cfg *YamlConfig `yaml:"major"`
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
	return y.Cfg.SrcDir != "" && y.Cfg.OutputDir != "", nil
}

// Mutant describes a line in the mutants.log file
type Mutant struct {
	ID           int
	Operator     string
	OriginalType string
	NewType      string
	ClassPath    string
	LineNumber   int
	CharNumber   int
	Original     string
	Replacement  string
}

func (m *Mutant) file() string {
	return strings.ReplaceAll(m.ClassPath, ".", "/") + ".java"
}

func getType(t string) string {
	if re.MatchString(t) {
		return strings.Split(t, "(")[0]
	}
	return t
}

func marshalMutant(mutant string) (*Mutant, error) {
	mutant = strings.ReplaceAll(mutant, "\n", "")
	split := strings.Split(mutant, ":")

	id, err := strconv.ParseInt(split[0], 10, 32)
	if err != nil {
		return nil, err
	}

	class := strings.Split(split[4], "@")

	line, err := strconv.ParseInt(split[5], 10, 32)
	if err != nil {
		return nil, err
	}

	char, err := strconv.ParseInt(split[6], 10, 32)
	if err != nil {
		return nil, err
	}

	code := strings.Split(split[7], " |==> ")

	return &Mutant{
		ID:           int(id),
		Operator:     split[1],
		OriginalType: getType(split[2]),
		NewType:      getType(split[3]),
		ClassPath:    class[0],
		LineNumber:   int(line),
		CharNumber:   int(char),
		Original:     code[0],
		Replacement:  code[1],
	}, nil
}

func ReadMutants(file string) ([]*Mutant, error) {
	raw, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	ms := make([]*Mutant, 0)
	for line := range bytes.Lines(raw) {
		if bytes.Equal(line, []byte("")) {
			continue
		}
		m, err := marshalMutant(string(line))
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}
	return ms, nil
}

// Detail describes a line in the details.csv file
type Detail struct {
	ID     int
	Status string
}

func (d *Detail) status() mutations.Status {
	switch d.Status {
	case "FAIL":
		return mutations.Killed
	case "TIME":
		return mutations.Timeout
	case "EXC":
		return mutations.Crashed
	case "LIVE":
		return mutations.Survived
	default:
		return mutations.NoCoverage
	}
}

func marshalDetail(detail string) (*Detail, error) {
	split := strings.Split(detail, ",")
	id, err := strconv.ParseInt(split[0], 10, 32)
	if err != nil {
		return nil, err
	}
	return &Detail{ID: int(id), Status: split[1]}, nil
}

func ReadDetails(file string) ([]*Detail, error) {
	raw, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	ds := make([]*Detail, 0)
	// NOTE: the first line is the csv column titles
	for i, line := range strings.Split(string(raw), "\n") {
		if i == 0 || line == "" {
			continue
		}
		d, err := marshalDetail(line)
		if err != nil {
			return nil, err
		}
		ds = append(ds, d)
	}
	return ds, nil
}

type Major struct {
	yml     *YamlWrapper
	mutants []*Mutant
	details []*Detail
	ms      mutations.Mutations
}

func NewMajor() *Major {
	return &Major{yml: &YamlWrapper{}}
}

func (m *Major) Meta() *fwlib.Meta {
	return &meta
}

func (m *Major) Yaml() fwlib.FWConfig {
	return m.yml
}

func (m *Major) LoadResults() error {
	log.Info().Msgf("%s - loading results", m.Meta().Name)

	var err error
	m.mutants, err = ReadMutants(path.Join(m.yml.Cfg.OutputDir, "mutants.log"))
	if err != nil {
		return err
	}
	m.details, err = ReadDetails(path.Join(m.yml.Cfg.OutputDir, "details.csv"))
	if err != nil {
		return err
	}
	return nil
}

func (m *Major) TransformResults() error {
	log.Info().Msgf("%s - transforming results", m.Meta().Name)
	bar := fwlib.NewProgressbar(len(m.mutants), "transforming")

	m.ms = make(mutations.Mutations)
	for _, mutant := range m.mutants {
		m.ms.Append(mutant.file(), &mutations.Mutation{
			ID:                uuid.New(),
			FrameworkMutantID: strconv.Itoa(mutant.ID),
			Description:       fmt.Sprintf("Replaced `%s` with `%s`", mutant.OriginalType, mutant.NewType),
			Operation:         mutant.Operator,
			Start: &mutations.Range{
				Line: mutant.LineNumber - 1,
				Char: 0, // TODO:
			},
			End: &mutations.Range{
				Line: mutant.LineNumber - 1,
				Char: 1, // TODO:
			},
			Status:      m.details[mutant.ID-1].status(),
			Replacement: mutant.Replacement,
		})
	}

	fwlib.FinishProgressbar(bar)
	return nil
}

func (m *Major) Mutations() mutations.Mutations {
	return m.ms
}

func (m *Major) ReadLines(file string) ([]string, error) {
	return fio.ReadLines(path.Join(m.yml.Cfg.SrcDir, file))
}
