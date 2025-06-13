package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// solve aoc 2025 day 25
func solve(input string) int {
	blocks := strings.Split(strings.TrimSpace(input), "\n\n")

	var locks [][]int
	var keys [][]int
	const totalHeight = 5

	for _, block := range blocks {
		lines := strings.Split(block, "\n")
		if len(lines) < 2 {
			continue
		}

		firstLine := lines[0]
		lastLine := lines[len(lines)-1]

		isLock := !strings.Contains(firstLine, ".") && !strings.Contains(lastLine, "#")
		isKey := !strings.Contains(firstLine, "#") && !strings.Contains(lastLine, ".")

		if !isLock && !isKey {
			continue
		}

		// Calculate heights for the schematic (applies to both locks and keys)
		// The height is the number of '#' in the middle 5 rows for each column.
		width := len(lines[0])
		heights := make([]int, width)
		// We only care about the middle rows (1 to 5, since row 0 and 6 are borders)
		for col := 0; col < width; col++ {
			pinHeight := 0
			for row := 1; row <= totalHeight; row++ {
				if col < len(lines[row]) && lines[row][col] == '#' {
					pinHeight++
				}
			}
			heights[col] = pinHeight
		}

		if isLock {
			locks = append(locks, heights)
		} else if isKey {
			keys = append(keys, heights)
		}
	}

	// Count fitting pairs
	fitCount := 0
	for _, lock := range locks {
		for _, key := range keys {
			if checkFit(lock, key, totalHeight) {
				fitCount++
			}
		}
	}

	return fitCount
}

// checkFit determines if a key fits into a lock.
func checkFit(lock, key []int, totalHeight int) bool {
	if len(lock) != len(key) {
		return false // Should not happen with valid input
	}
	for i := 0; i < len(lock); i++ {
		if lock[i]+key[i] > totalHeight {
			return false // Overlap
		}
	}
	return true // It fits
}

func main() {
	data, err := os.ReadFile("input")
	if err != nil {
		log.Fatal("failed to read input file:", err)
	}

	result := solve(string(data))
	fmt.Println("Total fitting lock/key pairs:", result)
}
