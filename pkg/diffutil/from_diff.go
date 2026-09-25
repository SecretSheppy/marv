package diffutil

import (
	"errors"
	"strings"
)

type DiffType string

const (
	Removed  DiffType = "REMOVED"
	Equal    DiffType = "EQUAL"
	Inserted DiffType = "INSERTED"

	NilLineIndex = -1
)

func lineType(line string) DiffType {
	switch line[:1] {
	case "-":
		return Removed
	case "+":
		return Inserted
	default:
		return Equal
	}
}

type DiffLine struct {
	Number int
	Type   DiffType
	Text   string
}

type DiffLines []*DiffLine

func (d DiffLines) Get(number int) *DiffLine {
	for _, line := range d {
		if line.Number == number {
			return line
		}
	}
	return nil
}

func (d DiffLines) LinesByType(diffType DiffType) DiffLines {
	dls := make(DiffLines, 0)
	for _, line := range d {
		if line.Type == diffType {
			dls = append(dls, line)
		}
	}
	return dls
}

func (d DiffLines) LineChanges() (removed, inserted DiffLines) {
	return d.LinesByType(Removed), d.LinesByType(Inserted)
}

func (d DiffLines) StringLines() []string {
	lines := make([]string, len(d))
	for i, line := range d {
		lines[i] = line.Text
	}
	return lines
}

type DiffConfig struct {
	PrefixLines            int // number of prefix metadata lines to be trimmed from diff
	SuffixLines            int // number of suffix metadata lines to be trimmed from diff
	FirstRemovedLineNumber int // line number associated with the first removed line
}

type FormattedDiff struct {
	config    *DiffConfig
	diffLines DiffLines
}

func FromFormattedDiff(diff string, config *DiffConfig) *FormattedDiff {
	var (
		lines     = strings.Split(diff, "\n")
		diffLines = make(DiffLines, 0)
		from      = config.PrefixLines
		to        = len(lines) - config.SuffixLines
	)
	for _, line := range lines[from:to] {
		diffLines = append(diffLines, &DiffLine{Type: lineType(line), Text: line[1:]})
	}
	return &FormattedDiff{
		config:    config,
		diffLines: diffLines,
	}
}

// returns the first removed line from the text diff (i.e. the first DiffLine where DiffLine.Type == Removed)
func (f *FormattedDiff) firstRemovedLineIndex() int {
	for i, line := range f.diffLines {
		if line.Type == Removed {
			return i
		}
	}
	return NilLineIndex
}

// returns an array of indexes for each removed line.
func (f *FormattedDiff) removedLineIndexes() []int {
	removed := make([]int, 0)
	for i, line := range f.diffLines {
		if line.Type == Removed {
			removed = append(removed, i)
		}
	}
	return removed
}

func (f *FormattedDiff) linesSynced(lines []string, removedIndexes []int, sourceIndex int) bool {
	firstRemoved := removedIndexes[0]
	for _, index := range removedIndexes {
		diffLine := strings.TrimSpace(f.diffLines[index].Text)
		// NOTE: sourceIndex is adjusted to match the same line as in the diff by using the current index (which will
		// always be bigger than the first removed index) and subtracting the first removed index.
		sourceLine := strings.TrimSpace(lines[sourceIndex+(index-firstRemoved)])
		if diffLine != sourceLine {
			return false
		}
	}
	return true
}

// SyncLineNumbers syncs the DiffLine numbers with the original source code. This is used for cases where the produced
// text diff may be slightly different to the actual change made (i.e. an extra deleted blank line is present in the
// diff before the expected first line.)
func (f *FormattedDiff) SyncLineNumbers(lines []string) error {
	removed := f.removedLineIndexes()
	if len(removed) == 0 {
		return errors.New("no removed lines in text diff")
	}

	if !f.linesSynced(lines, removed, f.config.FirstRemovedLineNumber) {
		for i := -3; i <= 3; i++ {
			if f.linesSynced(lines, removed, f.config.FirstRemovedLineNumber+i) {
				f.config.FirstRemovedLineNumber = f.config.FirstRemovedLineNumber + i
				break
			}
		}
	}

	number := f.config.FirstRemovedLineNumber - removed[0]
	for _, line := range f.diffLines {
		if line.Type == Inserted {
			line.Number = NilLineIndex
			continue
		}
		line.Number = number
		number++
	}
	return nil
}

func (f *FormattedDiff) SyncLineFormatting(lines []string) {
	text := lines[f.config.FirstRemovedLineNumber]
	trim := strings.TrimSpace(text)
	truePadding := len(text) - len(trim)

	diffLineText := f.diffLines[f.firstRemovedLineIndex()].Text
	diffLineTrim := strings.TrimSpace(diffLineText)
	diffPadding := len(diffLineText) - len(diffLineTrim)

	padding := truePadding - diffPadding

	for _, line := range f.diffLines {
		line.Text = strings.Repeat(" ", padding) + line.Text
	}
}

func (f *FormattedDiff) Lines() DiffLines {
	return f.diffLines
}
