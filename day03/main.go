package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func readInputAndSumMul(filename string) (int, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return 0, err
	}

	re := regexp.MustCompile(`mul\((\d+),(\d+)\)`)
	matches := re.FindAllStringSubmatch(string(data), -1)
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
