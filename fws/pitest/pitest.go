package pitest

import (
	"encoding/xml"
	"os"

	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/pkg/mutations"
	"gopkg.in/yaml.v3"
)

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

func (m *Mutation) sourceFilePath() string {
	return ""
}

func (m *Mutation) mutatedFilePath() string {
	return ""
}

type PitXML struct {
	XMLName   xml.Name    `xml:"mutations"`
	Mutations []*Mutation `xml:"mutation"`
}

// pitestYamlWrapper used to load the pitest configuration from the .marv.yml file.
type pitestYamlWrapper struct {
	Cfg *pitestYamlCfg `yaml:"pitest"`
}

type pitestYamlCfg struct {
	Run            string `yaml:"run"`
	XmlPath        string `yaml:"xml-path"`
	SrcClassesPath string `yaml:"src-classes-path"`
	MutClassesPath string `yaml:"mut-classes-path"`
}

func (p pitestYamlCfg) IsPopulated() bool {
	return p.XmlPath != "" || p.SrcClassesPath != "" || p.MutClassesPath != ""
}

// Pitest is the Framework object for the pitest library.
type Pitest struct {
	cfg  *pitestYamlCfg
	muts []*Mutation
}

func (p *Pitest) Meta() *fwlib.Meta {
	return &fwlib.Meta{
		Name:      "Pitest",
		Extension: "java",
		URL:       "https://pitest.org/",
	}
}

func (p *Pitest) LoadYamlCfg(yml []byte) (bool, error) {
	wrapper := &pitestYamlWrapper{}
	if err := yaml.Unmarshal(yml, wrapper); err != nil {
		return false, err
	}
	p.cfg = wrapper.Cfg
	return p.cfg.IsPopulated(), nil
}

func (p *Pitest) Init() error {
	rawxml, err := os.ReadFile(p.cfg.XmlPath)
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

func (p *Pitest) Mutations() (mutations.Mutations, error) {
	ms := mutations.Mutations{}
	for _, mu := range p.muts {
		_ = mu.streamline()

	}
	return ms, nil
}
