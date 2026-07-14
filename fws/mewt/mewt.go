package mewt

import (
	"strconv"
	"strings"

	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/internal/mutations"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var meta = fwlib.Meta{
	Name: "mewt",
	URL:  "https://github.com/trailofbits/mewt",
}

type YamlConfig struct {
	SQLiteDB string `yaml:"sqlite-path"`
}

type YamlWrapper struct {
	Cfg *YamlConfig `yaml:"mewt"`
}

func (y *YamlWrapper) Init() interface{} {
	return &YamlWrapper{Cfg: &YamlConfig{}}
}

func (y *YamlWrapper) Load(yml []byte) (bool, error) {
	if err := yaml.Unmarshal(yml, y); err != nil {
		return false, err
	}
	if y.Cfg == nil {
		return false, nil
	}
	return y.Cfg.SQLiteDB != "", nil
}

type Target struct {
	TargetID int `gorm:"primary_key;column:id"`
	Path     string
	Text     string
	Mutants  []*Mutant
}

type Mutant struct {
	MutantID     int `gorm:"primary_key;column:id"`
	TargetID     int `gorm:"foreign_key:TargetID"`
	ByteOffset   int
	LineOffset   int
	OldText      string
	NewText      string
	MutationSlug string
	Outcome      *Outcome
}

func (m *Mutant) status() mutations.Status {
	switch m.Outcome.Status {
	case "TestFail":
		return mutations.Killed
	case "Skipped":
		return mutations.Ignored
	case "Timeout":
		return mutations.Timeout
	default: // "Uncaught"
		return mutations.Survived
	}
}

func (m *Mutant) operator() (string, string) {
	switch m.MutationSlug {
	case "ER":
		return "Error Replacement", "Replaced statement with an error"
	case "CR":
		return "Comment Replacement", "Replaced statement with an in-line comment"
	case "IF":
		return "If False", "Hardcoded an if condition to false"
	case "IT":
		return "If True", "Hardcoded an if condition to true"
	case "WF":
		return "While False", "Hardcoded while condition to false"
	case "AS":
		return "Argument Swap", "Swaped pairs of adjacent arguments"
	case "LC":
		return "Loop Control", "Swaped break and continue statements"
	case "BL":
		return "Boolean Literal Flip", "Flipped boolean literal true <-> false"
	case "AOS":
		return "Arithmetic Operator Shuffle", "Replaced arithmetic operators (+, -, *, /)"
	case "AAOS":
		return "Arithmetic Assignment Operator Shuffle", "Replaced arithmetic assignment operators (+=, -=, *=, /=)"
	case "BOS":
		return "Bitwise Operator Shuffle", "Replaced bitwise operators (&, |, ^)"
	case "BAOS":
		return "Bitwise Assignment Operator Shuffle", "Replaced bitwise assignment operators (&=, |=, ^=)"
	case "LOS":
		return "Logical Operator Shuffle", "Replaced logical operators (&&, ||)"
	case "COS":
		return "Comparison Operator Shuffle", "Replaced comparison operators (==, !=, <, <=, >, >=)"
	case "SOS":
		return "Shift Operator Shuffle", "Replaced shift operators (<<, >>)"
	case "SAOS":
		return "Shift Assignment Operator Shuffle", "Replaced shift assignment operators (<<=, >>=)"
	case "NR":
		return "Negation Removal", "Removed logical negation operator (!x -> x)"
	default:
		return m.MutationSlug, "Marv: Unknown Operator"
	}
}

type Outcome struct {
	MutantID int `gorm:"primary_key;foreign_key:MutantID"`
	Status   string
}

type Mewt struct {
	yml   *YamlWrapper
	mdb   *gorm.DB
	ms    mutations.Mutations
	files map[string][]string
}

func NewMewt() *Mewt {
	return &Mewt{yml: &YamlWrapper{}}
}

func (m *Mewt) Meta() *fwlib.Meta {
	return &meta
}

func (m *Mewt) Yaml() fwlib.FWConfig {
	return m.yml
}

func (m *Mewt) LoadResults() error {
	log.Info().Msgf("%s - loading results", m.Meta().Name)

	var err error
	m.mdb, err = gorm.Open(sqlite.Open(m.yml.Cfg.SQLiteDB), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	return err
}

func (m *Mewt) getOriginalCharacterOffset(byteOffset int, file []byte) int {
	count := -1
	for i := byteOffset; i > 0; i-- {
		if file[i] == '\n' {
			return count
		}
		count++
	}
	return count
}

func (m *Mewt) TransformResults() error {
	log.Info().Msgf("%s - transforming results", m.Meta().Name)

	var targets []Target
	if err := m.mdb.Preload("Mutants").Preload("Mutants.Outcome").Find(&targets).Error; err != nil {
		return err
	}

	size := 0
	for _, target := range targets {
		size += len(target.Mutants)
	}
	bar := fwlib.NewProgressbar(size, "transforming")

	m.ms = make(mutations.Mutations)
	m.files = make(map[string][]string)
	for _, file := range targets {
		m.addFile(file.Path, file.Text)
		byteFile := []byte(file.Text)

		for _, mutant := range file.Mutants {
			op, opDesc := mutant.operator()
			lines := strings.Split(mutant.OldText, "\n")
			singleLineCharOffset := m.getOriginalCharacterOffset(mutant.ByteOffset, byteFile)
			endCharOffset := len(lines[len(lines)-1])
			if len(lines) == 1 {
				endCharOffset += singleLineCharOffset
			}

			m.ms.Append(file.Path, &mutations.Mutation{
				ID:                uuid.New(),
				FrameworkMutantID: strconv.Itoa(mutant.MutantID),
				Description:       opDesc,
				Operation:         op,
				Start: &mutations.Range{
					Line: mutant.LineOffset,
					Char: singleLineCharOffset,
				},
				End: &mutations.Range{
					Line: mutant.LineOffset + len(lines) - 1,
					Char: endCharOffset,
				},
				Status:      mutant.status(),
				Replacement: mutant.NewText,
			})
			bar.Add(1)
		}
	}
	fwlib.FinishProgressbar(bar)
	return nil
}

func (m *Mewt) addFile(path, text string) {
	lines := make([]string, 0)
	for line := range strings.Lines(text) {
		lines = append(lines, strings.ReplaceAll(line, "\n", ""))
	}
	m.files[path] = lines
}

func (m *Mewt) Mutations() mutations.Mutations {
	return m.ms
}

func (m *Mewt) ReadLines(file string) ([]string, error) {
	return m.files[file], nil
}
