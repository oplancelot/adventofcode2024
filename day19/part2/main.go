package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// CountWays 计算一个设计可以由一组毛巾模式拼接而成的方法总数。
//
// 参数:
//
//	design - 目标设计字符串。
//	patterns - 包含所有可用毛巾模式的字符串切片。
//
// 返回:
//
//	可以构成该设计的不同方法的总数。
//
// 实现思路:
// 动态规划。dp[i] 存储构成 design[:i] 的方法总数。
// dp[0] = 1 (构成空字符串有1种方法)。
// 状态转移方程: dp[i] = Σ dp[j]，对于所有 j < i 且 design[j:i] 是有效模式的情况。
func CountWays(design string, patterns []string) int {
	patternSet := make(map[string]bool)
	for _, p := range patterns {
		patternSet[p] = true
	}

	// dp[i] 存储构成前 i 个字符的方法数
	dp := make([]int, len(design)+1)
	dp[0] = 1 // 基础案例：构成空字符串有 1 种方法

	for i := 1; i <= len(design); i++ {
		// 遍历所有可能的分割点 j
		for j := 0; j < i; j++ {
			// 如果子字符串 design[j:i] 是一个有效的毛巾模式
			if patternSet[design[j:i]] {
				// 将构成 design[:j] 的方法数累加到 dp[i]
				dp[i] += dp[j]
			}
		}
	}

	return dp[len(design)]
}

func main() {
	const filename = "input"

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("错误：无法打开文件 '%s': %v", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// 解析模式（包含上一部分的空白清理逻辑）
	if !scanner.Scan() {
		log.Fatal("错误：输入文件为空或无法读取第一行。")
	}
	rawPatterns := strings.Split(scanner.Text(), ",")
	patterns := make([]string, len(rawPatterns))
	for i, p := range rawPatterns {
		patterns[i] = strings.TrimSpace(p)
	}
	fmt.Printf("已加载并清理了 %d 个毛巾模式。\n", len(patterns))

	// 跳过空行
	if !scanner.Scan() {
		log.Fatal("错误：文件中缺少空行分隔符。")
	}

	// 计算所有设计的方法总和
	totalWays := 0
	designCount := 0
	fmt.Println("\n开始计算每个设计的方法数...")

	for scanner.Scan() {
		design := scanner.Text()
		if design == "" {
			continue
		}
		designCount++

		// 调用新的计数函数
		ways := CountWays(design, patterns)
		if ways > 0 {
			fmt.Printf("- 设计 #%d: ✅ 有 %d 种方法\n", designCount, ways)
			totalWays += ways
		} else {
			fmt.Printf("- 设计 #%d: ❌ 不可能 (0 种方法)\n", designCount)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("读取文件时发生错误: %v", err)
	}

	// 打印最终结果
	fmt.Println("\n--- 计算完成 ---")
	fmt.Printf("所有可行的设计的拼接方法总数为: %d\n", totalWays)
}
