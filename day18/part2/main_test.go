package main

import (
	"reflect"
	"testing"
)

func TestFindBlockingByte(t *testing.T) {
	// 谜题中提供的完整示例字节列表
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
		{1, 2},
		{5, 5},
		{2, 5},
		{6, 5},
		{1, 4},
		{0, 4},
		{6, 4},
		{1, 1},
		{6, 1},
		{1, 0},
		{0, 5},
		{1, 6},
		{2, 0},
	}

	tests := []struct {
		name          string
		width         int
		height        int
		bytePositions []Point
		wantPoint     Point
		wantFound     bool
	}{
		{
			name:          "Example from puzzle part two",
			width:         7,
			height:        7,
			bytePositions: exampleBytes,
			wantPoint:     Point{6, 1}, // 根据描述，这个字节是第一个阻塞路径的
			wantFound:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPoint, gotFound := findBlockingByte(tt.width, tt.height, tt.bytePositions)
			if !reflect.DeepEqual(gotPoint, tt.wantPoint) {
				t.Errorf("findBlockingByte() gotPoint = %v, want %v", gotPoint, tt.wantPoint)
			}
			if gotFound != tt.wantFound {
				t.Errorf("findBlockingByte() gotFound = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}
