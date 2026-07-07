package colour

import "testing"

func TestIsBright(t *testing.T) {
	tests := []struct {
		Name, Hex string
		Expected  bool
	}{
		{
			"full white is bright",
			"#ffffff",
			true,
		},
		{
			"full black is dark",
			"#000000",
			false,
		},
		{
			"rose pink is bright",
			"#FFD0C7",
			true,
		},
		{
			"red-pink is dark",
			"#D14747",
			false,
		},
		{
			"red-brown is dark",
			"#614C4C",
			false,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			if isBright, _ := IsBright(test.Hex); isBright != test.Expected {
				t.Errorf("expected %t, got %t", test.Expected, isBright)
			}
		})
	}
}
