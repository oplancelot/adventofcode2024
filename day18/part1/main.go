package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Point 定义了坐标点
type Point struct {
	X, Y int
}

// QueueItem 用于 BFS 队列，包含坐标点和到达该点的步数
type QueueItem struct {
	Point Point
	Dist  int
}

// findShortestPath 使用广度优先搜索 (BFS) 寻找最短路径
func findShortestPath(width, height, byteCount int, bytePositions []Point) int {
	// 创建网格并标记障碍物
	grid := make([][]bool, height)
	for i := range grid {
		grid[i] = make([]bool, width)
	}

	for i := 0; i < byteCount && i < len(bytePositions); i++ {
		p := bytePositions[i]
		if p.Y >= 0 && p.Y < height && p.X >= 0 && p.X < width {
			grid[p.Y][p.X] = true // true 表示被破坏
		}
	}

	start := Point{0, 0}
	end := Point{width - 1, height - 1}

	// 如果起点或终点是障碍物，则无解
	if grid[start.Y][start.X] || grid[end.Y][end.X] {
		return -1
	}

	// 初始化访问记录
	visited := make([][]bool, height)
	for i := range visited {
		visited[i] = make([]bool, width)
	}

	queue := []QueueItem{{Point: start, Dist: 0}}
	visited[start.Y][start.X] = true

	// 定义四个移动方向：上、下、左、右
	dx := []int{0, 0, -1, 1}
	dy := []int{-1, 1, 0, 0}

	for len(queue) > 0 {
		currentItem := queue[0]
		queue = queue[1:]

		p := currentItem.Point

		// 到达终点
		if p == end {
			return currentItem.Dist
		}

		// 探索相邻节点
		for i := 0; i < 4; i++ {
			nextX, nextY := p.X+dx[i], p.Y+dy[i]

			// 检查边界、是否为障碍物、是否已访问
			if nextX >= 0 && nextX < width && nextY >= 0 && nextY < height &&
				!grid[nextY][nextX] && !visited[nextY][nextX] {
				visited[nextY][nextX] = true
				queue = append(queue, QueueItem{Point: Point{nextX, nextY}, Dist: currentItem.Dist + 1})
			}
		}
	}

	return -1 // 未找到路径
}

func main() {
	file, err := os.Open("input")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var bytePositions []Point
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		if len(parts) != 2 {
			continue
		}
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])
		bytePositions = append(bytePositions, Point{x, y})
	}

	// 实际谜题要求
	gridSize := 71
	byteCount := 1024
	result := findShortestPath(gridSize, gridSize, byteCount, bytePositions)

	fmt.Printf("The minimum number of steps needed is: %d\n", result)
}
