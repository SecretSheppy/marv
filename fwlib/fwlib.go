package fwlib

import (
	"github.com/SecretSheppy/marv/pkg/mutations"
)

type Meta struct {
	Name      string
	Extension string
	TSLang    string // TODO: should be replaced with type of tree sitter language, haven't yet decided which tree sitter package to use
	URL       string
}

// Runnable interface describes Framework instances that have the ability to re-run the framework to generate a
// new report.
type Runnable interface {
	Run()
}

// Decompiling interfaces describes Framework instances that have to decompile binaries in order to extract mutants.
type Decompiling interface {
	// SetDecompiler sets the decompiler that is being used.
	SetDecompiler()
}

type FWConfig interface {
	Init() interface{}
	Load(yml []byte) (bool, error)
	SourceCodeDir() string
}

// Framework defines what methods an extension must have in order to interact with the marv system.
type Framework interface {
	// Meta returns the Framework's Meta information.
	Meta() *Meta
	// Yaml returns the YAML configuration struct.
	Yaml() FWConfig
	// LoadResults loads the Framework's output data.
	LoadResults() error
	// TransformResults transforms the Framework's output data into the marv format.
	TransformResults() error
	// Mutations returns the mutations in the marv format. Returns nil if TransformResults has not been called.
	Mutations() mutations.Mutations
}
