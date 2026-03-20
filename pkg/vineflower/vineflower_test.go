package vineflower

import (
	"os"
	"testing"
)

func TestHelp(t *testing.T) {
	if os.Getenv("SKIP_WRAPPER_TESTS") == "true" {
		t.Skip("skipping due to SKIP_WRAPPER_TESTS environment variable being set")
	}
	if _, err := Help(); err != nil {
		t.Fatal(err)
	}
}
