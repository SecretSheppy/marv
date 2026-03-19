package pitest

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/pkg/mutations"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

// YamlConfig represents Pitest's yml config data.
type YamlConfig struct {
	XmlPath         string `yaml:"xml-path"`
	SrcPath         string `yaml:"src-path"`
	SrcBytecodePath string `yaml:"src-bytecode-path"`
	MutBytecodePath string `yaml:"mut-bytecode-path"`
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
	return y.Cfg.XmlPath != "" || y.Cfg.SrcPath != "" || y.Cfg.SrcBytecodePath != "" || y.Cfg.MutBytecodePath != "", nil
}

// Mutation is a struct that can accept the pitest xml output.
type Mutation struct {
	Detected          bool   `xml:"detected,attr"`
	Status            string `xml:"status,attr"`
	NumTestsRun       int    `xml:"numberOfTestsRun,attr"`
	SourceFile        string `xml:"sourceFile"`
	MutatedClass      string `xml:"mutatedClass"`
	MutatedMethod     string `xml:"mutatedMethod"`
	MethodDescription string `xml:"methodDescription"`
	LineNumber        int    `xml:"lineNumber"`
	Mutator           string `xml:"mutator"`
	KillingTest       string `xml:"killingTest"`
	Description       string `xml:"description"`
}

func (m *Mutation) streamline() *mutations.Mutation {
	return &mutations.Mutation{
		Name:   m.Description,
		OpDesc: m.Mutator,
		Status: mutations.Status(m.Status),
		Starts: &mutations.Range{
			Line: m.LineNumber,
		},
		Ends: &mutations.Range{
			Line: m.LineNumber,
		},
	}
}

type PitXML struct {
	XMLName   xml.Name    `xml:"mutations"`
	Mutations []*Mutation `xml:"mutation"`
}

// Pitest is the Framework object for the pitest library.
type Pitest struct {
	yml  *YamlWrapper
	muts []*Mutation
	ms   mutations.Mutations
}

func NewPitest() *Pitest {
	return &Pitest{yml: &YamlWrapper{}}
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

	ms := mutations.Mutations{}
	for _, mu := range p.muts {
		_ = mu.streamline()

		fmt.Print(".")
	}
	fmt.Println()

	p.ms = ms
	return nil
}

func (p *Pitest) Mutations() mutations.Mutations {
	return p.ms
}
