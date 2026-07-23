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
	PrefixLines, SuffixLines, FirstRemovedLineNumber int
	IgnoreBlankLines                                 bool
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

func (f *FormattedDiff) firstRemovedLineIndex() int {
	for i, line := range f.diffLines {
		if line.Type == Removed {
			return i
		}
	}
	return NilLineIndex
}

func (f *FormattedDiff) Number() error {
	firstRemoved := f.firstRemovedLineIndex()
	if firstRemoved == NilLineIndex {
		return errors.New("must set DiffConfig.FirstRemovedLineNumber to number diff lines")
	}
	number := f.config.FirstRemovedLineNumber - firstRemoved
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

func (f *FormattedDiff) SyncLineFormatting(source []string) {
	text := source[f.config.FirstRemovedLineNumber]
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
