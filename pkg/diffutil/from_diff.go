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

var ErrNoRemovedLines = errors.New("no removed lines in text diff")

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

func FromFormattedDiff(sourceLines []string, diff string, config *DiffConfig) (*FormattedDiff, error) {
	var (
		lines     = strings.Split(diff, "\n")
		diffLines = make(DiffLines, 0)
		from      = config.PrefixLines
		to        = len(lines) - config.SuffixLines
	)

	for _, line := range lines[from:to] {
		diffLines = append(diffLines, &DiffLine{Type: lineType(line), Text: line[1:]})
	}

	fDiff := &FormattedDiff{
		config:    config,
		diffLines: diffLines,
	}
	if err := fDiff.syncLineNumbers(sourceLines); err != nil {
		return nil, err
	}
	if err := fDiff.syncLineFormatting(sourceLines); err != nil {
		return nil, err
	}
	return fDiff, nil
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

// verifies whether the diff line content matches the source line content. the indent formatting from both lines is
// removed in the process to ensure that any formatting differences from the text diff do not affect the outcome.
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

// syncs the DiffLine numbers with the original source code. This is used for cases where the produced
// text diff may be slightly different to the actual change made (i.e. an extra deleted blank line is present in the
// diff before the expected first line.)
func (f *FormattedDiff) syncLineNumbers(lines []string) error {
	removed := f.removedLineIndexes()
	if len(removed) == 0 {
		return ErrNoRemovedLines
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

// syncs the line formatting from the source lines into the diff.
func (f *FormattedDiff) syncLineFormatting(lines []string) error {
	text := lines[f.config.FirstRemovedLineNumber]
	trim := strings.TrimSpace(text)
	truePadding := len(text) - len(trim)

	removed := f.removedLineIndexes()
	if len(removed) == 0 {
		return ErrNoRemovedLines
	}
	first := removed[0]

	diffLineText := f.diffLines[first].Text
	diffLineTrim := strings.TrimSpace(diffLineText)
	diffPadding := len(diffLineText) - len(diffLineTrim)

	padding := truePadding - diffPadding

	for _, line := range f.diffLines {
		line.Text = strings.Repeat(" ", padding) + line.Text
	}
	return nil
}

func (f *FormattedDiff) Lines() DiffLines {
	return f.diffLines
}
