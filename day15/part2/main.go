package main

import (
	"fmt"
	"os"
	"strings"
)

// Point 结构体表示地图上的一个坐标 (物理字符坐标)
type Point struct {
	R, C int
}

// expandMap 函数将原始地图放大 (生成物理字符地图)
func expandMap(originalMapStr string) [][]rune {
// 	// (expandMap 函数代码保持不变，此处省略)
// 	specialSmallMapRaw := `
// #######
// #...#.#
// #.....#
// #..OO@#
// #..O..#
// #.....#
// #######
// `
// 	if strings.TrimSpace(originalMapStr) == strings.TrimSpace(specialSmallMapRaw) {
// 		hardcodedMap := make([][]rune, 7)
// 		hardcodedMap[0] = []rune{'#', '#', '#', '#', '#', '#', '#', '#', '#', '#', '#', '#', '#', '#'}
// 		hardcodedMap[1] = []rune{'#', '#', '.', '.', '.', '.', '.', '.', '#', '#', '.', '.', '#', '#'}
// 		hardcodedMap[2] = []rune{'#', '#', '.', '.', '.', '.', '.', '.', '.', '.', '.', '.', '#', '#'}
// 		hardcodedMap[3] = []rune{'#', '#', '.', '.', '.', '.', '[', ']', '[', ']', '@', '.', '#', '#'}
// 		hardcodedMap[4] = []rune{'#', '#', '.', '.', '.', '.', '[', ']', '.', '.', '.', '.', '#', '#'}
// 		hardcodedMap[5] = []rune{'#', '#', '.', '.', '.', '.', '.', '.', '.', '.', '.', '.', '#', '#'}
// 		hardcodedMap[6] = []rune{'#', '#', '#', '#', '#', '#', '#', '#', '#', '#', '#', '#', '#', '#'}
// 		return hardcodedMap
// 	}
	originalLines := strings.Split(strings.TrimSpace(originalMapStr), "\n")
	expandedMap := make([][]rune, 0)
	for _, line := range originalLines {
		expandedRow := make([]rune, 0, len(line)*2)
		for _, char := range line {
			switch char {
			case '#':
				expandedRow = append(expandedRow, '#', '#')
			case 'O':
				expandedRow = append(expandedRow, '[', ']')
			case '.':
				expandedRow = append(expandedRow, '.', '.')
			case '@':
				expandedRow = append(expandedRow, '@', '.')
			}
		}
		expandedMap = append(expandedMap, expandedRow)
	}
	return expandedMap
}

// printMap 辅助函数，用于打印当前地图状态
func printMap(m [][]rune, robot Point) {
	// (printMap 函数代码保持不变，此处省略)
	fmt.Println("--- Current Map State ---")
	for _, row := range m {
		for c := 0; c < len(row); c++ {
			char := row[c]
			if char == '[' {
				fmt.Print("[]")
				c++
			} else if char == '.' && c+1 < len(row) && row[c+1] == '.' {
				fmt.Print("..")
				c++
			} else if char == '#' && c+1 < len(row) && row[c+1] == '#' {
				fmt.Print("##")
				c++
			} else if char == '@' && c+1 < len(row) && row[c+1] == '.' {
				fmt.Print("@.")
				c++
			} else {
				fmt.Printf("%c", char)
			}
		}
		fmt.Println()
	}
	fmt.Println("Robot at:", robot)
	fmt.Println("-------------------------")
}

// solvePart2 模拟机器人和宽箱子在放大仓库中的移动
func solvePart2(warehouseMapStr, movesStr string) int {
	warehouseMap := expandMap(warehouseMapStr)

	fmt.Println("--- DIAGNOSIS: Map content immediately after expandMap (RAW RUNES) ---")
	for r, row := range warehouseMap {
		fmt.Printf("R%d: %s\n", r, string(row))
	}
	fmt.Println("-------------------------------------------------------------------")

	moves := strings.ReplaceAll(movesStr, "\n", "")
	rows := len(warehouseMap)
	if rows == 0 {
		return 0
	}
	cols := len(warehouseMap[0])

	var robotPos Point
	foundRobot := false
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if warehouseMap[r][c] == '@' {
				robotPos = Point{r, c}
				foundRobot = true
				break
			}
		}
		if foundRobot {
			break
		}
	}
	if !foundRobot {
		fmt.Println("ERROR: Robot not found in the map!")
		return 0
	}

	fmt.Println("--- Initial State ---")
	printMap(warehouseMap, robotPos)

	directions := map[rune]Point{
		'^': {-1, 0}, '<': {0, -1}, 'v': {1, 0}, '>': {0, 1},
	}

	for i, move := range moves {
		fmt.Printf("--- Move %d: %c ---\n", i+1, move)
		fmt.Printf("  Current robot pos: %v\n", robotPos)
		dr, dc := directions[move].R, directions[move].C
		nextRobotR, nextRobotC := robotPos.R+dr, robotPos.C+dc
		oldRobotPos := robotPos

		if nextRobotR < 0 || nextRobotR >= rows || nextRobotC < 0 || nextRobotC >= cols ||
			warehouseMap[nextRobotR][nextRobotC] == '#' {
			fmt.Printf("  Robot cannot move to (%d,%d): out of bounds or wall.\n", nextRobotR, nextRobotC)
			printMap(warehouseMap, robotPos)
			continue
		}

		if warehouseMap[nextRobotR][nextRobotC] == '[' || warehouseMap[nextRobotR][nextRobotC] == ']' {
			fmt.Println("  Robot attempting to push wide box(es).")

			firstBoxHitR, firstBoxHitC := nextRobotR, nextRobotC
			if warehouseMap[firstBoxHitR][firstBoxHitC] == ']' {
				firstBoxHitC--
			}
			fmt.Printf("  First box identified at (%d,%d) and (%d,%d).\n", firstBoxHitR, firstBoxHitC, firstBoxHitR, firstBoxHitC+1)

			// --- NEW BFS-based chain building logic ---
			pushChain := []Point{}
			queue := []Point{{R: firstBoxHitR, C: firstBoxHitC}} // Queue stores the '[' coords of boxes to process
			processedInChain := make(map[Point]bool)             // To avoid adding/processing the same box multiple times

			for len(queue) > 0 {
				currentBoxStart := queue[0]
				queue = queue[1:]

				if processedInChain[currentBoxStart] {
					continue
				}

				cbR, cbC := currentBoxStart.R, currentBoxStart.C
				if cbR < 0 || cbR >= rows || cbC < 0 || cbC+1 >= cols ||
					warehouseMap[cbR][cbC] != '[' || warehouseMap[cbR][cbC+1] != ']' {
					continue // Invalid box from queue
				}

				pushChain = append(pushChain, currentBoxStart)
				processedInChain[currentBoxStart] = true

				// Determine cells this currentBoxStart would try to push into
				var cellsToInvestigate []Point

				if dr != 0 { // Vertical push: checks two cells in front
					cellsToInvestigate = append(cellsToInvestigate, Point{R: cbR + dr, C: cbC})     // Cell above/below current box's '['
					cellsToInvestigate = append(cellsToInvestigate, Point{R: cbR + dr, C: cbC + 1}) // Cell above/below current box's ']'
				} else if dc != 0 { // Horizontal push
					if dc > 0 { // Pushing right: check cell to the right of current box's ']'
						cellsToInvestigate = append(cellsToInvestigate, Point{R: cbR, C: cbC + 1 + dc})
					} else { // Pushing left (dc < 0): check cell to the left of current box's '['
						cellsToInvestigate = append(cellsToInvestigate, Point{R: cbR, C: cbC + dc})
					}
				}

				for _, cell := range cellsToInvestigate {
					probeR, probeC := cell.R, cell.C
					actualNextBoxR, actualNextBoxC := -1, -1

					if probeR < 0 || probeR >= rows || probeC < 0 || probeC >= cols {
						continue // Probe point out of bounds
					}

					charAttProbe := warehouseMap[probeR][probeC]
					if charAttProbe == '[' { // Direct hit on a box's start
						actualNextBoxR, actualNextBoxC = probeR, probeC
					} else if charAttProbe == ']' { // Hit the right part of a box
						if probeC-1 >= 0 && warehouseMap[probeR][probeC-1] == '[' {
							actualNextBoxR, actualNextBoxC = probeR, probeC-1
						}
					}

					if actualNextBoxR != -1 {
						nextBoxPoint := Point{R: actualNextBoxR, C: actualNextBoxC}
						// Check if the identified next box is valid '[]' and not yet processed
						if !processedInChain[nextBoxPoint] &&
							actualNextBoxC+1 < cols && warehouseMap[actualNextBoxR][actualNextBoxC+1] == ']' {
							queue = append(queue, nextBoxPoint)
						}
					}
				}
			}
			// --- End of BFS-based chain building ---

			if len(pushChain) == 0 {
				fmt.Println("  No valid box chain found (this is unexpected if robot hit a box).")
				printMap(warehouseMap, robotPos)
				continue
			}
			fmt.Printf("  Identified push chain (length %d): %v\n", len(pushChain), pushChain)

			// Collision detection and movement logic (remains the same)
			totalShiftPhysicalR, totalShiftPhysicalC := dr, dc
			canPushChain := true

			for _, box := range pushChain { // Check wall/boundary collision for each box in chain
				newBoxR, newBoxC := box.R+totalShiftPhysicalR, box.C+totalShiftPhysicalC
				if newBoxR < 0 || newBoxR >= rows || newBoxC < 0 || newBoxC+1 >= cols ||
					warehouseMap[newBoxR][newBoxC] == '#' || warehouseMap[newBoxR][newBoxC+1] == '#' {
					fmt.Printf("  Chain collision: Wall/boundary for box (%d,%d) at new pos (%d,%d)-(%d,%d).\n", box.R, box.C, newBoxR, newBoxC, newBoxR, newBoxC+1)
					canPushChain = false
					break
				}
			}

			if canPushChain { // Check inter-box collision (non-chain boxes)
				tempMap := make([][]rune, rows)
				for r_copy := range warehouseMap {
					tempMap[r_copy] = make([]rune, cols)
					copy(tempMap[r_copy], warehouseMap[r_copy])
				}
				for _, box := range pushChain { // Clear chain boxes from tempMap
					tempMap[box.R][box.C] = '.'
					if box.C+1 < cols {
						tempMap[box.R][box.C+1] = '.'
					}
				}

				for _, box := range pushChain {
					newBoxR, newBoxC := box.R+totalShiftPhysicalR, box.C+totalShiftPhysicalC
					if (tempMap[newBoxR][newBoxC] == '[' || tempMap[newBoxR][newBoxC] == ']') ||
						(newBoxC+1 < cols && (tempMap[newBoxR][newBoxC+1] == '[' || tempMap[newBoxR][newBoxC+1] == ']')) {
						fmt.Printf("  Chain collision: Other box at new pos (%d,%d)-(%d,%d) for chain box originally at (%d,%d).\n", newBoxR, newBoxC, newBoxR, newBoxC+1, box.R, box.C)
						canPushChain = false
						break
					}
				}
			}

			if !canPushChain {
				fmt.Println("  Chain of boxes cannot be moved due to collision.")
				printMap(warehouseMap, robotPos)
				continue
			}

			fmt.Println("  Chain of boxes and Robot successfully moved.")
			warehouseMap[oldRobotPos.R][oldRobotPos.C] = '.' // Clear robot's old '@'

			for _, box := range pushChain { // Clear old positions of chain boxes
				warehouseMap[box.R][box.C] = '.'
				if box.C+1 < cols {
					warehouseMap[box.R][box.C+1] = '.'
				}
			}
			for _, box := range pushChain { // Place chain boxes in new positions
				newBoxR, newBoxC := box.R+totalShiftPhysicalR, box.C+totalShiftPhysicalC
				warehouseMap[newBoxR][newBoxC] = '['
				if newBoxC+1 < cols {
					warehouseMap[newBoxR][newBoxC+1] = ']'
				}
			}

			robotPos = Point{nextRobotR, nextRobotC}   // Move robot
			warehouseMap[robotPos.R][robotPos.C] = '@' // Place new robot '@'

		} else { // Robot moving to empty space
			fmt.Println("  Robot moving to empty space.")
			warehouseMap[oldRobotPos.R][oldRobotPos.C] = '.'
			robotPos = Point{nextRobotR, nextRobotC}
			warehouseMap[robotPos.R][robotPos.C] = '@'
			fmt.Println("  Robot successfully moved.")
		}
		printMap(warehouseMap, robotPos)
	}

	fmt.Println("--- Final State ---")
	printMap(warehouseMap, robotPos)
	totalGPSCoordinates := 0
	fmt.Println("--- Calculating Final GPS Sum ---")
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if warehouseMap[r][c] == '[' {
				gps := (r * 100) + c
				fmt.Printf("  Box at (%d,%d) ([) has GPS: %d\n", r, c, gps)
				totalGPSCoordinates += gps
			}
		}
	}
	fmt.Printf("--- Total GPS Sum: %d ---\n", totalGPSCoordinates)
	return totalGPSCoordinates
}

func main() {
	// 读取文件 "input"
	data, err := os.ReadFile("input") // 根据你的偏好，文件名是 "input"
	if err != nil {
		fmt.Println("读取文件错误:", err)
		return
	}

	// 将文件内容转换为字符串
	fileContent := string(data)

	// 按空行分割地图和移动指令
	// strings.SplitN 最多分N-1次，所以用2确保只在第一个空行分割
	parts := strings.SplitN(fileContent, "\n\n", 2)

	if len(parts) < 2 {
		fmt.Println("输入文件格式不正确：请确保地图和移动指令之间至少有一个空行。")
		return
	}

	warehouseInput := parts[0]
	movesInput := parts[1]

	// 调用 solvePart2 并打印结果
	// 注意: 为了让 main 函数能够调用 solvePart2, expandMap, printMap, Point,
	// 它们都需要在同一个包（本例中是 "main" 包）并且在同一个文件或者在同一个目录下被编译。
	// 如果你的 solvePart2 等函数在另一个文件中但属于同一个包，这没有问题。
	fmt.Println("--- Running Simulation from input file ---")
	result := solvePart2(warehouseInput, movesInput)
	fmt.Printf("Final GPS sum from input file: %d\n", result)
}
