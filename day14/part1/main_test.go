package main

import (
	"testing"
)

func TestCalculateSafetyFactor(t *testing.T) {
	// 题目中提供的示例输入数据
	exampleInput := `p=0,4 v=3,-3
p=6,3 v=-1,-3
p=10,3 v=-1,2
p=2,0 v=2,-1
p=0,0 v=1,3
p=3,0 v=-2,-2
p=7,6 v=-1,-3
p=3,0 v=-1,-2
p=9,3 v=2,3
p=7,3 v=-1,2
p=2,4 v=2,-3
p=9,5 v=-3,-3`

	// 定义测试用例
	tests := []struct {
		name           string // 测试名称
		input          string // 输入数据
		width          int    // 空间宽度
		height         int    // 空间高度
		simulationTime int    // 模拟时间
		expected       int    // 期望的安全系数
	}{
		{
			name:           "示例测试",
			input:          exampleInput,
			width:          11,
			height:         7,
			simulationTime: 100,
			expected:       12, // 根据题目示例，期望结果为12
		},
	}

	// 遍历并运行所有测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 解析输入数据
			robots, err := ParseInput(tt.input)
			if err != nil {
				t.Fatalf("解析输入失败: %v", err)
			}

			// 调用核心计算函数
			got := CalculateSafetyFactor(robots, tt.width, tt.height, tt.simulationTime)
			// 检查结果是否符合预期
			if got != tt.expected {
				t.Errorf("CalculateSafetyFactor() 得到 %d, 期望 %d", got, tt.expected)
			}
		})
	}
}
