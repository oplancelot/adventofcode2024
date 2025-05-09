package main

import (
	"fmt"
	"os"
	"strings"
)

// readInput reads the input file and builds the character grid.
// 读取输入文件并构建字符网格。
func readInput(filename string) ([][]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	grid := make([][]string, len(lines))

	for i, line := range lines {
		line = strings.TrimSuffix(line, "\r") // Handle potential CR characters from Windows-style line endings
		grid[i] = strings.Split(line, "")
	}

	return grid, nil
}

// findXMASPatterns counts the number of X-MAS patterns in the grid.
// An X-MAS pattern is two "MAS" strings (forwards or backwards) intersecting at 'A'
// in an X shape. Example: M.S
//                          .A.
//                          M.S
func findXMASPatterns(grid [][]string) int {
	count := 0
	rows := len(grid)
	if rows < 3 { // Need at least 3 rows for a pattern
		return 0
	}

	cols := 0
	if len(grid[0]) > 0 { // Assuming a rectangular grid, get column count from the first row
		cols = len(grid[0])
	}
	if cols < 3 { // Need at least 3 columns for a pattern
		return 0
	}

	// Iterate through each cell that could be the center 'A' of an X-MAS pattern.
	// This means r and c must allow for r-1, r+1, c-1, c+1 accesses.
	for r := 1; r < rows-1; r++ {
		for c := 1; c < cols-1; c++ {
			if grid[r][c] == "A" {
				// Potential center 'A' found.
				// Get diagonal characters. Loop bounds ensure these are valid for rectangular grid.
				tl := grid[r-1][c-1] // Top-Left
				tr := grid[r-1][c+1] // Top-Right
				bl := grid[r+1][c-1] // Bottom-Left
				br := grid[r+1][c+1] // Bottom-Right

				// Check first diagonal (Top-Left to Bottom-Right) for "MAS" or "SAM"
				diag1Valid := (tl == "M" && br == "S") || (tl == "S" && br == "M")

				// Check second diagonal (Top-Right to Bottom-Left) for "MAS" or "SAM"
				diag2Valid := (tr == "M" && bl == "S") || (tr == "S" && bl == "M")

				if diag1Valid && diag2Valid {
					count++
				}
			}
		}
	}
	return count
}

func main() {
	const inputFile = "/home/ubuntu/adventofcode/day04/input"
	grid, err := readInput(inputFile)
	if err != nil {
		fmt.Printf("Failed to read input (%s): %v\n", inputFile, err)
		return
	}

	// --- Part Two: Find X-MAS patterns ---
	totalXMASPatterns := findXMASPatterns(grid)
	fmt.Printf("Number of X-MAS patterns found: %d\n", totalXMASPatterns)
}
