package main

import (
	"strings"
	"testing"
)

func TestSolve(t *testing.T) {
	// Example from the puzzle description.
	// We use backticks for the multi-line string and then remove the leading spaces
	// that are added for indentation in the source code.
	exampleInput := `
#####
.####
.####
.####
.#.#.
.#...
.....

#####
##.##
.#.##
...##
...#.
...#.
.....

.....
#....
#....
#...#
#.#.#
#.###
#####

.....
.....
#.#..
###..
###.#
###.#
#####

.....
.....
.....
#....
#.#..
#.#.#
#####`

	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "example",
			input: strings.TrimSpace(exampleInput),
			want:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := solve(tt.input)
			if got != tt.want {
				t.Errorf("solve() = %v, want %v", got, tt.want)
			}
		})
	}
}
