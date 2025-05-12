package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// 读取输入
func readInput(filename string) ([][2]int, [][]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	var rules [][2]int
	var updates [][]int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.Contains(line, "|") {
			parts := strings.Split(line, "|")
			if len(parts) != 2 {
				fmt.Println("Invalid rule line:", line)
				continue
			}
			a, err1 := strconv.Atoi(parts[0])
			b, err2 := strconv.Atoi(parts[1])
			if err1 != nil || err2 != nil {
				fmt.Println("Invalid rule integers:", line)
				continue
			}
			rules = append(rules, [2]int{a, b})
		} else if strings.Contains(line, ",") {
			parts := strings.Split(line, ",")
			var update []int
			for _, p := range parts {
				num, err := strconv.Atoi(strings.TrimSpace(p))
				if err != nil {
					fmt.Println("Invalid update integer:", p)
					continue
				}
				update = append(update, num)
			}
			updates = append(updates, update)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return rules, updates, nil
}

// 检查单个更新是否符合所有规则
func isValidUpdate(update []int, rules [][2]int) bool {
	positions := make(map[int]int)
	for i, v := range update {
		positions[v] = i
	}

	for _, rule := range rules {
		a, b := rule[0], rule[1]
		posA, hasA := positions[a]
		posB, hasB := positions[b]

		// 只有在 A 和 B 都存在时，才检查是否 A 在 B 前面
		if hasA && hasB {
			if posA >= posB {
				return false
			}
		}
	}
	return true
}

func topologicalSort(update []int, rules [][2]int) []int {
	// 仅对 update 中的元素进行排序
	graph := make(map[int][]int)
	inDegree := make(map[int]int)
	inUpdate := make(map[int]bool)

	// 先标记哪些元素在 update 中
	for _, v := range update {
		inUpdate[v] = true
		inDegree[v] = 0
	}

	// 构建图
	for _, rule := range rules {
		a, b := rule[0], rule[1]
		if inUpdate[a] && inUpdate[b] {
			graph[a] = append(graph[a], b)
			inDegree[b]++
		}
	}

	// 拓扑排序（Kahn 算法）
	queue := []int{}
	for _, v := range update {
		if inDegree[v] == 0 {
			queue = append(queue, v)
		}
	}

	sorted := []int{}
	used := make(map[int]bool)

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		sorted = append(sorted, node)
		used[node] = true

		for _, neighbor := range graph[node] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	// 如果有遗漏节点（没有在拓扑图中参与），按原顺序补上
	for _, v := range update {
		if !used[v] {
			sorted = append(sorted, v)
		}
	}
	return sorted
}

func processUpdates(updates [][]int, rules [][2]int) int {
	sum := 0
	for _, update := range updates {
		if !isValidUpdate(update, rules) {
			sorted := topologicalSort(update, rules)
			mid := sorted[len(sorted)/2]
			sum += mid
		}
	}
	return sum
}

func main() {
	const inputFile = "input"

	rules, updates, err := readInput(inputFile)
	if err != nil {
		fmt.Printf("Failed to read input (%s): %v\n", inputFile, err)
		return
	}

	sums := processUpdates(updates, rules)
	fmt.Println("sums of middle elements \n", sums)
}
