package mutations

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestStatusToTextConversion(t *testing.T) {
	tests := []struct {
		Status   Status
		Expected string
	}{
		{Killed, "killed"},
		{Survived, "survived"},
		{Crashed, "crashed"},
		{Timeout, "timeout"},
		{NoCoverage, "no_coverage"},
		{Status("OTHER_STATUS"), "other_status"},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s.Text() == \"%s\"?", test.Status, test.Expected), func(t *testing.T) {
			result := test.Status.Text()
			if result != test.Expected {
				t.Errorf("expected %s but got %s", test.Expected, result)
			}
		})
	}
}

func TestStatusToIconWithTextConversion(t *testing.T) {
	tests := []struct {
		Status         Status
		ExpectedSubstr string
	}{
		{Killed, "<p class=\"status-text\">killed</p>"},
		{Survived, "<p class=\"status-text\">survived</p>"},
		{Crashed, "<p class=\"status-text\">crashed</p>"},
		{Timeout, "<p class=\"status-text\">timeout</p>"},
		{NoCoverage, "<p class=\"status-text\">no_coverage</p>"},
		// looking for part of the no_coverage icon
		{Status("OTHER_STATUS"), "M256 512a256 256 0 1 0 0-512 256 256 0 1 0 0 512zm0-336c-17.7 0-32 14.3"},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s.IconWithText() contains \"%s\"?", test.Status, test.ExpectedSubstr), func(t *testing.T) {
			result := test.Status.IconWithText()
			if !strings.Contains(result, test.ExpectedSubstr) {
				t.Errorf("expected %s but got %s", test.ExpectedSubstr, result)
			}
		})
	}
}

func TestRangeLessThanComparisons(t *testing.T) {
	// NOTE: test will always do RangeA.LessThan(RangeB).
	tests := []struct {
		RangeA       Range
		RangeB       Range
		ExpectRangeA bool // true if test is expecting RangeA to be less than RangeB
	}{
		{Range{Line: 0, Char: 0}, Range{Line: 1, Char: 0}, true},
		{Range{Line: 100, Char: 0}, Range{Line: 1, Char: 0}, false},
		{Range{Line: 1, Char: 2}, Range{Line: 1, Char: 0}, false},
		{Range{Line: 1, Char: 2}, Range{Line: 1, Char: 15}, true},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s.LessThan(%s) == %v?", test.RangeA, test.RangeB, test.ExpectRangeA), func(t *testing.T) {
			result := test.RangeA.LessThan(&test.RangeB)
			if result != test.ExpectRangeA {
				t.Errorf("expected %v but got %v", test.ExpectRangeA, result)
			}
		})
	}
}

func TestMutationReturnsOperationWhenNoDescription(t *testing.T) {
	tests := []struct {
		Name     string
		Mutation Mutation
		Expected string
	}{
		{"Mutation returns description", Mutation{Description: "deleted `x`", Operation: "deletion_mutator"}, "deleted `x`"},
		{"Mutation returns operation in place of empty description", Mutation{Operation: "deletion_mutator"}, "deletion_mutator"},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result := test.Mutation.GetDescription()
			if result != test.Expected {
				t.Errorf("expected %s but got %s", test.Expected, result)
			}
		})
	}
}

func TestMutationBrokenCheck(t *testing.T) {
	tests := []struct {
		Name            string
		Mutation        Mutation
		ExpectingBroken bool
	}{
		{"Mutation is not broken", Mutation{Start: &Range{0, 0}, End: &Range{1, 0}}, false},
		{"Mutation is broken as start range > end range", Mutation{Start: &Range{10, 0}, End: &Range{1, 0}}, true},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result := test.Mutation.IsBroken()
			if result != test.ExpectingBroken {
				t.Errorf("expected %v but got %v", test.ExpectingBroken, result)
			}
		})
	}
}

func TestConflictCreationCreatesCorrectRegion(t *testing.T) {
	m := Mutation{Start: &Range{10, 35}, End: &Range{10, 45}}
	c := NewConflict(&m)
	if c.StartLine != m.Start.Line {
		t.Errorf("expected start line to be %d but got %d", m.Start.Line, c.StartLine)
	}
	if c.EndLine != m.End.Line {
		t.Errorf("expected end line to be %d but got %d", m.End.Line, c.EndLine)
	}
}

func TestConflictCorrectlyReportsMutationsOverlap(t *testing.T) {
	c := &Conflict{StartLine: 10, EndLine: 25}
	tests := []struct {
		Name              string
		Mutation          Mutation
		ExpectingConflict bool
	}{
		{"Mutation outside of conflict zone", Mutation{Start: &Range{0, 0}, End: &Range{1, 0}}, false},
		{"Mutation overlaps lower boundary of conflict zone", Mutation{Start: &Range{8, 0}, End: &Range{12, 0}}, true},
		{"Mutation overlaps upper boundary of conflict zone", Mutation{Start: &Range{22, 0}, End: &Range{28, 0}}, true},
		{"Mutation wrapped by boundaries of conflict zone", Mutation{Start: &Range{14, 0}, End: &Range{20, 0}}, true},
		{"Mutation surrounds boundaries of conflict zone from outside the zone", Mutation{Start: &Range{4, 0}, End: &Range{28, 0}}, true},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result := c.ConflictsWithMutation(&test.Mutation)
			if result != test.ExpectingConflict {
				t.Errorf("expected mutation conflict == %v but got %v", test.ExpectingConflict, result)
			}
		})
	}
}

func TestConflictBoundaryExpansionWhenAppendingNewMutation(t *testing.T) {
	// tests conducted with new conflict created inside of runner loop. conflict always starts with StartLine: 10
	// and EndLine: 25
	tests := []struct {
		Name                      string
		Mutation                  Mutation
		ExpectedConflictStartLine int
		ExpectedConflictEndLine   int
	}{
		{"Mutation expands lower boundary of conflict zone", Mutation{Start: &Range{8, 0}, End: &Range{12, 0}}, 8, 25},
		{"Mutation expands upper boundary of conflict zone", Mutation{Start: &Range{22, 0}, End: &Range{28, 0}}, 10, 28},
		{"Mutation does not expand boundaries of conflict zone", Mutation{Start: &Range{14, 0}, End: &Range{20, 0}}, 10, 25},
		{"Mutation expands both boundaries of conflict zone", Mutation{Start: &Range{4, 0}, End: &Range{28, 0}}, 4, 28},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			c := &Conflict{StartLine: 10, EndLine: 25}
			c.Append(&test.Mutation)
			if c.StartLine != test.ExpectedConflictStartLine {
				t.Errorf("expected start line to be %d but got %d", test.ExpectedConflictStartLine, c.StartLine)
			}
			if c.EndLine != test.ExpectedConflictEndLine {
				t.Errorf("expected end line to be %d but got %d", test.ExpectedConflictEndLine, c.EndLine)
			}
		})
	}
}

func TestSortCorrectlyOrdersConflictsFromFirstToLastBasedOnStartLine(t *testing.T) {
	c := Conflicts{
		&Conflict{StartLine: 1201},
		&Conflict{StartLine: 102},
		&Conflict{StartLine: 400},
		&Conflict{StartLine: 0},
	}
	c.Sort()
	if c[0].StartLine != 0 {
		t.Errorf("expected first conflict to start from line 0 but got %d", c[0].StartLine)
	}
	if c[1].StartLine != 102 {
		t.Errorf("expected second conflict to start from line 102 but got %d", c[1].StartLine)
	}
	if c[2].StartLine != 400 {
		t.Errorf("expected third conflict to start from line 400 but got %d", c[2].StartLine)
	}
	if c[3].StartLine != 1201 {
		t.Errorf("expected last conflict to start from line 1201 but got %d", c[3].StartLine)
	}
}

func TestConflictsGetMutant(t *testing.T) {
	m1ID := uuid.New()
	m2ID := uuid.New()

	m1 := Mutation{ID: m1ID}
	m2 := Mutation{ID: m2ID}

	c := Conflicts{
		&Conflict{StartLine: 45, Mutations: []*Mutation{&m1}},
		&Conflict{StartLine: 65, Mutations: []*Mutation{&m2}},
	}

	tests := []struct {
		ID                uuid.UUID
		Mutation          *Mutation
		ExpectedStartLine int
	}{
		{m1ID, &m1, 45},
		{m2ID, &m2, 65},
		{uuid.Nil, nil, 0},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("correctly retrieves mutant %s", test.ID), func(t *testing.T) {
			conflict, mutant := c.GetMutant(test.ID)
			if test.ID == uuid.Nil && conflict == nil && mutant == nil {
				return // this passes the test as if id is nil then it is not looking for anything
			}
			if conflict.StartLine != test.ExpectedStartLine {
				t.Fatalf("expected conflict start line to be %d but got %d", test.ExpectedStartLine, conflict.StartLine)
			}
			if mutant != test.Mutation {
				t.Fatal("got an incorrect mutation reference")
			}
		})
	}
}

func TestMergingMutationMaps(t *testing.T) {
	m1 := Mutations{"file/path/1.lang": Conflicts{}}
	m2 := Mutations{"file/path/2.lang": Conflicts{}}
	m1.Merge(m2)
	if m1["file/path/1.lang"] == nil {
		t.Fatal("m1 lost contents of file/path/1.lang during merge operation")
	}
	if m1["file/path/2.lang"] == nil {
		t.Fatal("m1 did not gain contents of file/path/2.lang during merge operation")
	}
}

// TODO: add Append test

func TestExtractingBrokenMutantsFromMutationsMap(t *testing.T) {
	validMut1 := &Mutation{Start: &Range{1, 0}, End: &Range{2, 0}}
	validMut2 := &Mutation{Start: &Range{1, 0}, End: &Range{2, 0}}
	validMut3 := &Mutation{Start: &Range{1, 0}, End: &Range{2, 0}}
	validMut4 := &Mutation{Start: &Range{1, 0}, End: &Range{2, 0}}
	brokenMut1 := &Mutation{Start: &Range{100, 0}, End: &Range{0, 0}}
	brokenMut2 := &Mutation{Start: &Range{1, 400}, End: &Range{1, 0}}
	con1 := &Conflict{Mutations: []*Mutation{brokenMut1, validMut1}}
	con2 := &Conflict{Mutations: []*Mutation{validMut3, validMut4, brokenMut2}}
	m := Mutations{
		"file/path/1.lang": Conflicts{con1, &Conflict{Mutations: []*Mutation{validMut2}}},
		"file/path/2.lang": Conflicts{con2},
	}
	broken := m.ExtractBrokenMutations()
	if len(broken) != 2 {
		t.Errorf("expected 2 broken mutations but got %d", len(broken))
	}
	if !slices.Contains(broken, brokenMut1) {
		t.Errorf("slice of broken mutations does not contain brokenMut1")
	}
	if slices.Contains(con1.Mutations, brokenMut1) {
		t.Errorf("conflict1 of correct mutations contains brokenMut1")
	}
	if !slices.Contains(broken, brokenMut2) {
		t.Errorf("slice of broken mutations does not contain brokenMut2")
	}
	if slices.Contains(con2.Mutations, brokenMut2) {
		t.Errorf("conflict2 of correct mutations contains brokenMut2")
	}
}

func TestEnsureGenerateIDsAssignsAnIDToAllMutations(t *testing.T) {
	m := Mutations{
		"file1.lang": Conflicts{
			&Conflict{Mutations: []*Mutation{{Status: Killed}}},
			&Conflict{Mutations: []*Mutation{{Status: Killed}}},
		},
		"file2.lang": Conflicts{
			&Conflict{Mutations: []*Mutation{{Status: Killed}}},
			&Conflict{Mutations: []*Mutation{{Status: Killed}}},
		},
		"file3.lang": Conflicts{
			&Conflict{Mutations: []*Mutation{{Status: Killed}}},
			&Conflict{Mutations: []*Mutation{{Status: Killed}}},
		},
		"file4.lang": Conflicts{
			&Conflict{Mutations: []*Mutation{{Status: Killed}}},
			&Conflict{Mutations: []*Mutation{{Status: Killed}}},
		},
	}
	m.GenerateIDs()
	for _, conflicts := range m {
		for _, conflict := range conflicts {
			for _, mutation := range conflict.Mutations {
				if mutation.ID == uuid.Nil {
					t.Errorf("mutation was not assigned an id")
				}
			}
		}
	}
}

func TestGetStatisticsFromMutationsMap(t *testing.T) {
	m := Mutations{
		// KILLED: 4, SURVIVED: 3, TIMEOUT: 1
		"root/file1.lang": Conflicts{
			&Conflict{Mutations: []*Mutation{
				{Status: Killed}, {Status: Killed}, {Status: Killed}, {Status: Survived},
				{Status: Survived}, {Status: Killed}, {Status: Timeout}, {Status: Survived},
			}},
		},
		// KILLED: 8, SURVIVED: 6, TIMEOUT: 2
		"root/path/file2.lang": Conflicts{
			&Conflict{Mutations: []*Mutation{
				{Status: Killed}, {Status: Killed}, {Status: Killed}, {Status: Survived},
				{Status: Survived}, {Status: Killed}, {Status: Timeout}, {Status: Survived},
				{Status: Killed}, {Status: Killed}, {Status: Killed}, {Status: Survived},
				{Status: Survived}, {Status: Killed}, {Status: Timeout}, {Status: Survived},
			}},
		},
		// KILLED: 2, CRASHED: 1
		"root/path/long/file.lang": Conflicts{
			&Conflict{Mutations: []*Mutation{
				{Status: Killed}, {Status: Killed}, {Status: Crashed},
			}},
		},
	}
	rootStats := m.StatisticsFrom("root/")
	if rootStats.StatusCounts[Killed] != 14 {
		t.Errorf("expected 14 killed mutations on root/ but got %f", rootStats.StatusCounts[Killed])
	}
	pathStats := m.StatisticsFrom("root/path/")
	if pathStats.StatusCounts[Killed] != 10 {
		t.Errorf("expected 10 killed mutations on root/path/ but got %f", pathStats.StatusCounts[Killed])
	}
	fileStats := m.StatisticsFrom("root/path/long/file.lang")
	if fileStats.StatusCounts[Crashed] != 1 {
		t.Errorf("expected 1 crashed mutation on root/path/long/file.lang but got %f", fileStats.StatusCounts[Crashed])
	}
}

func TestStatisticsCalculations(t *testing.T) {
	var (
		killedCount     float64 = 2102
		survivedCount   float64 = 3042
		crashedCount    float64 = 3
		timeoutCount    float64 = 60
		noCoverageCount float64 = 1344
		count                   = killedCount + survivedCount + crashedCount + timeoutCount + noCoverageCount
	)
	s := Statistics{
		Count: count,
		StatusCounts: map[Status]float64{
			Killed:     killedCount,
			Survived:   survivedCount,
			Crashed:    crashedCount,
			Timeout:    timeoutCount,
			NoCoverage: noCoverageCount,
		},
	}
	if s.Detected() != killedCount+timeoutCount {
		t.Errorf("incorrect detected calculation")
	}
	if s.Undetected() != survivedCount+noCoverageCount {
		t.Errorf("incorrect undetected calculation")
	}
	if s.Covered() != killedCount+timeoutCount+survivedCount {
		t.Errorf("incorrect covered calculation")
	}
	if s.Valid() != killedCount+timeoutCount+survivedCount+noCoverageCount {
		t.Errorf("incorrect valid calculation")
	}
	if s.Invalid() != crashedCount {
		t.Errorf("incorrect invalid calculation")
	}
	if s.Score() != (killedCount+timeoutCount)/(killedCount+timeoutCount+survivedCount+noCoverageCount)*100 {
		t.Errorf("incorrect score calculation")
	}
	if s.ScoreOfCovered() != (killedCount+timeoutCount)/(killedCount+timeoutCount+survivedCount)*100 {
		t.Errorf("incorrect score of covered calculation")
	}
	if s.Coverage() != (killedCount+timeoutCount+survivedCount)/count*100 {
		t.Errorf("incorrect coverage calculation")
	}
}
