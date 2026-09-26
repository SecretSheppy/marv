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

// trims the provided line indexes if the contents
func (d DiffLines) trimLinesIfBothBlank(removed, inserted int) DiffLines {
	removedLine := strings.TrimSpace(d[removed].Text)
	insertedLine := strings.TrimSpace(d[inserted].Text)
	trim := removedLine == "" && insertedLine == ""
	trimmed := make(DiffLines, 0)
	for i, line := range d {
		if trim && (i == removed || i == inserted) {
			continue
		}
		trimmed = append(trimmed, line)
	}
	return trimmed
}

// trims and leading blank lines.
func (d DiffLines) trimLeadingLines() DiffLines {
	removed, inserted := NilLineIndex, NilLineIndex

	for i, line := range d {
		if line.Type == Removed && removed == NilLineIndex {
			removed = i
		}
		if line.Type == Inserted && inserted == NilLineIndex {
			inserted = i
		}
	}

	return d.trimLinesIfBothBlank(removed, inserted)
}

// trims any trailing blank lines.
func (d DiffLines) trimTrailingLines() DiffLines {
	var removed, inserted int

	for i, line := range d {
		switch line.Type {
		case Removed:
			removed = i
		case Inserted:
			inserted = i
		}
	}

	return d.trimLinesIfBothBlank(removed, inserted)
}

// trims leading and trailing blank lines that are sometimes erroneously added to diffs.
func (d DiffLines) trimLeadingAndTrailingLines() DiffLines {
	removed, inserted := d.LineChanges()
	if len(removed) <= 1 || len(inserted) <= 1 {
		return d
	}
	trimmed := d.trimLeadingLines()
	removed, inserted = trimmed.LineChanges()
	if len(removed) <= 1 || len(inserted) <= 1 {
		return trimmed
	}
	return trimmed.trimTrailingLines()
}

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
		diffLines: diffLines.trimLeadingAndTrailingLines(),
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
