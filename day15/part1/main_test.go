package main

import (
	"testing"
)

// TestCase 结构体定义了一个测试用例的输入和期望输出。
type TestCase struct {
	name      string // 测试用例的名称
	warehouse string // 仓库地图的字符串表示
	moves     string // 机器人移动指令的字符串表示
	expected  int    // 期望的箱子GPS坐标总和
}

// printDebugMap 辅助函数：用于在测试调试时打印当前地图状态。
// 它会根据机器人和箱子的当前位置动态构建并打印地图。
func printDebugMap(t *testing.T, baseGrid [][]rune, robotPos point, boxes []point, moveIdx int, moveRune rune) {
	t.Helper() // 标记为辅助函数，错误报告会指向调用它的地方

	// 克隆一个只包含墙壁的基础网格
	tempGrid := cloneGrid(baseGrid)

	// 在临时网格上放置箱子
	for _, box := range boxes {
		// 确保箱子位置在地图范围内，避免panic
		if box.row >= 0 && box.row < len(tempGrid) && box.col >= 0 && box.col < len(tempGrid[0]) {
			tempGrid[box.row][box.col] = 'O'
		}
	}
	// 在临时网格上放置机器人
	// 确保机器人位置在地图范围内
	if robotPos.row >= 0 && robotPos.row < len(tempGrid) && robotPos.col >= 0 && robotPos.col < len(tempGrid[0]) {
		tempGrid[robotPos.row][robotPos.col] = '@'
	}

	// 打印当前移动信息和地图状态
	if moveRune == ' ' { // 初始状态时，moveRune 为空格
		t.Logf("--- Initial State ---")
	} else {
		t.Logf("--- After Move %d (%c) ---", moveIdx, moveRune)
	}
	for _, row := range tempGrid {
		t.Logf("%s", string(row))
	}
	t.Log("\n") // 打印空行以分隔不同状态
}

// TestSolveWarehouse 函数用于运行所有预定义的测试用例。
func TestSolveWarehouse(t *testing.T) {
	testCases := []TestCase{
		{
			name: "Small Example",
			warehouse: `########
#..O.O.#
##@.O..#
#...O..#
#.#.O..#
#...O..#
#......#
########`,
			moves:    `<^^>>>vv<v>>v<<`,
			expected: 2028, // 经过手动验证，此期望值与题目描述和图示一致
		},
		{
			name: "Large Example",
			warehouse: `##########
#..O..O.O#
#......O.#
#.OO..O.O#
#..O@..O.#
#O#..O...#
#O..O..O.#
#.OO.O.OO#
#....O...#
##########`,
			// 为了确保复制粘贴的准确性，建议直接从 Advent of Code 网站复制这段指令。
			moves: `<vv>^<v^>v>^vv^v>v<>v^v<v<^vv<<<^><<><>>v<vvv<>^v^>^<<<><<v<<<v^vv^v>^
vvv<<^>^v^^><<>>><>^<<><^vv^^<>vvv<>><^^v>^>vv<>v<<<<v<^v>^<^^>>>^<v<v
><>vv>v^v^<>><>>>><^^>vv>v<^^^>>v^v^<^^>v^^>v^<^v>v<>>v^v^<v>v^^<^^vv<
<<v<^>>^^^^>>>v^<>vvv^><v<<<>^^^vv^<vvv>^>v<^^^^v<>^>vvvv><>>v^<<^^^^^
^><^><>>><>^^<<^^v>>><^<v>^<vv>>v>>>^v><>^v><<<<v>>v<v<v>vvv>^<><<>^><
^>><>^v<><^vvv<^^<><v<<<<<><^v<<<><<<^^<v<^^^><^>>^<v^><<<^>>^v<v^v<v^
>^>>^v>vv>^<<^v<>><<><<v<<v><>v<^vv<<<>^^v^>^^>>><<^v>>v^v><^^>>^<>vv^
<><^^>^^^<><vvvvv^v<v<<>^v<v>v<<^><<><<><<<^^<<<^<<>><<><^^^>^^<>^>v<>
^^>vv<^v^v<vv>^<><v<^v>^^^>>>^^vvv^>vvv<>>>^<^>>>>>^<<^v>^vvv<>^<><<v>
v^^>>><<^^<>>^v^<v^vv<>v^<<>^<^v^v><^<<<><<^<v><v<>vv>>v><v^<vv<>v^<<^`,
			expected: 10092, // 题目提供的期望值
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 调用核心模拟函数，并传递 *testing.T 实例用于打印调试信息。
			actual := solveWarehouse(tc.warehouse, tc.moves)
			if actual != tc.expected {
				t.Errorf("For test case %s: expected %d, got %d", tc.name, tc.expected, actual)
				// 此时，详细的调试日志会通过 t.Logf 显示出来，帮助定位问题。
			}
		})
	}
}
