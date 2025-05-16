package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sort"
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

// compactDiskWholeFiles 模拟通过移动整个文件进行碎片整理的过程
func compactDiskWholeFiles(diskMap []int) []int {
	// 1. 构建初始磁盘状态
	var expandedDisk []int // -1表示空闲空间，非负整数表示文件ID
	filePositions := make(map[int]int) // 记录每个文件ID在磁盘上的起始位置
	fileSizes := make(map[int]int)     // 记录每个文件ID的大小
	
	// 填充扩展磁盘并记录文件信息
	isFile := true
	fileID := 0
	pos := 0
	
	for _, size := range diskMap {
		if isFile {
			// 记录文件信息
			filePositions[fileID] = pos
			fileSizes[fileID] = size
			
			// 填充文件块
			for j := 0; j < size; j++ {
				expandedDisk = append(expandedDisk, fileID)
			}
			fileID++
		} else {
			// 填充空闲空间
			for j := 0; j < size; j++ {
				expandedDisk = append(expandedDisk, -1)
			}
		}
		pos += size
		isFile = !isFile
	}
	
	// 2. 获取所有文件ID并按降序排序
	var fileIDs []int
	for id := range fileSizes {
		fileIDs = append(fileIDs, id)
	}
	sort.Slice(fileIDs, func(i, j int) bool {
		return fileIDs[i] > fileIDs[j] // 降序排序
	})
	
	// 3. 按ID从大到小尝试移动每个文件
	for _, id := range fileIDs {
		fileSize := fileSizes[id]
		currentPos := -1
		
		// 找到文件当前位置
		for i := 0; i < len(expandedDisk); i++ {
			if expandedDisk[i] == id {
				currentPos = i
				break
			}
		}
		
		if currentPos == -1 {
			continue // 找不到文件，跳过
		}
		
		// 寻找左侧最近的足够大的连续空闲空间
		bestPos := -1
		for i := 0; i < currentPos; i++ {
			if expandedDisk[i] == -1 { // 找到空闲空间起始点
				// 检查是否有足够的连续空闲空间
				j := i
				freeCount := 0
				
				while: for ; j < currentPos && expandedDisk[j] == -1; j++ {
					freeCount++
					if freeCount >= fileSize {
						bestPos = i
						break while
					}
				}
				
				// 如果找到了足够的空间，就不再继续寻找
				if bestPos != -1 {
					break
				}
				
				// 跳过已检查的空闲空间
				i = j - 1
			}
		}
		
		// 如果找到合适的空闲空间，移动整个文件
		if bestPos != -1 {
			// 复制文件块到新位置
			for i := 0; i < fileSize; i++ {
				expandedDisk[bestPos+i] = id
			}
			
			// 清除原位置
			for i := 0; i < fileSize; i++ {
				expandedDisk[currentPos+i] = -1
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

	// 模拟整个文件移动的碎片整理过程
	compactedDisk := compactDiskWholeFiles(diskMap)

	// 计算校验和
	checksum := calculateChecksum(compactedDisk)

	fmt.Printf("碎片整理后的文件系统校验和: %d\n", checksum)
}
