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

// FileRange represents a range of blocks belonging to a file
type FileRange struct {
	ID       int
	StartPos int
	Size     int
}

// FreeRange represents a range of free space blocks
type FreeRange struct {
	StartPos int
	Size     int
}

// readInput reads the disk map from input file
func readInput(filename string) ([]int, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	// Trim whitespace and get the single line of input
	line := strings.TrimSpace(string(data))

	// Validate input is not empty
	if len(line) == 0 {
		return nil, fmt.Errorf("empty input file: %s", filename)
	}

	// Convert each character to an integer
	var diskMap []int
	for i, char := range line {
		num, err := strconv.Atoi(string(char))
		if err != nil {
			return nil, fmt.Errorf("invalid character at position %d: %c", i, char)
		}
		diskMap = append(diskMap, num)
	}

	return diskMap, nil
}

// createExpandedDisk creates an expanded representation of the disk
// where each element represents a single block
func createExpandedDisk(diskMap []int) []int {
	// Estimate total disk size to pre-allocate memory
	totalSize := 0
	for _, size := range diskMap {
		totalSize += size
	}

	expandedDisk := make([]int, 0, totalSize) // Pre-allocate memory

	isFile := true
	fileID := 0

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

	return expandedDisk
}

// compactExpandedDisk performs the compaction process on the expanded disk
// by moving file blocks from right to left
func compactExpandedDisk(expandedDisk []int) []int {
	// Create a copy of the expanded disk to avoid modifying the original
	compactedDisk := make([]int, len(expandedDisk))
	copy(compactedDisk, expandedDisk)

	// Track free space positions for optimization
	var freeSpaces []int
	for i, block := range compactedDisk {
		if block == -1 {
			freeSpaces = append(freeSpaces, i)
		}
	}

	// Process each free space from left to right
	for _, freePos := range freeSpaces {
		// Skip if this position is no longer free (already filled by a previous move)
		if compactedDisk[freePos] != -1 {
			continue
		}

		// Find the rightmost file block
		rightmostFilePos := -1
		for j := len(compactedDisk) - 1; j > freePos; j-- {
			if compactedDisk[j] != -1 {
				rightmostFilePos = j
				break
			}
		}

		if rightmostFilePos != -1 {
			// Move the file block to the free space
			compactedDisk[freePos] = compactedDisk[rightmostFilePos]
			compactedDisk[rightmostFilePos] = -1
		} else {
			// No more file blocks to move
			break
		}
	}

	return compactedDisk
}

// calculateChecksum computes the checksum based on file positions
func calculateChecksum(expandedDisk []int) int {
	checksum := 0

	// Calculate checksum: sum of (position * fileID) for each block
	for pos, fileID := range expandedDisk {
		if fileID != -1 { // Skip free space
			checksum += pos * fileID
		}
	}

	return checksum
}

// compactDisk is a wrapper function that handles the entire compaction process
func compactDisk(diskMap []int) ([]int, error) {
	if len(diskMap) == 0 {
		return nil, fmt.Errorf("empty disk map")
	}

	expandedDisk := createExpandedDisk(diskMap)
	compactedDisk := compactExpandedDisk(expandedDisk)

	return compactedDisk, nil
}

func main() {
	// Parse command line arguments
	const inputFile = "input"

	// Read and process the input
	diskMap, err := readInput(inputFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Read disk map with %d entries\n", len(diskMap))

	// Compact the disk
	compactedDisk, err := compactDisk(diskMap)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Calculate the checksum
	checksum := calculateChecksum(compactedDisk)

	fmt.Printf("Filesystem checksum after compaction: %d\n", checksum)
}
