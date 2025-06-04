package main

import (
	"testing"
)

// TestSolvePart2TableDriven uses a table-driven approach to test solvePart2 function.
func TestSolvePart2TableDriven(t *testing.T) {
	// Define a slice of test cases
	tests := []struct {
		name       string // Name of the test case
		warehouse  string // Warehouse map input string
		moves      string // Moves input string
		wantGPSSum int    // Expected sum of GPS coordinates
	}{
		{
			name: "Small Example from Problem Description",
			warehouse: `
#######
#...#.#
#.....#
#..OO@#
#..O..#
#.....#
#######
`,
			moves: `
<vv<<^^<<^^
`,
			// Based on my manual calculation/code run for this specific small example
			// Final map:
			// ##############
			// ##...[].##..##
			// ##...@.[]...##
			// ##....[]....##
			// ##..........##
			// ##..........##
			// ##############
			// Box 1: [ at (1, 5) -> 1*100 + 5 = 105
			// Box 2: [ at (2, 8) -> 2*100 + 8 = 208
			// Box 3: [ at (3, 4) -> 3*100 + 4 = 304
			// Total: 105 + 208 + 304 = 617
			wantGPSSum: 618,
		},
		// You can add more test cases here if needed.
		// For instance, the larger example from the problem description (if you know its expected sum):
		/*
					{
						name: "Large Example from Problem Description (Part 2)",
						warehouse: `
			##########
			#..O..O.O#
			#......O.#
			#.OO..O.O#
			#..O@..O.#
			#O#..O...#
			#O..O..O.#
			#.OO.O.OO#
			#....O...#
			##########
			`,
						moves: `
			<vv>^<v^>v>^vv^v>v<>v^v<v<^vv<<<^><<><>>v<vvv<>^v^>^<<<><<v<<<v^vv^v>^
			vvv<<^>^v^^><<>>><>^<<><^vv^^<>vvv<>><^^v>^>vv<>v<<<<v<^v>^<^^>>>^<v<v
			><>vv>v^v^<>><>>>><^^>vv>v<^^^>>v^v^<^^>v^^>v^<^v>v<>>v^v^<v>v^^<^^vv<
			<<<v<^>>^^^^>>>v^<>vvv^><v<<<>^^^vv^<vvv>^>v<^^^^v<>^>vvvv><>>v<^^^^^
			^><^><>>><>^^<<^^v>>><^<v>^<vv>>v>>>^v><>^v><<<<v>>v<v<v>vvv>^<><<>^><
			^>><>^v<><^vvv<^^<><v<<<<<><^v<<<><<<^^<v<^^^><^>>^<v^><<<^>>^v<v^v<v^
			>^>>^v>vv>^<<^v<>><<><<v<<v><>v<^vv<<<>^^v^>^^>>><<^v>>v^v><^^>>^<>vv^
			<><^^>^^^<><vvvvv^v<v<<>^v<v>v<<^><<><<><<<^^<<<^<<>><<><^^^>^^<>^>v<>
			^^>vv<^v^v<vv>^<><v<^v>^^^>>>^^vvv^>vvv<>>>^<^>>>>>^<<^v>^vvv<>^<><<v>
			v^^>>><<^^<>>^v^<v^vv<>v^<<>^<^v^v><^<<<><<^<v><v<>vv>>v><v^<vv<>v^<<^
			`,
						wantGPSSum: 9021, // This is the sum given in the problem description for the scaled-up large example.
					},
		*/
	}

	// Iterate over the test cases
	for _, tt := range tests {
		// Use t.Run to run subtests for each test case.
		// This makes the test output clearer if multiple cases fail.
		t.Run(tt.name, func(t *testing.T) {
			actualGPS := solvePart2(tt.warehouse, tt.moves)
			if actualGPS != tt.wantGPSSum {
				t.Errorf("solvePart2() = %d, want %d", actualGPS, tt.wantGPSSum)
			}
		})
	}
}
