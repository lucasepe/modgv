package text

import "testing"

func TestFind(t *testing.T) {
	tests := []struct {
		set  []string
		el   string
		want bool
	}{
		{
			[]string{"filled", "rounded", "striped"},
			"filled",
			true,
		},
		{
			[]string{"filled", "rounded", "striped"},
			"slashed",
			false,
		},
		{
			[]string{"mela", "banana", "caffè"},
			"caffè",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.el, func(t *testing.T) {
			if _, got := Find(tt.set, tt.el); got != tt.want {
				t.Errorf("got [%v] want [%v]", got, tt.want)
			}
		})
	}
}
