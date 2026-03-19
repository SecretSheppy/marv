package config

type Features struct {
	ToolRunner   bool `yaml:"tool-runner"`
	FileWatchers bool `yaml:"file-watchers"`
}

type Marv struct {
	Port      int      `yaml:"port"`
	ReviewDir string   `yaml:"review-dir"`
	Features  Features `yaml:"features"`
}

type Config struct {
	Marv Marv `yaml:"marv"`
}

// Init returns the default .marv.yml config for creating the default .marv.yml file.
func Init() *Config {
	return &Config{
		Marv{
			Port:      8080,
			ReviewDir: ".marv/reviews",
		},
	}
}
