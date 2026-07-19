package mutant

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/internal/mutations"
	"github.com/SecretSheppy/marv/pkg/fio"
	"github.com/SecretSheppy/marv/pkg/pathutil"
	"github.com/aymanbagabas/go-udiff"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var meta = fwlib.Meta{
	Name: "mutant",
	URL:  "https://github.com/mbj/mutant",
}

type YamlConfig struct {
	RootDir    string `yaml:"root-dir"`
	ResultsDir string `yaml:"results-dir"`
	Session    string `yaml:"results-session,omitempty"`
}

type YamlWrapper struct {
	Cfg *YamlConfig `yaml:"mutant"`
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
	return y.Cfg.RootDir != "" && y.Cfg.ResultsDir != "", nil
}

type Exception struct{}

type ProcessStatus struct {
	ExitStatus int `json:"exitstatus"`
}

type Timeout struct{}

type Value struct {
	Passed bool `json:"passed"`
}

// IsolationResult uses the presence of an Exception or Timeout to determine whether the status should be
// CRASHED or TIMEOUT.
type IsolationResult struct {
	Exception     *Exception     `json:"exception"`
	ProcessStatus *ProcessStatus `json:"process_status"`
	Timeout       *Timeout       `json:"timeout"`
	Value         *Value         `json:"value"`
}

type MutationResult struct {
	IsolationResult        *IsolationResult `json:"isolation_result"`
	MutationSource         string           `json:"mutation_source"`
	MutationType           string           `json:"mutation_type"`
	MutationIdentification string           `json:"mutation_identification"`
}

func (m *MutationResult) getMarvStatus() mutations.Status {
	if m.IsolationResult.Timeout != nil {
		return mutations.Timeout
	}
	if m.IsolationResult.Exception != nil || m.IsolationResult.ProcessStatus.ExitStatus != 0 {
		return mutations.Crashed
	}
	if m.IsolationResult.Value.Passed {
		return mutations.Survived
	}
	return mutations.Killed
}

func (m *MutationResult) toMarvMutation(lines []string, originalCode string, startLine int) (*mutations.Mutation, error) {
	edits := udiff.Strings(originalCode, m.MutationSource)
	diff, err := udiff.ToUnifiedDiff("old", "new", originalCode, edits, 0)
	if err != nil {
		return nil, err
	}
	var (
		endLine     = startLine
		endChar     int
		replacement string
	)
	// NOTE: there should only ever be one hunk here
	for _, hunk := range diff.Hunks {
		startLine += hunk.FromLine - 2
		endLine += hunk.FromLine - 3

		trimmed := strings.TrimSpace(lines[startLine])
		padding := len(lines[startLine]) - len(trimmed)

		for i, line := range hunk.Lines {
			switch line.Kind {
			case udiff.Equal:
				endLine++
			case udiff.Delete:
				endLine++
				endChar = len(lines[startLine+i])
			case udiff.Insert:
				replacement += strings.Repeat(" ", padding) + strings.TrimSpace(line.Content)
			}
		}
	}

	desc := fmt.Sprintf("Replaced ```ruby\n%s\n``` with ```ruby\n%s\n```",
		strings.Join(lines[startLine:endLine+1], "\n"),
		strings.TrimSpace(replacement))
	return &mutations.Mutation{
		ID:                uuid.New(),
		FrameworkMutantID: m.MutationIdentification,
		Description:       desc,
		Operation:         mutations.UnrecoverableOperator,
		Start: &mutations.Range{
			Line: startLine,
			Char: 0,
		},
		End: &mutations.Range{
			Line: endLine,
			Char: endChar,
		},
		Status:      m.getMarvStatus(),
		Replacement: replacement,
	}, nil
}

type CoverageResult struct {
	MutationResult *MutationResult `json:"mutation_result"`
}

type SubjectResult struct {
	AmountMutations int               `json:"amount_mutations"`
	CoverageResults []*CoverageResult `json:"coverage_results"`
	Identification  string            `json:"identification"`
	Source          string            `json:"source"`
	SourcePath      string            `json:"source_path"`
}

func (s *SubjectResult) startLine() (int64, error) {
	split := strings.Split(s.Identification, ":")
	sln := split[len(split)-1]
	return strconv.ParseInt(sln, 10, 32)
}

type Results struct {
	SubjectResults []*SubjectResult `json:"subject_results"`
}

func (r *Results) mutationsCount() int {
	count := 0
	for _, sr := range r.SubjectResults {
		count += sr.AmountMutations
	}
	return count
}

type Mutant struct {
	yml     *YamlWrapper
	results Results
	ms      mutations.Mutations
	files   map[string][]string
}

func NewMutant() *Mutant {
	return &Mutant{yml: &YamlWrapper{}}
}

func (m *Mutant) Meta() *fwlib.Meta {
	return &meta
}

func (m *Mutant) Yaml() fwlib.FWConfig {
	return m.yml
}

// If the session name is specified then that particular session JSON is returned. Otherwise, the most recently created
// JSON with a UUID name is returned.
func (m *Mutant) getResultsPath() (string, error) {
	if m.yml.Cfg.Session != "" {
		return path.Join(m.yml.Cfg.ResultsDir, m.yml.Cfg.Session+".json"), nil
	}
	entries, err := os.ReadDir(m.yml.Cfg.ResultsDir)
	if err != nil {
		return "", err
	}
	var (
		rtime time.Time
		rfile string
		uid   uuid.UUID
	)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, _ := entry.Info()
		fext := path.Ext(info.Name())
		fnme := strings.TrimSuffix(info.Name(), fext)
		uid, err = uuid.Parse(fnme)
		if info.ModTime().After(rtime) && fext == ".json" && err == nil && uid != uuid.Nil {
			rtime = info.ModTime()
			rfile = info.Name()
		}
	}
	return path.Join(m.yml.Cfg.ResultsDir, rfile), nil
}

func (m *Mutant) LoadResults() error {
	log.Info().Msgf("%s - loading results", m.Meta().Name)

	file, err := m.getResultsPath()
	if err != nil {
		return err
	}
	log.Info().Msgf("%s - loading %s", m.Meta().Name, file)

	raw, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, &m.results)
}

func (m *Mutant) TransformResults() error {
	log.Info().Msgf("%s - transforming results", m.Meta().Name)

	wd, _ := os.Getwd()
	root := path.Base(path.Join(wd, m.yml.Cfg.RootDir))
	var absPathPrefix string
	if len(m.results.SubjectResults) > 0 {
		for _, dir := range pathutil.Split(m.results.SubjectResults[0].SourcePath) {
			absPathPrefix += "/" + dir
			if dir == root {
				absPathPrefix += "/"
				break
			}
		}
	}
	log.Info().Msgf("%s - removing absolute path %s", m.Meta().Name, absPathPrefix)

	bar := fwlib.NewProgressbar(m.results.mutationsCount(), "transforming")
	m.ms = make(mutations.Mutations)
	m.files = make(map[string][]string)

	for _, sr := range m.results.SubjectResults {
		sourcePath := strings.TrimPrefix(sr.SourcePath, absPathPrefix)
		lines, err := fio.ReadLines(path.Join(m.yml.Cfg.RootDir, sourcePath))
		if err != nil {
			return err
		}
		m.files[sourcePath] = lines

		for _, cr := range sr.CoverageResults {
			// NOTE: Marv does not support rendering original unmutated source code inside a mutation block. If the
			// replacement field is blank then Marv marks the last line of the block as deleted. This is not
			// deliberate behavior, but it causes no issues in any other frameworks and has not been fixed for
			// this framework. Neutral "mutations" are therefore be ignored.
			if cr.MutationResult.MutationType == "neutral" {
				continue
			}
			sl, err := sr.startLine()
			if err != nil {
				return err
			}
			mu, err := cr.MutationResult.toMarvMutation(lines, sr.Source, int(sl))
			if err != nil {
				return err
			}
			m.ms.Append(sourcePath, mu)
			bar.Add(1)
		}
	}

	fwlib.FinishProgressbar(bar)
	return nil
}

func (m *Mutant) Mutations() mutations.Mutations {
	return m.ms
}

func (m *Mutant) ReadLines(file string) ([]string, error) {
	return m.files[file], nil
}
