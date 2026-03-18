package fws

import (
	"encoding/json"
	"errors"
	"os"
	"path"

	"github.com/SecretSheppy/marv/pkg/mutations"
)

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

var meta = &Meta{
	Name:   "mutest-rs",
	Lang:   "rs",
	URL:    "https://github.com/zalanlevai/mutest-rs",
	RunStr: "cargo mutest",
}

// MutestRS wraps the evaluation.json and mutations.json objects into a single struct.
type MutestRS struct {
	eval *Evaluation
	muts *Mutations
}

func (m *MutestRS) Meta() *Meta {
	return meta
}

func (m *MutestRS) Init(p string) error {
	eval, err := os.ReadFile(path.Join(p, "evaluation.json"))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(eval, &m.eval); err != nil {
		return err
	}
	muts, err := os.ReadFile(path.Join(p, "mutations.json"))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(muts, &m.muts); err != nil {
		return err
	}
	return nil
}

func (m *MutestRS) Mutations() (mutations.Mutations, error) {
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
