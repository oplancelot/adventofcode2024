package main

import (
	"testing" // 引入 Go 的测试包
)

// TestCanFormDesign 是一个表驱动测试，用于验证 CanFormDesign 函数的正确性。
func TestCanFormDesign(t *testing.T) {
	// 定义所有可用的毛巾模式
	availablePatterns := []string{"r", "wr", "b", "g", "bwu", "rb", "gb", "br"}

	// 定义测试用例表
	// 每个条目都包含一个名字(name)、要测试的设计(design)和期望的结果(expected)。
	testCases := []struct {
		name     string
		design   string
		expected bool
	}{
		// --- 这里是根据谜题描述提供的所有例子 ---
		{
			name:     "brwrr - 成功",
			design:   "brwrr",
			expected: true, // 可由 "br" + "wr" + "r" 组成
		},
		{
			name:     "bggr - 成功",
			design:   "bggr",
			expected: true, // 可由 "b" + "g" + "g" + "r" 组成
		},
		{
			name:     "gbbr - 成功",
			design:   "gbbr",
			expected: true, // 可由 "gb" + "br" 组成
		},
		{
			name:     "rrbgbr - 成功",
			design:   "rrbgbr",
			expected: true, // 可由 "r" + "rb" + "g" + "br" 组成
		},
		{
			name:     "ubwu - 失败",
			design:   "ubwu",
			expected: false, // 包含 "u"，但没有以 "u" 开头的毛巾
		},
		{
			name:     "bwurrg - 成功",
			design:   "bwurrg",
			expected: true, // 可由 "bwu" + "r" + "r" + "g" 组成
		},
		{
			name:     "brgr - 成功",
			design:   "brgr",
			expected: true, // 可由 "br" + "g" + "r" 组成
		},
		{
			name:     "bbrgwb - 失败",
			design:   "bbrgwb",
			expected: false, // 最后的 "wb" 无法由任何毛巾模式组成
		},
		// --- 你可以在这里添加更多的测试用例 ---
		{
			name:     "空设计字符串",
			design:   "",
			expected: true, // 空字符串应该返回 true
		},
		{
			name:     "单个完全匹配的毛巾",
			design:   "bwu",
			expected: true,
		},
		{
			name:     "无法匹配的单字符",
			design:   "x",
			expected: false,
		},
	}

	// 遍历所有测试用例
	for _, tc := range testCases {
		// t.Run 会为每个测试用例创建一个独立的子测试
		// 这使得在大量测试失败时，结果报告更清晰
		t.Run(tc.name, func(t *testing.T) {
			// 调用我们正在测试的函数
			actual := CanFormDesign(tc.design, availablePatterns)

			// 检查实际结果是否与期望结果相符
			if actual != tc.expected {
				// 如果不符，t.Errorf 会报告一个错误，但不会停止整个测试过程
				t.Errorf("对于设计 '%s', 期望得到 %v, 但实际得到 %v", tc.design, tc.expected, actual)
			}
		})
	}
}
