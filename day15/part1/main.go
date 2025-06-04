package main

import (
	"fmt"
	"os"
	"strings"
)

// point 结构体用于表示二维坐标 (行, 列)。
type point struct {
	row, col int
}

// parseMap 辅助函数：将多行字符串地图转换为二维 rune 切片。
func parseMap(s string) [][]rune {
	lines := strings.Split(s, "\n")
	var grid [][]rune
	for _, line := range lines {
		if line == "" {
			continue // 跳过空行
		}
		grid = append(grid, []rune(line))
	}
	return grid
}

// cleanMoves 辅助函数：移除移动指令字符串中的所有换行符。
func cleanMoves(moves string) string {
	return strings.ReplaceAll(moves, "\n", "")
}

// getRobotAndBoxes 辅助函数：从初始地图中解析出机器人和所有箱子的位置。
func getRobotAndBoxes(grid [][]rune) (point, []point) {
	var robotPos point
	var boxPos []point
	for r, row := range grid {
		for c, char := range row {
			if char == '@' {
				robotPos = point{r, c}
			} else if char == 'O' {
				boxPos = append(boxPos, point{r, c})
			}
		}
	}
	return robotPos, boxPos
}

// cloneGrid 辅助函数：创建一个二维 rune 切片的深拷贝。
// 这样在模拟过程中修改地图时不会影响到原始地图数据。
func cloneGrid(grid [][]rune) [][]rune {
	newGrid := make([][]rune, len(grid))
	for i := range grid {
		newGrid[i] = make([]rune, len(grid[i]))
		copy(newGrid[i], grid[i])
	}
	return newGrid
}

// calculateGPSCoordinates 辅助函数：计算所有箱子的GPS坐标之和。
// GPS坐标 = 100 * (行) + (列)。
func calculateGPSCoordinates(boxes []point) int {
	totalGPS := 0
	for _, box := range boxes {
		totalGPS += 100*box.row + box.col
	}
	return totalGPS
}

// findBoxAt 检查某个位置是否有箱子，并返回箱子在切片中的索引。
// `ignoreIndex` 参数用于在查找时忽略特定索引的箱子。
// 在链式推动的场景下，ignoreIndex 通常设置为 -1，因为我们需要找到目标位置的任何箱子。
func findBoxAt(boxes []point, r, c int, ignoreIndex int) (int, bool) {
	for i, box := range boxes {
		if i == ignoreIndex { // 如果是需要忽略的箱子，跳过
			continue
		}
		if box.row == r && box.col == c {
			return i, true
		}
	}
	return -1, false // 未找到箱子
}

// solveWarehouse 是核心模拟函数。
// 它接收初始地图字符串和原始移动指令字符串，模拟机器人和箱子的移动，
// 并返回最终箱子的GPS坐标总和。
// 注意：在 main.go 中，我们不传递 *testing.T，所以移除了调试打印。
func solveWarehouse(initialMapStr string, rawMoves string) int {
	initialGrid := parseMap(initialMapStr)
	grid := cloneGrid(initialGrid) // 使用地图的深拷贝进行操作
	robotPos, boxes := getRobotAndBoxes(grid)
	moves := cleanMoves(rawMoves)

	// 清理初始地图显示，将 @ 和 O 的位置变成 .
	// 这样做是为了在后续的碰撞检测中，grid 只反映墙壁，方便判断。
	// 确保清理位置在地图范围内，避免panic
	if robotPos.row >= 0 && robotPos.row < len(grid) && robotPos.col >= 0 && robotPos.col < len(grid[0]) {
		grid[robotPos.row][robotPos.col] = '.'
	}
	for _, box := range boxes {
		if box.row >= 0 && box.row < len(grid) && box.col >= 0 && box.col < len(grid[0]) {
			grid[box.row][box.col] = '.'
		}
	}

	for _, move := range moves {
		dr, dc := 0, 0 // 机器人移动的行和列增量
		switch move {
		case '^': // 向上
			dr = -1
		case 'v': // 向下
			dr = 1
		case '<': // 向左
			dc = -1
		case '>': // 向右
			dc = 1
		}

		nextRobotR, nextRobotC := robotPos.row+dr, robotPos.col+dc

		// 检查机器人是否会移动到地图边界外
		if nextRobotR < 0 || nextRobotR >= len(grid) || nextRobotC < 0 || nextRobotC >= len(grid[0]) {
			continue // 机器人试图移出地图，不移动
		}

		// 检查机器人目标位置是否是墙
		if grid[nextRobotR][nextRobotC] == '#' {
			continue // 机器人撞墙，不移动
		}

		// 检查机器人目标位置是否有箱子（通过箱子列表判断，而不是清理后的grid）
		_, hasBox := findBoxAt(boxes, nextRobotR, nextRobotC, -1) // 初始查找不忽略任何箱子

		if hasBox { // 机器人尝试推动箱子
			// 收集要推动的箱子链
			var pushChain []int // 存储箱子在 boxes 切片中的索引
			currentPushR, currentPushC := nextRobotR, nextRobotC

			for {
				foundBoxIndex, isBox := findBoxAt(boxes, currentPushR, currentPushC, -1) // 查找当前位置的箱子
				if !isBox {
					break // 遇到空地，链条结束
				}
				pushChain = append(pushChain, foundBoxIndex)

				// 检查链条的下一个位置
				currentPushR += dr
				currentPushC += dc

				// 如果链条末端会撞墙或出界，则整个链条不移动
				if currentPushR < 0 || currentPushR >= len(grid) || currentPushC < 0 || currentPushC >= len(grid[0]) || grid[currentPushR][currentPushC] == '#' {
					goto NextMove // 跳到下一个循环迭代 (整个推动失败)
				}
			}

			// 如果链条可以移动，则反向更新所有箱子的位置
			// 从链条末端开始移动，避免覆盖
			for j := len(pushChain) - 1; j >= 0; j-- {
				idx := pushChain[j]
				boxes[idx].row += dr
				boxes[idx].col += dc
			}

			// 机器人移动到新位置
			robotPos = point{nextRobotR, nextRobotC}

		} else { // 机器人移动到空地
			robotPos = point{nextRobotR, nextRobotC}
		}

	NextMove: // 跳转标签，用于跳过当前循环的剩余部分
		// 生产代码中不再打印调试信息
	}

	return calculateGPSCoordinates(boxes)
}

func main() {
	// 1. 读取 input.txt 文件
	input, err := os.ReadFile("input")
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		return
	}

	// 2. 将输入内容分割成地图和移动指令
	// 假定地图和移动指令之间用一个或多个空行分隔
	parts := strings.SplitN(string(input), "\n\n", 2)
	if len(parts) != 2 {
		fmt.Println("Error: Invalid input format. Expected map and moves separated by a double newline.")
		// 尝试处理只有单行分隔的情况
		parts = strings.SplitN(string(input), "\n", 2)
		if len(parts) != 2 {
			fmt.Println("Error: Invalid input format. Could not split map and moves.")
			return
		}
	}

	warehouseMapStr := parts[0]
	movesStr := parts[1]

	// 3. 调用 solveWarehouse 函数进行模拟
	result := solveWarehouse(warehouseMapStr, movesStr)

	// 4. 打印最终结果
	fmt.Printf("The sum of all boxes' GPS coordinates after the robot finishes moving is: %d\n", result)
}
