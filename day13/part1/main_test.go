package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// ClawMachine represents the configuration and prize location for a single claw machine.
// type ClawMachine struct {
// 	MoveAX, MoveAY   int // Button A movement
// 	CostA            int // Button A token cost
// 	MoveBX, MoveBY   int // Button B movement
// 	CostB            int // Button B token cost
// 	TargetX, TargetY int // Prize target location
// }

// TODO: The actual implementation of CalculateMinTokens will go here.
// func CalculateMinTokens(machine ClawMachine) int { ... }

func TestCalculateMinTokens(t *testing.T) {
	tests := []struct {
		name     string
		machine  ClawMachine
		expected int // Expected minimum tokens, or -1 if impossible
	}{
		{
			name: "Example 1: Solvable Case",
			machine: ClawMachine{
				MoveAX: 94, MoveAY: 34, CostA: 3,
				MoveBX: 22, MoveBY: 67, CostB: 1,
				TargetX: 8400, TargetY: 5400,
			},
			expected: 280, // 80*3 + 40*1 = 280
		},
		{
			name: "Example 2: Unsolvable Case",
			machine: ClawMachine{
				MoveAX: 26, MoveAY: 66, CostA: 3,
				MoveBX: 67, MoveBY: 21, CostB: 1,
				TargetX: 12748, TargetY: 12176,
			},
			expected: -1, // No solution
		},
		{
			name: "Example 3: Another Solvable Case",
			machine: ClawMachine{
				MoveAX: 17, MoveAY: 86, CostA: 3,
				MoveBX: 84, MoveBY: 37, CostB: 1,
				TargetX: 7870, TargetY: 6450,
			},
			expected: 200, // 38*3 + 86*1 = 200
		},
		{
			name: "Example 4: Unsolvable Case",
			machine: ClawMachine{
				MoveAX: 69, MoveAY: 23, CostA: 3,
				MoveBX: 27, MoveBY: 71, CostB: 1,
				TargetX: 18641, TargetY: 10279,
			},
			expected: -1, // No solution
		},
		{
			name: "Zero Target - Should be 0 tokens",
			machine: ClawMachine{
				MoveAX: 10, MoveAY: 10, CostA: 3,
				MoveBX: 5, MoveBY: 5, CostB: 1,
				TargetX: 0, TargetY: 0,
			},
			expected: 0,
		},
		{
			name: "One Button Only - Solvable by A",
			machine: ClawMachine{
				MoveAX: 10, MoveAY: 10, CostA: 3,
				MoveBX: 0, MoveBY: 0, CostB: 1, // B button does nothing
				TargetX: 100, TargetY: 100,
			},
			expected: 30, // 10*3 = 30
		},
		{
			name: "One Button Only - Solvable by B",
			machine: ClawMachine{
				MoveAX: 0, MoveAY: 0, CostA: 3, // A button does nothing
				MoveBX: 5, MoveBY: 5, CostB: 1,
				TargetX: 50, TargetY: 50,
			},
			expected: 10, // 10*1 = 10
		},
		{
			name: "One Button Only - Unsolvable (B does nothing, A wrong target)",
			machine: ClawMachine{
				MoveAX: 10, MoveAY: 10, CostA: 3,
				MoveBX: 0, MoveBY: 0, CostB: 1,
				TargetX: 101, TargetY: 101,
			},
			expected: -1,
		},
		{
			name: "Co-linear movements, but no solution",
			machine: ClawMachine{
				MoveAX: 10, MoveAY: 20, CostA: 3,
				MoveBX: 20, MoveBY: 40, CostB: 1, // B is 2*A
				TargetX: 15, TargetY: 30,
			},
			expected: -1, // Can only reach multiples of (10,20)
		},
		{
			name: "Large numbers - Solvable",
			machine: ClawMachine{
				MoveAX: 1, MoveAY: 1000, CostA: 1,
				MoveBX: 1000, MoveBY: 1, CostB: 1,
				TargetX: 1001, TargetY: 1001,
			},
			expected: 2, // 1*1 + 1*1 = 2 (1 A, 1 B)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function we are testing
			actual := CalculateMinTokens(tt.machine) // This function is not yet implemented

			// Assert the result using testify/require
			require.Equal(t, tt.expected, actual, "Test case: %s", tt.name)
		})
	}
}
