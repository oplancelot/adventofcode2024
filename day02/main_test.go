package main

import (
	"testing"
)

func TestIsSafe(t *testing.T) {
	tests := []struct {
		input []int
		want  bool
	}{
		{[]int{1, 2, 3}, true},
		{[]int{3, 2, 1}, true},
		{[]int{1, 2, 6}, false},      // diff > 3
		{[]int{1, 1, 2}, false},      // diff == 0
		{[]int{1, 2, 1}, false},      // 折返
		{[]int{5}, false},           // 长度不足
		{[]int{1, 3, 5}, true},      // 差值都是 2
		{[]int{5, 3, 1}, true},      // 递减有效
		{[]int{5, 2, -1}, true},    // diff 超过 3
	}
	for _, tt := range tests {
		got := isSafe(tt.input)
		if got != tt.want {
			t.Errorf("isSafe(%v) = %v; want %v", tt.input, got, tt.want)
		}
	}
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		line string
		want LineEvaluation
	}{
		{"1 2 3", EvalSafe},
		{"3 2 1", EvalSafe},
		{"1 2 6", EvalSafe},
		{"1 2 6 3", EvalSafe},     // 去掉 6 或 3 后可以变安全
		{"5", EvalUnsafe},
		{"a b", EvalUnsafe},       // 非数字
		{"", EvalUnsafe},          // 空行
	}

	for _, tt := range tests {
		got, _ := parseLine(tt.line)
		if got != tt.want {
			t.Errorf("parseLine(%q) = %v; want %v", tt.line, got, tt.want)
		}
	}
}
