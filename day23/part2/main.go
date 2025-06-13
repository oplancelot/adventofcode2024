package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// solve 负责解决谜题
func solve(input string) string {
	adj := make(map[string]map[string]bool)
	var nodes []string
	nodesSet := make(map[string]struct{})

	// 1. 解析输入，构建图
	lines := strings.Split(strings.TrimSpace(input), "\n")
	for _, line := range lines {
		parts := strings.Split(line, "-")
		u, v := parts[0], parts[1]

		if adj[u] == nil {
			adj[u] = make(map[string]bool)
		}
		if adj[v] == nil {
			adj[v] = make(map[string]bool)
		}
		adj[u][v] = true
		adj[v][u] = true
		nodesSet[u] = struct{}{}
		nodesSet[v] = struct{}{}
	}
	for node := range nodesSet {
		nodes = append(nodes, node)
	}
	sort.Strings(nodes)

	// 2. 迭代寻找最大团
	// 从大小为1的团（每个节点自身）开始
	cliques := make([][]string, 0)
	for _, node := range nodes {
		cliques = append(cliques, []string{node})
	}

	for {
		nextCliques := make([][]string, 0)
		nextCliquesMap := make(map[string]bool)

		for _, clique := range cliques {
			for _, node := range nodes {
				// 检查节点是否可以扩展当前团
				isExtendable := true
				
				// 检查新节点是否已在团中
				inClique := false
				for _, member := range clique {
					if member == node {
						inClique = true
						break
					}
				}
				if inClique {
					continue
				}

				// 检查新节点是否与团中所有成员都相连
				for _, member := range clique {
					if !adj[node][member] {
						isExtendable = false
						break
					}
				}

				if isExtendable {
					// 创建新团并保持排序以生成唯一键
					newClique := append([]string{}, clique...)
					newClique = append(newClique, node)
					sort.Strings(newClique)
					key := strings.Join(newClique, ",")
					
					if !nextCliquesMap[key] {
						nextCliques = append(nextCliques, newClique)
						nextCliquesMap[key] = true
					}
				}
			}
		}

		if len(nextCliques) == 0 {
			// 无法找到更大的团，当前 cliques 是最大的
			break
		}
		cliques = nextCliques
	}
	
	// 3. 格式化输出
	// 此时, cliques[0] 就是最大团之一（题目暗示唯一）
	largestClique := cliques[0]
	sort.Strings(largestClique) // 确保最终排序
	return strings.Join(largestClique, ",")
}

func main() {
	data, err := os.ReadFile("input")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(solve(string(data)))
}