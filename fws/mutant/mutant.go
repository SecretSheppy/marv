package mutant

import (
	"encoding/json"
	"os"
	"path"
	"strings"
	"time"

	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/internal/mutations"
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
	MutationIdentification string           `json:"mutation_identification"`
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

type Results struct {
	SubjectResults []*SubjectResult `json:"subject_results"`
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
	_ = fwlib.NewProgressbar(0, "transforming") // TODO

	// TODO: use m.yml.Cfg.RootDir to generate relative paths for all files. need to join it with the workdir and then
	//  take the last directory name and use that to strip out the absolute paths from the mutant results.
	//  - only need to generate the head path once, and can just trim prefix for all others
	return nil
}

func (m *Mutant) Mutations() mutations.Mutations {
	return m.ms
}

func (m *Mutant) ReadLines(file string) ([]string, error) {
	return m.files[file], nil
}

// TODO: Diff MutationResult.MutationSource and SubjectResult.Source to get the actual mutation and its lines etc...

// TODO: schema very useful:
//  https://github.com/mbj/mutant/blob/main/docs/session-json-schema.yml

// TODO: mutation_type can be either evil (not test killed it) or neutral (killed by a test) or noop

// TODO: will have to extract replacement lines beginning with + and deleted lines with -
// TODO: (contd) will have to diff the original and replacement to produce descriptions as well as actual replacements

// TODO: operators will have to be defined by marv and then determined based off of this list: (? maybe, this could be very difficult)
//  https://github.com/mbj/mutant/blob/59517844547eef3d67b71a3c736f05bb3c2376da/ruby/lib/mutant/mutation/operators.rb

// TODO: exit_status != 0 || exception == STATUS CRASHED

// TODO: if value:passed == STATUS SURVIVED else STATUS KILLED
