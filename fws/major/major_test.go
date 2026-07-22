package major

import (
	"reflect"
	"testing"
)

func TestMarhsallingMutants(t *testing.T) {
	tests := []struct {
		Name      string
		MutantStr string
		ExpMutant *Mutant
	}{
		{
			Name:      "Parse mutant with regular types",
			MutantStr: "1:ROR:<=(int,int):<(int,int):triangle.Triangle@classify(int,int,int):11:191:a <= 0 |==> a < 0\n",
			ExpMutant: &Mutant{1, "ROR", "<=", "<", "triangle.Triangle", 11, 191, "a <= 0", "a < 0"},
		},
		{
			Name:      "Parse mutant with <RETURN> and <NO-OP> types",
			MutantStr: "18:STD:<RETURN>:<NO-OP>:triangle.Triangle@classify(int,int,int):12:222:return Type.INVALID; |==> <NO-OP>\n",
			ExpMutant: &Mutant{18, "STD", "<RETURN>", "<NO-OP>", "triangle.Triangle", 12, 222, "return Type.INVALID;", "<NO-OP>"},
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			m, err := marshalMutant(test.MutantStr)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(m, test.ExpMutant) {
				t.Errorf("Expected:\n\n%v\n\nGot:\n\n%v\n\n", test.ExpMutant, m)
			}
		})
	}
}

func TestMarshallingDetail(t *testing.T) {
	detail := "1,FAIL"
	exp := &Detail{1, "FAIL"}
	got, err := marshalDetail(detail)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(exp, got) {
		t.Errorf("Expected:\n\n%v\n\nGot:\n\n%v\n\n", exp, got)
	}
}

func TestMutantClassPathToFilePathConversion(t *testing.T) {
	m := &Mutant{ClassPath: "test.class.path.Class"}
	exp := "test/class/path/Class.java"
	got := m.file()
	if got != exp {
		t.Errorf("Expected:\n\n%v\n\nGot:\n\n%v\n\n", exp, got)
	}
}
