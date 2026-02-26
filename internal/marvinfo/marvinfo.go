package marvinfo

import (
	"os"

	"gopkg.in/yaml.v3"
)

// MarvInfo is the struct representation of .marvinfo.yml which just stores metadata about marvs current state.
type MarvInfo struct {
	Version string `yaml:"version"`
}

func Default() *MarvInfo {
	return &MarvInfo{
		Version: "not detected, check for .marvinfo.yml",
	}
}

func Get() *MarvInfo {
	info := &MarvInfo{}
	file, err := os.ReadFile(".marvinfo.yml")
	if err != nil {
		return Default()
	}
	if err = yaml.Unmarshal(file, &info); err != nil {
		return Default()
	}
	return info
}
