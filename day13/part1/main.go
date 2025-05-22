package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
)

// ClawMachine represents the configuration and prize location for a single claw machine.
type ClawMachine struct {
	MoveAX, MoveAY   int // Button A movement
	CostA            int // Button A token cost (given as 3)
	MoveBX, MoveBY   int // Button B movement
	CostB            int // Button B token cost (given as 1)
	TargetX, TargetY int // Prize target location
}

// CalculateMinTokens calculates the minimum tokens required to win a prize for a given machine.
// Returns -1 if the prize cannot be won.
// (This is the function we developed and tested earlier)
func CalculateMinTokens(machine ClawMachine) int {
	minTokens := math.MaxInt // Initialize with a very large number

	// Set a reasonable upper bound for numA and numB iterations.
	// Based on problem examples, target values are up to ~18000.
	// Smallest move is 1. So, maximum individual presses can be ~18000.
	// A safe upper bound for sum of numA and numB could be around 20000-50000.
	// Here, we iterate numA up to a generous limit.
	// If performance is critical for very large inputs,
	// more advanced math (e.g., extended Euclidean algorithm for Diophantine equations)
	// would be necessary. For typical AoC constraints, this loop range is usually fine.
	upperBound := 20000 // A sufficiently large upper bound for numA iterations

	// Handle the special case where target is (0,0)
	if machine.TargetX == 0 && machine.TargetY == 0 {
		return 0
	}

	// If all moves are 0, but target is not (0,0), it's unsolvable.
	if machine.MoveAX == 0 && machine.MoveAY == 0 && machine.MoveBX == 0 && machine.MoveBY == 0 {
		return -1
	}

	for numA := 0; numA <= upperBound; numA++ {
		currentX := numA * machine.MoveAX
		currentY := numA * machine.MoveAY

		remainingX := machine.TargetX - currentX
		remainingY := machine.TargetY - currentY

		// Optimization: if remaining is negative and all moves are positive, this path won't work.
		// Assumes MoveBX, MoveBY are always positive, which they are in this problem.
		if remainingX < 0 || remainingY < 0 {
			continue
		}

		// Calculate numB required for X-axis
		var numBX int
		if machine.MoveBX != 0 {
			if remainingX%machine.MoveBX != 0 {
				continue // numBX is not an integer
			}
			numBX = remainingX / machine.MoveBX
		} else { // MoveBX is 0
			if remainingX != 0 {
				continue // Cannot reach remainingX target if MoveBX is 0 and remainingX is not 0
			}
			numBX = 0 // If MoveBX is 0 and remainingX is 0, numB can be anything, but we look for min non-negative.
		}

		// Calculate numB required for Y-axis
		var numBY int
		if machine.MoveBY != 0 {
			if remainingY%machine.MoveBY != 0 {
				continue // numBY is not an integer
			}
			numBY = remainingY / machine.MoveBY
		} else { // MoveBY is 0
			if remainingY != 0 {
				continue // Cannot reach remainingY target if MoveBY is 0 and remainingY is not 0
			}
			numBY = 0 // Similar logic for Y-axis
		}

		// Check if numB values are non-negative
		if numBX < 0 || numBY < 0 {
			continue
		}

		// Crucial check: numBX and numBY must be equal to satisfy both axes
		// unless one of the MoveB values is 0, in which case we check the other.
		if (machine.MoveBX != 0 && machine.MoveBY != 0 && numBX == numBY) ||
			(machine.MoveBX != 0 && machine.MoveBY == 0 && remainingY == 0) || // B only moves X, Y already fulfilled by A or is 0
			(machine.MoveBX == 0 && machine.MoveBY != 0 && remainingX == 0) || // B only moves Y, X already fulfilled by A or is 0
			(machine.MoveBX == 0 && machine.MoveBY == 0 && remainingX == 0 && remainingY == 0) { // B moves nothing, X,Y already fulfilled by A

			// Determine the actual numB to use.
			// If one of MoveB is 0, the other's numB determines the actual button presses.
			// If both MoveB are 0, numB is 0.
			var actualNumB int
			if machine.MoveBX != 0 {
				actualNumB = numBX
			} else if machine.MoveBY != 0 {
				actualNumB = numBY
			} else { // Both MoveBX and MoveBY are 0
				actualNumB = 0
			}

			// Ensure that if one move is 0, the other axis is also consistent.
			// This covers cases like: B moves X, but Y is already at target.
			// Or B moves Y, but X is already at target.
			// The crucial part is that the calculated numB (actualNumB)
			// should not create new misalignments if one move is 0.
			// If MoveBX != 0 and remainingY != actualNumB * machine.MoveBY: this is a mismatch
			// If MoveBY != 0 and remainingX != actualNumB * machine.MoveBX: this is a mismatch

			// This check is implicitly handled by the previous remainingX/Y and numBX/numBY calculations.
			// The most robust check is simply that if numBX and numBY were calculated (i.e. corresponding MoveB != 0), they must match.
			// If one of them is 0 (MoveB is 0), then the other needs to be valid.
			// The conditions above for `if (machine.MoveBX != 0 && machine.MoveBY != 0 && numBX == numBY)` cover the main case.
			// The subsequent `else if` cases for when one move is 0 also work.
			// Let's refine the logic for combining numBX and numBY.

			// A simpler approach to combine: if any of the non-zero move axes
			// lead to different numB, then it's not a solution for this numA.

			isValidCombination := true
			if machine.MoveBX != 0 && machine.MoveBY != 0 {
				if numBX != numBY {
					isValidCombination = false
				}
			} else if machine.MoveBX != 0 { // MoveBX is non-zero, MoveBY is zero
				// Must ensure that remainingY is 0, because B cannot affect Y.
				if remainingY != 0 {
					isValidCombination = false
				}
				actualNumB = numBX
			} else if machine.MoveBY != 0 { // MoveBY is non-zero, MoveBX is zero
				// Must ensure that remainingX is 0, because B cannot affect X.
				if remainingX != 0 {
					isValidCombination = false
				}
				actualNumB = numBY
			} else { // Both MoveBX and MoveBY are zero
				// Must ensure both remainingX and remainingY are 0.
				if remainingX != 0 || remainingY != 0 {
					isValidCombination = false
				}
				actualNumB = 0 // No presses of B button needed
			}

			if isValidCombination {
				currentTokens := numA*machine.CostA + actualNumB*machine.CostB
				if currentTokens < minTokens {
					minTokens = currentTokens
				}
			}
		}
	}

	if minTokens == math.MaxInt {
		return -1 // No solution found
	}
	return minTokens
}

// parseInput parses the input string into a slice of ClawMachine structs.
func parseInput(input string) []ClawMachine {
	var machines []ClawMachine
	// Regex to match each block of machine data
	// Example: Button A: X+94, Y+34\nButton B: X+22, Y+67\nPrize: X=8400, Y=5400
	re := regexp.MustCompile(`Button A: X\+(\d+), Y\+(\d+)\nButton B: X\+(\d+), Y\+(\d+)\nPrize: X=(\d+), Y=(\d+)`)
	matches := re.FindAllStringSubmatch(input, -1)

	for _, m := range matches {
		if len(m) != 7 { // Expect 6 captured groups + full match
			log.Fatalf("Input parsing error: unexpected number of matches for a block: %v", m)
		}
		// Convert captured strings to integers
		moveAX, _ := strconv.Atoi(m[1])
		moveAY, _ := strconv.Atoi(m[2])
		moveBX, _ := strconv.Atoi(m[3])
		moveBY, _ := strconv.Atoi(m[4])
		targetX, _ := strconv.Atoi(m[5])
		targetY, _ := strconv.Atoi(m[6])

		// Costs are fixed for this problem: A=3, B=1
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
func SolveClawContraption(machines []ClawMachine) int {
	winnablePrizes := 0
	totalMinTokensForWinnable := 0

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

	// The problem asks for "the smallest number of tokens you would have to spend
	// to win as many prizes as possible".
	// This implies we sum the minimum tokens for all solvable machines.
	// If no prizes are winnable, totalMinTokensForWinnable will be 0.
	// If only 1 prize is winnable, it's just that prize's min tokens.
	// The problem states: "So, the most prizes you could possibly win is two;
	// the minimum tokens you would have to spend to win all (two) prizes is 480."
	// This confirms we just sum up the minimums for solvable machines.

	if winnablePrizes == 0 {
		fmt.Println("No prizes are winnable.")
		return 0 // Or -1, depending on specific problem interpretation if no prizes can be won at all.
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
	machines := parseInput(string(inputData))

	// Renamed to SolveClawContraption to differentiate from single machine calculation
	totalFewestTokens := SolveClawContraption(machines)
	fmt.Println("最终结果: ", totalFewestTokens)
}
