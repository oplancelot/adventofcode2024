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

// buyerData holds the generated price and change history for a single buyer.
type buyerData struct {
	prices  []int
	changes []int
}

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

// generateBuyerHistory creates the full list of prices and changes for one buyer.
func generateBuyerHistory(initialSecret int) buyerData {
	prices := make([]int, numNewSecrets+1)
	changes := make([]int, numNewSecrets)

	currentSecret := initialSecret
	prices[0] = currentSecret % 10

	for i := 1; i <= numNewSecrets; i++ {
		currentSecret = nextSecret(currentSecret)
		prices[i] = currentSecret % 10
		changes[i-1] = prices[i] - prices[i-1]
	}

	return buyerData{prices: prices, changes: changes}
}

// solve calculates the maximum number of bananas obtainable.
func solve(input string) int {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var buyersHistory []buyerData
	for _, line := range lines {
		if line == "" {
			continue
		}
		initialSecret, _ := strconv.Atoi(line)
		buyersHistory = append(buyersHistory, generateBuyerHistory(initialSecret))
	}

	// Collect all unique 4-change sequences that occurred.
	// We use a map where the key is a 4-integer array to find unique sequences.
	candidateSequences := make(map[[changeSeqLength]int]bool)
	for _, buyer := range buyersHistory {
		for i := 0; i <= len(buyer.changes)-changeSeqLength; i++ {
			var seq [changeSeqLength]int
			copy(seq[:], buyer.changes[i:i+changeSeqLength])
			candidateSequences[seq] = true
		}
	}

	maxBananas := 0

	// Test each candidate sequence to see how many bananas it yields.
	for seq := range candidateSequences {
		currentBananas := 0
		for _, buyer := range buyersHistory {
			// Find the first time this sequence occurs for the buyer.
			for i := 0; i <= len(buyer.changes)-changeSeqLength; i++ {
				match := true
				for j := 0; j < changeSeqLength; j++ {
					if buyer.changes[i+j] != seq[j] {
						match = false
						break
					}
				}

				if match {
					// Sale is made at the price corresponding to the end of the sequence.
					// The change at index `i+3` is `price[i+4] - price[i+3]`.
					// So, the price of the sale is `price[i+4]`.
					priceAtSale := buyer.prices[i+changeSeqLength]
					currentBananas += priceAtSale
					goto nextBuyer // Move to the next buyer after the first sale.
				}
			}
		nextBuyer:
		}
		if currentBananas > maxBananas {
			maxBananas = currentBananas
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
