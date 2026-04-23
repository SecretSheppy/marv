package mtelib

import (
	"encoding/json"
	"os"
	"sort"
	"strings"

	"github.com/SecretSheppy/marv/internal/mutations"
)

// Mutation Testing Elements Library
//
// A library that provides structs and methods to unmarshal and marshal mutations from the Mutation Testing Elements
// JSON format. Built off of mutation testing report schema version 2.0.1

// MutationTestResult represents the main Mutation Testing Elements JSON file.
type MutationTestResult struct {
	SchemaVersion string               `json:"schemaVersion"`
	Files         FileResultDictionary `json:"files"`
}

// FileResultDictionary is a dictionary that stores FileResults against their string file paths.
type FileResultDictionary map[string]FileResult

// FileResult contains the files language, mutants and unedited source code.
type FileResult struct {
	Language string         `json:"language"`
	Mutants  []MutantResult `json:"mutants"`
	Source   string         `json:"source"`
}

// MutantResult contains the data about a specific mutant.
type MutantResult struct {
	ID          string    `json:"id"`
	MutatorName string    `json:"mutatorName"`
	Replacement string    `json:"replacement"`
	Location    Location  `json:"location"`
	Status      MTEStatus `json:"status"`
	Description string    `json:"description"`
}

func (m *MutantResult) toMarvMutation() *mutations.Mutation {
	return &mutations.Mutation{
		FrameworkMutantID: m.ID,
		Description:       m.Description,
		Operation:         m.MutatorName,
		Start:             m.Location.Start.toMarvRange(),
		End:               m.Location.End.toMarvRange(),
		Status:            m.Status.toMarvStatus(),
		Replacement:       m.Replacement,
	}
}

type MTEStatus string

func (m MTEStatus) toMarvStatus() mutations.Status {
	switch string(m) {
	case "Survived":
		return mutations.Survived
	case "Killed":
		return mutations.Killed
	case "RuntimeError", "CompileError":
		return mutations.Crashed
	case "Timeout":
		return mutations.Timeout
	case "Pending":
		return mutations.Pending
	case "Ignored":
		return mutations.Ignored
	default:
		return mutations.NoCoverage
	}
}

// Location describes a range within the source code. Start is inclusive, end is exclusive.
type Location struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// Position describes a single position within the source code. Both line and column start at one.
type Position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

func (p *Position) toMarvRange() *mutations.Range {
	return &mutations.Range{
		Line: p.Line - 1,
		Char: p.Column - 1,
	}
}

type MTE struct {
	result    MutationTestResult
	mutations mutations.Mutations
	files     map[string]string
}

func NewMTE(file string) (*MTE, error) {
	raw, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	mte := &MTE{}
	if err := json.Unmarshal(raw, &mte.result); err != nil {
		return nil, err
	}
	return mte, nil
}

func (m *MTE) Transform() {
	// TODO: progress bar
	m.mutations = make(mutations.Mutations)
	m.files = make(map[string]string)

	for file, fileResult := range m.result.Files {
		if strings.HasPrefix(file, "/") {
			file = file[1:]
		}
		m.files[file] = fileResult.Source
		SortMutantsByRange(fileResult.Mutants)
		for _, mutant := range fileResult.Mutants {
			m.mutations.Append(file, mutant.toMarvMutation())
		}
	}

	m.result = MutationTestResult{}
}

// SortMutantsByRange sorts in-place: largest ranges first.
func SortMutantsByRange(ms []MutantResult) {
	sort.Slice(ms, func(i, j int) bool {
		ri := rangeSize(ms[i])
		rj := rangeSize(ms[j])
		if ri != rj {
			return ri > rj // larger range first
		}
		// tie-breaker: compare column span
		ci := columnSpan(ms[i])
		cj := columnSpan(ms[j])
		if ci != cj {
			return ci > cj
		}
		// final tie-breaker: stable deterministic ordering by ID
		return ms[i].ID < ms[j].ID
	})
}

func rangeSize(m MutantResult) int {
	return m.Location.End.Line - m.Location.Start.Line
}

func columnSpan(m MutantResult) int {
	return m.Location.End.Column - m.Location.Start.Column
}

func (m *MTE) Mutations() mutations.Mutations {
	return m.mutations
}

func (m *MTE) ReadLines(file string) []string {
	return getLinesFromString(m.files[file])
}

func getLinesFromString(str string) []string {
	lines := make([]string, 0)
	for line := range strings.Lines(str) {
		lines = append(lines, strings.ReplaceAll(line, "\n", ""))
	}
	return lines
}
