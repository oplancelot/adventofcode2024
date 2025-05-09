package main

import (
	"testing"
)

func TestFindXMASPatterns(t *testing.T) {
	tests := []struct {
		name     string
		grid     [][]string
		expected int
	}{
		{
			name:     "empty grid",
			grid:     [][]string{},
			expected: 0,
		},
		{
			name:     "nil grid",
			grid:     nil,
			expected: 0,
		},
		{
			name:     "grid with empty rows",
			grid:     [][]string{{}, {}, {}},
			expected: 0,
		},
		{
			name:     "grid 2x2 (too small)",
			grid:     [][]string{{"M", "A"}, {"S", "X"}},
			expected: 0,
		},
		{
			name:     "grid 3x2 (cols < 3)",
			grid:     [][]string{{"M", "A"}, {"S", "X"}, {"M", "S"}},
			expected: 0,
		},
		{
			name:     "grid 2x3 (rows < 3)",
			grid:     [][]string{{"M", "A", "S"}, {"S", "X", "M"}},
			expected: 0,
		},
		{
			name:     "grid 3x3 no A",
			grid:     [][]string{{"M", "X", "S"}, {"X", "X", "X"}, {"M", "X", "S"}},
			expected: 0,
		},
		{
			name:     "grid 3x3 with A but no surrounding M/S pattern",
			grid:     [][]string{{"X", "X", "X"}, {"X", "A", "X"}, {"X", "X", "X"}},
			expected: 0,
		},
		{
			name: "pattern: diag1=(M,S), diag2=(M,S)", // TL=M, BR=S; TR=M, BL=S
			grid: [][]string{
				{"M", "X", "M"},
				{"X", "A", "X"},
				{"S", "X", "S"},
			},
			expected: 1,
		},
		{
			name: "pattern: diag1=(S,M), diag2=(S,M)", // TL=S, BR=M; TR=S, BL=M
			grid: [][]string{
				{"S", "X", "S"},
				{"X", "A", "X"},
				{"M", "X", "M"},
			},
			expected: 1,
		},
		{
			name: "pattern: diag1=(M,S), diag2=(S,M)", // TL=M, BR=S; TR=S, BL=M
			grid: [][]string{
				{"M", "X", "S"},
				{"X", "A", "X"},
				{"M", "X", "S"},
			},
			expected: 1,
		},
		{
			name: "pattern: diag1=(S,M), diag2=(M,S)", // TL=S, BR=M; TR=M, BL=S
			grid: [][]string{
				{"S", "X", "M"},
				{"X", "A", "X"},
				{"S", "X", "M"},
			},
			expected: 1,
		},
		{
			name: "non-pattern: specific M,S,A arrangement that is not XMAS",
			grid: [][]string{
				{"M", "X", "S"}, // TL=M, TR=S
				{"X", "A", "X"},
				{"S", "X", "M"}, // BL=S, BR=M
			},
			// tl=M, tr=S, bl=S, br=M
			// diag1: (M,M) or (S,S) -> false
			// diag2: (S,S) or (M,M) -> false
			expected: 0,
		},
		{
			name: "Only one diagonal is MAS/SAM, other is not",
			grid: [][]string{
				{"M", "X", "X"}, // TL=M, TR=X
				{"X", "A", "X"}, // A
				{"X", "X", "S"}, // BL=X, BR=S
			},
			// tl=M, tr=X, bl=X, br=S
			// diag1: (M,S) -> true
			// diag2: (X,X) -> false
			expected: 0,
		},
		{
			name: "Multiple patterns in a larger grid",
			grid: [][]string{
				{"M", "X", "S", "X", "S", "X", "M"}, // Row 0
				{"X", "A", "X", "A", "X", "A", "X"}, // Row 1
				{"M", "X", "S", "X", "M", "X", "S"}, // Row 2
				{"X", "X", "X", "X", "X", "X", "X"}, // Row 3
				{"S", "X", "M", "X", "M", "X", "S"}, // Row 4
				{"X", "A", "X", "A", "X", "A", "X"}, // Row 5
				{"M", "X", "S", "X", "S", "X", "M"}, // Row 6
			},
			// A at (1,1): TL(M),BR(S) [diag1:MAS]; TR(S),BL(M) [diag2:SAM]. Pattern.
			// A at (1,3): TL(S),BR(M) [diag1:SAM]; TR(S),BL(S) [diag2:S.S]. No.
			// A at (1,5): TL(S),BR(S) [diag1:S.S]. No.
			// A at (5,1): TL(S),BR(S) [diag1:S.S]. No.
			// A at (5,3): TL(M),BR(S) [diag1:MAS]; TR(M),BL(S) [diag2:MAS]. Pattern.
			// A at (5,5): TL(M),BR(M) [diag1:M.M]. No.
			expected: 2,
		},
		{
			name: "Overlapping potential (but distinct centers)",
			grid: [][]string{
				{"M", "X", "S", "X", "M"},
				{"X", "A", "X", "A", "X"},
				{"M", "X", "S", "X", "S"},
				{"X", "A", "X", "A", "X"},
				{"S", "X", "M", "X", "M"},
			},
			// A at (1,1): TL(M),BR(S) [MAS]; TR(S),BL(M) [SAM]. Yes.
			// A at (1,3): TL(S),BR(S) [S.S]. No.
			// A at (3,1): TL(M),BR(M) [M.M]. No.
			// A at (3,3): TL(S),BR(M) [SAM]; TR(S),BL(M) [SAM]. Yes.
			expected: 2,
		},
		{
			name: "Grid with only one possible center A, but not forming pattern",
			grid: [][]string{
				{"M", "A", "S"},
				{"A", "A", "A"},
				{"S", "A", "M"},
			},
			// Center A is (1,1)
			// TL(M), BR(M) -> diag1 M.M (invalid)
			expected: 0,
		},
		{
			name: "Valid pattern with other chars in filler spots",
			grid: [][]string{
				{"M", "Z", "S"}, // TL=M, TR=S
				{"Y", "A", "W"}, // A
				{"M", "V", "S"}, // BL=M, BR=S
			},
			// tl=M, tr=S, bl=M, br=S
			// diag1: (M,S) -> valid
			// diag2: (S,M) -> valid
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findXMASPatterns(tt.grid)
			if got != tt.expected {
				// For small grids, printing the grid is helpful. For very large ones, it might be too much.
				// These test grids are reasonably small.
				t.Errorf("findXMASPatterns() for grid %v\n  got %v, want %v", tt.grid, got, tt.expected)
			}
		})
	}
}
