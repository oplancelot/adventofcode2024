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

func main() {
	const inputFile = "input"

	rules, updates, err := readInput(inputFile)
	if err != nil {
		fmt.Printf("Failed to read input (%s): %v\n", inputFile, err)
		return
	}

	total := 0
	for _, update := range updates {
		if isValidUpdate(update, rules) {
			// 打印符合规则的 update
			fmt.Printf("Valid update: %v\n", update)

			// 取中间页码并累加
			if len(update) > 0 {
				mid := update[len(update)/2]
				total += mid
			}
		}
	}

	fmt.Printf("Sum of middle page numbers from valid updates: %d\n", total)
}
