// main_test.go
package main

import (
	"testing"
)

func TestFindLowestScore(t *testing.T) {
	tests := []struct {
		name          string
		maze          []string
		expectedScore int
	}{
		{
			name: "谜题示例 1",
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
			expectedScore: 7036,
		},
		{
			name: "谜题示例 2",
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
			expectedScore: 11048,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 调用 main.go 中的 findLowestScore 函数
			got := findLowestScore(tt.maze)
			if got != tt.expectedScore {
				t.Errorf("findLowestScore(%s) 测试失败: 得到 %v, 期望 %v", tt.name, got, tt.expectedScore)
			}
		})
	}
}
