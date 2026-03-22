package mutest_rs

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/pkg/mutations"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
	"gopkg.in/yaml.v3"
)

// YamlConfig represents Mutest-RS yml config data.
type YamlConfig struct {
	Run     string `yaml:"run"`
	Src     string `yaml:"src"`
	JsonDir string `yaml:"json-dir"`
}

// YamlWrapper used to load the mutest-rs configuration from the .marv.yml file.
type YamlWrapper struct {
	Cfg *YamlConfig `yaml:"mutest-rs"`
}

func (y *YamlWrapper) Init() interface{} {
	return &YamlWrapper{Cfg: &YamlConfig{}}
}

func (y *YamlWrapper) Load(yml []byte) (bool, error) {
	if err := yaml.Unmarshal(yml, y); err != nil {
		return false, err
	}
	return y.Cfg.Src != "" || y.Cfg.JsonDir != "", nil
}

func (y *YamlWrapper) SourceCodeDir() string {
	return y.Cfg.Src
}

// Evaluation marshals to evaluation.json from the mutest output data
type Evaluation struct {
	MutationRuns []*MutationRuns `json:"mutation_runs"`
}

type MutationRuns struct {
	DetectionMatrix *DetectionMatrix `json:"mutation_detection_matrix"`
}

type DetectionMatrix struct {
	OverallDetections string `json:"overall_detections"`
}

// Mutations marshals to mutations.json from the mutest output data
type Mutations struct {
	Mutations []*Mutation `json:"mutations"`
}

type Mutation struct {
	MutationID    int             `json:"mutation_id"`
	Location      *Location       `json:"origin_span"`
	MutationOp    string          `json:"mutation_op"`
	DisplayName   string          `json:"display_name"`
	Substitutions []*Substitution `json:"substs"`
}

type Substitution struct {
	Substitution *Substitute `json:"substitute"`
}

type Substitute struct {
	Replacement string `json:"replacement"`
}

type Location struct {
	Path  string `json:"path"`
	Begin []int  `json:"begin"`
	End   []int  `json:"end"`
}

// MutestRS wraps the evaluation.json and mutations.json objects into a single struct.
type MutestRS struct {
	yml  *YamlWrapper
	eval *Evaluation
	muts *Mutations
	ms   mutations.Mutations
}

func NewMutestRS() *MutestRS {
	return &MutestRS{yml: &YamlWrapper{}}
}

func (m *MutestRS) Meta() *fwlib.Meta {
	return &fwlib.Meta{
		Name:      "mutest-rs",
		Extension: "rs",
		URL:       "https://github.com/zalanlevai/mutest-rs",
	}
}

func (m *MutestRS) Yaml() fwlib.FWConfig {
	return m.yml
}

func (m *MutestRS) LoadResults() error {
	log.Info().Msgf("%s - loading results", m.Meta().Name)

	eval, err := os.ReadFile(path.Join(m.yml.Cfg.JsonDir, "evaluation.json"))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(eval, &m.eval); err != nil {
		return err
	}

	muts, err := os.ReadFile(path.Join(m.yml.Cfg.JsonDir, "mutations.json"))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(muts, &m.muts); err != nil {
		return err
	}

	return nil
}

func (m *MutestRS) TransformResults() error {
	log.Info().Msgf("%s - transforming results", m.Meta().Name)

	bar := progressbar.NewOptions(
		len(m.muts.Mutations),
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionSetDescription("transforming"),
		progressbar.OptionSetRenderBlankState(true))

	ms, err := m.transformResults(bar)
	if err != nil {
		return err
	}
	bar.Finish()
	fmt.Println()

	m.ms = ms
	return nil
}

func (m *MutestRS) transformResults(bar *progressbar.ProgressBar) (mutations.Mutations, error) {
	ms := mutations.Mutations{}
	for _, mu := range m.muts.Mutations {
		sl, err := streamlineMutation(mu)
		if err != nil {
			return nil, err
		}

		sl.Status, err = getMutationStatus(mu.MutationID, m.eval)
		if err != nil {
			return nil, err
		}

		added := false
		for _, c := range ms[mu.Location.Path] {
			if c.Conflicts(sl) {
				c.Append(sl)
				added = true
				break
			}
		}

		if !added {
			ms[mu.Location.Path] = append(ms[mu.Location.Path], mutations.NewConflict(sl))
		}

		bar.Add(1)
	}
	return ms, nil
}

func streamlineMutation(m *Mutation) (*mutations.Mutation, error) {
	if len(m.Location.Begin) != 2 {
		return nil, errors.New("plugin mutest_rs: Mutation.Location.Begin does not have two positions")
	}
	if len(m.Location.End) != 2 {
		return nil, errors.New("plugin mutest_rs: Mutation.Location.End does not have two positions")
	}
	if len(m.Substitutions) == 0 {
		return nil, errors.New("plugin mutest_rs: Mutation.Location.Substitutions is empty")
	}
	return &mutations.Mutation{
		ID:       m.MutationID - 1, // NOTE: mutest-rs mutation id start from 1, but they are used to index an array from 0
		IDOffset: 1,
		Name:     m.DisplayName,
		OpDesc:   m.MutationOp,
		Starts: &mutations.Range{
			Line: m.Location.Begin[0],
			Char: m.Location.Begin[1],
		},
		Ends: &mutations.Range{
			Line: m.Location.End[0],
			Char: m.Location.End[1],
		},
		Type:   mutations.Replacement, // for now all mutest does is replacement
		Source: m.Substitutions[0].Substitution.Replacement,
	}, nil
}

func getMutationStatus(id int, ev *Evaluation) (mutations.Status, error) {
	if len(ev.MutationRuns) == 0 {
		return "", errors.New("Mutation.Runs is empty")
	}
	status := mutations.Survived
	switch ev.MutationRuns[0].DetectionMatrix.OverallDetections[id-1] {
	case 'D':
		status = mutations.Killed
	case 'T':
		status = mutations.Timeout
	case 'C':
		status = mutations.Crashed
	}
	return status, nil
}

func (m *MutestRS) Mutations() mutations.Mutations {
	return m.ms
}
