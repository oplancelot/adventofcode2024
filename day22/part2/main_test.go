package main

import (
	"strings"
	"testing"
)

func TestSolvePart2(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name: "Part 2 Puzzle Example",
			input: `
1
2
3
2024
`,
			want: 23,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := strings.TrimSpace(tt.input)
			if got := solve(input); got != tt.want {
				t.Errorf("solve() = %v, want %v", got, tt.want)
			}
		})
	}
}
