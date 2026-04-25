package mtelib

import "testing"

func TestSortingMutantsByRange(t *testing.T) {
	mutants := []MutantResult{
		{ID: "1", Location: Location{Start: Position{Line: 72, Column: 13}, End: Position{Line: 102, Column: 46}}},
		{ID: "2", Location: Location{Start: Position{Line: 995, Column: 4}, End: Position{Line: 1000, Column: 28}}},
		{ID: "3", Location: Location{Start: Position{Line: 311, Column: 98}, End: Position{Line: 343, Column: 120}}},
		{ID: "4", Location: Location{Start: Position{Line: 9, Column: 0}, End: Position{Line: 21, Column: 7}}},
		{ID: "5", Location: Location{Start: Position{Line: 680, Column: 150}, End: Position{Line: 709, Column: 177}}},
		{ID: "6", Location: Location{Start: Position{Line: 257, Column: 33}, End: Position{Line: 287, Column: 66}}},
		{ID: "7", Location: Location{Start: Position{Line: 423, Column: 5}, End: Position{Line: 448, Column: 52}}},
		{ID: "8", Location: Location{Start: Position{Line: 56, Column: 199}, End: Position{Line: 91, Column: 200}}},
		{ID: "9", Location: Location{Start: Position{Line: 814, Column: 22}, End: Position{Line: 846, Column: 45}}},
		{ID: "10", Location: Location{Start: Position{Line: 137, Column: 0}, End: Position{Line: 166, Column: 12}}},
		{ID: "11", Location: Location{Start: Position{Line: 499, Column: 77}, End: Position{Line: 528, Column: 110}}},
		{ID: "12", Location: Location{Start: Position{Line: 921, Column: 3}, End: Position{Line: 946, Column: 9}}},
		{ID: "13", Location: Location{Start: Position{Line: 203, Column: 44}, End: Position{Line: 236, Column: 94}}},
		{ID: "14", Location: Location{Start: Position{Line: 35, Column: 12}, End: Position{Line: 64, Column: 60}}},
		{ID: "15", Location: Location{Start: Position{Line: 742, Column: 185}, End: Position{Line: 772, Column: 200}}},
		{ID: "16", Location: Location{Start: Position{Line: 608, Column: 0}, End: Position{Line: 638, Column: 25}}},
		{ID: "17", Location: Location{Start: Position{Line: 482, Column: 140}, End: Position{Line: 512, Column: 178}}},
		{ID: "18", Location: Location{Start: Position{Line: 120, Column: 66}, End: Position{Line: 149, Column: 99}}},
		{ID: "19", Location: Location{Start: Position{Line: 871, Column: 9}, End: Position{Line: 894, Column: 37}}},
		{ID: "20", Location: Location{Start: Position{Line: 34, Column: 2}, End: Position{Line: 65, Column: 53}}},
	}
	expectedIDOrder := []string{"8", "13", "9", "3", "20", "17", "1", "6", "16", "15", "14", "11", "18", "5", "10", "7", "12", "19", "4", "2"}
	sortMutantsByRange(mutants)
	for i, mutant := range mutants {
		if mutant.ID != expectedIDOrder[i] {
			t.Errorf("expected id in position %d to be \"%s\" but got \"%s\"", i, expectedIDOrder[i], mutant.ID)
		}
	}
}
