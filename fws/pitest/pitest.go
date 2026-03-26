package pitest

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"

	"github.com/SecretSheppy/marv/decompilers"
	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/pkg/mutations"
	"github.com/aymanbagabas/go-udiff"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
	"gopkg.in/yaml.v3"
)

// YamlConfig represents Pitest's yml config data.
type YamlConfig struct {
	XmlPath      string `yaml:"xml-path"`
	SrcCodePath  string `yaml:"src-code-path"`
	SrcClassPath string `yaml:"src-class-path"`
	MutClassPath string `yaml:"mut-class-path"`
	Decompiler   string `yaml:"decompiler"`
}

// YamlWrapper used to load the pitest configuration from the .marv.yml file.
type YamlWrapper struct {
	Cfg *YamlConfig `yaml:"pitest"`
}

func (y *YamlWrapper) Init() interface{} {
	return &YamlWrapper{Cfg: &YamlConfig{}}
}

func (y *YamlWrapper) Load(yml []byte) (bool, error) {
	if err := yaml.Unmarshal(yml, y); err != nil {
		return false, err
	}
	return y.Cfg.XmlPath != "" || y.Cfg.SrcCodePath != "" || y.Cfg.SrcClassPath != "" || y.Cfg.MutClassPath != "", nil
}

// Mutation is a struct that can accept the pitest xml output.
type Mutation struct {
	Detected          bool             `xml:"detected,attr"`
	Status            mutations.Status `xml:"status,attr"`
	NumTestsRun       int              `xml:"numberOfTestsRun,attr"`
	SourceFile        string           `xml:"sourceFile"`
	MutatedClass      string           `xml:"mutatedClass"`
	MutatedMethod     string           `xml:"mutatedMethod"`
	MethodDescription string           `xml:"methodDescription"`
	LineNumber        int              `xml:"lineNumber"`
	Mutator           string           `xml:"mutator"`
	KillingTest       string           `xml:"killingTest"`
	Description       string           `xml:"description"`
	MutationIndex     int              // NOTE: used to determine which mutants to decompile. Not present in XML.
}

func (m *Mutation) SourceCodePath() string {
	base := strings.ReplaceAll(m.MutatedClass, ".", "/")
	return path.Join(path.Dir(base), m.SourceFile)
}

func (m *Mutation) SourceClassPath() string {
	return strings.ReplaceAll(m.MutatedClass, ".", "/") + ".class"
}

func (m *Mutation) MutantExportDir() string {
	base := strings.ReplaceAll(m.MutatedClass, ".", "/")
	return path.Join(base, "mutants")
}

func (m *Mutation) MutatedClassPath() string {
	file := fmt.Sprintf("%d/%s.class", m.MutationIndex, m.MutatedClass)
	return path.Join(m.MutantExportDir(), file)
}

type PitXML struct {
	XMLName   xml.Name    `xml:"mutations"`
	Mutations []*Mutation `xml:"mutation"`
}

// Pitest is the Framework object for the pitest library.
type Pitest struct {
	yml   *YamlWrapper
	muts  []*Mutation
	ms    mutations.Mutations
	dcomp decompilers.Decompiler
}

func NewPitest() *Pitest {
	return &Pitest{yml: &YamlWrapper{}}
}

func (p *Pitest) SetDecompiler() {
	p.dcomp = decompilers.JavaDecompiler(p.yml.Cfg.Decompiler)
}

func (p *Pitest) Meta() *fwlib.Meta {
	return &fwlib.Meta{
		Name:      "Pitest",
		Extension: "java",
		URL:       "https://pitest.org/",
	}
}

func (p *Pitest) Yaml() fwlib.FWConfig {
	return p.yml
}

func (p *Pitest) LoadResults() error {
	log.Warn().Msgf("%s - experimental framwork that relies on decompilation of binary files. results may not be 100%% accurate!", p.Meta().Name)
	log.Info().Msgf("%s - loading results", p.Meta().Name)

	rawxml, err := os.ReadFile(p.yml.Cfg.XmlPath)
	if err != nil {
		return err
	}

	pitxml := &PitXML{}
	if err := xml.Unmarshal(rawxml, pitxml); err != nil {
		return err
	}
	p.muts = pitxml.Mutations

	return nil
}

func (p *Pitest) TransformResults() error {
	log.Info().Msgf("%s - transforming results", p.Meta().Name)
	log.Warn().Msgf("%s - mutations of status NO_COVERAGE will not be decompiled", p.Meta().Name)

	groupBar := progressbar.NewOptions(
		len(p.muts),
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionSetDescription("[1/3]     grouping"),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionShowCount())
	msMap := p.groupMutations(groupBar)
	groupBar.Finish()
	fmt.Println()

	indexBar := progressbar.NewOptions(
		len(p.muts),
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionSetDescription("[2/3]     indexing"),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionShowCount())
	err := p.indexMutations(msMap, indexBar)
	if err != nil {
		return err
	}
	indexBar.Finish()
	fmt.Println()

	log.Info().Msgf("%s - using %s", p.Meta().Name, p.dcomp)

	if err := p.dcomp.Setup(); err != nil {
		return err
	}
	transformBar := progressbar.NewOptions(
		len(msMap),
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionSetDescription("[3/3] transforming"),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionShowCount())
	p.transformMutations(msMap, transformBar)
	transformBar.Finish()
	fmt.Println()
	return p.dcomp.Teardown()
}

// groups mutations by class name (so all mutations of the same file will be grouped together)
func (p *Pitest) groupMutations(bar *progressbar.ProgressBar) map[string][]*Mutation {
	msMap := make(map[string][]*Mutation)
	for _, m := range p.muts {
		msMap[m.SourceCodePath()] = append(msMap[m.SourceCodePath()], m)
		bar.Add(1)
	}
	return msMap
}

func (p *Pitest) indexMutations(msMap map[string][]*Mutation, bar *progressbar.ProgressBar) error {
	for _, ms := range msMap {
		for _, m := range ms {
			for i := 0; i < len(ms); i++ {
				details := fmt.Sprintf("%d/details.txt", i)
				txtpth := path.Join(p.yml.Cfg.MutClassPath, m.MutantExportDir(), details)
				data, err := os.ReadFile(txtpth)
				if err != nil {
					return err
				}
				vals := fmt.Sprintf("lineNumber=%d, description=%s", m.LineNumber, m.Description)
				if strings.Contains(string(data), vals) {
					m.MutationIndex = i
					break
				}
			}
			bar.Add(1)
		}
	}
	return nil
}

type TransformWorkerJob struct {
	Mutations []*Mutation
}

type TransformWorkerResult struct {
	Mutations mutations.Mutations
	Error     error
}

func transformMutationsWorker(jobs <-chan TransformWorkerJob, results chan<- TransformWorkerResult, wg *sync.WaitGroup, cfg *YamlConfig, bar *progressbar.ProgressBar, dcomp decompilers.Decompiler) {
	defer wg.Done()

	for job := range jobs {

		// NOTE: Errors will be printed into stderr but will not interrupt the main process. Files for which this
		// process fails will just be left out of the visualizations.
		srcCodePath := path.Join(cfg.SrcCodePath, job.Mutations[0].SourceCodePath())
		rawSrcCode, err := os.ReadFile(srcCodePath)
		if err != nil {
			results <- TransformWorkerResult{
				Error: errors.New("failed to process mutations for " + srcCodePath),
			}
			bar.Add(1)
			continue
		}

		var srcCodeLines []string
		lines := strings.Lines(string(rawSrcCode))
		for line := range lines {
			srcCodeLines = append(srcCodeLines, line)
		}

		// Cache the decompiled class files where they can be reused.
		decompiled := make(map[string]string)

		ms := make(mutations.Mutations)
		// Adds the mutation to the mutations map.
		appendMutation := func(pth string, mu *mutations.Mutation) {
			added := false
			for _, c := range ms[pth] {
				if c.Conflicts(mu) {
					c.Append(mu)
					added = true
					break
				}
			}

			if !added {
				ms[pth] = append(ms[pth], mutations.NewConflict(mu))
			}
		}

		for _, m := range job.Mutations {
			srcLine := srcCodeLines[m.LineNumber-1]

			// NOTE: Ignore mutations with status NO_COVERAGE in order to save time. Marv is not useful where there is
			// no coverage, so this does not affect the quality of its visualizations.
			if m.Status == mutations.NoCoverage {
				appendMutation(m.SourceCodePath(), streamlineMutation(
					m,
					&mutations.Range{Line: m.LineNumber - 1},
					&mutations.Range{Line: m.LineNumber - 1, Char: len(srcLine) - 1}))
				continue
			}

			srcClassPath := path.Join(cfg.SrcClassPath, m.SourceClassPath())
			if decompiled[srcClassPath] == "" {
				dcomp, err := dcomp.Decompile(srcClassPath)
				if err != nil {
					// FIXME: handle this error
					panic(err)
				}
				decompiled[srcClassPath] = string(dcomp)
			}

			mutClassPath := path.Join(cfg.MutClassPath, m.MutatedClassPath())
			mutated, err := dcomp.Decompile(mutClassPath)
			if err != nil {
				// FIXME: handle this error
				panic(err)
			}

			edits := udiff.Strings(decompiled[srcClassPath], string(mutated))
			d, err := udiff.ToUnifiedDiff("old", "new", decompiled[srcClassPath], edits, 0)
			if err != nil {
				// FIXME: handle this error
				panic(err)
			}

			rmlines := 0
			builder := strings.Builder{}
			for _, h := range d.Hunks {
				if strings.Contains(h.Lines[0].Content, "import") {
					continue
				}
				for _, l := range h.Lines {
					switch l.Kind {
					case udiff.Delete:
						rmlines++
					case udiff.Insert, udiff.Equal:
						builder.WriteString(l.Content)
					}
				}
			}

			// NOTE: m.LineNumber is technically the first line removed in rmlines, so hence the -2 here and -1 in
			// the below mutations.Range.
			srcEndLine := srcCodeLines[m.LineNumber+rmlines-2]
			mutant := streamlineMutation(
				m,
				&mutations.Range{Line: m.LineNumber - 1},
				&mutations.Range{Line: m.LineNumber + rmlines - 2, Char: len(srcEndLine) - 1})
			mutant.Source = builder.String()
			appendMutation(m.SourceCodePath(), mutant)
		}

		results <- TransformWorkerResult{Mutations: ms}
		bar.Add(1)
	}
}

func (p *Pitest) transformMutations(ms map[string][]*Mutation, bar *progressbar.ProgressBar) {
	numWorkers := runtime.NumCPU()
	jobs := make(chan TransformWorkerJob, len(ms))
	results := make(chan TransformWorkerResult, len(ms))

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go transformMutationsWorker(jobs, results, &wg, p.yml.Cfg, bar, p.dcomp)
	}

	for _, fileMutations := range ms {
		jobs <- TransformWorkerJob{Mutations: fileMutations}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	p.ms = make(mutations.Mutations)
	for result := range results {
		if result.Error != nil {
			// FIXME: this is not printing errors
			log.Error().Err(result.Error)
		}
		for k, v := range result.Mutations {
			p.ms[k] = v
		}
	}
}

func streamlineMutation(m *Mutation, starts, ends *mutations.Range) *mutations.Mutation {
	return &mutations.Mutation{
		Name:   m.Description,
		OpDesc: m.Mutator,
		Status: m.Status,
		Starts: starts,
		Ends:   ends,
	}
}

func (p *Pitest) Mutations() mutations.Mutations {
	return p.ms
}
