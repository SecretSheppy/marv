package pitest

import (
	"encoding/xml"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/rs/zerolog"
)

// NOTE: taken from results of running pitest on guava
const rawxml = `<?xml version="1.0" encoding="UTF-8"?>
<mutations partial="true">
    <mutation detected='false' status='SURVIVED' numberOfTestsRun='48'>
        <sourceFile>MapInterfaceTest.java</sourceFile>
        <mutatedClass>com.google.common.collect.testing.MapInterfaceTest</mutatedClass>
        <mutatedMethod>assertEntrySetNotContainsString</mutatedMethod>
        <methodDescription>(Ljava/util/Set;)V</methodDescription>
        <lineNumber>266</lineNumber>
        <mutator>org.pitest.mutationtest.engine.gregor.mutators.VoidMethodCallMutator</mutator>
        <indexes>
            <index>6</index>
        </indexes>
        <blocks>
            <block>1</block>
        </blocks>
        <killingTest/>
        <description>removed call to com/google/common/collect/testing/MapInterfaceTest::assertFalse</description>
    </mutation>
</mutations>`

func TestPitestXmlUnmarshalling(t *testing.T) {
	pitxml := &PitXML{}
	if err := xml.Unmarshal([]byte(rawxml), pitxml); err != nil {
		t.Fatal(err)
	}
	if pitxml.Mutations == nil {
		t.Fatal("pitxml parsed no mutations")
	}
}

func TestMutationPaths(t *testing.T) {
	tests := []struct {
		Name                     string
		Mutation                 *Mutation
		ExpectedSourceCodePath   string
		ExpectedSourceClassPath  string
		ExpectedMutatedClassPath string
	}{
		{
			"simple mutation",
			&Mutation{
				SourceFile:    "TestClass.java",
				MutatedClass:  "com.example.testing.TestClass",
				MutationIndex: 0,
			},
			"com/example/testing/TestClass.java",
			"com/example/testing/TestClass.class",
			"com/example/testing/TestClass/mutants/0/com.example.testing.TestClass.class",
		},
		{
			"more complex mutation",
			&Mutation{
				SourceFile:    "TestClassGenerator.java",
				MutatedClass:  "com.example.testing.TestClass$FirstGenerationOfClass",
				MutationIndex: 44,
			},
			"com/example/testing/TestClassGenerator.java",
			"com/example/testing/TestClass$FirstGenerationOfClass.class",
			"com/example/testing/TestClass$FirstGenerationOfClass/mutants/44/com.example.testing.TestClass$FirstGenerationOfClass.class",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			if test.Mutation.SourceCodePath() != test.ExpectedSourceCodePath {
				t.Errorf("expected %s but got %s", test.Mutation.SourceCodePath(), test.ExpectedSourceCodePath)
			}
			if test.Mutation.SourceClassPath() != test.ExpectedSourceClassPath {
				t.Errorf("expected %s but got %s", test.Mutation.SourceClassPath(), test.ExpectedSourceClassPath)
			}
			if test.Mutation.MutatedClassPath() != test.ExpectedMutatedClassPath {
				t.Errorf("expected %s but got %s", test.Mutation.MutatedClassPath(), test.ExpectedMutatedClassPath)
			}
		})
	}
}

func TestYamlWrapper_Load(t *testing.T) {
	tests := []struct {
		Name     string
		Wrapper  *YamlWrapper
		Yml      []byte
		Expected bool
	}{
		{
			"YAML that provides all required fields should load successfully",
			&YamlWrapper{},
			[]byte("pitest:\n    xml-path: a\n    src-code-path: b\n    src-class-path: c\n    mut-class-path: d"),
			true,
		},
		{
			"YAML that provides no required fields should not load successfully",
			&YamlWrapper{},
			[]byte("pitest:\n    xml-path:\n    src-code-path:\n    src-class-path:\n    mut-class-path:"),
			false,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			load, err := test.Wrapper.Load(test.Yml)
			if err != nil {
				t.Error(err)
			}
			if load != test.Expected {
				t.Errorf("expected %t but got %t", test.Expected, load)
			}
		})
	}
}

func TestPitest_LoadResults(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)

	dir := t.TempDir()
	xmlPath := path.Join(dir, "mutations.xmlPath")
	yml := fmt.Sprintf("pitest:\n    xml-path: %s\n    src-code-path: b\n    src-class-path: c\n    mut-class-path: d", xmlPath)

	t.Run("tries to load file that does not exist", func(t *testing.T) {
		pt := NewPitest()
		if _, err := pt.Yaml().Load([]byte(yml)); err != nil {
			t.Fatal(err)
		}
		if err := pt.LoadResults(); err == nil {
			t.Fatal("managed to read file that does not exist")
		}
	})

	if err := os.WriteFile(xmlPath, []byte(rawxml), 0644); err != nil {
		t.Fatal(err)
	}

	t.Run("successfully reads xml file and parses it", func(t *testing.T) {
		pt := NewPitest()
		if _, err := pt.Yaml().Load([]byte(yml)); err != nil {
			t.Fatal(err)
		}
		if err := pt.LoadResults(); err != nil {
			t.Error(err)
		}
	})
}
