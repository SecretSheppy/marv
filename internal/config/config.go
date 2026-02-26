package config

type Web struct {
	Port int `yaml:"port"`
}

type Paths struct {
	Sources string `yaml:"src"`
	Reviews string `yaml:"review"`
}

type Features struct {
	ToolRunner   bool `yaml:"tool-runner"`
	FileWatchers bool `yaml:"file-watchers"`
}

type Framework struct {
	Name string `yaml:"name"`
	Fomi string `yaml:"fomi"`
	Run  string `yaml:"run"`
}

type Config struct {
	Web        Web         `yaml:"web"`
	Paths      Paths       `yaml:"paths"`
	Features   Features    `yaml:"config"`
	Frameworks []Framework `yaml:"frameworks"`
}

// Init returns the default .marv.yml config for creating the default .marv.yml file.
func Init() *Config {
	return &Config{
		Web:   Web{Port: 8080},
		Paths: Paths{Sources: ".", Reviews: "./marv/reviews"},
	}
}
