package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// solve 负责解决谜题
func solve(input string) int {
	// 使用邻接表表示图
	adj := make(map[string][]string)
	// 使用 map 快速获取所有唯一的节点名
	nodesSet := make(map[string]struct{})

	// 1. 解析输入，构建图
	lines := strings.Split(strings.TrimSpace(input), "\n")
	for _, line := range lines {
		parts := strings.Split(line, "-")
		u, v := parts[0], parts[1]
		// 连接是双向的
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
		nodesSet[u] = struct{}{}
		nodesSet[v] = struct{}{}
	}

	// 辅助函数：检查两个节点是否相连
	isConnected := func(u, v string) bool {
		for _, neighbor := range adj[u] {
			if neighbor == v {
				return true
			}
		}
		return false
	}

	// 2. 获取并排序所有节点，以避免重复计算
	var nodes []string
	for node := range nodesSet {
		nodes = append(nodes, node)
	}
	sort.Strings(nodes)

	// 3. 遍历所有三个节点的组合以查找三角形
	count := 0
	nodeCount := len(nodes)
	for i := 0; i < nodeCount; i++ {
		for j := i + 1; j < nodeCount; j++ {
			for k := j + 1; k < nodeCount; k++ {
				nodeA, nodeB, nodeC := nodes[i], nodes[j], nodes[k]

				// 检查是否为三角形
				if isConnected(nodeA, nodeB) && isConnected(nodeB, nodeC) && isConnected(nodeA, nodeC) {
					// 4. 应用过滤条件
					if strings.HasPrefix(nodeA, "t") || strings.HasPrefix(nodeB, "t") || strings.HasPrefix(nodeC, "t") {
						count++
					}
				}
			}
		}
	}

	return count
}

func main() {
	data, err := os.ReadFile("input")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(solve(string(data)))
}
