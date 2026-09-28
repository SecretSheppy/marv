package littledarwin

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/internal/mutations"
	"github.com/SecretSheppy/marv/pkg/fio"
	"github.com/aymanbagabas/go-udiff"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var meta = fwlib.Meta{
	Name: "littledarwin",
	URL:  "https://github.com/aliparsai/LittleDarwin",
}

type YamlConfig struct {
	LittleDarwinResultsDir string `yaml:"results-dir"`
}

type YamlWrapper struct {
	Cfg *YamlConfig `yaml:"littledarwin"`
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
	return y.Cfg.LittleDarwinResultsDir != "", nil
}

type Mutation struct {
	SourceFilePath  string
	MutatedFilePath string
	Status          mutations.Status
}

var re = regexp.MustCompile("(.+):\\ssurvived\\s\\(\\d+/\\d+\\)\\s->\\s(\\[[0-9'.,jav ]*])\\s-\\skilled\\s\\(\\d+/\\d+\\)\\s->\\s(\\[[0-9'.,jav ]*])")

type LittleDarwin struct {
	yml   *YamlWrapper
	muts  []*Mutation
	ms    mutations.Mutations
	files map[string][]string
}

func NewLittleDarwin() *LittleDarwin {
	return &LittleDarwin{yml: &YamlWrapper{}}
}

func (l *LittleDarwin) Meta() *fwlib.Meta {
	return &meta
}

func (l *LittleDarwin) Yaml() fwlib.FWConfig {
	return l.yml
}

func (l *LittleDarwin) addMutationEntries(sourcePath string, mutationFiles []string, status mutations.Status) {
	for _, m := range mutationFiles {
		l.muts = append(l.muts, &Mutation{
			SourceFilePath:  sourcePath,
			MutatedFilePath: path.Join(sourcePath, m),
			Status:          status,
		})
	}
}

func (l *LittleDarwin) LoadResults() error {
	log.Info().Msgf("%s - loading results", l.Meta().Name)

	file, err := os.ReadFile(path.Join(l.yml.Cfg.LittleDarwinResultsDir, "report.txt"))
	if err != nil {
		return err
	}

	matches := re.FindAllStringSubmatch(string(file), -1)
	for _, match := range matches {
		sourcePath := match[1]
		survivedRaw := strings.ReplaceAll(match[2], "'", "\"")
		killedRaw := strings.ReplaceAll(match[3], "'", "\"")

		var survived, killed []string
		if err = json.Unmarshal([]byte(survivedRaw), &survived); err != nil {
			return err
		}
		if err = json.Unmarshal([]byte(killedRaw), &killed); err != nil {
			return err
		}

		l.addMutationEntries(sourcePath, survived, mutations.Survived)
		l.addMutationEntries(sourcePath, killed, mutations.Killed)
	}
	return nil
}

func (l *LittleDarwin) TransformResults() error {
	log.Info().Msgf("%s - transforming results", l.Meta().Name)
	bar := fwlib.NewProgressbar(len(l.muts), "transforming")
	l.ms = make(mutations.Mutations)
	l.files = make(map[string][]string)

	for _, mutation := range l.muts {
		if l.files[mutation.SourceFilePath] == nil {
			lines, err := fio.ReadLines(path.Join(l.yml.Cfg.LittleDarwinResultsDir, mutation.SourceFilePath, "original.java"))
			if err != nil {
				return err
			}
			l.files[mutation.SourceFilePath] = lines
		}

		file, err := os.Open(path.Join(l.yml.Cfg.LittleDarwinResultsDir, mutation.MutatedFilePath))
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lines := make([]string, 0)
		for i := 0; i < 5; i++ {
			scanner.Scan()
			lines = append(lines, scanner.Text())
		}

		operator := strings.TrimPrefix(lines[1], "mutant type: ")
		before := strings.TrimPrefix(lines[2], "----> before: ")
		after := strings.TrimPrefix(lines[3], "----> after: ")
		strLineNum := strings.TrimPrefix(lines[4], "----> line number in original file: ")
		lineNum, err := strconv.Atoi(strLineNum)
		if err != nil {
			return err
		}

		// NOTE: if there is no difference between before and after, then the line the mutation is on is very long, and
		// has been broken across several lines in the file, despite being treated as one line by the parser. To find
		// the mutant, we do the extra content below.
		if before == after {
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}
			source := strings.Join(l.files[mutation.SourceFilePath], "\n")
			mutant := strings.Join(lines[9:], "\n")
			edits := udiff.Strings(source, mutant)
			diff, err := udiff.ToUnifiedDiff("original", "new", source, edits, 0)
			if err != nil {
				return err
			}
			hunk := diff.Hunks[0]
			lineNum = hunk.FromLine
			before = hunk.Lines[0].Content
			after = hunk.Lines[1].Content
		}

		edits := udiff.Strings(before, after)
		edit := edits[0]

		m := &mutations.Mutation{
			FrameworkMutantID: mutation.MutatedFilePath,
			Description:       fmt.Sprintf("Replaced `%s` with `%s`", before[edit.Start:edit.End], edit.New),
			Operation:         operator,
			Start: &mutations.Range{
				Line: lineNum - 1,
				Char: edit.Start,
			},
			End: &mutations.Range{
				Line: lineNum - 1,
				Char: edit.End,
			},
			Status:      mutation.Status,
			Replacement: edit.New,
		}

		l.ms.Append(mutation.SourceFilePath, m)
		bar.Add(1)
	}

	fwlib.FinishProgressbar(bar)
	return nil
}

func (l *LittleDarwin) Mutations() mutations.Mutations {
	return l.ms
}

func (l *LittleDarwin) ReadLines(file string) ([]string, error) {
	return l.files[file], nil
}
