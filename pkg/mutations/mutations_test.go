package mutations

import "testing"

var m1 = &Mutation{
	FrameworkMutantID: "0",
	Description:       "test mutation 1",
	Operation:         "bla bla bla",
	Start: &Range{
		Line: 44,
		Char: 4,
	},
	End: &Range{
		Line: 70,
		Char: 92,
	},
	Status:      Killed,
	Replacement: "line that was replaced ...",
}
var m2 = &Mutation{
	FrameworkMutantID: "1",
	Description:       "test mutation 2",
	Operation:         "bla bla bla",
	Start: &Range{
		Line: 60,
		Char: 40,
	},
	End: &Range{
		Line: 1006,
		Char: 30,
	},
	Status:      Killed,
	Replacement: "line that was replaced ...",
}
var m3 = &Mutation{
	FrameworkMutantID: "2",
	Description:       "test mutation 3",
	Operation:         "bla bla bla",
	Start: &Range{
		Line: 20,
		Char: 0,
	},
	End: &Range{
		Line: 43,
		Char: 400,
	},
	Status:      Killed,
	Replacement: "line that was replaced ...",
}
var m4 = &Mutation{
	FrameworkMutantID: "2",
	Description:       "test mutation 3",
	Operation:         "bla bla bla",
	Start: &Range{
		Line: 20,
		Char: 0,
	},
	End: &Range{
		Line: 50,
		Char: 400,
	},
	Status:      Killed,
	Replacement: "line that was replaced ...",
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
		expConflictEndLine: m2.End.Line,
	},
	{
		name:               "conflict with start line before conflict start line, but end line between conflict start and end line",
		conflict:           NewConflict(m1),
		mutation:           m4,
		expConflict:        true,
		expConflictEndLine: m1.End.Line,
	},
}

func assertConflicts(t *testing.T, c conflictsTestCase) {
	conf := c.conflict.Conflicts(c.mutation)
	if conf != c.expConflict && c.expConflict {
		t.Errorf("conflict between Conflict{%s} and Mutation{%s} was not detected",
			c.conflict.Mutations[0].FrameworkMutantID, c.mutation.FrameworkMutantID)
	}
	if conf != c.expConflict && !c.expConflict {
		t.Errorf("unexpected conflict detected between Conflict{%s} and Mutation{%s}",
			c.conflict.Mutations[0].FrameworkMutantID, c.mutation.FrameworkMutantID)
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
		t.Errorf("conflict Conflict{%s} expected end line to be %d, got %d",
			c.conflict.Mutations[0].FrameworkMutantID, c.expConflictEndLine, c.conflict.EndLine)
	}
}

func TestConflict_Append(t *testing.T) {
	for _, c := range conflictsTestCases {
		t.Run(c.name, func(t *testing.T) {
			assertAppend(t, c)
		})
	}
}
