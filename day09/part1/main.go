package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type File struct {
	ID   int
	Size int
	Pos  int // Starting position on disk
}

// readInput reads the disk map from input file
func readInput(filename string) ([]int, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Trim whitespace and get the single line of input
	line := strings.TrimSpace(string(data))

	// Convert each character to an integer
	var diskMap []int
	for _, char := range line {
		num, err := strconv.Atoi(string(char))
		if err != nil {
			return nil, fmt.Errorf("invalid character in input: %c", char)
		}
		diskMap = append(diskMap, num)
	}

	return diskMap, nil
}

// compactDisk simulates the compaction process and returns the final positions of each file
func compactDisk(diskMap []int) []int {
	// Create expanded disk representation
	var expandedDisk []int // -1 for free space, file ID for file blocks

	isFile := true
	fileID := 0

	// 构建初始磁盘状态
	for _, size := range diskMap {
		for i := 0; i < size; i++ {
			if isFile {
				expandedDisk = append(expandedDisk, fileID)
			} else {
				expandedDisk = append(expandedDisk, -1) // -1 represents free space
			}
		}
		if isFile {
			fileID++
		}
		isFile = !isFile
	}

	// 模拟碎片整理过程
	for i := 0; i < len(expandedDisk); i++ {
		if expandedDisk[i] == -1 { // Found free space
			// Find the rightmost file block
			rightmostFilePos := -1
			for j := len(expandedDisk) - 1; j > i; j-- {
				if expandedDisk[j] != -1 {
					rightmostFilePos = j
					break
				}
			}

			if rightmostFilePos != -1 {
				// Move the file block to the free space
				expandedDisk[i] = expandedDisk[rightmostFilePos]
				expandedDisk[rightmostFilePos] = -1
			} else {
				// No more file blocks to move
				break
			}
		}
	}

	return expandedDisk
}

// calculateChecksum computes the checksum based on file positions
func calculateChecksum(expandedDisk []int) int {
	checksum := 0

	// 计算校验和：每个文件块的(位置 * 文件ID)之和
	for pos, fileID := range expandedDisk {
		if fileID != -1 { // 跳过空闲空间
			checksum += pos * fileID
		}
	}

	return checksum
}

func main() {
	const inputFile = "input"
	diskMap, err := readInput(inputFile)
	if err != nil {
		fmt.Printf("读取输入文件失败 (%s): %v\n", inputFile, err)
		return
	}

	fmt.Printf("读取到磁盘映射，共 %d 个条目\n", len(diskMap))

	// 模拟碎片整理过程
	compactedDisk := compactDisk(diskMap)

	// Calculate the checksum
	checksum := calculateChecksum(compactedDisk)

	fmt.Printf("碎片整理后的文件系统校验和: %d\n", checksum)
}
