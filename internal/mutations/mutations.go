package mutations

import (
	"sort"

	"github.com/google/uuid"
)

// Status represents the outcome of a mutation.
type Status string

const (
	Killed     Status = "KILLED"
	Survived   Status = "SURVIVED"
	Crashed    Status = "CRASHED"
	Timeout    Status = "TIMEOUT"
	NoCoverage Status = "NO_COVERAGE"
)

// Range holds a line and char index.
type Range struct {
	Line int
	Char int
}

// Mutation represents a single mutation.
type Mutation struct {
	ID                uuid.UUID
	FrameworkMutantID string // FrameworkMutantID is an optional identifier provided by a framework.
	Description       string // Description functions as the mutations title, it is displayed to the user when they preview the mutation.
	Operation         string // Operation is an optional (short) description of the operation.
	Start             *Range
	End               *Range
	Status            Status
	Replacement       string
}

// Conflict represents all mutations that would conflict with each other if they were displayed simultaneously.
type Conflict struct {
	ID        uuid.UUID
	StartLine int
	EndLine   int
	Mutations []*Mutation
}

func NewConflict(m *Mutation) *Conflict {
	return &Conflict{
		StartLine: m.Start.Line,
		EndLine:   m.End.Line,
		Mutations: []*Mutation{m},
	}
}

func (c *Conflict) Conflicts(m *Mutation) bool {
	return m.Start.Line >= c.StartLine && m.Start.Line <= c.EndLine ||
		m.End.Line >= c.StartLine && m.End.Line <= c.EndLine
}

func (c *Conflict) Append(m *Mutation) {
	if m.End.Line > c.EndLine {
		c.EndLine = m.End.Line
	}
	c.Mutations = append(c.Mutations, m)
}

// Conflicts is a slice of Conflict instances.
type Conflicts []*Conflict

func (c Conflicts) Sort() {
	sort.Slice(c, func(i, j int) bool {
		return c[i].StartLine < c[j].StartLine
	})
}

// Mutations is a map of file names to groups of conflicting mutations.
type Mutations map[string]Conflicts

func (m Mutations) Merge(b Mutations) {
	for k, v := range b {
		m[k] = v
	}
}

func (m Mutations) Append(file string, mutation *Mutation) {
	added := false
	for _, c := range m[file] {
		if c.Conflicts(mutation) {
			c.Append(mutation)
			added = true
			break
		}
	}

	if !added {
		m[file] = append(m[file], NewConflict(mutation))
	}
}

// GenerateIDs generates UUIDs for all conflicts and mutations
func (m Mutations) GenerateIDs() {
	for _, conflicts := range m {
		for _, conflict := range conflicts {
			conflict.ID = uuid.New()
			for _, mutation := range conflict.Mutations {
				mutation.ID = uuid.New()
			}
		}
	}
}
