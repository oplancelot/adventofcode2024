package main

import (
	"fmt"
	"os"
	"strings"
)

// Direction types, defined using iota. (使用 iota 定义的方向类型。)
const (
	DirectionRight = iota // Right (向右)
	DirectionLeft         // Left (向左)
	DirectionUp           // Up (向上)
	DirectionDown         // Down (向下)
	DirectionUpRight      // Up-Right (向右上)
	DirectionUpLeft       // Up-Left (向左上)
	DirectionDownRight    // Down-Right (向右下)
	DirectionDownLeft     // Down-Left (向左下)
)

// Defines the (dx, dy) offsets for the 8 directions. (定义了8个方向的 (dx, dy) 偏移量。)
// dx: change in column, dy: change in row. (dx: 列变化, dy: 行变化。)
var directions = [][2]int{
	{1, 0},   // 向右 (Right)
	{-1, 0},  // 向左 (Left)
	{0, -1},  // 向上 (Up)
	{0, 1},   // 向下 (Down)
	{1, -1},  // 向右上 (Up-Right)
	{-1, -1}, // 向左上 (Up-Left)
	{1, 1},   // 向右下 (Down-Right)
	{-1, 1},  // 向左下 (Down-Left)
}

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

// findCharPositions finds all coordinates of a given character in the grid.
// 查找给定字符在网格中的所有坐标。
func findCharPositions(grid [][]string, char string) [][2]int {
	var positions [][2]int
	for r, row := range grid {
		for c := range row {
			// Ensure row is not empty before accessing grid[r][c]
			// Although current input structure implies non-empty rows of same length.
			// For robustness, a check like `if len(row) > c && grid[r][c] == char` could be added if row lengths vary.
			if len(row) > 0 && grid[r][c] == char { // Assuming c will be within bounds if row is not empty.
				positions = append(positions, [2]int{r, c})
			}
		}
	}
	return positions
}

// searchWordRecursive recursively searches for the word starting from a given character.
// (r, c) are the coordinates of the previously matched character (word[idx-1]).
// It attempts to match word[idx] at the next position in the direction (dx, dy).
// 从给定字符开始递归地搜索单词。
// (r, c) 是前一个匹配字符 (word[idx-1]) 的坐标。
// 它尝试在 (dx, dy) 方向上的下一个位置匹配 word[idx]。
func searchWordRecursive(grid [][]string, r, c int, word string, idx int, dx, dy int) bool {
	// Base case: If idx reaches the length of the word,
	// it means all characters word[0]...word[len(word)-1] have been matched.
	if idx == len(word) {
		// Successfully matched all characters from word[0] to word[len(word)-1]
		return true
	}

	// Calculate new coordinates (计算新的坐标)
	newR := r + dy
	newC := c + dx

	// Check if out of bounds (检查是否越界)
	if newR < 0 || newR >= len(grid) || newC < 0 || (len(grid[newR]) == 0 || newC >= len(grid[newR])) {
		return false
	}

	// Check if the current character matches (检查当前字符是否匹配)
	if grid[newR][newC] != string(word[idx]) {
		return false
	}

	// Recursively search for the next character (递归查找下一个字符)
	return searchWordRecursive(grid, newR, newC, word, idx+1, dx, dy)
}

// findWordInGrid finds the number of occurrences of the word in the grid.
// 查找单词在网格中出现的次数。
func findWordInGrid(grid [][]string, word string) int {
	if len(word) == 0 {
		return 0
	}
	if len(grid) == 0 || len(grid[0]) == 0 { // Basic check for empty grid
		return 0
	}

	count := 0
	charPositions := findCharPositions(grid, string(word[0])) // Find all positions of the first letter (查找第一个字母的所有位置)

	// For each starting position, try to find the complete word in all directions.
	// 对于每个起始位置，尝试各个方向查找完整单词。
	for _, pos := range charPositions {
		r, c := pos[0], pos[1]
		// (r,c) is the position of word[0].
		for _, dir := range directions {
			// Start recursive search for word[1] from (r,c) in direction dir.
			// searchWordRecursive will try to match word[idx=1] at (r+dy, c+dx).
			if searchWordRecursive(grid, r, c, word, 1, dir[0], dir[1]) {
				count++
			}
		}
	}

	return count
}

func main() {
	const inputFile = "input"
	grid, err := readInput(inputFile)
	if err != nil {
		fmt.Println("Failed to read input (读取输入失败):", err)
		return
	}

	word := "XMAS"
	total := findWordInGrid(grid, word)
	fmt.Printf("“%s”出现的次数为: %d\n", word, total)
	// Note: For a single-character word, this logic would count it 8 times for each occurrence.
}
