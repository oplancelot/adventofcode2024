package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type Point struct {
	R, C int
}

type Interval struct {
	Start, End int
}

func isValid(r, c, numRows, numCols int) bool {
	return r >= 0 && r < numRows && c >= 0 && c < numCols
}

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

func findRegionsAndCalculateTotalPrice(grid [][]rune) int {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return 0
	}
	numRows := len(grid)
	numCols := len(grid[0])

	// 定义新的扩展网格尺寸
	// expandedGrid[2r][2c] 存储原始单元格类型
	// expandedGrid[2r+1][2c] 存储垂直栅栏状态
	// expandedGrid[2r][2c+1] 存储水平栅栏状态
	// expandedGrid[2r+1][2c+1] 存储对角线交点栅栏状态
	expandedNumRows := numRows*2 - 1 // 考虑边界，实际只有 (N-1) 个中间点
	expandedNumCols := numCols*2 - 1
	if numRows == 1 {
		expandedNumRows = 1
	} // 单行网格，没有中间行
	if numCols == 1 {
		expandedNumCols = 1
	} // 单列网格，没有中间列

	// 如果原始网格为空，或者只有一行/一列
	if numRows == 0 || numCols == 0 {
		return 0
	}
	if numRows == 1 && numCols == 1 { // 1x1 网格
		return 1 * 4 // 1个面积，4条边
	}
	// 如果是 1xn 或 nx1，需要调整 expandedNumRows/Cols
	if numRows == 1 && numCols > 1 {
		expandedNumRows = 1
	} else if numRows > 1 {
		expandedNumRows = numRows*2 - 1
	}
	if numCols == 1 && numRows > 1 {
		expandedNumCols = 1
	} else if numCols > 1 {
		expandedNumCols = numCols*2 - 1
	}

	expGrid := make([][]rune, expandedNumRows)
	for i := range expGrid {
		expGrid[i] = make([]rune, expandedNumCols)
		// 初始填充，0 表示空或分隔，其他表示植物类型
		for j := range expGrid[i] {
			expGrid[i][j] = 0 // 默认分隔
		}
	}

	// 填充原始植物类型
	for r := 0; r < numRows; r++ {
		for c := 0; c < numCols; c++ {
			expGrid[r*2][c*2] = grid[r][c]
		}
	}

	// 处理水平和垂直连接
	for r := 0; r < numRows; r++ {
		for c := 0; c < numCols; c++ {
			// 右侧连接
			if c < numCols-1 {
				if grid[r][c] == grid[r][c+1] {
					expGrid[r*2][c*2+1] = grid[r][c] // 如果同类型，则连接
				} else {
					expGrid[r*2][c*2+1] = ' ' // 否则为边界，用空格或其他非零但非植物字符表示
				}
			}
			// 下方连接
			if r < numRows-1 {
				if grid[r][c] == grid[r+1][c] {
					expGrid[r*2+1][c*2] = grid[r][c] // 如果同类型，则连接
				} else {
					expGrid[r*2+1][c*2] = ' ' // 否则为边界
				}
			}
		}
	}

	// 处理对角线交点
	// expandedGrid[2r+1][2c+1] 对应原始网格的 (r,c), (r,c+1), (r+1,c), (r+1,c+1) 的中心
	for r := 0; r < numRows-1; r++ {
		for c := 0; c < numCols-1; c++ {
			topLeft := grid[r][c]
			topRight := grid[r][c+1]
			bottomLeft := grid[r+1][c]
			bottomRight := grid[r+1][c+1]

			// 根据题目说明：如果对角线上的单元格类型相同，但与另一对角线上的单元格类型不同，则栅栏不连接。
			// 此时，中心点应为分隔，即使是相同类型的植物，也不能通过此点连接。
			if topLeft == bottomRight && topRight == bottomLeft && topLeft != topRight {
				expGrid[r*2+1][c*2+1] = ' ' // ' ' 表示不可通过，是栅栏
			} else {
				// 如果不是 X 型交界，且四角都是相同类型，那么中心可以连接
				if topLeft == topRight && topLeft == bottomLeft && topLeft == bottomRight {
					expGrid[r*2+1][c*2+1] = topLeft // 相同类型，可以连接
				} else {
					// 如果有混合类型，中心点也是分隔
					expGrid[r*2+1][c*2+1] = ' '
				}
			}
		}
	}

	visited := make([][]bool, expandedNumRows)
	for i := range visited {
		visited[i] = make([]bool, expandedNumCols)
	}
	totalPrice := 0

	// 邻居的四个方向：上, 下, 左, 右 (只考虑垂直和水平移动)
	dr := []int{-1, 1, 0, 0}
	dc := []int{0, 0, -1, 1}

	for r := 0; r < expandedNumRows; r++ {
		for c := 0; c < expandedNumCols; c++ {
			// 只对原始单元格对应的点 (2r, 2c) 进行 BFS
			if r%2 == 0 && c%2 == 0 && !visited[r][c] && expGrid[r][c] != ' ' && expGrid[r][c] != 0 {
				currentPlantType := expGrid[r][c]
				currentArea := 0
				regionHedges := make(map[int][]Interval) // y -> list of x-intervals
				regionVedges := make(map[int][]Interval) // x -> list of y-intervals

				q := []Point{{R: r, C: c}}
				visited[r][c] = true
				head := 0

				for head < len(q) {
					curr := q[head]
					head++

					// 如果当前是原始单元格的点 (2r, 2c)，则面积+1
					if curr.R%2 == 0 && curr.C%2 == 0 {
						currentArea++
					}

					for i := 0; i < 4; i++ {
						nr, nc := curr.R+dr[i], curr.C+dc[i]
						if isValid(nr, nc, expandedNumRows, expandedNumCols) {
							// 检查邻居类型是否相同且未访问
							if expGrid[nr][nc] == currentPlantType && !visited[nr][nc] {
								visited[nr][nc] = true
								q = append(q, Point{R: nr, C: nc})
							} else if expGrid[nr][nc] != currentPlantType && expGrid[nr][nc] != 0 { // 边界
								// 只有当邻居是' ' (栅栏) 或不同类型的植物时，才算作边界
								// 记录边界段
								// 注意：边界是在 expGrid 中计算的
								if dr[i] == -1 { // 邻居在上方，即当前点的上边缘
									y := curr.R // expGrid 中的行
									interval := Interval{Start: curr.C, End: curr.C + 1}
									regionHedges[y] = append(regionHedges[y], interval)
								} else if dr[i] == 1 { // 邻居在下方，即当前点的下边缘
									y := curr.R + 1 // expGrid 中的行
									interval := Interval{Start: curr.C, End: curr.C + 1}
									regionHedges[y] = append(regionHedges[y], interval)
								} else if dc[i] == -1 { // 邻居在左侧，即当前点的左边缘
									x := curr.C // expGrid 中的列
									interval := Interval{Start: curr.R, End: curr.R + 1}
									regionVedges[x] = append(regionVedges[x], interval)
								} else { // 邻居在右侧，即当前点的右边缘
									x := curr.C + 1 // expGrid 中的列
									interval := Interval{Start: curr.R, End: curr.R + 1}
									regionVedges[x] = append(regionVedges[x], interval)
								}
							}
						} else { // 超出边界，也算作边界
							if dr[i] == -1 {
								y := curr.R
								interval := Interval{Start: curr.C, End: curr.C + 1}
								regionHedges[y] = append(regionHedges[y], interval)
							} else if dr[i] == 1 {
								y := curr.R + 1
								interval := Interval{Start: curr.C, End: curr.C + 1}
								regionHedges[y] = append(regionHedges[y], interval)
							} else if dc[i] == -1 {
								x := curr.C
								interval := Interval{Start: curr.R, End: curr.R + 1}
								regionVedges[x] = append(regionVedges[x], interval)
							} else {
								x := curr.C + 1
								interval := Interval{Start: curr.R, End: curr.R + 1}
								regionVedges[x] = append(regionVedges[x], interval)
							}
						}
					}
				}

				numberOfSides := 0
				for _, segments := range regionHedges {
					numberOfSides += countAndMergeSegments(segments)
				}
				for _, segments := range regionVedges {
					numberOfSides += countAndMergeSegments(segments)
				}

				regionPrice := currentArea * numberOfSides
				totalPrice += regionPrice

				fmt.Printf("区域 %c: 面积=%d, 边数=%d, 价格=%d\n",
					currentPlantType, currentArea, numberOfSides, regionPrice)
			}
		}
	}
	return totalPrice
}

func countAndMergeSegments(segments []Interval) int {
	if len(segments) == 0 {
		return 0
	}

	// 按起点排序
	sort.Slice(segments, func(i, j int) bool {
		return segments[i].Start < segments[j].Start
	})

	// 合并相邻或重叠的段
	mergedCount := 1 // 计入第一个段
	currentEnd := segments[0].End

	for i := 1; i < len(segments); i++ {
		current := segments[i]

		if current.Start <= currentEnd {
			// 合并段
			if current.End > currentEnd {
				currentEnd = current.End
			}
		} else {
			// 添加新段
			mergedCount++
			currentEnd = current.End
		}
	}

	return mergedCount
}

func main() {
	const inputFile = "input"
	inputData, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("从 %s 读取谜题输入失败: %v", inputFile, err)
	}
	grid := parseInput(string(inputData))

	totalPrice := findRegionsAndCalculateTotalPrice(grid)
	fmt.Println(totalPrice)
}
