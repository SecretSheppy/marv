package fwlib

import (
	"github.com/SecretSheppy/marv/internal/languages"
	"github.com/SecretSheppy/marv/internal/mutations"
)

type Meta struct {
	Name     string
	Language *languages.Language
	URL      string
}

// Runnable interface describes Framework instances that have the ability to re-run the framework to generate a
// new report.
type Runnable interface {
	Run()
}

// Decompiling interface describes Framework instances that have to decompile binaries in order to extract mutants.
type Decompiling interface {
	// SetDecompiler sets the decompiler that is being used.
	SetDecompiler()
}

// FWConfig interface describes objects that are used to read custom Framework configurations from the .marv.yml file.
type FWConfig interface {
	// Init returns the default configuration for use with the init command.
	Init() interface{}
	// Load unmarshals the yml data into the struct if it exists, and returns true if any configuration was loaded.
	Load(yml []byte) (bool, error)
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
	// ReadLines returns the lines of the specified file
	ReadLines(file string) ([]string, error)
}
