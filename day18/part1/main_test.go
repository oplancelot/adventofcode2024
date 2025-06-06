package main

import (
	"testing"
)

func TestFindShortestPath(t *testing.T) {
	// 谜题中提供的示例数据
	exampleBytes := []Point{
		{5, 4},
		{4, 2},
		{4, 5},
		{3, 0},
		{2, 1},
		{6, 3},
		{2, 4},
		{1, 5},
		{0, 6},
		{3, 3},
		{2, 6},
		{5, 1},
	}

	// 定义测试用例表
	tests := []struct {
		name          string
		width         int
		height        int
		byteCount     int
		bytePositions []Point
		want          int
	}{
		{
			name:          "Example from puzzle",
			width:         7,  // 网格大小为 0-6，即宽度 7
			height:        7,  // 网格大小为 0-6，即高度 7
			byteCount:     12, // 模拟前 12 个字节
			bytePositions: exampleBytes,
			want:          22, // 期望的最短路径长度
		},
	}

	// 遍历并执行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findShortestPath(tt.width, tt.height, tt.byteCount, tt.bytePositions); got != tt.want {
				t.Errorf("findShortestPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
