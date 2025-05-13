package main

import (
	"fmt"
	"os"
	"strings"
)

// Direction represents the four cardinal directions
type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

// Position represents a coordinate on the grid
type Position struct {
	row, col int
}

// readInput reads the input file and builds the character grid.
func readInput(filename string) ([][]string, Position, Direction, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, Position{}, 0, err
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	grid := make([][]string, len(lines))
	
	var guardPos Position
	var guardDir Direction

	for i, line := range lines {
		line = strings.TrimSuffix(line, "\r") // Handle potential CR characters from Windows-style line endings
		grid[i] = strings.Split(line, "")
		
		// Find the guard's starting position and direction
		for j, char := range grid[i] {
			switch char {
			case "^":
				guardPos = Position{i, j}
				guardDir = Up
				grid[i][j] = "." // Replace with empty space for tracking
			case ">":
				guardPos = Position{i, j}
				guardDir = Right
				grid[i][j] = "."
			case "v":
				guardPos = Position{i, j}
				guardDir = Down
				grid[i][j] = "."
			case "<":
				guardPos = Position{i, j}
				guardDir = Left
				grid[i][j] = "."
			}
		}
	}

	return grid, guardPos, guardDir, nil
}

// isInBounds checks if a position is within the grid boundaries
func isInBounds(grid [][]string, pos Position) bool {
	return pos.row >= 0 && pos.row < len(grid) && pos.col >= 0 && pos.col < len(grid[0])
}

// getNextPosition returns the position in front of the current position based on direction
func getNextPosition(pos Position, dir Direction) Position {
	switch dir {
	case Up:
		return Position{pos.row - 1, pos.col}
	case Right:
		return Position{pos.row, pos.col + 1}
	case Down:
		return Position{pos.row + 1, pos.col}
	case Left:
		return Position{pos.row, pos.col - 1}
	}
	return pos // Should never happen
}

// turnRight returns the direction after turning right 90 degrees
func turnRight(dir Direction) Direction {
	return (dir + 1) % 4
}

// findDistinctPositions counts the number of distinct positions visited by the guard
func findDistinctPositions(grid [][]string, startPos Position, startDir Direction) int {
	// Create a map to track visited positions
	visited := make(map[Position]bool)
	
	// Start with the guard's initial position
	pos := startPos
	dir := startDir
	
	// Mark the starting position as visited
	visited[pos] = true
	
	// Continue until the guard leaves the mapped area
	for {
		// Check what's in front
		nextPos := getNextPosition(pos, dir)
		
		// If out of bounds, the guard has left the area
		if !isInBounds(grid, nextPos) {
			break
		}
		
		// If there's an obstacle in front, turn right
		if grid[nextPos.row][nextPos.col] == "#" {
			dir = turnRight(dir)
		} else {
			// Otherwise, move forward
			pos = nextPos
			visited[pos] = true
		}
	}
	
	return len(visited)
}

func main() {
	const inputFile = "input"
	grid, guardPos, guardDir, err := readInput(inputFile)
	if err != nil {
		fmt.Printf("Failed to read input (%s): %v\n", inputFile, err)
		return
	}

	distinctPositions := findDistinctPositions(grid, guardPos, guardDir)
	fmt.Printf("Number of distinct positions visited by the guard: %d\n", distinctPositions)
}
