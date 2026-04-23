package mutations

import (
	"fmt"
	"sort"
	"strings"

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

var Statuses = []Status{Killed, Survived, Crashed, Timeout, NoCoverage}
var statusPaths = map[Status]string{
	Killed: "<path d=\"M256 512a256 256 0 1 1 0-512 256 256 0 1 1 0 512zM374 145.7c-10.7-7.8-25.7-5.4-33.5 5.3L221" +
		".1 315.2 169 263.1c-9.4-9.4-24.6-9.4-33.9 0s-9.4 24.6 0 33.9l72 72c5 5 11.8 7.5 18.8 7s13.4-4.1 17.5-9.8L" +
		"379.3 179.2c7.8-10.7 5.4-25.7-5.3-33.5z\"/>",
	Survived: "<path d=\"M256 512a256 256 0 1 0 0-512 256 256 0 1 0 0 512zM167 167c9.4-9.4 24.6-9.4 33.9 0l55 55 5" +
		"5-55c9.4-9.4 24.6-9.4 33.9 0s9.4 24.6 0 33.9l-55 55 55 55c9.4 9.4 9.4 24.6 0 33.9s-24.6 9.4-33.9 0l-55-55" +
		"-55 55c-9.4 9.4-24.6 9.4-33.9 0s-9.4-24.6 0-33.9l55-55-55-55c-9.4-9.4-9.4-24.6 0-33.9z\"/>",
	Crashed: "<path d=\"M256 512a256 256 0 1 1 0-512 256 256 0 1 1 0 512zm0-192a32 32 0 1 0 0 64 32 32 0 1 0 0-64z" +
		"m0-192c-18.2 0-32.7 15.5-31.4 33.7l7.4 104c.9 12.6 11.4 22.3 23.9 22.3 12.6 0 23-9.7 23.9-22.3l7.4-104c1." +
		"3-18.2-13.1-33.7-31.4-33.7z\"/>",
	Timeout: "<path d=\"M256 0a256 256 0 1 1 0 512 256 256 0 1 1 0-512zM232 120l0 136c0 8 4 15.5 10.7 20l96 64c11 " +
		"7.4 25.9 4.4 33.3-6.7s4.4-25.9-6.7-33.3L280 243.2 280 120c0-13.3-10.7-24-24-24s-24 10.7-24 24z\"/>",
	NoCoverage: "<path d=\"M256 512a256 256 0 1 0 0-512 256 256 0 1 0 0 512zm0-336c-17.7 0-32 14.3-32 32 0 13.3-10" +
		".7 24-24 24s-24-10.7-24-24c0-44.2 35.8-80 80-80s80 35.8 80 80c0 47.2-36 67.2-56 74.5l0 3.8c0 13.3-10.7 24" +
		"-24 24s-24-10.7-24-24l0-8.1c0-20.5 14.8-35.2 30.1-40.2 6.4-2.1 13.2-5.5 18.2-10.3 4.3-4.2 7.7-10 7.7-19.6" +
		" 0-17.7-14.3-32-32-32zM224 368a32 32 0 1 1 64 0 32 32 0 1 1 -64 0z\"/>",
}

func (s Status) Text() string {
	return strings.ToLower(string(s))
}

// Icon returns the icon belonging to the respective Status.
func (s Status) Icon() string {
	svgPath, exists := statusPaths[s]
	if !exists {
		svgPath = statusPaths[NoCoverage]
	}
	return fmt.Sprintf("<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 512 512\" class=\"status-icon %s\">"+
		"<!--!Font Awesome Free v7.2.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/"+
		"license/free Copyright 2026 Fonticons, Inc.-->%s</svg>", s.Text(), svgPath)
}

func (s Status) IconWithText() string {
	return fmt.Sprintf("<div class=\"status-wrapper %s\">%s<p class=\"status-text\">%s</p></div>", s.Text(), s.Icon(), s.Text())
}

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

func (c *Conflict) ConflictsWithMutation(m *Mutation) bool {
	return m.Start.Line <= c.EndLine && m.End.Line >= c.StartLine
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

func (c Conflicts) GetMutant(ID uuid.UUID) (*Conflict, *Mutation) {
	for _, conflict := range c {
		for _, mutation := range conflict.Mutations {
			if mutation.ID == ID {
				return conflict, mutation
			}
		}
	}
	return nil, nil
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

func (m Mutations) StatisticsFrom(prefix string) Statistics {
	s := Statistics{
		StatusCounts: make(map[Status]float64),
	}
	for path, conflicts := range m {
		if strings.HasPrefix(path, prefix) {
			for _, conflict := range conflicts {
				for _, mutation := range conflict.Mutations {
					s.Count++
					s.StatusCounts[mutation.Status]++
				}
			}
		}
	}
	return s
}

type Statistics struct {
	Count        float64
	StatusCounts map[Status]float64
}

func (s Statistics) covered() float64 {
	return s.Count - s.StatusCounts[NoCoverage]
}

func (s Statistics) Coverage() float64 {
	return s.covered() / s.Count * 100
}

func (s Statistics) Score() float64 {
	return s.StatusCounts[Killed] / s.Count * 100
}

func (s Statistics) ScoreOfCovered() float64 {
	return s.StatusCounts[Killed] / s.covered() * 100
}
