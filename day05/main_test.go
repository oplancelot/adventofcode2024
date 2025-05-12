package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// Helper function to create a temporary file with given content for testing readInput
func createTempInputFile(t *testing.T, content string) string {
	t.Helper()
	// Create a temporary directory for the test file, which will be cleaned up automatically
	tempDir := t.TempDir()
	tmpFile, err := os.Create(filepath.Join(tempDir, "test_input.txt"))
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	if _, err := tmpFile.WriteString(content); err != nil {
		tmpFile.Close()
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}
	return tmpFile.Name()
}

func TestReadInput(t *testing.T) {
	tests := []struct {
		name          string
		fileContent   string // Content to write to a temp file
		nonExistent   bool   // If true, try to read a non-existent file
		expectedRules [][2]int
		expectedUpdates [][]int
		expectError   bool
	}{
		{
			name: "valid input with rules and updates",
			fileContent: `1|2
3|4

10,20,30
40,50
`,
			expectedRules:   [][2]int{{1, 2}, {3, 4}},
			expectedUpdates: [][]int{{10, 20, 30}, {40, 50}},
			expectError:     false,
		},
		{
			name:            "empty file",
			fileContent:     "",
			expectedRules:   nil, // Or [][2]int{} depending on initialization
			expectedUpdates: nil, // Or [][]int{}
			expectError:     false,
		},
		{
			name:          "file not found",
			nonExistent:   true,
			expectedRules: nil,
			expectedUpdates: nil,
			expectError:   true,
		},
		{
			name: "rules only",
			fileContent: `5|6
7|8
`,
			expectedRules:   [][2]int{{5, 6}, {7, 8}},
			expectedUpdates: nil,
			expectError:     false,
		},
		{
			name:        "updates only",
			fileContent: `1,2,3`,
			expectedRules:   nil,
			expectedUpdates: [][]int{{1, 2, 3}},
			expectError:     false,
		},
		{
			name: "invalid rule format - too many parts",
			fileContent: `1|2|3
4|5
`, // The first rule is skipped, second is processed
			expectedRules:   [][2]int{{4, 5}},
			expectedUpdates: nil,
			expectError:     false, // readInput prints and continues
		},
		{
			name: "invalid rule format - not integer",
			fileContent: `a|2
3|4
`, // First rule skipped
			expectedRules:   [][2]int{{3, 4}},
			expectedUpdates: nil,
			expectError:     false, // readInput prints and continues
		},
		{
			name: "invalid update format - not integer",
			fileContent: `10,b,30
40,50
`, // "b" is skipped in the first update line
			expectedRules:   nil,
			expectedUpdates: [][]int{{10, 30}, {40, 50}},
			expectError:     false, // readInput prints and continues
		},
		{
			name: "mixed valid and invalid lines",
			fileContent: `1|2
invalid_rule_line
10,20
another_invalid_line_for_update_or_rule
3|4
30,x,40
`,
			expectedRules:   [][2]int{{1, 2}, {3, 4}},
			expectedUpdates: [][]int{{10, 20}, {30, 40}},
			expectError:     false,
		},
		{
			name:            "input with only blank lines and spaces",
			fileContent:     "\n   \n\n  \t\n",
			expectedRules:   nil,
			expectedUpdates: nil,
			expectError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var inputFile string
			if tt.nonExistent {
				inputFile = filepath.Join(t.TempDir(), "non_existent_file.txt")
			} else {
				inputFile = createTempInputFile(t, tt.fileContent)
			}

			rules, updates, err := readInput(inputFile)

			if tt.expectError {
				if err == nil {
					t.Errorf("readInput() expected an error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("readInput() returned an unexpected error: %v", err)
				}
				if !reflect.DeepEqual(rules, tt.expectedRules) {
					t.Errorf("readInput() rules = %v, want %v", rules, tt.expectedRules)
				}
				if !reflect.DeepEqual(updates, tt.expectedUpdates) {
					t.Errorf("readInput() updates = %v, want %v", updates, tt.expectedUpdates)
				}
			}
		})
	}
}

func TestIsValidUpdate(t *testing.T) {
	tests := []struct {
		name   string
		update []int
		rules  [][2]int
		want   bool
	}{
		{
			name:   "empty update, empty rules",
			update: []int{},
			rules:  [][2]int{},
			want:   true,
		},
		{
			name:   "update with elements, empty rules",
			update: []int{1, 2, 3},
			rules:  [][2]int{},
			want:   true,
		},
		{
			name:   "simple valid: a before b, both present",
			update: []int{1, 2},
			rules:  [][2]int{{1, 2}},
			want:   true,
		},
		{
			name:   "simple invalid: b before a, both present",
			update: []int{2, 1},
			rules:  [][2]int{{1, 2}},
			want:   false,
		},
		{
			name:   "a present, b not present - valid by current logic",
			update: []int{1, 3},
			rules:  [][2]int{{1, 2}},
			want:   true, // Rule {1,2}: hasA=true, hasB=false. Condition (hasA && hasB) is false.
		},
		{
			name:   "b present, a not present - valid by current logic",
			update: []int{3, 2},
			rules:  [][2]int{{1, 2}},
			want:   true, // Rule {1,2}: hasA=false, hasB=true. Condition (hasA && hasB) is false.
		},
		{
			name:   "neither a nor b present - valid by current logic",
			update: []int{3, 4},
			rules:  [][2]int{{1, 2}},
			want:   true, // Rule {1,2}: hasA=false, hasB=false. Condition (hasA && hasB) is false.
		},
		{
			name:   "multiple rules, all valid",
			update: []int{1, 2, 3, 4, 7},
			rules:  [][2]int{{1, 2}, {3, 4}, {5, 6}}, // {5,6} not in update
			want:   true,
		},
		{
			name:   "multiple rules, one invalid (a after b, both present)",
			update: []int{2, 1, 3, 4},
			rules:  [][2]int{{1, 2}, {3, 4}}, // {1,2} is invalid
			want:   false,
		},
		{
			name:   "multiple rules, one invalid, others fine",
			update: []int{10, 20, 5, 1},
			rules:  [][2]int{{10, 20}, {1, 5}}, // {1,5} is invalid because 5 (pos 2) is not after 1 (pos 3)
			want:   false,
		},
		{
			name:   "rule with same numbers, a before b (posA >= posB is true)",
			update: []int{1, 1}, // positions map: {1:1} (last occurrence)
			rules:  [][2]int{{1, 1}},
			want:   false, // posA=1, posB=1. posA >= posB is true.
		},
		{
			name:   "update with duplicate numbers, rule applies to last occurrences due to map population",
			update: []int{1, 2, 1, 3}, // positions: {1:2, 2:1, 3:3}
			rules:  [][2]int{{1, 2}}, // Rule 1 before 2. posA=2 (for 1), posB=1 (for 2). posA >= posB.
			want:   false,
		},
		{
			name:   "update with duplicate numbers, valid case with last occurrences",
			update: []int{2, 1, 3, 1}, // positions: {2:0, 1:3, 3:2}
			rules:  [][2]int{{2, 1}}, // Rule 2 before 1. posA=0 (for 2), posB=3 (for 1). posA < posB.
			want:   true,
		},
		{
			name:   "complex case with multiple rules and varied presence",
			update: []int{10, 20, 30, 40, 50},
			rules: [][2]int{
				{10, 20}, // valid
				{30, 20}, // invalid
				{60, 70}, // valid (neither present)
				{40, 80}, // valid (b not present)
			},
			want: false, // because {30,20} is invalid
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidUpdate(tt.update, tt.rules); got != tt.want {
				t.Errorf("isValidUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}