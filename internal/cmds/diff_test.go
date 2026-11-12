package cmds

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/SecretSheppy/marv/internal/mutations"
	mapset "github.com/deckarep/golang-set/v3"
)

func TestDiffResults(t *testing.T) {
	oldResults, err := os.ReadFile("../../test/cmds/diff_export/old.json")
	if err != nil {
		t.Fatal(err)
	}

	newResults, err := os.ReadFile("../../test/cmds/diff_export/new.json")
	if err != nil {
		t.Fatal(err)
	}

	expected, err := os.ReadFile("../../test/cmds/diff_export/expected.json")
	if err != nil {
		t.Fatal(err)
	}

	var or, nr mutations.Mutations
	if err = json.Unmarshal(oldResults, &or); err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(newResults, &nr); err != nil {
		t.Fatal(err)
	}

	exp := &diff{}
	if err = json.Unmarshal(expected, exp); err != nil {
		t.Fatal(err)
	}

	got := diffResults(or, nr)

	if !reflect.DeepEqual(exp.Removed, got.Removed) {
		t.Fatalf("diff removed: expected %v, got %v", exp.Removed, got.Removed)
	}

	if !reflect.DeepEqual(exp.Added, got.Added) {
		t.Fatalf("diff added: expected %v, got %v", exp.Added, got.Added)
	}

	expChanges := mapset.NewSet[change](exp.Changed...)
	gotChanges := mapset.NewSet[change](got.Changed...)

	if !expChanges.Equal(gotChanges) {
		t.Fatalf("diff changes - expected:\n%v\ngot:\n%v\n", expChanges, gotChanges)
	}
}
