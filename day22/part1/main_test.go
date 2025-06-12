package main

import (
	"strings"
	"testing"
)

// TestNextSecretSequence verifies the example sequence starting with 123.
func TestNextSecretSequence(t *testing.T) {
	initialSecret := 123
	expectedSequence := []int{
		15887950, 16495136, 527345, 704524, 1553684,
		12683156, 11100544, 12249484, 7753432, 5908254,
	}

	currentSecret := initialSecret
	for i, expected := range expectedSequence {
		currentSecret = nextSecret(currentSecret)
		if currentSecret != expected {
			t.Errorf("Sequence step %d: got %d, want %d", i+1, currentSecret, expected)
		}
	}
}

// TestSolve uses a table-driven test to check the main solving logic with the puzzle's example.
func TestSolve(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name: "Puzzle Example",
			input: `
1
10
100
2024
`,
			want: 37327623,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Trim leading/trailing whitespace for cleaner test cases.
			input := strings.TrimSpace(tt.input)
			if got := solve(input); got != tt.want {
				t.Errorf("solve() = %v, want %v", got, tt.want)
			}
		})
	}
}
