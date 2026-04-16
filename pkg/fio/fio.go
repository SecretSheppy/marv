package fio

import (
	"bytes"
	"os"
	"strings"
)

func ReadLines(file string) ([]string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	lines := make([]string, 0)
	for line := range bytes.Lines(data) {
		lines = append(lines, strings.ReplaceAll(string(line), "\n", ""))
	}
	return lines, nil
}
