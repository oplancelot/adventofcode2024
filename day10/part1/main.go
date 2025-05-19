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

// findTrailheadScore calculates the score for a single trailhead
func findTrailheadScore(grid [][]int, start Position) int {
	// Set to track unique 9s we've reached
	reachableNines := make(map[Position]bool)
	
	// Use BFS to find all paths
	type QueueItem struct {
		pos Position
		height int
	}
	
	queue := []QueueItem{{pos: start, height: 0}}
	visited := make(map[string]bool)
	
	// Mark the starting position as visited
	key := fmt.Sprintf("%d,%d,%d", start.row, start.col, 0)
	visited[key] = true
	
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		
		// If we've reached a 9, mark it as reachable
		if grid[current.pos.row][current.pos.col] == 9 {
			reachableNines[current.pos] = true
			continue
		}
		
		// Try all four directions
		for _, dir := range directions {
			next := Position{current.pos.row + dir.row, current.pos.col + dir.col}
			
			if isWithinBounds(grid, next) {
				nextHeight := current.height + 1
				
				// We can only move to positions with height exactly one more than current
				if grid[next.row][next.col] == grid[current.pos.row][current.pos.col] + 1 {
					// Create a unique key for this state
					key := fmt.Sprintf("%d,%d,%d", next.row, next.col, nextHeight)
					
					if !visited[key] {
						visited[key] = true
						queue = append(queue, QueueItem{
							pos:    next,
							height: nextHeight,
						})
					}
				}
			}
		}
	}
	
	// Return the number of unique 9s reached
	return len(reachableNines)
}

// calculateTotalScore calculates the sum of scores for all trailheads
func calculateTotalScore(grid [][]int) int {
	totalScore := 0
	
	// Find all trailheads (positions with height 0)
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			if grid[i][j] == 0 {
				score := findTrailheadScore(grid, Position{i, j})
				fmt.Printf("Trailhead at (%d,%d) has score: %d\n", i, j, score)
				totalScore += score
			}
		}
	}
	
	return totalScore
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

	totalScore := calculateTotalScore(grid)
	fmt.Printf("Sum of all trailhead scores: %d\n", totalScore)
}
