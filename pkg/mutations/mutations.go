package mutations

// Status represents the outcome of a mutation.
type Status string

const (
	Killed     Status = "KILLED"
	Survived   Status = "SURVIVED"
	Crashed    Status = "CRASHED"
	Timeout    Status = "TIMEOUT"
	NoCoverage Status = "NO_COVERAGE"
)

// Modification represents the type of mutation.
type Modification string

const (
	Replacement Modification = "REPLACEMENT"
	Deletion    Modification = "DELETION"
	Insertion   Modification = "INSERTION"
	Swap        Modification = "SWAP"
	Negation    Modification = "NEGATION"
	Reorder     Modification = "REORDER"
)

// Range holds a line and char index.
type Range struct {
	Line int
	Char int
}

// Mutation represents a single mutation.
type Mutation struct {
	ID       int
	IDOffset int    // IDOffset used to specify an offset that should be used to acquire the real mutation ID.
	Name     string // Name functions as the mutations title, it is displayed to the user when they preview the mutation.
	OpDesc   string // OpDesc is an optional (short) description of the operation.
	Starts   *Range
	Ends     *Range
	Status   Status
	Type     Modification
	Source   string
}

// TrueID returns the mutation id + its id offset. This method exists as marv expects all Mutation.ID's to be indexed
// from 0, however this altered mutation id would not reflect the true mutation id that was created by the framework.
func (m *Mutation) TrueID() int {
	return m.ID + m.IDOffset
}

// Conflict represents all mutations that would conflict with each other if they were displayed simultaneously.
type Conflict struct {
	StartLine int
	EndLine   int
	Mutations []*Mutation
}

func NewConflict(m *Mutation) *Conflict {
	return &Conflict{
		StartLine: m.Starts.Line,
		EndLine:   m.Ends.Line,
		Mutations: []*Mutation{m},
	}
}

func (c *Conflict) Conflicts(m *Mutation) bool {
	return m.Starts.Line >= c.StartLine && m.Starts.Line <= c.EndLine ||
		m.Ends.Line >= c.StartLine && m.Ends.Line <= c.EndLine
}

func (c *Conflict) Append(m *Mutation) {
	if m.Ends.Line > c.EndLine {
		c.EndLine = m.Ends.Line
	}
	c.Mutations = append(c.Mutations, m)
}

// Conflicts is a slice of Conflict instances.
type Conflicts []*Conflict

// Mutations is a map of file names to groups of conflicting mutations.
type Mutations map[string]Conflicts

func (m Mutations) Merge(b Mutations) {
	for k, v := range b {
		m[k] = v
	}
}
