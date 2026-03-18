package fws

import (
	"github.com/SecretSheppy/marv/pkg/mutations"
)

type Meta struct {
	Name   string
	Lang   string
	TSLang string // TODO: should be replaced with type of tree sitter language, haven't yet decided which tree sitter package to use
	URL    string
	RunStr string
}

// Runnable interface describes Framework instances that have the ability to re-run the framework to generate a
// new report.
type Runnable interface {
	Run()
}

// Framework defines what methods an extension must have in order to interact with the marv system.
type Framework interface {
	// Meta returns the Meta information for the Framework
	Meta() *Meta
	// LoadYamlCfg unmarshals the frameworks YAML config if it exists, and returns whether the YAML config existed for
	// that Framework. This should be used to initialize all the Framework instances and to filter out all Framework
	// instances that are not in use in the current project.
	LoadYamlCfg(yml []byte) (bool, error)
	// Init initializes the framework plugin by parsing all the mutations and formatting them into the marv internal
	// format.
	Init() error
	// Mutations returns the mutations in the marv internal format.
	Mutations() (mutations.Mutations, error)
}

func Frameworks() []Framework {
	return []Framework{
		&MutestRS{},
	}
}
