package config

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
}

type Config struct {
	Marv Marv `yaml:"marv"`
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
