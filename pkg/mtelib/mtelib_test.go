package mtelib

import (
	"testing"
)

func TestMTE(t *testing.T) {
	mte, err := NewMTE("")
	if err != nil {
		t.Fatal(err)
	}
	mte.Transform()
}
