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
