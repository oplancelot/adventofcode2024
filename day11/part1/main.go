package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// applyRules applies the transformation rules to a single stone
func applyRules(stone int) []int {
	// Convert to string to work with digits
	stoneStr := strconv.Itoa(stone)
	
	// Rule 1: If the stone is 0, replace with 1
	if stone == 0 {
		return []int{1}
	}
	
	// Rule 2: If the stone has an even number of digits, split it
	if len(stoneStr) % 2 == 0 {
		midpoint := len(stoneStr) / 2
		leftHalf := stoneStr[:midpoint]
		rightHalf := stoneStr[midpoint:]
		
		leftNum, _ := strconv.Atoi(leftHalf)
		rightNum, _ := strconv.Atoi(rightHalf)
		
		return []int{leftNum, rightNum}
	}
	
	// Rule 3: Multiply by 2024
	return []int{stone * 2024}
}

// simulateBlink simulates one blink transformation on all stones
func simulateBlink(stones []int) []int {
	var newStones []int
	
	for _, stone := range stones {
		transformedStones := applyRules(stone)
		newStones = append(newStones, transformedStones...)
	}
	
	return newStones
}

// readInput reads the initial stone arrangement from a file
func readInput(filename string) ([]int, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	
	// Split by whitespace to get individual stone values
	stoneStrs := strings.Fields(strings.TrimSpace(string(data)))
	stones := make([]int, len(stoneStrs))
	
	for i, stoneStr := range stoneStrs {
		stone, err := strconv.Atoi(stoneStr)
		if err != nil {
			return nil, fmt.Errorf("invalid stone value: %s", stoneStr)
		}
		stones[i] = stone
	}
	
	return stones, nil
}

func main() {
	const inputFile = "input"
	stones, err := readInput(inputFile)
	if err != nil {
		fmt.Printf("Failed to read input file (%s): %v\n", inputFile, err)
		return
	}
	
	fmt.Printf("Initial arrangement: %v\n", stones)
	
	// Simulate 25 blinks
	for i := 1; i <= 25; i++ {
		stones = simulateBlink(stones)
		if i <= 10 || i == 25 {
			fmt.Printf("After %d blinks: %d stones\n", i, len(stones))
		}
	}
	
	fmt.Printf("Final count: %d stones\n", len(stones))
}
