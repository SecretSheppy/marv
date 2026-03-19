package mutations

import "testing"

var m1 = &Mutation{
	ID:     0,
	Name:   "test mutation 1",
	OpDesc: "bla bla bla",
	Starts: &Range{
		Line: 44,
		Char: 4,
	},
	Ends: &Range{
		Line: 70,
		Char: 92,
	},
	Status: Killed,
	Type:   Replacement,
	Source: "line that was replaced ...",
}
var m2 = &Mutation{
	ID:     1,
	Name:   "test mutation 2",
	OpDesc: "bla bla bla",
	Starts: &Range{
		Line: 60,
		Char: 40,
	},
	Ends: &Range{
		Line: 1006,
		Char: 30,
	},
	Status: Killed,
	Type:   Replacement,
	Source: "line that was replaced ...",
}
var m3 = &Mutation{
	ID:     2,
	Name:   "test mutation 3",
	OpDesc: "bla bla bla",
	Starts: &Range{
		Line: 20,
		Char: 0,
	},
	Ends: &Range{
		Line: 43,
		Char: 400,
	},
	Status: Killed,
	Type:   Replacement,
	Source: "line that was replaced ...",
}
var m4 = &Mutation{
	ID:     2,
	Name:   "test mutation 3",
	OpDesc: "bla bla bla",
	Starts: &Range{
		Line: 20,
		Char: 0,
	},
	Ends: &Range{
		Line: 50,
		Char: 400,
	},
	Status: Killed,
	Type:   Replacement,
	Source: "line that was replaced ...",
}

type conflictsTestCase struct {
	name               string
	conflict           *Conflict
	mutation           *Mutation
	expConflict        bool
	expConflictEndLine int
}

var conflictsTestCases = []conflictsTestCase{
	{
		name:        "no conflict",
		conflict:    NewConflict(m1),
		mutation:    m3,
		expConflict: false,
	},
	{
		name:               "conflict with start line contained between conflict start and end line",
		conflict:           NewConflict(m1),
		mutation:           m2,
		expConflict:        true,
		expConflictEndLine: m2.Ends.Line,
	},
	{
		name:               "conflict with start line before conflict start line, but end line between conflict start and end line",
		conflict:           NewConflict(m1),
		mutation:           m4,
		expConflict:        true,
		expConflictEndLine: m1.Ends.Line,
	},
}

func assertConflicts(t *testing.T, c conflictsTestCase) {
	conf := c.conflict.Conflicts(c.mutation)
	if conf != c.expConflict && c.expConflict {
		t.Errorf("conflict between Conflict{%d} and Mutation{%d} was not detected",
			c.conflict.Mutations[0].ID, c.mutation.ID)
	}
	if conf != c.expConflict && !c.expConflict {
		t.Errorf("unexpected conflict detected between Conflict{%d} and Mutation{%d}",
			c.conflict.Mutations[0].ID, c.mutation.ID)
	}
}

func TestConflict_Conflicts(t *testing.T) {
	for _, c := range conflictsTestCases {
		t.Run(c.name, func(t *testing.T) {
			assertConflicts(t, c)
		})
	}
}

func assertAppend(t *testing.T, c conflictsTestCase) {
	if !c.expConflict {
		return
	}
	c.conflict.Append(c.mutation)
	if c.conflict.EndLine != c.expConflictEndLine {
		t.Errorf("conflict Conflict{%d} expected end line to be %d, got %d",
			c.conflict.Mutations[0].ID, c.expConflictEndLine, c.conflict.EndLine)
	}
}

func TestConflict_Append(t *testing.T) {
	for _, c := range conflictsTestCases {
		t.Run(c.name, func(t *testing.T) {
			assertAppend(t, c)
		})
	}
}
