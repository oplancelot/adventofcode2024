package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Point 代表网格中的一个坐标 (行, 列)
type Point struct {
	R, C int
}

// parseInput 将地图的字符串表示形式转换为 [][]rune 网格。
// 它会处理输入块周围和每行末尾可能存在的空白字符。
func parseInput(input string) [][]rune {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		return [][]rune{} // 处理空输入或只有空行的输入
	}
	grid := make([][]rune, len(lines))
	for i, line := range lines {
		grid[i] = []rune(strings.TrimSpace(line))
	}
	return grid
}

// isValid 检查点 (r, c) 是否在网格边界内。
func isValid(r, c, numRows, numCols int) bool {
	return r >= 0 && r < numRows && c >= 0 && c < numCols
}

// findRegionsAndCalculateTotalPrice 处理网格以找到所有区域，
// 并返回它们价格的总和。
func findRegionsAndCalculateTotalPrice(grid [][]rune) int {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return 0 // 处理空网格
	}

	numRows := len(grid)
	numCols := len(grid[0])
	visited := make([][]bool, numRows)
	for i := range visited {
		visited[i] = make([]bool, numCols)
	}

	totalPrice := 0

	// 邻居的四个方向：上, 下, 左, 右
	dr := []int{-1, 1, 0, 0}
	dc := []int{0, 0, -1, 1}

	for r := 0; r < numRows; r++ {
		for c := 0; c < numCols; c++ {
			if !visited[r][c] {
				// 开始对新区域进行 BFS
				currentPlantType := grid[r][c]
				currentArea := 0
				currentPerimeter := 0

				q := []Point{{R: r, C: c}} // BFS 队列
				visited[r][c] = true

				head := 0 // BFS 队列的头部指针
				for head < len(q) {
					curr := q[head]
					head++

					currentArea++

					// 检查4个邻居以计算周长和确定区域中的下一个单元格
					for i := 0; i < 4; i++ {
						nr, nc := curr.R+dr[i], curr.C+dc[i]

						if !isValid(nr, nc, numRows, numCols) {
							// 邻居超出边界，对周长有贡献
							currentPerimeter++
						} else {
							// 邻居在边界内
							if grid[nr][nc] != currentPlantType {
								// 邻居是不同类型的植物，对周长有贡献
								currentPerimeter++
							} else {
								// 邻居是相同类型的植物
								if !visited[nr][nc] {
									visited[nr][nc] = true
									q = append(q, Point{R: nr, C: nc})
								}
							}
						}
					}
				}

				// 区域已找到，计算其价格
				regionPrice := currentArea * currentPerimeter
				totalPrice += regionPrice

				// // 如果需要，可以取消注释以打印每个区域的详细信息
				// fmt.Printf("找到区域: 类型 %c, 面积 %d, 周长 %d, 价格 %d\n",
				// currentPlantType, currentArea, currentPerimeter, regionPrice)
			}
		}
	}

	return totalPrice
}

func main() {
	// 通过命令行参数提供了文件名
	const inputFile = "input"
	inputData, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("从 %s 读取谜题输入失败: %v", inputFile, err)
	}
	puzzleGrid := parseInput(string(inputData))
	finalPrice := findRegionsAndCalculateTotalPrice(puzzleGrid)
	fmt.Printf("文件 %s 的总价格: %d\n", inputFile, finalPrice)
}
