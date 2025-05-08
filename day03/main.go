package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

// removeDoDontSections 删除 do() 和 don't() 之间的内容，以及从 don't() 到文件结束的内容
// 删除 don't() 到下一个 do() 之间的内容（包含 don't 和 do）
func removeDontDoSections(data string) string {
	re := regexp.MustCompile(`(?s)don't\(\).*?(do\(\)|\z)`)
	return re.ReplaceAllString(data, "")
}
func readInputAndSumMul(filename string) (int, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return 0, err
	}
	// 先删除 do() 和 don't() 之间的内容，以及从 don't() 到文件结束的内容
	cleanData := removeDontDoSections(string(data))
	re := regexp.MustCompile(`mul\((\d+),(\d+)\)`)
	matches := re.FindAllStringSubmatch(string(cleanData), -1)
	sum := 0
	for _, match := range matches {
		x, _ := strconv.Atoi(match[1])
		y, _ := strconv.Atoi(match[2])
		product := x * y
		sum += product
	}
	return sum, nil
}

func main() {
	const inputFile = "input"
	total, err := readInputAndSumMul(inputFile)
	if err != nil {
		fmt.Println("读取输入失败:", err)
		return
	}

	fmt.Printf("所有乘积之和为: %d\n", total)
}
