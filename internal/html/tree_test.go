package html

import (
	"bytes"
	"testing"

	"github.com/SecretSheppy/marv/fws/mockfw"
	"github.com/SecretSheppy/marv/internal/languages"
)

var paths = []string{
	"test/files/location/test1.lang",
	"test/files/location/test2.lang",
	"test/files/location/subtest/subtest1.lang",
	"test/files/newlocation/newlocationtest1.lang",
	"test/files2/location/test1.lang",
}

func TestTreeStructure(t *testing.T) {
	root := PathNode{}
	for _, path := range paths {
		root.AddFile(path)
	}
	location := root.
		ChildNode("test").
		ChildNode("files").
		ChildNode("location")
	if location == nil {
		t.Fatal("test/files/location directory not found")
	}
	if len(location.children) != 3 {
		t.Errorf("expected 3 children in test/files/location but got %d", len(location.children))
	}
	test1File := location.ChildNode("test1.lang")
	if test1File == nil {
		t.Fatal("test/files/location/test1.lang file not found")
	}
	if test1File.Type != File {
		t.Errorf("expected test/files/location/test1.lang type to be %d but got %d", File, test1File.Type)
	}
	files2files := root.
		ChildNode("test").
		ChildNode("files2").
		ChildNode("location").
		Children()
	if len(files2files) != 1 {
		t.Errorf("expected test/files2/location to have 1 child but got %d", len(files2files))
	}
	if files2files[0].Name != "test1.lang" {
		t.Errorf("expected test/files2/location/test1.lang but got test/files2/location/%s", files2files[0].Name)
	}
}

func TestTreeSorting(t *testing.T) {
	root := PathNode{}
	for _, path := range paths {
		root.AddFile(path)
	}
	root.SortChildren()
	location := root.
		ChildNode("test").
		ChildNode("files").
		ChildNode("location")
	if location == nil {
		t.Fatal("test/files/location directory not found")
	}
	children := location.Children()
	if children[0].Name != "subtest" {
		t.Errorf("expected first child to be subtest but got %s", children[0].Name)
	}
	if children[1].Name != "test1.lang" {
		t.Errorf("expected second child to be test1.lang but got %s", children[1].Name)
	}
	if children[2].Name != "test2.lang" {
		t.Errorf("expected third child to be test2.lang but got %s", children[2].Name)
	}
}

func TestTreeRendering(t *testing.T) {
	root := PathNode{}
	for _, path := range paths {
		root.AddFile(path)
	}
	root.SortChildren()
	var buff bytes.Buffer
	root.Render(&buff, &mockfw.MockFW{}, &languages.Language{})
	// TODO: check this somehow
	//fmt.Println(buff.String())
}
