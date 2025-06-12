package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	pruneMod        = 16777216
	numNewSecrets   = 2000 // Each buyer generates 2000 new secrets
	changeSeqLength = 4    // The monkey looks for a sequence of 4 changes
)

// nextSecret generates the next secret number in the sequence (reused from Part 1).
func nextSecret(secret int) int {
	res1 := secret * 64
	secret = (secret ^ res1) % pruneMod
	res2 := secret / 32
	secret = (secret ^ res2) % pruneMod
	res3 := secret * 2048
	secret = (secret ^ res3) % pruneMod
	return secret
}

// solve calculates the maximum number of bananas obtainable using an optimized approach.
func solve(input string) int {
	lines := strings.Split(strings.TrimSpace(input), "\n")

	// This map will act as a global scoreboard.
	// Key: a 4-change sequence. Value: total bananas earned for this sequence.
	sequenceScores := make(map[[changeSeqLength]int]int)

	// Process each buyer one by one.
	for _, line := range lines {
		if line == "" {
			continue
		}
		initialSecret, _ := strconv.Atoi(line)

		// --- Generate this buyer's history ---
		prices := make([]int, numNewSecrets+1)
		changes := make([]int, numNewSecrets)
		currentSecret := initialSecret
		prices[0] = currentSecret % 10

		for i := 1; i <= numNewSecrets; i++ {
			currentSecret = nextSecret(currentSecret)
			prices[i] = currentSecret % 10
			changes[i-1] = prices[i] - prices[i-1]
		}

		// --- Find first sale for each unique sequence for THIS buyer ---
		// Use a temporary map to ensure we only record the FIRST sale per sequence for this buyer.
		firstSaleForBuyer := make(map[[changeSeqLength]int]int)
		for i := 0; i <= len(changes)-changeSeqLength; i++ {
			var seq [changeSeqLength]int
			copy(seq[:], changes[i:i+changeSeqLength])

			// If we haven't seen this sequence for this buyer yet, record the sale.
			if _, ok := firstSaleForBuyer[seq]; !ok {
				priceAtSale := prices[i+changeSeqLength]
				firstSaleForBuyer[seq] = priceAtSale
			}
		}

		// --- Add this buyer's contributions to the global scoreboard ---
		for seq, price := range firstSaleForBuyer {
			sequenceScores[seq] += price
		}
	}

	// Find the highest score on the global scoreboard.
	maxBananas := 0
	for _, score := range sequenceScores {
		if score > maxBananas {
			maxBananas = score
		}
	}

	return maxBananas
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
