package config

import (
	"os"
	"path"
)

const (
	DefaultPort  = 8080
	DefaultPath  = ".marv"
	DefaultMerge = false
)

type Output struct {
	Path  string
	Merge bool
}

type Marv struct {
	Port   int    `yaml:"port"`
	Output Output `yaml:"output"`
	Debug  bool   `yaml:"debug,omitempty"`
	Theme  string `yaml:"-"`
}

type Config struct {
	Marv Marv `yaml:"marv"`
}

func (c *Config) LoadPersistentData() {
	if c.Marv.Theme == "" {
		c.Marv.Theme = GetPersistentData(PersistentTheme)
	}
}

// Init returns the default .marv.yml config for creating the default .marv.yml file.
func Init() *Config {
	return &Config{
		Marv{
			Port: DefaultPort,
			Output: Output{
				Path:  DefaultPath,
				Merge: DefaultMerge,
			},
		},
	}
}

const (
	DefaultTheme    = "darcula"
	PersistentTheme = "marv.theme"
)

func GetPersistentData(name string) string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	raw, err := os.ReadFile(path.Join(dir, "marv", name))
	if err != nil {
		return "" // NOTE: this likely means that the user has yet to set the relevant persistent data
	}
	return string(raw)
}

func SetPersistentData(name, value string) error {
	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	if err = os.Mkdir(path.Join(dir, "marv"), 0777); err != nil && !os.IsExist(err) {
		return err
	}
	return os.WriteFile(path.Join(dir, "marv", name), []byte(value), 0777)
}
