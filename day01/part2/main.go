// part2/main.go
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// parseInput 将输入数据解析为左右两个整数切片
func parseInput(data []byte) ([]int, []int, error) {
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	var leftList, rightList []int

	for scanner.Scan() {
		line := scanner.Text()
		// strings.Fields可以处理一个或多个空格/制表符分隔的情况
		parts := strings.Fields(line)
		if len(parts) != 2 {
			// 可以选择忽略格式错误的行或返回错误
			continue
		}

		leftNum, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, nil, fmt.Errorf("无法解析左侧数字 '%s': %w", parts[0], err)
		}
		rightNum, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, nil, fmt.Errorf("无法解析右侧数字 '%s': %w", parts[1], err)
		}

		leftList = append(leftList, leftNum)
		rightList = append(rightList, rightNum)
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("扫描输入时出错: %w", err)
	}

	return leftList, rightList, nil
}

// solvePart2 计算相似度分数
func solvePart2(leftList, rightList []int) int {
	// 为右侧列表创建一个频率映射
	rightCounts := make(map[int]int)
	for _, num := range rightList {
		rightCounts[num]++
	}

	var totalScore int
	// 遍历左侧列表并计算分数
	for _, num := range leftList {
		// 如果数字不在映射中，其计数值将为0
		count := rightCounts[num]
		totalScore += num * count
	}

	return totalScore
}

func main() {
	// 从 "input" 文件读取数据
	file, err := os.ReadFile("input")
	if err != nil {
		log.Fatalf("无法读取输入文件: %s", err)
	}

	leftList, rightList, err := parseInput(file)
	if err != nil {
		log.Fatalf("解析输入失败: %s", err)
	}

	result := solvePart2(leftList, rightList)
	fmt.Printf("相似度分数为: %d\n", result)
}
