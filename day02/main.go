package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// readInput 读取并解析输入文件，返回两个整数切片
func readInput(filename string) (int, int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	var safe, unsafe int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result, err := parseLine(scanner.Text())
		if err != nil {
			fmt.Println(err)
			continue
		}
		if result == "safe" {
			safe = safe + 1
		} else {
			unsafe = unsafe + 1
		}

	}
	if err := scanner.Err(); err != nil {
		return 0, 0, err
	}
	return safe, unsafe, nil
}

// parseLine 解析每一行，返回左右两个整数
func parseLine(line string) (string, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return "", fmt.Errorf("not enough numbers to evaluate")
	}

	nums := make([]int, len(fields))
	for i, f := range fields {
		n, err := strconv.Atoi(f)
		if err != nil {
			return "", fmt.Errorf("invalid number: %s", f)
		}
		nums[i] = n
	}

	direction := 0 // 0 = unknown, 1 = increasing, -1 = decreasing

	for i := 1; i < len(nums); i++ {
		diff := nums[i] - nums[i-1]

		// 差值必须在 1 到 3 或 -1 到 -3 之间
		if diff == 0 || abs(diff) > 3 {
			return "unsafe", nil
		}

		// 初始化方向
		if direction == 0 {
			if diff > 0 {
				direction = 1
			} else {
				direction = -1
			}
		} else {
			// 检查方向一致性
			if (direction == 1 && diff < 0) || (direction == -1 && diff > 0) {
				return "unsafe", nil
			}
		}
	}

	return "safe", nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	const inputFile = "./input"
	safe, unsafe, err := readInput(inputFile)
	if err != nil {
		fmt.Println("读取输入失败:", err)
		return
	}
	fmt.Println("safe:", safe, "unsafe:", unsafe)

}
