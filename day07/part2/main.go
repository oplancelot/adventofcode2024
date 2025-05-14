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
	// Create a copy of the numbers to work with
	nums := make([]int, len(numbers))
	copy(nums, numbers)
	
	// Process operators from left to right
	for i := 0; i < len(operators); i++ {
		// Apply the operator between nums[0] and nums[1]
		if operators[i] == "+" {
			nums[0] += nums[1]
		} else if operators[i] == "*" {
			nums[0] *= nums[1]
		} else if operators[i] == "||" {
			// Concatenation operator
			// Convert both numbers to strings, concatenate, then convert back to int
			numStr1 := strconv.Itoa(nums[0])
			numStr2 := strconv.Itoa(nums[1])
			concatenated, _ := strconv.Atoi(numStr1 + numStr2)
			nums[0] = concatenated
		}
		
		// Shift the remaining numbers left
		for j := 1; j < len(nums)-1; j++ {
			nums[j] = nums[j+1]
		}
	}
	
	return nums[0]
}

// canSolveEquation checks if the equation can be solved with any combination of +, *, and || operators
func canSolveEquation(eq Equation) bool {
	// If there's only one number, check if it equals the test value
	if len(eq.numbers) == 1 {
		return eq.numbers[0] == eq.testValue
	}

	// Generate all possible combinations of operators
	numOperators := len(eq.numbers) - 1
	// Now we have 3 operators (+, *, ||), so we need 3^numOperators combinations
	
	// Helper function to generate all possible operator combinations
	var checkCombinations func(int, []string) bool
	checkCombinations = func(pos int, ops []string) bool {
		if pos == numOperators {
			// We have a complete set of operators, evaluate the expression
			result := evaluateExpression(eq.numbers, ops)
			return result == eq.testValue
		}
		
		// Try each operator at the current position
		for _, op := range []string{"+", "*", "||"} {
			ops[pos] = op
			if checkCombinations(pos+1, ops) {
				return true
			}
		}
		
		return false
	}
	
	operators := make([]string, numOperators)
	return checkCombinations(0, operators)
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