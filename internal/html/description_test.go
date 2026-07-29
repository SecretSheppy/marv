package html

import "testing"

func TestDescriptionFormatter(t *testing.T) {
	tests := []struct {
		Name, Description, Expected string
	}{
		{
			Name:        "test no backticks yields no codeblocks",
			Description: "replaced return statement with null",
			Expected:    "replaced return statement with null",
		},
		{
			Name:        "test single backticks yields a codeblock",
			Description: "replaced `45` with null",
			Expected:    `replaced <span class="inline-code">45</span> with null`,
		},
		{
			Name:        "test multiple sets of single backticks yields multiple codeblocks",
			Description: "replaced `45` with `null`",
			Expected:    `replaced <span class="inline-code">45</span> with <span class="inline-code">null</span>`,
		},
		{
			Name:        "test triple backticks yields a codeblock",
			Description: "replaced ```45``` with null",
			Expected:    `replaced <span class="inline-code">45</span> with null`,
		},
		{
			Name:        "test multiple sets of triple backticks yields multiple codeblocks",
			Description: "replaced ```45``` with ```null```",
			Expected:    `replaced <span class="inline-code">45</span> with <span class="inline-code">null</span>`,
		},
		{
			Name:        "test mixed codeblock generation with single backticks and triple backticks",
			Description: "replaced ```45``` with `null`",
			Expected:    `replaced <span class="inline-code">45</span> with <span class="inline-code">null</span>`,
		},
		{
			Name:        "test single backticks inside tripple backticks",
			Description: "replaced ```def `x` with``` with null",
			Expected:    "replaced <span class=\"inline-code\">def `x` with</span> with null",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			formatted := formatDescription(test.Description)
			if formatted != test.Expected {
				t.Errorf("got %q, want %q", formatted, test.Expected)
			}
		})
	}
}
