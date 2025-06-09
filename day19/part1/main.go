package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// CanFormDesign 函数保持不变
func CanFormDesign(design string, patterns []string) bool {
	patternSet := make(map[string]bool)
	for _, p := range patterns {
		patternSet[p] = true
	}

	dp := make([]bool, len(design)+1)
	dp[0] = true

	for i := 1; i <= len(design); i++ {
		for j := 0; j < i; j++ {
			if dp[j] && patternSet[design[j:i]] {
				dp[i] = true
				break
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

	// --- 2. 解析毛巾模式 (已修复) ---
	if !scanner.Scan() {
		log.Fatal("错误：输入文件为空或无法读取第一行。")
	}
	patternsLine := scanner.Text()
	// 首先，像之前一样用逗号分割
	rawPatterns := strings.Split(patternsLine, ",")

	// 然后，创建一个新的切片来存放清理后的模式
	// 并且遍历每一个原始模式，使用 TrimSpace 清理空白
	patterns := make([]string, len(rawPatterns))
	for i, p := range rawPatterns {
		patterns[i] = strings.TrimSpace(p)
	}

	// 为了验证，我们可以打印清理后的模式
	fmt.Printf("已加载并清理了 %d 个毛巾模式。\n", len(patterns))

	// --- 3. 跳过空行 ---
	if !scanner.Scan() {
		log.Fatal("错误：文件中缺少空行分隔符。")
	}
	_ = scanner.Text()

	// --- 4. 统计可行的设计数量 ---
	possibleDesignsCount := 0
	designCount := 0
	fmt.Println("\n开始检查设计...")

	for scanner.Scan() {
		design := scanner.Text()
		if design == "" { // 如果遇到空行，则跳过
			continue
		}
		designCount++
		if CanFormDesign(design, patterns) {
			possibleDesignsCount++
			fmt.Printf("- 设计 #%d: ✅ 可能\n", designCount)
		} else {
			fmt.Printf("- 设计 #%d: ❌ 不可能\n", designCount)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("读取文件时发生错误: %v", err)
	}

	// --- 5. 打印最终结果 ---
	fmt.Println("\n--- 检查完成 ---")
	fmt.Printf("在 %d 个总设计中，有 %d 个是可行的。\n", designCount, possibleDesignsCount)
}
