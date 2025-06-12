package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	// Prune a secret number by taking it modulo 16777216 (2^24)
	pruneMod = 16777216
	// In a single day, buyers generate 2000 new secret numbers
	iterations = 2000
)

// nextSecret generates the next secret number in the sequence.
func nextSecret(secret int) int {
	// 1. Multiply by 64, mix, and prune.
	res1 := secret * 64
	secret = (secret ^ res1) % pruneMod

	// 2. Divide by 32, mix, and prune.
	res2 := secret / 32 // Integer division rounds down.
	secret = (secret ^ res2) % pruneMod

	// 3. Multiply by 2048, mix, and prune.
	res3 := secret * 2048
	secret = (secret ^ res3) % pruneMod

	return secret
}

// generateFinalSecret simulates the process for a given number of iterations.
func generateFinalSecret(initialSecret int, iters int) int {
	currentSecret := initialSecret
	for i := 0; i < iters; i++ {
		currentSecret = nextSecret(currentSecret)
	}
	return currentSecret
}

// solve processes the entire input and returns the sum of the 2000th secret numbers.
func solve(input string) int {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	totalSum := 0

	for _, line := range lines {
		if line == "" {
			continue
		}
		initialSecret, err := strconv.Atoi(line)
		if err != nil {
			panic(fmt.Sprintf("failed to parse number: %v", err))
		}
		finalSecret := generateFinalSecret(initialSecret, iterations)
		totalSum += finalSecret
	}

	return totalSum
}

func main() {
	data, err := os.ReadFile("input")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading input file: %v\n", err)
		os.Exit(1)
	}
	result := solve(string(data))
	fmt.Println(result)
}
