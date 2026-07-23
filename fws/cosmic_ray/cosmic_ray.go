package cosmic_ray

import (
	"fmt"
	"path"
	"strings"

	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/internal/mutations"
	"github.com/SecretSheppy/marv/pkg/diffutil"
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

func (m *MutationResult) StartLine() int {
	return m.StartPosRow - 1
}

func (m *MutationResult) StartChar() int {
	return m.StartPosCol
}

func (m *MutationResult) EndLine() int {
	return m.EndPosRow - 1
}

func (m *MutationResult) EndChar() int {
	return m.EndPosCol
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
		lines, err := c.getOrAddFile(result.ModulePath)
		if err != nil {
			return err
		}

		diff := diffutil.FromFormattedDiff(result.Diff, &diffutil.DiffConfig{
			PrefixLines:            4,
			FirstRemovedLineNumber: result.StartLine(),
			IgnoreBlankLines:       true,
		})
		if err = diff.Number(); err != nil {
			return err
		}
		diff.SyncLineFormatting(lines)

		removed, inserted := diff.Lines().LineChanges()
		prefix := removed.Get(result.StartLine()).Text[:result.StartChar()]
		suffix := removed.Get(result.EndLine()).Text[result.EndChar():]
		start := len(prefix)
		end := len(suffix)

		rem := strings.Join(removed.StringLines(), "\n")
		ins := strings.Join(inserted.StringLines(), "\n")

		original := rem[start : len(rem)-end]
		replacement := ins[start : len(ins)-end]

		c.ms.Append(result.ModulePath, &mutations.Mutation{
			ID:                uuid.New(),
			FrameworkMutantID: result.JobID.String(),
			Description:       fmt.Sprintf("Replaced `%s` with `%s`", original, replacement),
			Operation:         result.OperatorName,
			Start: &mutations.Range{
				Line: result.StartLine(),
				Char: result.StartChar(),
			},
			End: &mutations.Range{
				Line: result.EndLine(),
				Char: result.EndChar(),
			},
			Status:      result.status(),
			Replacement: replacement,
		})
	}

	fwlib.FinishProgressbar(bar)
	return nil
}

func (c *CosmicRay) Mutations() mutations.Mutations {
	return c.ms
}

func (c *CosmicRay) getOrAddFile(file string) ([]string, error) {
	if c.files[file] == nil {
		fp := path.Join(c.yml.Cfg.CRWorkDir, file)
		lines, err := fio.ReadLines(fp)
		if err != nil {
			return nil, err
		}
		c.files[file] = lines
	}
	return c.ReadLines(file)
}

func (c *CosmicRay) ReadLines(file string) ([]string, error) {
	return c.files[file], nil
}
