package main

import (
	"reflect"
	"testing"
)

func TestParseInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  [][]rune
	}{
		{
			name:  "empty string",
			input: "",
			want:  [][]rune{},
		},
		{
			name:  "only newlines and spaces",
			input: "\n  \n ",
			want:  [][]rune{},
		},
		{
			name:  "simple input",
			input: "A",
			want:  [][]rune{{'A'}},
		},
		{
			name:  "multiline input",
			input: "AB\nCD",
			want:  [][]rune{{'A', 'B'}, {'C', 'D'}},
		},
		{
			name:  "input with leading/trailing spaces on lines",
			input: "  AB  \n  CD  ",
			want:  [][]rune{{'A', 'B'}, {'C', 'D'}},
		},
		{
			name:  "input with leading/trailing spaces around block",
			input: "  \n  AB\nCD  \n  ",
			want:  [][]rune{{'A', 'B'}, {'C', 'D'}},
		},
		{
			name:  "input with empty lines in between",
			input: "AB\n\nCD", // strings.Split will produce an empty string for the middle line
			// parseInput trims lines, so an effectively empty line becomes an empty []rune
			// However, the current parseInput splits by "\n" then trims.
			// "AB\n\nCD" -> lines: ["AB", "", "CD"]
			// grid[0] = []rune{'A','B'}
			// grid[1] = []rune{} // from ""
			// grid[2] = []rune{'C','D'}
			// This might be an edge case to consider if rows must have consistent length
			// For now, assuming parseInput handles it as per its implementation.
			// The current problem implies a rectangular grid, so empty inner lines might be problematic
			// for `findRegionsAndCalculateTotalPrice` if not handled (e.g. grid[0] access).
			// Let's assume valid rectangular or empty grid inputs for `findRegionsAndCalculateTotalPrice`.
			// For `parseInput` itself, this is what it would produce:
			want: [][]rune{{'A', 'B'}, {}, {'C', 'D'}},
		},
		{
			name:  "single line with spaces",
			input: " A B ",
			want:  [][]rune{{'A', ' ', 'B'}}, // strings.TrimSpace(line)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseInput(tt.input); !reflect.DeepEqual(got, tt.want) {
				// For easier debugging of [][]rune differences
				t.Errorf("parseInput() got = %v, want %v", got, tt.want)
				if len(got) != len(tt.want) {
					t.Errorf("parseInput() length mismatch: got len %d, want len %d", len(got), len(tt.want))
				} else {
					for i := range got {
						if string(got[i]) != string(tt.want[i]) {
							t.Errorf("parseInput() row %d mismatch: got %s, want %s", i, string(got[i]), string(tt.want[i]))
						}
					}
				}
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		name    string
		r, c    int
		numRows int
		numCols int
		want    bool
	}{
		{"top-left corner", 0, 0, 2, 2, true},
		{"bottom-right corner", 1, 1, 2, 2, true},
		{"inside grid", 1, 0, 3, 3, true},
		{"r too small", -1, 0, 2, 2, false},
		{"c too small", 0, -1, 2, 2, false},
		{"r too large", 2, 0, 2, 2, false},
		{"c too large", 0, 2, 2, 2, false},
		{"r at boundary", 1, 1, 2, 2, true}, // (numRows-1)
		{"c at boundary", 1, 1, 2, 2, true}, // (numCols-1)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValid(tt.r, tt.c, tt.numRows, tt.numCols); got != tt.want {
				t.Errorf("isValid(%d, %d, %d, %d) = %v, want %v", tt.r, tt.c, tt.numRows, tt.numCols, got, tt.want)
			}
		})
	}
}

func TestFindRegionsAndCalculateTotalPrice(t *testing.T) {
	tests := []struct {
		name string
		grid [][]rune
		want int
	}{
		{
			name: "Simple example1",
			grid: parseInput("EEEEE\nEXXXX\nEEEEE\nEXXXX\nEEEEE"),
			want: 236, // As per the problem description
		},
		{
			name: "Simple example2",
			grid: parseInput("AAAA\nBBCD\nBBCC\nEEEC"),
			want: 80,
		},
		{
			name: "Simple example3",
			grid: parseInput("AAAAAA\nAAABBA\nAAABBA\nABBAAA\nABBAAA\nAAAAAA"),
			want: 368,
		},

		{
			name: "Simple example4",
			grid: parseInput("RRRRIICCFF\nRRRRIICCCF\nVVRRRCCFFF\nVVRCCCJFFF\nVVVVCJJCFE\nVVIVCCJJEE\nVVIIICJJEE\nMIIIIIJJEE\nMIIISIJEEE\nMMMISSJEEE"),
			want: 1206,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findRegionsAndCalculateTotalPrice(tt.grid); got != tt.want {
				t.Errorf("findRegionsAndCalculateTotalPrice() for grid %v = %v, want %v", tt.grid, got, tt.want)
			}
		})
	}
}
