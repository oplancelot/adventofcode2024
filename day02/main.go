package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type LineEvaluation int

const (
	EvalSafe LineEvaluation = iota
	EvalUnsafe
)

func (e LineEvaluation) String() string {
	switch e {
	case EvalSafe:
		return "safe"
	case EvalUnsafe:
		return "unsafe"
	default:
		return "unknown"
	}
}

// readInput 读取并分析输入文件，返回 safe 和 unsafe 行的数量
func readInput(filename string) (int, int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	var safeCount, unsafeCount int
	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		result, err := parseLine(scanner.Text())
		if err != nil {
			fmt.Printf("line %d: %v\n", lineNum, err)
			continue
		}
		if result == EvalSafe {
			safeCount++
		} else {
			unsafeCount++
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, 0, err
	}
	return safeCount, unsafeCount, nil
}

// parseLine 解析并评估一行数据
func parseLine(line string) (LineEvaluation, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return EvalUnsafe, fmt.Errorf("not enough numbers to evaluate")
	}

	nums := make([]int, len(fields))
	for i, f := range fields {
		n, err := strconv.Atoi(f)
		if err != nil {
			return EvalUnsafe, fmt.Errorf("invalid number: %s", f)
		}
		nums[i] = n
	}

	if isSafe(nums) {
		return EvalSafe, nil
	}

	// 尝试删除一个数字来看看是否变为安全
	for i := 0; i < len(nums); i++ {
		tmp := append([]int{}, nums[:i]...)
		tmp = append(tmp, nums[i+1:]...)
		if len(tmp) >= 2 && isSafe(tmp) {
			return EvalSafe, nil
		}
	}

	return EvalUnsafe, nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// isSafe 判断一个数列是否符合安全条件
func isSafe(nums []int) bool {
	if len(nums) < 2 {
		return false
	}

	direction := 0 // 0 = 未确定, 1 = 递增, -1 = 递减

	for i := 1; i < len(nums); i++ {
		diff := nums[i] - nums[i-1]
		if diff == 0 || abs(diff) > 3 {
			return false
		}
		if direction == 0 {
			if diff > 0 {
				direction = 1
			} else {
				direction = -1
			}
		} else {
			if (direction == 1 && diff < 0) || (direction == -1 && diff > 0) {
				return false
			}
		}
	}
	return true
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
