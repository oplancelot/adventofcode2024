// main.go
package main

import (
	"bufio"
	"container/heap" // 用于优先队列
	"fmt"
	"log"
	"os"
)

// --- 核心寻路逻辑 ---

// 方向常量
const (
	East  = 0 // 东
	South = 1 // 南
	West  = 2 // 西
	North = 3 // 北
)

// 对应方向的行、列变化量: dr[East], dc[East] ...
var (
	dr = []int{0, 1, 0, -1} // 行变化: 东, 南, 西, 北
	dc = []int{1, 0, -1, 0} // 列变化: 东, 南, 西, 北
)

// State 代表搜索过程中的一个状态：位置 (R, C)，朝向 (Dir)，以及到达此状态的当前分数 (Score)
type State struct {
	R, C  int
	Dir   int
	Score int
}

// StateKey 用作 minScores 哈希表的键，唯一标识一个状态 (不包含分数)
type StateKey struct {
	R, C, Dir int
}

// PriorityQueue 是一个最小堆，根据 State 的 Score 排序
type PriorityQueue []*State

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// 我们希望得到分数最小的 State，所以是最小堆
	return pq[i].Score < pq[j].Score
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*State)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // 防止内存泄漏
	*pq = old[0 : n-1]
	return item
}

// findLowestScore 函数计算从 'S' 到 'E' 的最低分数
func findLowestScore(maze []string) int {
	rows := len(maze)
	if rows == 0 {
		return -1 // 空迷宫
	}
	cols := len(maze[0])
	if cols == 0 {
		return -1 // 空迷宫
	}

	var startR, startC int = -1, -1
	// 寻找起点 'S'
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if maze[r][c] == 'S' {
				startR, startC = r, c
				break
			}
		}
		if startR != -1 {
			break
		}
	}

	if startR == -1 {
		log.Println("错误: 未在迷宫中找到起点 'S'")
		return -1
	}

	// minScores 记录到达状态 (R, C, Dir) 的已知最低分数
	minScores := make(map[StateKey]int)

	// 初始化优先队列
	pq := &PriorityQueue{}
	heap.Init(pq)

	// 初始状态：在 'S' 位置，朝向东，分数为 0
	initialState := &State{R: startR, C: startC, Dir: East, Score: 0}
	initialStateKey := StateKey{R: startR, C: startC, Dir: East}

	heap.Push(pq, initialState)
	minScores[initialStateKey] = 0

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*State)
		r, c, dir, score := current.R, current.C, current.Dir, current.Score

		// 如果当前路径的分数比已记录到达此状态 (r,c,dir) 的最小分数还要高，则跳过
		// 这是因为一个更优的路径已经被处理或已在队列中
		currentKey := StateKey{R: r, C: c, Dir: dir}
		if recordedMinScore, ok := minScores[currentKey]; ok && score > recordedMinScore {
			continue
		}

		// 如果到达 'E' (终点)，则返回当前分数，因为Dijkstra保证这是最短路径
		if maze[r][c] == 'E' {
			return score
		}

		// 尝试操作1: 前进一步
		nr, nc := r+dr[dir], c+dc[dir] // 根据当前方向计算新位置
		newScoreMove := score + 1

		// 检查新位置是否有效 (界内、非墙)
		if nr >= 0 && nr < rows && nc >= 0 && nc < cols && maze[nr][nc] != '#' {
			moveStateKey := StateKey{R: nr, C: nc, Dir: dir}
			// 如果新路径更优 (或首次到达)，则更新分数并加入队列
			if val, ok := minScores[moveStateKey]; !ok || newScoreMove < val {
				minScores[moveStateKey] = newScoreMove
				heap.Push(pq, &State{R: nr, C: nc, Dir: dir, Score: newScoreMove})
			}
		}

		// 尝试操作2: 顺时针旋转90度
		newDirCW := (dir + 1) % 4 // (0E -> 1S -> 2W -> 3N -> 0E)
		newScoreRotate := score + 1000
		rotateCWStateKey := StateKey{R: r, C: c, Dir: newDirCW}
		if val, ok := minScores[rotateCWStateKey]; !ok || newScoreRotate < val {
			minScores[rotateCWStateKey] = newScoreRotate
			heap.Push(pq, &State{R: r, C: c, Dir: newDirCW, Score: newScoreRotate})
		}

		// 尝试操作3: 逆时针旋转90度
		newDirCCW := (dir - 1 + 4) % 4 // `+4` 确保结果为正 (0E -> 3N -> 2W -> 1S -> 0E)
		// 旋转分数是相同的 newScoreRotate
		rotateCCWStateKey := StateKey{R: r, C: c, Dir: newDirCCW}
		if val, ok := minScores[rotateCCWStateKey]; !ok || newScoreRotate < val {
			minScores[rotateCCWStateKey] = newScoreRotate
			heap.Push(pq, &State{R: r, C: c, Dir: newDirCCW, Score: newScoreRotate})
		}
	}

	return -1 // 如果队列为空仍未到达 'E'，说明无法到达
}

// main 函数保持不变，用于读取输入并调用 findLowestScore
func main() {
	filePath := "input" // 默认输入文件名
	if len(os.Args) > 1 {
		filePath = os.Args[1] // 允许通过命令行参数指定
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("错误：无法打开文件 '%s': %v", filePath, err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("错误：读取文件 '%s' 时发生错误: %v", filePath, err)
	}

	if len(lines) == 0 {
		log.Fatalf("错误：文件 '%s' 为空或无法读取内容。", filePath)
	}

	fmt.Printf("正在分析地图: %s\n", filePath)
	scoreResult := findLowestScore(lines)
	if scoreResult == -1 {
		fmt.Println("未能找到到达终点 'E' 的路径。")
	} else {
		fmt.Printf("可以获得的最低分数是: %d\n", scoreResult)
	}
}
