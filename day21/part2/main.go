package main

import (
	"bufio"
	"container/list"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// Keypad layout definitions
var numKeypad = [][]rune{
	{'7', '8', '9'},
	{'4', '5', '6'},
	{'1', '2', '3'},
	{' ', '0', 'A'},
}

var dirKeypad = [][]rune{
	{' ', '^', 'A'},
	{'<', 'v', '>'},
}

// computeSequences uses BFS to find all shortest paths between any two keys on a given keypad.
// This new version removes the `visited` map and uses a more aggressive pruning strategy
// to more closely match the reference Rust implementation.
func computeSequences(keypad [][]rune) map[[2]rune][]string {
	pos := make(map[rune]struct{ r, c int })
	for r, row := range keypad {
		for c, b := range row {
			if b != ' ' {
				pos[b] = struct{ r, c int }{r, c}
			}
		}
	}

	sequences := make(map[[2]rune][]string)
	for fromKey := range pos {
		for toKey := range pos {
			if fromKey == toKey {
				sequences[[2]rune{fromKey, toKey}] = []string{"A"}
				continue
			}

			type bfsState struct {
				pos  struct{ r, c int }
				path string
			}

			q := list.New()
			q.PushBack(bfsState{pos: pos[fromKey], path: ""})

			shortestPaths := []string{}
			minLen := math.MaxInt32 // Use a large integer to represent infinity

			for q.Len() > 0 {
				elem := q.Front()
				q.Remove(elem)
				curr := elem.Value.(bfsState)

				// Pruning: if the current path is already long enough that its neighbors
				// would exceed the shortest path found so far, we can stop.
				if len(curr.path)+1 > minLen {
					continue
				}

				dr := []int{-1, 1, 0, 0}
				dc := []int{0, 0, -1, 1}
				moveChars := []rune{'^', 'v', '<', '>'}

				for i := 0; i < 4; i++ {
					nr, nc := curr.pos.r+dr[i], curr.pos.c+dc[i]

					if nr >= 0 && nr < len(keypad) && nc >= 0 && nc < len(keypad[0]) && keypad[nr][nc] != ' ' {
						nextPos := struct{ r, c int }{nr, nc}
						nextPath := curr.path + string(moveChars[i])

						if keypad[nr][nc] == toKey {
							fullPath := nextPath + "A"
							if len(fullPath) < minLen {
								minLen = len(fullPath)
								shortestPaths = []string{fullPath}
							} else if len(fullPath) == minLen {
								shortestPaths = append(shortestPaths, fullPath)
							}
						} else {
							// Only add to queue if this path isn't doomed to be longer than the current shortest.
							if len(nextPath)+1 < minLen {
								q.PushBack(bfsState{pos: nextPos, path: nextPath})
							}
						}
					}
				}
			}
			sequences[[2]rune{fromKey, toKey}] = shortestPaths
		}
	}
	return sequences
}

// cartesianProduct is a helper to generate all combinations of sequences.
func cartesianProduct(sets [][]string) []string {
	if len(sets) == 0 {
		return []string{""}
	}
	result := []string{}
	subProduct := cartesianProduct(sets[1:])
	for _, item := range sets[0] {
		for _, sub := range subProduct {
			result = append(result, item+sub)
		}
	}
	return result
}

// generateAllNumSequences creates all possible D-Pad sequences for a given numeric code.
func generateAllNumSequences(code string, sequences map[[2]rune][]string) []string {
	var options [][]string
	prevChar := 'A'
	for _, char := range code {
		options = append(options, sequences[[2]rune{prevChar, char}])
		prevChar = char
	}
	return cartesianProduct(options)
}

// Cache for the recursive computeMinLength function.
var computeCache = make(map[string]uint64)

// computeMinLength recursively calculates the minimum cost to type a sequence at a given depth.
func computeMinLength(seq string, depth int, dirSequences map[[2]rune][]string, dirLengths map[[2]rune]int) uint64 {
	if depth == 1 {
		var length uint64
		prevChar := 'A'
		for _, char := range seq {
			length += uint64(dirLengths[[2]rune{prevChar, char}])
			prevChar = char
		}
		return length
	}

	cacheKey := fmt.Sprintf("%s|%d", seq, depth)
	if length, ok := computeCache[cacheKey]; ok {
		return length
	}

	var totalMinLength uint64
	prevChar := 'A'
	for _, char := range seq {
		minSubLength := uint64(math.MaxUint64)
		options := dirSequences[[2]rune{prevChar, char}]
		for _, subSeq := range options {
			subLength := computeMinLength(subSeq, depth-1, dirSequences, dirLengths)
			if subLength < minSubLength {
				minSubLength = subLength
			}
		}
		totalMinLength += minSubLength
		prevChar = char
	}

	computeCache[cacheKey] = totalMinLength
	return totalMinLength
}

func solve(input []string, depth int) uint64 {
	fmt.Println("Pre-computing sequences for numeric keypad...")
	numSequences := computeSequences(numKeypad)
	fmt.Println("Pre-computing sequences for directional keypad...")
	dirSequences := computeSequences(dirKeypad)

	// Create a simple length map for the base case (depth=1) of the recursion.
	dirLengths := make(map[[2]rune]int)
	for k, v := range dirSequences {
		dirLengths[k] = len(v[0])
	}

	// Clear cache for a new run.
	computeCache = make(map[string]uint64)

	var totalComplexity uint64
	for _, line := range input {
		// Generate all possible D-pad sequences for the numeric code.
		possibleSequences := generateAllNumSequences(line, numSequences)

		minOverallLength := uint64(math.MaxUint64)

		// For each possible sequence, compute its minimum cost at the target depth.
		for _, s := range possibleSequences {
			length := computeMinLength(s, depth, dirSequences, dirLengths)
			if length < minOverallLength {
				minOverallLength = length
			}
		}

		numPart, _ := strconv.ParseUint(strings.TrimSuffix(line, "A"), 10, 64)
		totalComplexity += minOverallLength * numPart
	}

	return totalComplexity
}

func main() {
	file, err := os.Open("input")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not open 'input' file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var puzzleCodes []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if txt := strings.TrimSpace(scanner.Text()); txt != "" {
			puzzleCodes = append(puzzleCodes, txt)
		}
	}

	if len(puzzleCodes) == 0 {
		fmt.Println("'input' file is empty. Cannot proceed.")
		return
	}

	fmt.Println("--- Calculating Part 2 ---")
	// Following the reference Rust code's interpretation, Part 2 has a depth of 25.
	result := solve(puzzleCodes, 25)
	fmt.Printf("\nFinal Sum of complexities for Part 2: %d\n", result)
}
