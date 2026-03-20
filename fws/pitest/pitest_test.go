package pitest

import (
	"encoding/xml"
	"testing"
)

func TestPitestXmlUnmarshalling(t *testing.T) {
	rawxml := `<?xml version="1.0" encoding="UTF-8"?>
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
