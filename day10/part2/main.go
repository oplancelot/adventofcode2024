package main

import (
	"fmt"
	"os"
	"strings"
)

// Position represents a location in the grid
type Position struct {
	row, col int
}

// Direction represents possible movement directions
var directions = []Position{
	{-1, 0}, // Up
	{1, 0},  // Down
	{0, -1}, // Left
	{0, 1},  // Right
}

func readInput(filename string) ([][]int, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	grid := make([][]int, len(lines))

	for i, line := range lines {
		grid[i] = make([]int, len(line))
		for j, char := range line {
			grid[i][j] = int(char - '0')
		}
	}

	return grid, nil
}

// isWithinBounds checks if a position is within the grid boundaries
func isWithinBounds(grid [][]int, pos Position) bool {
	return pos.row >= 0 && pos.row < len(grid) && pos.col >= 0 && pos.col < len(grid[0])
}

// countDistinctTrails calculates the number of distinct hiking trails from a trailhead
func countDistinctTrails(grid [][]int, start Position) int {
    // Create a memoization cache
    memo := make(map[string]int)
    
    // Define a recursive function to count paths
    var countPaths func(pos Position, height int) int
    countPaths = func(pos Position, height int) int {
        // If we've reached height 9, we've found a complete trail
        if height == 9 {
            return 1
        }
        
        // Create a key for memoization
        key := fmt.Sprintf("%d,%d,%d", pos.row, pos.col, height)
        
        // Check if we've already computed this
        if count, exists := memo[key]; exists {
            return count
        }
        
        // Count paths from this position
        count := 0
        
        // Try all four directions
        for _, dir := range directions {
            next := Position{pos.row + dir.row, pos.col + dir.col}
            
            if isWithinBounds(grid, next) {
                // We can only move to positions with height exactly one more than current
                if grid[next.row][next.col] == height + 1 {
                    count += countPaths(next, height + 1)
                }
            }
        }
        
        // Store result in memo
        memo[key] = count
        return count
    }
    
    // Start counting from the trailhead
    return countPaths(start, grid[start.row][start.col])
}

// calculateTotalRating calculates the sum of ratings for all trailheads
func calculateTotalRating(grid [][]int) int {
	totalRating := 0
	
	// Find all trailheads (positions with height 0)
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			if grid[i][j] == 0 {
				rating := countDistinctTrails(grid, Position{i, j})
				fmt.Printf("Trailhead at (%d,%d) has rating: %d\n", i, j, rating)
				totalRating += rating
			}
		}
	}
	
	return totalRating
}

func main() {
	const inputFile = "input"
	grid, err := readInput(inputFile)
	if err != nil {
		fmt.Printf("Failed to read input file (%s): %v\n", inputFile, err)
		return
	}

	// Count trailheads
	trailheadCount := 0
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			if grid[i][j] == 0 {
				trailheadCount++
			}
		}
	}
	fmt.Printf("Found %d trailheads\n", trailheadCount)

	totalRating := calculateTotalRating(grid)
	fmt.Printf("Sum of all trailhead ratings: %d\n", totalRating)
}
