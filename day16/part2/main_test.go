// main_test.go
package main

import (
	"testing"
)

// TestCountTilesOnBestPath (Part 2 测试)
func TestCountTilesOnBestPath(t *testing.T) {
	tests := []struct {
		name          string
		maze          []string
		expectedCount int
	}{
		{
			name: "Part 2 Example 1",
			maze: []string{
				"###############",
				"#.......#....E#",
				"#.#.###.#.###.#",
				"#.....#.#...#.#",
				"#.###.#####.#.#",
				"#.#.#.......#.#",
				"#.#.#####.###.#",
				"#...........#.#",
				"###.#.#####.#.#",
				"#...#.....#.#.#",
				"#.#.#.###.#.#.#",
				"#.....#...#.#.#",
				"#.###.#.#.#.#.#",
				"#S..#.....#...#",
				"###############",
			},
			expectedCount: 45,
		},
		{
			name: "Part 2 Example 2",
			maze: []string{
				"#################",
				"#...#...#...#..E#",
				"#.#.#.#.#.#.#.#.#",
				"#.#.#.#...#...#.#",
				"#.#.#.#.###.#.#.#",
				"#...#.#.#.....#.#",
				"#.#.#.#.#.#####.#",
				"#.#...#.#.#.....#",
				"#.#.#####.#.###.#",
				"#.#.#.......#...#",
				"#.#.###.#####.###",
				"#.#.#...#.....#.#",
				"#.#.#.#####.###.#",
				"#.#.#.........#.#",
				"#.#.#.#########.#",
				"#S#.............#",
				"#################",
			},
			expectedCount: 64,
		},
		{
			name: "Simple S.E path",
			maze: []string{
				"#####",
				"#S.E#",
				"#####",
			},
			expectedCount: 3, // S, ., E
		},
		{
			name: "S to E direct",
			maze: []string{
				"###",
				"#SE#",
				"###",
			},
			expectedCount: 2,
		},
		{
			name: "No path",
			maze: []string{
				"S#E",
			},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countTilesOnBestPath(tt.maze)
			if got != tt.expectedCount {
				t.Errorf("countTilesOnBestPath(%s) = %v, want %v", tt.name, got, tt.expectedCount)
			}
		})
	}
}
