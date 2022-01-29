package filter

import (
	"testing"
)

func TestPass(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		include  bool
		input    string
		want     bool
	}{
		{
			name:     "empty list, empty input - include",
			patterns: []string{},
			include:  true,
			input:    "",
			want:     false,
		},
		{
			name:     "empty list, empty input - exclude",
			patterns: []string{},
			include:  false,
			input:    "",
			want:     true,
		},

		{
			name:     "literals list, found - exclude",
			patterns: []string{"hello", "world"},
			include:  false,
			input:    "hello",
			want:     false,
		},
		{
			name:     "literals list, partial found - exclude",
			patterns: []string{"hello", "world"},
			include:  false,
			input:    "hell",
			want:     true,
		},
		{
			name:     "regex list, found - exclude",
			patterns: []string{"hell.*", "world"},
			include:  false,
			input:    "hello",
			want:     false,
		},
		{
			name:     "regex list, partial found - exclude",
			patterns: []string{"hell.*", "world"},
			include:  false,
			input:    "hell",
			want:     false,
		},
		{
			name:     "literals list, found - include",
			patterns: []string{"hello", "world"},
			include:  true,
			input:    "hello",
			want:     true,
		},
		{
			name:     "literals list, partial found - include",
			patterns: []string{"hello", "world"},
			include:  true,
			input:    "hell",
			want:     false,
		},
		{
			name:     "regex list, found - include",
			patterns: []string{"hell.*", "world"},
			include:  true,
			input:    "hello",
			want:     true,
		},
		{
			name:     "regex list, partial found - include",
			patterns: []string{"hell.*", "world"},
			include:  true,
			input:    "hell",
			want:     true,
		},
		{
			name:     "literals list, empty input - include",
			patterns: []string{},
			include:  false,
			input:    "",
			want:     true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			filter, err := New(tt.patterns, tt.include)
			if err != nil {
				t.Fatalf("failed to create filter: %v", err)
			}

			got := filter.Pass(tt.input)
			if got != tt.want {
				t.Fatalf("Pass(%q)=%v want: %v", tt.input, got, tt.want)
			}

		})
	}
}
