package pitest

import (
	"encoding/xml"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/SecretSheppy/marv/decompilers"
	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/internal/mutations"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
	"gopkg.in/yaml.v3"
)

const FWName = "Pitest"

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
	if y.Cfg == nil {
		return false, nil
	}
	return y.Cfg.XmlPath != "" || y.Cfg.SrcCodePath != "" || y.Cfg.SrcClassPath != "" || y.Cfg.MutClassPath != "", nil
}

func (y *YamlWrapper) SourceCodeDir() string {
	return y.Cfg.SrcCodePath
}

// FileMutations stores all mutants of the same file against the string of the full file path.
type FileMutations map[string][]*Mutation

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

func (p *Pitest) config() *YamlConfig {
	return p.yml.Cfg
}

func (p *Pitest) decompiler() decompilers.Decompiler {
	return p.dcomp
}

func (p *Pitest) SetDecompiler() {
	p.dcomp = decompilers.JavaDecompiler(p.yml.Cfg.Decompiler)
}

func (p *Pitest) Meta() *fwlib.Meta {
	return &fwlib.Meta{
		Name:      FWName,
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

	groupBar := progressbar.NewOptions(
		len(p.muts),
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionSetDescription("[1/3]     grouping"),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionShowCount())
	fileMutations := p.groupMutants(groupBar)
	groupBar.Finish()
	fmt.Println()

	indexBar := progressbar.NewOptions(
		len(p.muts),
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionSetDescription("[2/3]     indexing"),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionShowCount())
	err := p.indexMutants(fileMutations, indexBar)
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
		len(fileMutations),
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionSetDescription("[3/3] transforming"),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionShowCount())
	var errs []error
	p.ms, errs = transform(p, fileMutations, transformBar)
	// NOTE: perform stdout cleanup before printing errors.
	transformBar.Finish()
	fmt.Println()
	if len(errs) > 0 {
		for _, err := range errs {
			err.(*transformError).log()
		}
	}
	return p.dcomp.Teardown()
}

func (p *Pitest) Mutations() mutations.Mutations {
	return p.ms
}

func streamlineMutation(m *Mutation, starts, ends *mutations.Range) *mutations.Mutation {
	return &mutations.Mutation{
		Description: m.Description,
		Operation:   m.Mutator,
		Status:      m.Status,
		Start:       starts,
		End:         ends,
	}
}
