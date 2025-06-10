package main

import (
	"fmt"
	"os"
	"strings"
)

type Pos struct {
	X, Y int
}

func solve(input string) int {
	grid, start, end, trackTiles := parseInput(input)

	distFromS := bfs(start, grid)
	distToE := bfs(end, grid)

	baseTime, ok := distFromS[end]
	if !ok {
		return 0 // No path from S to E
	}

	cheats := make(map[[2]Pos]bool)

	for i := 0; i < len(trackTiles); i++ {
		for j := i + 1; j < len(trackTiles); j++ {
			pStart := trackTiles[i]
			pEnd := trackTiles[j]

			cheatDuration := abs(pStart.X-pEnd.X) + abs(pStart.Y-pEnd.Y)

			if cheatDuration <= 20 {
				timeWithCheat := distFromS[pStart] + cheatDuration + distToE[pEnd]
				if baseTime-timeWithCheat >= 100 {
					// Canonicalize the cheat to avoid duplicates (A,B) vs (B,A)
					// The loop structure (j=i+1) already prevents this.
					cheats[[2]Pos{pStart, pEnd}] = true
				}

				// Also check the reverse direction path
				timeWithCheatReverse := distFromS[pEnd] + cheatDuration + distToE[pStart]
				if baseTime-timeWithCheatReverse >= 100 {
					cheats[[2]Pos{pStart, pEnd}] = true
				}
			}
		}
	}

	return len(cheats)
}

func parseInput(input string) ([][]rune, Pos, Pos, []Pos) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	grid := make([][]rune, len(lines))
	var start, end Pos
	var trackTiles []Pos

	for y, line := range lines {
		grid[y] = []rune(line)
		for x, r := range line {
			if r == 'S' {
				start = Pos{x, y}
				grid[y][x] = '.' // Treat S as track
			} else if r == 'E' {
				end = Pos{x, y}
				grid[y][x] = '.' // Treat E as track
			}
		}
	}

	// After replacing S and E, find all track tiles
	for y, row := range grid {
		for x, r := range row {
			if r == '.' {
				trackTiles = append(trackTiles, Pos{x, y})
			}
		}
	}

	return grid, start, end, trackTiles
}

func bfs(start Pos, grid [][]rune) map[Pos]int {
	q := []Pos{start}
	dist := make(map[Pos]int)
	dist[start] = 0
	head := 0

	for head < len(q) {
		curr := q[head]
		head++

		for _, move := range []Pos{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
			next := Pos{curr.X + move.X, curr.Y + move.Y}
			if next.Y >= 0 && next.Y < len(grid) && next.X >= 0 && next.X < len(grid[0]) && grid[next.Y][next.X] == '.' {
				if _, visited := dist[next]; !visited {
					dist[next] = dist[curr] + 1
					q = append(q, next)
				}
			}
		}
	}
	return dist
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	data, err := os.ReadFile("input")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
		os.Exit(1)
	}
	result := solve(string(data))
	fmt.Println(result)
}
