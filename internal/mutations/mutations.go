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
	Pending    Status = "PENDING"
	Ignored    Status = "IGNORED"
)

var Statuses = []Status{Killed, Survived, Crashed, Timeout, NoCoverage, Ignored, Pending}
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
	// NOTE: Pending will just use the NoCoverage Image.
	Ignored: "<path d=\"M41-24.9c-9.4-9.4-24.6-9.4-33.9 0S-2.3-.3 7 9.1l528 528c9.4 9.4 24.6 9.4 33.9 0s9.4-24.6 0" +
		"-33.9l-96.4-96.4c2.7-2.4 5.4-4.8 8-7.2 46.8-43.5 78.1-95.4 93-131.1 3.3-7.9 3.3-16.7 0-24.6-14.9-35.7-46." +
		"2-87.7-93-131.1-47.1-43.7-111.8-80.6-192.6-80.6-56.8 0-105.6 18.2-146 44.2L41-24.9zM204.5 138.7c23.5-16.8" +
		" 52.4-26.7 83.5-26.7 79.5 0 144 64.5 144 144 0 31.1-9.9 59.9-26.7 83.5l-34.7-34.7c12.7-21.4 17-47.7 10.1-" +
		"73.7-13.7-51.2-66.4-81.6-117.6-67.9-8.6 2.3-16.7 5.7-24 10l-34.7-34.7zM325.3 395.1c-11.9 3.2-24.4 4.9-37." +
		"3 4.9-79.5 0-144-64.5-144-144 0-12.9 1.7-25.4 4.9-37.3L69.4 139.2c-32.6 36.8-55 75.8-66.9 104.5-3.3 7.9-3" +
		".3 16.7 0 24.6 14.9 35.7 46.2 87.7 93 131.1 47.1 43.7 111.8 80.6 192.6 80.6 37.3 0 71.2-7.9 101.5-20.6l-6" +
		"4.2-64.2z\"/>",
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

func (r Range) LessThan(rge *Range) bool {
	if r.Line < rge.Line {
		return true
	}
	if r.Line == rge.Line && r.Char < rge.Char {
		return true
	}
	return false
}

func (r Range) String() string {
	return fmt.Sprintf("L%dC%d", r.Line, r.Char)
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

func (m Mutation) GetDescription() string {
	if m.Description == "" {
		return m.Operation
	}
	return m.Description
}

func (m Mutation) IsBroken() bool {
	return m.End.LessThan(m.Start)
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

func (c *Conflict) ConflictsWithConflict(cb *Conflict) bool {
	return cb.StartLine <= c.EndLine && cb.EndLine >= c.StartLine
}

func (c *Conflict) Append(m *Mutation) {
	if m.Start.Line < c.StartLine {
		c.StartLine = m.Start.Line
	}
	if m.End.Line > c.EndLine {
		c.EndLine = m.End.Line
	}
	c.Mutations = append(c.Mutations, m)
}

func (c *Conflict) Merge(cb *Conflict) {
	if cb.StartLine < c.StartLine {
		c.StartLine = cb.StartLine
	}
	if cb.EndLine > c.EndLine {
		c.EndLine = cb.EndLine
	}
	c.Mutations = append(c.Mutations, cb.Mutations...)
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
	for _, c := range m[file] {
		if c.ConflictsWithMutation(mutation) {
			c.Append(mutation)
			return
		}
	}
	m[file] = append(m[file], NewConflict(mutation))
}

func (m Mutations) ExtractBrokenMutations() []*Mutation {
	broken := make([]*Mutation, 0)
	for _, conflicts := range m {
		for _, conflict := range conflicts {
			mutations := make([]*Mutation, 0, len(conflict.Mutations))
			for _, mutation := range conflict.Mutations {
				if mutation.IsBroken() {
					broken = append(broken, mutation)
					continue
				}
				mutations = append(mutations, mutation)
			}
			conflict.Mutations = mutations
		}
	}
	return broken
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

func (s Statistics) Detected() float64 {
	return s.StatusCounts[Killed] + s.StatusCounts[Timeout]
}

func (s Statistics) Undetected() float64 {
	return s.StatusCounts[Survived] + s.StatusCounts[NoCoverage]
}

func (s Statistics) Covered() float64 {
	return s.Detected() + s.StatusCounts[Survived]
}

func (s Statistics) Valid() float64 {
	return s.Detected() + s.Undetected()
}

func (s Statistics) Invalid() float64 {
	return s.StatusCounts[Crashed]
}

func (s Statistics) Score() float64 {
	return s.Detected() / s.Valid() * 100
}

func (s Statistics) ScoreOfCovered() float64 {
	return s.Detected() / s.Covered() * 100
}

func (s Statistics) Coverage() float64 {
	return s.Covered() / s.Count * 100
}
