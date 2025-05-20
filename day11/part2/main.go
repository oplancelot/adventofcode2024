package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Bag 统计每种石头值出现的次数
type Bag map[int]int

// countDigits 返回一个非负整数的位数
func countDigits(n int) int {
	if n == 0 {
		return 1
	}
	count := 0
	for n > 0 {
		n /= 10
		count++
	}
	return count
}

// 缓存规则结果以避免重复转换
var ruleCache = make(map[int][]int)

// applyRules 应用规则并缓存转换结果
func applyRules(stone int) []int {
	if cached, ok := ruleCache[stone]; ok {
		return cached
	}
	var result []int
	if stone == 0 {
		result = []int{1}
	} else {
		numDigits := countDigits(stone)
		if numDigits%2 == 0 {
			divisor := 1
			for i := 0; i < numDigits/2; i++ {
				divisor *= 10
			}
			left := stone / divisor
			right := stone % divisor
			result = []int{left, right}
		} else {
			result = []int{stone * 2024}
		}
	}
	ruleCache[stone] = result
	return result
}

// simulateBlinkMap 模拟一次 blink，返回新的 Bag，支持 trace 打印
func simulateBlinkMap(bag Bag, trace bool) Bag {
	newBag := make(Bag)
	for stone, count := range bag {
		newStones := applyRules(stone)
		if trace {
			fmt.Printf("  %d x %d → ", stone, count)
			for _, s := range newStones {
				fmt.Printf("%d x %d ", s, count)
			}
			fmt.Println()
		}
		for _, s := range newStones {
			newBag[s] += count
		}
	}
	return newBag
}

// readInput 读取输入文件并返回初始石头的 Bag
func readInput(filename string) (Bag, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	stoneStrs := strings.Fields(strings.TrimSpace(string(data)))
	bag := make(Bag)
	for _, stoneStr := range stoneStrs {
		stone, err := strconv.Atoi(stoneStr)
		if err != nil {
			return nil, fmt.Errorf("invalid stone value: %s", stoneStr)
		}
		bag[stone]++
	}
	return bag, nil
}

// totalCount 返回 Bag 中石头总数量
func totalCount(bag Bag) int {
	sum := 0
	for _, count := range bag {
		sum += count
	}
	return sum
}

func main() {
	const inputFile = "input"
	const maxTraceRounds = 14

	bag, err := readInput(inputFile)
	if err != nil {
		fmt.Printf("Failed to read input: %v\n", err)
		return
	}

	fmt.Printf("Initial count: %d stones\n", totalCount(bag))

	for i := 1; i <= 75; i++ {
		trace := i <= maxTraceRounds
		if trace {
			fmt.Printf("\n[Trace] Blink #%d\n", i)
		}
		bag = simulateBlinkMap(bag, trace)

		if i%5 == 0 || trace {
			fmt.Printf("After %d blinks: %d stones\n", i, totalCount(bag))
		}
	}

	fmt.Printf("Final stone count: %d\n", totalCount(bag))
}
