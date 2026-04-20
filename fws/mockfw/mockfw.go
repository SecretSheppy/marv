package mockfw

import (
	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/internal/languages"
	"github.com/SecretSheppy/marv/internal/mutations"
)

type YamlWrapper struct{}

func (y *YamlWrapper) Init() interface{} {
	return nil
}

func (y *YamlWrapper) Load(_ []byte) (bool, error) {
	return false, nil
}

// MockFW is a framework that is only for use when testing methods that require a fwlib.Framework as a parameter. It
// allows you to set the mutations that are returned and has a Meta value, but otherwise returns nil or false in all
// cases.
type MockFW struct {
	Muts mutations.Mutations
}

func (m *MockFW) Meta() *fwlib.Meta {
	return &fwlib.Meta{
		Name:     "mock-framework",
		Language: &languages.Language{},
		URL:      "https://github.com/SecretSheppy/marv/fws/mockfw",
	}
}

func (m *MockFW) Yaml() fwlib.FWConfig {
	return &YamlWrapper{}
}

func (m *MockFW) LoadResults() error {
	return nil
}

func (m *MockFW) TransformResults() error {
	return nil
}

func (m *MockFW) Mutations() mutations.Mutations {
	return m.Muts
}

func (m *MockFW) ReadLines(_ string) ([]string, error) {
	return nil, nil
}
