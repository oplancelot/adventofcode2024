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

// QueueItem 用于 BFS 队列
type QueueItem struct {
	Point Point
	Dist  int
}

// pathExists 使用广度优先搜索 (BFS) 检查路径是否存在。
// 它是一个精简版的寻路函数，因为我们只需要知道路径是否存在，而不需要长度。
func pathExists(width, height int, corrupted map[Point]bool) bool {
	start := Point{0, 0}
	end := Point{width - 1, height - 1}

	// 如果起点或终点一开始就被阻塞，则路径不存在
	if corrupted[start] || corrupted[end] {
		return false
	}

	visited := make(map[Point]bool)
	queue := []Point{start}
	visited[start] = true

	dx := []int{0, 0, -1, 1}
	dy := []int{-1, 1, 0, 0}

	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]

		if p == end {
			return true // 成功到达终点
		}

		for i := 0; i < 4; i++ {
			nextP := Point{X: p.X + dx[i], Y: p.Y + dy[i]}

			// 检查边界、是否为障碍物、是否已访问
			if nextP.X >= 0 && nextP.X < width && nextP.Y >= 0 && nextP.Y < height &&
				!corrupted[nextP] && !visited[nextP] {
				visited[nextP] = true
				queue = append(queue, nextP)
			}
		}
	}

	return false // 无法到达终点
}

// findBlockingByte 模拟字节坠落，找到第一个阻塞路径的字节
func findBlockingByte(width, height int, bytePositions []Point) (Point, bool) {
	corrupted := make(map[Point]bool)
	for _, p := range bytePositions {
		corrupted[p] = true
		// 在添加新字节后，检查路径是否还存在
		if !pathExists(width, height, corrupted) {
			// 此字节阻塞了路径
			return p, true
		}
	}
	// 所有字节坠落后路径依然存在
	return Point{-1, -1}, false
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

	gridSize := 71
	blockingByte, found := findBlockingByte(gridSize, gridSize, bytePositions)

	if found {
		fmt.Printf("%d,%d\n", blockingByte.X, blockingByte.Y)
	} else {
		fmt.Println("No byte was found that blocked the path.")
	}
}
