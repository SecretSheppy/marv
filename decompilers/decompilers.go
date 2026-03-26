package decompilers

import (
	"strings"

	"github.com/SecretSheppy/marv/decompilers/java/garlic"
	"github.com/SecretSheppy/marv/decompilers/java/vineflower"
	"github.com/SecretSheppy/marv/decompilers/java/vineflower_server"
)

// Decompiler describes an object that can be used to decompile a binary.
type Decompiler interface {
	// ExePath returns the executable path for the decompiler object.
	ExePath() string
	// Setup performs any necessary setup for the decompiler.
	Setup() error
	// Teardown performs any necessary teardown for the decompiler.
	Teardown() error
	// Decompile takes a string path to a file and returns the decompiled bytes of that file or an error.
	Decompile(path string) ([]byte, error)
}

func JavaDecompiler(name string) Decompiler {
	switch strings.ToLower(name) {
	case "garlic":
		return &garlic.Garlic{}
	case "vineflower":
		return &vineflower.Vineflower{}
	default:
		return &vineflower_server.VFServer{}
	}
}
