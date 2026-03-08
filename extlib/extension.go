package extlib

import "github.com/SecretSheppy/marv/pkg/mutations"

type Meta struct {
	Name   string
	Lang   string
	TSLang string // TODO: should be replaced with type of tree sitter language, haven't yet decided which tree sitter package to use
	URL    string
	RunStr string
}

// Extension defines what methods an extension must have in order to interact with the marv system.
type Extension interface {
	Meta() *Meta
	Init(path string) error
	Mutations() (mutations.Mutations, error)
}

// TODO: toggle bits of marv specification/functionality
