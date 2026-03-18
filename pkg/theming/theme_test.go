package theming

import "testing"

type CSSTestCase struct {
	Name     string
	Value    CSS
	Expected string
}

func RunTestCases(t *testing.T, tests []CSSTestCase) {
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			if test.Value.ToCSS() != test.Expected {
				t.Errorf("expected %s\ngot %s", test.Value.ToCSS(), test.Expected)
			}
		})
	}
}

func TestTextStyle_ToCSS(t *testing.T) {
	RunTestCases(t, []CSSTestCase{
		{
			"Correctly generates minimal css consisting of color and selection color",
			TextStyle{Color: "#fff", SelectionColor: "#000"},
			"color:#fff;&::selection{color:#000;}",
		},
		{
			"Correctly generates minimal css consisting of selection color, background and border",
			TextStyle{SelectionColor: "#000", SelectionBackground: "#fff", SelectionBorder: "none"},
			"&::selection{background:#fff;color:#000;border:none;}",
		},
	})
}

func TestCode_ToCSS(t *testing.T) {
	RunTestCases(t, []CSSTestCase{
		{
			"Correctly generates all classes with at least one piece of data",
			Code{
				Variable:  TextStyle{Color: "#fff", SelectionBackground: "#000"},
				Function:  TextStyle{Background: "#fff"},
				Keyword:   TextStyle{FontStyle: "italic"},
				String:    TextStyle{Color: "#aaafff"},
				Comment:   TextStyle{Background: "black", TextDecoration: "underline"},
				Number:    TextStyle{Color: "#fff"},
				Operator:  TextStyle{Color: "blue"},
				Type:      TextStyle{Color: "#fff"},
				Class:     TextStyle{Color: "black"},
				Interface: TextStyle{Color: "#fff"},
				Constant:  TextStyle{Color: "#fff"},
				Property:  TextStyle{Color: "#fff"},
			},
			".variable{color:#fff;&::selection{background:#000;}}.function{background:#fff;}" +
				".keyword{font-style:italic;}.string{color:#aaafff;}" +
				".comment{background:black;text-decoration:underline;}.number{color:#fff;}.operator{color:blue;}" +
				".type{color:#fff;}.class{color:black;}.interface{color:#fff;}.constant{color:#fff;}.property{color:#fff;}",
		},
	})
}
