package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Rule represents a dependency rule where element A must come before element B
type Rule struct {
	A int
	B int
}

// Update represents a sequence of integers that needs to be validated against rules
type Update []int

// InputData holds the parsed input data
type InputData struct {
	Rules   []Rule
	Updates []Update
}

// ReadInput parses the input file and returns the rules and updates
func ReadInput(filename string) (*InputData, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data := &InputData{
		Rules:   []Rule{},
		Updates: []Update{},
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.Contains(line, "|") {
			rule, err := parseRule(line)
			if err != nil {
				fmt.Printf("Warning: %v\n", err)
				continue
			}
			data.Rules = append(data.Rules, rule)
		} else if strings.Contains(line, ",") {
			update, err := parseUpdate(line)
			if err != nil {
				fmt.Printf("Warning: %v\n", err)
				continue
			}
			data.Updates = append(data.Updates, update)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return data, nil
}

// parseRule parses a rule line in the format "A|B"
func parseRule(line string) (Rule, error) {
	parts := strings.Split(line, "|")
	if len(parts) != 2 {
		return Rule{}, fmt.Errorf("invalid rule format: %s", line)
	}

	a, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return Rule{}, fmt.Errorf("invalid rule integer A: %s", parts[0])
	}

	b, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return Rule{}, fmt.Errorf("invalid rule integer B: %s", parts[1])
	}

	return Rule{A: a, B: b}, nil
}

// parseUpdate parses an update line in the format "a, b, c, ..."
func parseUpdate(line string) (Update, error) {
	parts := strings.Split(line, ",")
	update := make(Update, 0, len(parts))

	for _, p := range parts {
		num, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return nil, fmt.Errorf("invalid update integer: %s", p)
		}
		update = append(update, num)
	}

	return update, nil
}

// IsValid checks if an update satisfies all rules
func (u Update) IsValid(rules []Rule) bool {
	positions := make(map[int]int)
	for i, v := range u {
		positions[v] = i
	}

	for _, rule := range rules {
		posA, hasA := positions[rule.A]
		posB, hasB := positions[rule.B]

		// Only check if both elements exist in the update
		if hasA && hasB && posA >= posB {
			return false
		}
	}
	return true
}

// TopologicalSort sorts the update according to the rules
func (u Update) TopologicalSort(rules []Rule) Update {
	// Create a graph representation
	graph := make(map[int][]int)
	inDegree := make(map[int]int)
	inUpdate := make(map[int]bool)

	// Mark elements in the update
	for _, v := range u {
		inUpdate[v] = true
		inDegree[v] = 0
	}

	// Build the graph based on rules
	for _, rule := range rules {
		if inUpdate[rule.A] && inUpdate[rule.B] {
			graph[rule.A] = append(graph[rule.A], rule.B)
			inDegree[rule.B]++
		}
	}

	// Kahn's algorithm for topological sorting
	queue := []int{}
	for _, v := range u {
		if inDegree[v] == 0 {
			queue = append(queue, v)
		}
	}

	sorted := make(Update, 0, len(u))
	used := make(map[int]bool)

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		sorted = append(sorted, node)
		used[node] = true

		for _, neighbor := range graph[node] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	// Add any remaining elements that weren't part of the dependency graph
	for _, v := range u {
		if !used[v] {
			sorted = append(sorted, v)
		}
	}

	return sorted
}

// GetMiddleElement returns the middle element of the update
func (u Update) GetMiddleElement() int {
	if len(u) == 0 {
		return 0
	}
	return u[len(u)/2]
}

// ProcessUpdates processes all updates and returns the sum of middle elements of invalid updates
func ProcessUpdates(data *InputData) int {
	sum := 0

	for _, update := range data.Updates {
		if !update.IsValid(data.Rules) {
			sorted := update.TopologicalSort(data.Rules)
			sum += sorted.GetMiddleElement()
		}
	}

	return sum
}

func main() {
	const inputFile = "input"

	data, err := ReadInput(inputFile)
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return
	}

	result := ProcessUpdates(data)
	fmt.Println("Sum of middle elements:", result)
}
