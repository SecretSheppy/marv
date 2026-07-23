package cosmic_ray

import (
	"fmt"
	"path"
	"strings"

	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/internal/mutations"
	"github.com/SecretSheppy/marv/pkg/fio"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var meta = fwlib.Meta{
	Name: "cosmic-ray",
	URL:  "https://github.com/sixty-north/cosmic-ray",
}

type YamlConfig struct {
	SQLiteDB  string `yaml:"sqlite-path"`
	CRWorkDir string `yaml:"cr-work-dir"`
}

type YamlWrapper struct {
	Cfg *YamlConfig `yaml:"cosmic-ray"`
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
	return y.Cfg.SQLiteDB != "" && y.Cfg.CRWorkDir != "", nil
}

type MutationResult struct {
	ModulePath   string
	OperatorName string
	StartPosRow  int
	StartPosCol  int
	EndPosRow    int
	EndPosCol    int
	JobID        uuid.UUID
	TestOutcome  string
	Diff         string
}

func (m *MutationResult) status() mutations.Status {
	switch m.TestOutcome {
	case "SURVIVED":
		return mutations.Survived
	case "KILLED":
		return mutations.Killed
	default: // INCOMPETENT
		return mutations.Timeout
	}
}

func changes(diff string) (remove, insert string) {
	for _, line := range strings.Split(diff, "\n")[4:] {
		if line[:1] == "-" && line != "-" {
			remove = strings.TrimSpace(line[1:])
			continue
		}
		if line[:1] == "+" && line != "+" {
			insert = strings.TrimSpace(line[1:])
			continue
		}
	}
	return
}

type CosmicRay struct {
	yml     *YamlWrapper
	results []*MutationResult
	ms      mutations.Mutations
	files   map[string][]string
}

func NewCosmicRay() *CosmicRay {
	return &CosmicRay{yml: &YamlWrapper{}}
}

func (c *CosmicRay) Meta() *fwlib.Meta {
	return &meta
}

func (c *CosmicRay) Yaml() fwlib.FWConfig {
	return c.yml
}

func (c *CosmicRay) LoadResults() error {
	log.Info().Msgf("%s - loading results", c.Meta().Name)

	db, err := gorm.Open(sqlite.Open(c.yml.Cfg.SQLiteDB), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}
	return db.
		Table("mutation_specs").
		Joins("left join work_results on work_results.job_id = mutation_specs.job_id").
		Scan(&c.results).Error
}

func (c *CosmicRay) TransformResults() error {
	log.Info().Msgf("%s - transforming results", c.Meta().Name)

	c.ms = make(mutations.Mutations)
	c.files = make(map[string][]string)
	bar := fwlib.NewProgressbar(len(c.results), "transforming")

	for _, result := range c.results {
		if c.files[result.ModulePath] == nil {
			lines, err := fio.ReadLines(path.Join(c.yml.Cfg.CRWorkDir, result.ModulePath))
			if err != nil {
				return err
			}
			c.files[result.ModulePath] = lines
		}
		lines := c.files[result.ModulePath]
		remove, insert := changes(result.Diff)
		originalLine := lines[result.StartPosRow-1]
		padding := len(originalLine) - len(strings.TrimSpace(originalLine))
		fromStart := len(remove[:result.StartPosCol-padding])
		fromEnd := len(remove[result.EndPosCol-padding:])
		originalCode := remove[fromStart : len(remove)-fromEnd]
		newCode := insert[fromStart : len(insert)-fromEnd]

		c.ms.Append(result.ModulePath, &mutations.Mutation{
			ID:                uuid.New(),
			FrameworkMutantID: result.JobID.String(),
			Description:       fmt.Sprintf("Replaced `%s` with `%s`", originalCode, newCode),
			Operation:         result.OperatorName,
			Start: &mutations.Range{
				Line: result.StartPosRow - 1,
				Char: result.StartPosCol,
			},
			End: &mutations.Range{
				Line: result.EndPosRow - 1,
				Char: result.EndPosCol,
			},
			Status:      result.status(),
			Replacement: newCode,
		})
	}

	fwlib.FinishProgressbar(bar)
	return nil
}

func (c *CosmicRay) Mutations() mutations.Mutations {
	return c.ms
}

func (c *CosmicRay) ReadLines(file string) ([]string, error) {
	return c.files[file], nil
}
