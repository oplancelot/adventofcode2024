package main

import (
	"testing"
)

// TestCountWays 是一个表驱动测试，用于验证 CountWays 函数的正确性。
func TestCountWays(t *testing.T) {
	// 定义谜题示例中的可用毛巾模式
	availablePatterns := []string{"r", "wr", "b", "g", "bwu", "rb", "gb", "br"}

	// 定义测试用例表
	testCases := []struct {
		name     string
		design   string
		expected int // 期望值现在是整数
	}{
		// --- 这里是根据 Part Two 描述提供的所有例子 ---
		{name: "brwrr", design: "brwrr", expected: 2},
		{name: "bggr", design: "bggr", expected: 1},
		{name: "gbbr", design: "gbbr", expected: 4},
		{name: "rrbgbr", design: "rrbgbr", expected: 6},
		{name: "ubwu - 仍然不可能", design: "ubwu", expected: 0},
		{name: "bwurrg", design: "bwurrg", expected: 1},
		{name: "brgr", design: "brgr", expected: 2},
		{name: "bbrgwb - 仍然不可能", design: "bbrgwb", expected: 0},
		// --- 附加测试用例 ---
		{name: "空设计", design: "", expected: 1},
		{name: "不存在的模式", design: "xyz", expected: 0},
	}

	// 遍历所有测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 调用我们正在测试的函数
			actual := CountWays(tc.design, availablePatterns)

			// 检查实际结果是否与期望结果相符
			if actual != tc.expected {
				t.Errorf("对于设计 '%s', 期望得到 %d 种方法, 但实际得到 %d", tc.design, tc.expected, actual)
			}
		})
	}
}
