package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

// parseLine 解析每一行，返回左右两个整数
func parseLine(line string) (int, int, error) {
	parts := strings.Fields(strings.TrimSpace(line))
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("无效行（不是两个数字）: %s", line)
	}
	left, err1 := strconv.Atoi(parts[0])
	right, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return 0, 0, fmt.Errorf("数字解析失败: %s", line)
	}
	return left, right, nil
}

// readInput 读取并解析输入文件，返回两个整数切片
func readInput(filename string) ([]int, []int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	var lefts, rights []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		left, right, err := parseLine(scanner.Text())
		if err != nil {
			fmt.Println(err)
			continue
		}
		lefts = append(lefts, left)
		rights = append(rights, right)
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}
	return lefts, rights, nil
}



// combineAndPrint 输出排序后的左右列组合，并返回差值总和（绝对值）
func combineAndPrint(lefts, rights []int) int {
	sort.Ints(lefts)
	sort.Ints(rights)

	fmt.Println("左右列排序后重新组合：")
	total := 0
	for i := 0; i < len(lefts) && i < len(rights); i++ {
		fmt.Printf("%d %d\n", lefts[i], rights[i])
		diff := int(math.Abs(float64(rights[i] - lefts[i])))
		total += diff
	}
	return total
}


func main() {
	const inputFile = "input"

	lefts, rights, err := readInput(inputFile)
	if err != nil {
		fmt.Println("读取输入失败:", err)
		return
	}

	if len(lefts) != len(rights) {
		fmt.Println("左右列数量不一致，数据可能有误")
		return
	}

	total := combineAndPrint(lefts, rights)
	fmt.Printf("\n最终差值总和: %d\n", total)
}
