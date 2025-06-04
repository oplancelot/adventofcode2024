// main.go
package main

import (
	"bufio"
	"container/heap" // 用于优先队列
	"container/list" // 用于回溯时的BFS队列
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

// 对应方向的行、列变化量
var (
	dr = []int{0, 1, 0, -1} // 行变化: 东, 南, 西, 北
	dc = []int{1, 0, -1, 0} // 列变化: 东, 南, 西, 北
)

// State 代表搜索过程中的一个状态
type State struct {
	R, C  int
	Dir   int
	Score int
}

// StateKey 用作哈希表的键，唯一标识一个状态 (不包含分数)
type StateKey struct {
	R, C, Dir int
}

// Point 代表一个图块的坐标，用于标记在最佳路径上的图块
type Point struct {
	R, C int
}

// PriorityQueue 是一个最小堆，根据 State 的 Score 排序
type PriorityQueue []*State

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].Score < pq[j].Score }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(*State)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return item
}

// findLowestScoreAndFullDistMap 运行Dijkstra算法，
// 返回到达'E'的最低分数以及到达所有状态的最小分数映射。
// 1. minScoreToE: 到达任何'E'图块的最低分数 (-1 如果不可达)。
// 2. allDistances: 从'S'到每个StateKey的最小分数映射。
func findLowestScoreAndFullDistMap(maze []string) (int, map[StateKey]int) {
	rows := len(maze)
	if rows == 0 {
		return -1, nil
	}
	cols := len(maze[0])
	if cols == 0 {
		return -1, nil
	}

	var startR, startC int = -1, -1
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
		return -1, nil
	}

	allDistances := make(map[StateKey]int)
	pq := &PriorityQueue{}
	heap.Init(pq)

	initialState := &State{R: startR, C: startC, Dir: East, Score: 0}
	initialStateKey := StateKey{R: startR, C: startC, Dir: East}
	heap.Push(pq, initialState)
	allDistances[initialStateKey] = 0

	minScoreToE := -1

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*State)
		r, c, dir, score := current.R, current.C, current.Dir, current.Score

		currentKey := StateKey{R: r, C: c, Dir: dir}
		if recordedScore, ok := allDistances[currentKey]; !ok || score > recordedScore {
			continue
		}

		if maze[r][c] == 'E' {
			if minScoreToE == -1 || score < minScoreToE {
				minScoreToE = score
			}
		}

		// 尝试操作1: 前进一步
		nr, nc := r+dr[dir], c+dc[dir]
		newScoreMove := score + 1
		if nr >= 0 && nr < rows && nc >= 0 && nc < cols && maze[nr][nc] != '#' {
			moveStateKey := StateKey{R: nr, C: nc, Dir: dir}
			if val, ok := allDistances[moveStateKey]; !ok || newScoreMove < val {
				allDistances[moveStateKey] = newScoreMove
				heap.Push(pq, &State{R: nr, C: nc, Dir: dir, Score: newScoreMove})
			}
		}

		// 尝试操作2: 顺时针旋转
		newDirCW := (dir + 1) % 4
		newScoreRotate := score + 1000
		rotateCWStateKey := StateKey{R: r, C: c, Dir: newDirCW}
		if val, ok := allDistances[rotateCWStateKey]; !ok || newScoreRotate < val {
			allDistances[rotateCWStateKey] = newScoreRotate
			heap.Push(pq, &State{R: r, C: c, Dir: newDirCW, Score: newScoreRotate})
		}

		// 尝试操作3: 逆时针旋转
		newDirCCW := (dir - 1 + 4) % 4
		rotateCCWStateKey := StateKey{R: r, C: c, Dir: newDirCCW}
		if val, ok := allDistances[rotateCCWStateKey]; !ok || newScoreRotate < val {
			allDistances[rotateCCWStateKey] = newScoreRotate
			heap.Push(pq, &State{R: r, C: c, Dir: newDirCCW, Score: newScoreRotate})
		}
	}
	return minScoreToE, allDistances
}

// countTilesOnBestPath (Part 2 函数)
func countTilesOnBestPath(maze []string) int {
	actualMinScore, allMinScoresFromStart := findLowestScoreAndFullDistMap(maze)

	if actualMinScore == -1 {
		return 0 // 'E' 不可达
	}

	onBestPathTiles := make(map[Point]bool)
	tracebackQueue := list.New()                      // BFS 队列
	visitedTracebackStates := make(map[StateKey]bool) // 避免在回溯中重复处理状态

	rows, cols := len(maze), len(maze[0])

	// 初始化回溯队列：从所有以最低总分到达'E'的状态开始
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if maze[r][c] == 'E' {
				for dir := 0; dir < 4; dir++ { // 检查所有4个方向
					key := StateKey{R: r, C: c, Dir: dir}
					if score, ok := allMinScoresFromStart[key]; ok && score == actualMinScore {
						tracebackQueue.PushBack(State{R: r, C: c, Dir: dir, Score: score})
						visitedTracebackStates[key] = true
						onBestPathTiles[Point{R: r, C: c}] = true // 标记'E'图块
					}
				}
			}
		}
	}

	// BFS回溯
	for tracebackQueue.Len() > 0 {
		elem := tracebackQueue.Front()
		tracebackQueue.Remove(elem)
		current := elem.Value.(State)

		cr, cc, cdir, cscore := current.R, current.C, current.Dir, current.Score
		onBestPathTiles[Point{R: cr, C: cc}] = true // 标记当前图块

		// 尝试反转“前进一步”操作
		prMove := cr - dr[cdir]
		pcMove := cc - dc[cdir]
		if prMove >= 0 && prMove < rows && pcMove >= 0 && pcMove < cols && maze[prMove][pcMove] != '#' {
			prevMoveKey := StateKey{R: prMove, C: pcMove, Dir: cdir}
			if val, ok := allMinScoresFromStart[prevMoveKey]; ok && val == cscore-1 {
				if !visitedTracebackStates[prevMoveKey] {
					visitedTracebackStates[prevMoveKey] = true
					tracebackQueue.PushBack(State{R: prMove, C: pcMove, Dir: cdir, Score: cscore - 1})
				}
			}
		}

		// 尝试反转“顺时针旋转”操作
		pDirFromCwRot := (cdir - 1 + 4) % 4
		prevRotCwKey := StateKey{R: cr, C: cc, Dir: pDirFromCwRot}
		if val, ok := allMinScoresFromStart[prevRotCwKey]; ok && val == cscore-1000 {
			if !visitedTracebackStates[prevRotCwKey] {
				visitedTracebackStates[prevRotCwKey] = true
				tracebackQueue.PushBack(State{R: cr, C: cc, Dir: pDirFromCwRot, Score: cscore - 1000})
			}
		}

		// 尝试反转“逆时针旋转”操作
		pDirFromCcwRot := (cdir + 1) % 4
		prevRotCcwKey := StateKey{R: cr, C: cc, Dir: pDirFromCcwRot}
		if val, ok := allMinScoresFromStart[prevRotCcwKey]; ok && val == cscore-1000 {
			if !visitedTracebackStates[prevRotCcwKey] {
				visitedTracebackStates[prevRotCcwKey] = true
				tracebackQueue.PushBack(State{R: cr, C: cc, Dir: pDirFromCcwRot, Score: cscore - 1000})
			}
		}
	}
	return len(onBestPathTiles)
}

func main() {
	filePath := "input"
	if len(os.Args) > 1 {
		filePath = os.Args[1]
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

	fmt.Printf("--- Reindeer Maze: Part Two ---\n")
	tilesCount := countTilesOnBestPath(lines)
	fmt.Printf("最佳路径上的图块数量: %d\n", tilesCount)
}
