package main

import (
	"fmt"
	"os"
	"strings"
)

// Position 表示网格中的一个位置
type Position struct {
	row, col int
}

// readInput 读取输入文件并构建字符网格。
func readInput(filename string) ([][]string, Position, string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, Position{}, "", err
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	grid := make([][]string, len(lines))

	// 查找警卫的初始位置和方向
	var guardPos Position
	var guardDir string

	for i, line := range lines {
		line = strings.TrimSuffix(line, "\r") // 处理Windows风格的行结束符
		grid[i] = strings.Split(line, "")

		// 查找警卫的位置
		for j, char := range grid[i] {
			if char == "^" || char == ">" || char == "v" || char == "<" {
				guardPos = Position{i, j}
				guardDir = char
				// 不替换原始网格中的字符，保留警卫标记
			}
		}
	}

	return grid, guardPos, guardDir, nil
}

// checkLoop 检查警卫是否会进入循环
func checkLoop(grid [][]string, startPos Position, startDir string, maxSteps int) bool {
	// 记录警卫的状态 (位置+方向)
	type State struct {
		pos Position
		dir string
	}

	visited := make(map[State]int) // 状态 -> 步数

	// 当前位置和方向
	pos := startPos
	dir := startDir
	steps := 0

	for steps < maxSteps {
		// 记录当前状态
		currentState := State{pos, dir}
		if _, exists := visited[currentState]; exists {
			// 找到循环
			return true
		}
		visited[currentState] = steps

		var nextPos Position

		// 确定前方位置
		switch dir {
		case "^":
			nextPos = Position{pos.row - 1, pos.col}
		case ">":
			nextPos = Position{pos.row, pos.col + 1}
		case "v":
			nextPos = Position{pos.row + 1, pos.col}
		case "<":
			nextPos = Position{pos.row, pos.col - 1}
		}

		// 检查是否离开地图
		if nextPos.row < 0 || nextPos.row >= len(grid) || nextPos.col < 0 || nextPos.col >= len(grid[0]) {
			return false // 离开地图，没有循环
		}

		// 获取前方的内容
		nextCell := grid[nextPos.row][nextPos.col]

		// 如果前方有障碍物，向右转
		if nextCell == "#" {
			switch dir {
			case "^":
				dir = ">"
			case ">":
				dir = "v"
			case "v":
				dir = "<"
			case "<":
				dir = "^"
			}
		} else {
			// 否则，向前移动
			pos = nextPos
		}

		steps++
	}

	return false // 达到最大步数仍未找到循环
}
// countObstaclePositions 计算可以添加障碍物使警卫进入循环的位置数量
func countObstaclePositions(grid [][]string, guardPos Position, guardDir string) int {
	count := 0
	maxSteps := 10000 // 设置一个合理的最大步数限制

	// 创建一个网格副本
	copyGrid := make([][]string, len(grid))
	for i := range grid {
		copyGrid[i] = make([]string, len(grid[i]))
		copy(copyGrid[i], grid[i])
	}

	// 尝试在每个空位置添加障碍物
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			// 跳过已有障碍物或警卫的位置
			if grid[i][j] == "#" || (i == guardPos.row && j == guardPos.col) {
				continue
			}

			// 添加障碍物
			copyGrid[i][j] = "#"

			// 检查是否会形成循环
			if checkLoop(copyGrid, guardPos, guardDir, maxSteps) {
				count++
			}

			// 恢复原始状态
			copyGrid[i][j] = grid[i][j]
		}
	}

	return count
}

func main() {
	const inputFile = "input"
	grid, guardPos, guardDir, err := readInput(inputFile)
	if err != nil {
		fmt.Printf("Failed to read input (%s): %v\n", inputFile, err)
		return
	}

	// 计算可以添加障碍物使警卫进入循环的位置数量
	obstacleCount := countObstaclePositions(grid, guardPos, guardDir)
	fmt.Printf("Number of positions where adding an obstacle creates a loop: %d\n", obstacleCount)
}
