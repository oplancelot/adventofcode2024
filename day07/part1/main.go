package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Equation represents a calibration equation
type Equation struct {
	testValue int
	numbers   []int
}

// parseInput reads the input file and parses the equations
func parseInput(filename string) ([]Equation, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	equations := make([]Equation, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line format: %s", line)
		}

		testValue, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid test value: %s", parts[0])
		}

		numStrs := strings.Fields(strings.TrimSpace(parts[1]))
		numbers := make([]int, 0, len(numStrs))
		for _, numStr := range numStrs {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				return nil, fmt.Errorf("invalid number: %s", numStr)
			}
			numbers = append(numbers, num)
		}

		equations = append(equations, Equation{
			testValue: testValue,
			numbers:   numbers,
		})
	}

	return equations, nil
}

// evaluateExpression evaluates an expression with the given numbers and operators
func evaluateExpression(numbers []int, operators []string) int {
	result := numbers[0]
	for i := 0; i < len(operators); i++ {
		if operators[i] == "+" {
			result += numbers[i+1]
		} else if operators[i] == "*" {
			result *= numbers[i+1]
		}
	}
	return result
}

// canSolveEquation checks if the equation can be solved with any combination of + and * operators
func canSolveEquation(eq Equation) bool {
	// If there's only one number, check if it equals the test value
	if len(eq.numbers) == 1 {
		return eq.numbers[0] == eq.testValue
	}

	// Generate all possible combinations of operators
	numOperators := len(eq.numbers) - 1
	maxCombinations := 1 << numOperators // 2^numOperators

	for i := 0; i < maxCombinations; i++ {
		operators := make([]string, numOperators)
		for j := 0; j < numOperators; j++ {
			if (i & (1 << j)) == 0 {
				operators[j] = "+"
			} else {
				operators[j] = "*"
			}
		}

		// Evaluate the expression with this combination of operators
		result := evaluateExpression(eq.numbers, operators)
		if result == eq.testValue {
			return true
		}
	}

	return false
}

func main() {
	equations, err := parseInput("input")
	if err != nil {
		fmt.Printf("Error parsing input: %v\n", err)
		return
	}

	totalCalibration := 0
	for _, eq := range equations {
		if canSolveEquation(eq) {
			totalCalibration += eq.testValue
		}
	}

	fmt.Printf("Total calibration result: %d\n", totalCalibration)
}
