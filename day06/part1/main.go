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

// findDistinctPositions 计算警卫访问的不同位置数量。
func findDistinctPositions(grid [][]string, startPos Position, startDir string) int {
	// 使用map记录已访问的位置
	visited := make(map[Position]bool)

	// 当前位置和方向
	pos := startPos
	dir := startDir

	// 标记起始位置为已访问
	visited[pos] = true

	// 继续直到警卫离开地图区域
	for {
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
			break
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
			visited[pos] = true
		}
	}

	return len(visited)
}

func main() {
	const inputFile = "input"
	grid, guardPos, guardDir, err := readInput(inputFile)
	if err != nil {
		fmt.Printf("Failed to read input (%s): %v\n", inputFile, err)
		return
	}

	totalDistinctPositions := findDistinctPositions(grid, guardPos, guardDir)
	fmt.Printf("Number of distinct positions visited by the guard: %d\n", totalDistinctPositions)
}
