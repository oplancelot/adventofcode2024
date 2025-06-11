package main

import (
	"testing"
)

func TestCalculateComplexitySum(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected int
	}{
		{
			name: "Example from puzzle",
			input: []string{
				"029A",
				"980A",
				"179A",
				"456A",
				"379A",
			},
			expected: 126384,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateComplexitySum(tt.input)

			if result != tt.expected {
				t.Errorf("For input %v, expected %d, but got %d", tt.input, tt.expected, result)
			}
		})
	}
}

// 额外的测试来验证 findInitialAPos
func TestFindInitialAPos(t *testing.T) {
	// 测试数字键盘
	pos, err := findInitialAPos(numericKeypad[:])
	if err != nil {
		t.Errorf("findInitialAPos failed for numericKeypad: %v", err)
	}
	if pos.R != 3 || pos.C != 2 {
		t.Errorf("Expected numeric keypad A at (3,2), got (%d,%d)", pos.R, pos.C)
	}

	// 测试方向键盘
	pos, err = findInitialAPos(directionalKeypad[:])
	if err != nil {
		t.Errorf("findInitialAPos failed for directionalKeypad: %v", err)
	}
	if pos.R != 0 || pos.C != 2 {
		t.Errorf("Expected directional keypad A at (0,2), got (%d,%d)", pos.R, pos.C)
	}

	// // 测试找不到 A 的情况
	// 测试找不到 A 的情况
	// 创建一个符合函数签名的键盘，但其中不包含 'A'
	noAKeypad := [1][3]KeypadButton{{'X', 'Y', 'Z'}} // 确保是 [行][3]KeypadButton
	_, err = findInitialAPos(noAKeypad[:])
	if err == nil {
		t.Error("Expected error for keypad without A, got nil")
	}
}

// 额外的测试来验证 isValidMove 和 getButtonAtPos
func TestKeypadUtils(t *testing.T) {
	// 测试 numericKeypad
	if !isValidMove(Pos{0, 0}, numericKeypad[:]) {
		t.Errorf("Expected (0,0) to be valid on numericKeypad")
	}
	if isValidMove(Pos{3, 0}, numericKeypad[:]) { // 间隙
		t.Errorf("Expected (3,0) to be invalid (empty) on numericKeypad")
	}
	if getButtonAtPos(Pos{0, 0}, numericKeypad[:]) != '7' {
		t.Errorf("Expected (0,0) on numericKeypad to be '7'")
	}
	if getButtonAtPos(Pos{3, 0}, numericKeypad[:]) != Empty {
		t.Errorf("Expected (3,0) on numericKeypad to be Empty")
	}

	// 测试 directionalKeypad
	if !isValidMove(Pos{0, 1}, directionalKeypad[:]) { // '^'
		t.Errorf("Expected (0,1) to be valid on directionalKeypad")
	}
	if isValidMove(Pos{0, 0}, directionalKeypad[:]) { // 间隙
		t.Errorf("Expected (0,0) to be invalid (empty) on directionalKeypad")
	}
	if getButtonAtPos(Pos{0, 1}, directionalKeypad[:]) != '^' {
		t.Errorf("Expected (0,1) on directionalKeypad to be '^'")
	}
	if getButtonAtPos(Pos{0, 0}, directionalKeypad[:]) != Empty {
		t.Errorf("Expected (0,0) on directionalKeypad to be Empty")
	}
}
