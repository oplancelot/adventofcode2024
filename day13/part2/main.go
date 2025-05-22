package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

// ClawMachine represents the configuration and prize location for a single claw machine.
type ClawMachine struct {
	MoveAX, MoveAY   int64 // Button A movement
	CostA            int64 // Button A token cost
	MoveBX, MoveBY   int64 // Button B movement
	CostB            int64 // Button B token cost
	TargetX, TargetY int64 // Prize target location
}

// CalculateMinTokens calculates the minimum tokens required to win a prize for a given machine.
// Returns -1 if the prize cannot be won.
// This version uses a direct algebraic solution for numA and numB.
func CalculateMinTokens(machine ClawMachine) int64 {
	// Constants for easier reading
	a, b, c, d := machine.MoveAX, machine.MoveBX, machine.MoveAY, machine.MoveBY
	targetX, targetY := machine.TargetX, machine.TargetY

	// Handle the special case where target is (0,0)
	if targetX == 0 && targetY == 0 {
		return 0
	}

	// Calculate the determinant
	determinant := a*d - b*c

	// Case 1: Determinant is zero (linear dependence)
	if determinant == 0 {
		// If both buttons move nothing, but target is not (0,0), it's unsolvable.
		if a == 0 && c == 0 && b == 0 && d == 0 {
			return -1
		}

		// Handle cases where one button doesn't move.
		// If Button B does nothing (MoveBX=0, MoveBY=0)
		if b == 0 && d == 0 {
			// If Button A also does nothing, already handled above.
			// Only A button moves. targetX must be multiple of a, targetY must be multiple of c.
			if a != 0 && targetX%a != 0 {
				return -1
			}
			if c != 0 && targetY%c != 0 {
				return -1
			}

			var numAFromX int64 = -1
			if a != 0 {
				numAFromX = targetX / a
			} else if targetX != 0 { // TargetX must be 0 if a is 0
				return -1
			}

			var numAFromY int64 = -1
			if c != 0 {
				numAFromY = targetY / c
			} else if targetY != 0 { // TargetY must be 0 if c is 0
				return -1
			}

			// Combine results for numA
			var finalNumA int64
			if a != 0 && c != 0 {
				if numAFromX != numAFromY || numAFromX < 0 {
					return -1
				}
				finalNumA = numAFromX
			} else if a != 0 {
				if numAFromX < 0 {
					return -1
				}
				finalNumA = numAFromX
			} else if c != 0 {
				if numAFromY < 0 {
					return -1
				}
				finalNumA = numAFromY
			} else { // both a and c are 0, and target is not (0,0)
				return -1
			}
			return finalNumA * machine.CostA // Only A button costs
		}

		// Similarly for Button A doing nothing (MoveAX=0, MoveAY=0)
		if a == 0 && c == 0 { // Button A does nothing, Button B moves
			if b != 0 && targetX%b != 0 {
				return -1
			}
			if d != 0 && targetY%d != 0 {
				return -1
			}

			var numBFromX int64 = -1
			if b != 0 {
				numBFromX = targetX / b
			} else if targetX != 0 {
				return -1
			}

			var numBFromY int64 = -1
			if d != 0 {
				numBFromY = targetY / d
			} else if targetY != 0 {
				return -1
			}

			var finalNumB int64
			if b != 0 && d != 0 {
				if numBFromX != numBFromY || numBFromX < 0 {
					return -1
				}
				finalNumB = numBFromX
			} else if b != 0 {
				if numBFromX < 0 {
					return -1
				}
				finalNumB = numBFromX
			} else if d != 0 {
				if numBFromY < 0 {
					return -1
				}
				finalNumB = numBFromY
			} else { // both b and d are 0, and target is not (0,0)
				return -1
			}
			return finalNumB * machine.CostB
		}

		// General collinear case (determinant == 0 but neither button is totally useless)
		// For Advent of Code problems, if determinant is 0 and it's not one of the simple
		// single-button cases above, it often implies 'unsolvable' in this context,
		// as a full Diophantine equation solver is usually beyond typical AoC scope.
		return -1
	}

	// Case 2: Determinant is non-zero (unique solution)
	// Calculate numA and numB using Cramer's rule / algebraic solution
	numA_numerator := targetX*d - targetY*b
	numB_numerator := targetY*a - targetX*c

	// Check if numerators are perfectly divisible by the determinant
	if numA_numerator%determinant != 0 || numB_numerator%determinant != 0 {
		return -1 // No integer solution
	}

	numA := numA_numerator / determinant
	numB := numB_numerator / determinant

	// Check if solutions are non-negative
	if numA < 0 || numB < 0 {
		return -1 // No non-negative integer solution
	}

	return numA*machine.CostA + numB*machine.CostB
}

// parseInput parses the input string into a slice of ClawMachine structs.
func parseInput(input string) []ClawMachine {
	var machines []ClawMachine
	re := regexp.MustCompile(`Button A: X\+(\d+), Y\+(\d+)\nButton B: X\+(\d+), Y\+(\d+)\nPrize: X=(\d+), Y=(\d+)`)
	matches := re.FindAllStringSubmatch(input, -1)

	const prizeOffset int64 = 10000000000000 // 10 Trillion

	for _, m := range matches {
		if len(m) != 7 {
			log.Fatalf("Input parsing error: unexpected number of matches for a block: %v", m)
		}
		// Convert captured strings to int64 directly
		moveAX, _ := strconv.ParseInt(m[1], 10, 64) // ParseInt returns int64
		moveAY, _ := strconv.ParseInt(m[2], 10, 64)
		moveBX, _ := strconv.ParseInt(m[3], 10, 64)
		moveBY, _ := strconv.ParseInt(m[4], 10, 64)
		targetX, _ := strconv.ParseInt(m[5], 10, 64)
		targetY, _ := strconv.ParseInt(m[6], 10, 64)

		// Apply the prize offset for Part Two
		targetX += prizeOffset
		targetY += prizeOffset

		// Costs are fixed for this problem: A=3, B=1. Literal integers will be implicitly converted to int64.
		machines = append(machines, ClawMachine{
			MoveAX: moveAX, MoveAY: moveAY, CostA: 3,
			MoveBX: moveBX, MoveBY: moveBY, CostB: 1,
			TargetX: targetX, TargetY: targetY,
		})
	}
	return machines
}

// SolveClawContraption processes all machines and returns the minimum tokens needed
// to win as many prizes as possible.
func SolveClawContraption(machines []ClawMachine) int64 {
	var winnablePrizes int64 = 0
	var totalMinTokensForWinnable int64 = 0

	for i, machine := range machines {
		minTokens := CalculateMinTokens(machine)
		if minTokens != -1 {
			winnablePrizes++
			totalMinTokensForWinnable += minTokens
			fmt.Printf("Machine %d: Solvable with %d tokens\n", i+1, minTokens)
		} else {
			fmt.Printf("Machine %d: Unsolvable\n", i+1)
		}
	}

	if winnablePrizes == 0 {
		fmt.Println("No prizes are winnable.")
		return 0
	}

	fmt.Printf("Total winnable prizes: %d\n", winnablePrizes)
	return totalMinTokensForWinnable
}

func main() {
	const inputFile = "input"
	inputData, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("从 %s 读取谜题输入失败: %v", inputFile, err)
	}
	machines := parseInput(string(inputData)) // parseInput will now add the offset

	totalFewestTokens := SolveClawContraption(machines)
	fmt.Println("最终结果: ", totalFewestTokens)
}
