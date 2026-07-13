package mewt

import (
	"fmt"
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

func (m *Mewt) TransformResults() error {
	log.Info().Msgf("%s - transforming results", m.Meta().Name)

	var targets []Target
	if err := m.mdb.Preload("Mutants").Preload("Mutants.Outcome").Find(&targets).Error; err != nil {
		return err
	}

	// TODO: add progress bar

	m.ms = make(mutations.Mutations)
	m.files = make(map[string][]string)
	for _, file := range targets {
		m.addFile(file.Path, file.Text)

		// TODO: this whole section needs to be re thought through
		for _, mutant := range file.Mutants {
			originalLines := strings.Split(mutant.OldText, "\n")
			m.ms.Append(file.Path, &mutations.Mutation{
				ID:                uuid.New(),
				FrameworkMutantID: strconv.Itoa(mutant.MutantID),
				Description:       "", // TODO: generate description
				Operation:         mutant.MutationSlug,
				Start: &mutations.Range{
					Line: mutant.LineOffset,
					Char: 0, // TODO: might need to use byte offset to find actual position in line
				},
				End: &mutations.Range{
					Line: mutant.LineOffset + len(originalLines) - 1,
					Char: len(originalLines[len(originalLines)-1]) - 1, // TODO: this is wrong and doesn't account for start offset
				},
				Status:      mutant.status(),
				Replacement: mutant.NewText,
			})
		}
	}
	fmt.Println(m)
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
