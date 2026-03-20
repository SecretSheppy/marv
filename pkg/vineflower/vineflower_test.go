package vineflower

import (
	"testing"
)

func TestHelp(t *testing.T) {
	if _, err := Help(); err != nil {
		t.Fatal(err)
	}
}
