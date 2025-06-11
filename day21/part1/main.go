package main

import (
	"container/list"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// KeypadButton 表示键盘上的一个按钮，可能是数字，也可能是方向键，也可能是空
type KeypadButton rune

const (
	Empty KeypadButton = ' ' // 间隙
)

// 定义键盘布局
// 数字键盘 (Numeric Keypad)
// +---+---+---+
// | 7 | 8 | 9 |
// +---+---+---+
// | 4 | 5 | 6 |
// +---+---+---+
// | 1 | 2 | 3 |
// +---+---+---+
//
//	| 0 | A |
//	+---+---+
var numericKeypad = [4][3]KeypadButton{
	{'7', '8', '9'},
	{'4', '5', '6'},
	{'1', '2', '3'},
	{Empty, '0', 'A'}, // 'A' 是激活键
}

// 方向键盘 (Directional Keypad)
//
//	+---+---+
//	| ^ | A |
//
// +---+---+---+
// | < | v | > |
// +---+---+---+
var directionalKeypad = [2][3]KeypadButton{
	{Empty, '^', 'A'}, // 'A' 是激活键
	{'<', 'v', '>'},
}

// Pos 表示键盘上的位置 (Row, Col)
type Pos struct {
	R int
	C int
}

// State 表示 BFS 过程中的一个完整状态
type State struct {
	YourRobotPos    Pos    // 你的机器人在你的方向键盘上的位置
	Robot2Pos       Pos    // 第二个机器人在它的方向键盘上的位置
	Robot1Pos       Pos    // 第一个机器人在数字键盘上的位置
	TargetCharIndex int    // 目标代码中当前要按的字符的索引
	TotalSteps      int    // 你在你的键盘上按下的总次数 (这是我们要最小化的值)
	Path            string // 用于调试，记录按键路径
}

// StateKey 用于 visited map 的键，只包含需要避免重复搜索的状态部分
type StateKey struct {
	YourRobotR      int
	YourRobotC      int
	Robot2R         int
	Robot2C         int
	Robot1R         int
	Robot1C         int
	TargetCharIndex int
}

// getButtonAtPos 获取指定键盘在某个位置的按钮
func getButtonAtPos(p Pos, keypad [][3]KeypadButton) KeypadButton {
	if p.R < 0 || p.R >= len(keypad) || p.C < 0 || p.C >= len(keypad[0]) {
		return Empty // 超出边界
	}
	return keypad[p.R][p.C]
}

// isValidMove 检查移动是否有效 (在键盘范围内且不是空位)
func isValidMove(p Pos, keypad [][3]KeypadButton) bool {
	if p.R < 0 || p.R >= len(keypad) || p.C < 0 || p.C >= len(keypad[0]) {
		return false // 超出边界
	}
	return keypad[p.R][p.C] != Empty
}

// findInitialAPos 找到键盘上的 'A' 键的初始位置
func findInitialAPos(keypad [][3]KeypadButton) (Pos, error) {
	for r := 0; r < len(keypad); r++ {
		for c := 0; c < len(keypad[0]); c++ {
			if keypad[r][c] == 'A' {
				return Pos{r, c}, nil
			}
		}
	}
	return Pos{-1, -1}, fmt.Errorf("A key not found in keypad")
}

// FindShortestSequenceLength 计算到达目标代码所需的最少按键次数。
func FindShortestSequenceLength(targetCode string) int {
	// fmt.Printf("\n--- Calculating for code: %s ---\n", targetCode) // Debug info

	yourRobotInitialPos, err := findInitialAPos(directionalKeypad[:])
	if err != nil {
		panic(err)
	}
	robot2InitialPos, err := findInitialAPos(directionalKeypad[:])
	if err != nil {
		panic(err)
	}
	robot1InitialPos, err := findInitialAPos(numericKeypad[:])
	if err != nil {
		panic(err)
	}

	queue := list.New()
	initialState := State{
		YourRobotPos:    yourRobotInitialPos,
		Robot2Pos:       robot2InitialPos,
		Robot1Pos:       robot1InitialPos,
		TargetCharIndex: 0,
		TotalSteps:      0,
		Path:            "",
	}
	queue.PushBack(initialState)

	visited := make(map[StateKey]bool)
	visited[StateKey{
		YourRobotR:      initialState.YourRobotPos.R,
		YourRobotC:      initialState.YourRobotPos.C,
		Robot2R:         initialState.Robot2Pos.R,
		Robot2C:         initialState.Robot2Pos.C,
		Robot1R:         initialState.Robot1Pos.R,
		Robot1C:         initialState.Robot1Pos.C,
		TargetCharIndex: initialState.TargetCharIndex,
	}] = true

	dr := []int{-1, 1, 0, 0}
	dc := []int{0, 0, -1, 1}
	moveChars := []rune{'^', 'v', '<', '>'}

	for queue.Len() > 0 {
		e := queue.Front()
		queue.Remove(e)
		currentState := e.Value.(State)

		// 打印当前状态的关键信息（可选，如果输出过多可以注释掉）
		// fmt.Printf("Current: Steps=%d, Path='%s', YR=(%d,%d)%c, R2=(%d,%d)%c, R1=(%d,%d)%c, TargetIdx=%d\n",
		// 	currentState.TotalSteps, currentState.Path,
		// 	currentState.YourRobotPos.R, currentState.YourRobotPos.C, getButtonAtPos(currentState.YourRobotPos, directionalKeypad[:]),
		// 	currentState.Robot2Pos.R, currentState.Robot2Pos.C, getButtonAtPos(currentState.Robot2Pos, directionalKeypad[:]),
		// 	currentState.Robot1Pos.R, currentState.Robot1Pos.C, getButtonAtPos(currentState.Robot1Pos, numericKeypad[:]),
		// 	currentState.TargetCharIndex)

		if currentState.TargetCharIndex == len(targetCode) {
			// fmt.Printf("--- Found path for %s! Steps: %d, Path: %s ---\n", targetCode, currentState.TotalSteps, currentState.Path)
			return currentState.TotalSteps
		}

		// --- 你的机器人可以按方向键 (移动自己的臂) ---
		for i := 0; i < 4; i++ {
			newYourRobotPos := Pos{R: currentState.YourRobotPos.R + dr[i], C: currentState.YourRobotPos.C + dc[i]}
			if isValidMove(newYourRobotPos, directionalKeypad[:]) {
				newState := State{
					YourRobotPos:    newYourRobotPos,
					Robot2Pos:       currentState.Robot2Pos,
					Robot1Pos:       currentState.Robot1Pos,
					TargetCharIndex: currentState.TargetCharIndex,
					TotalSteps:      currentState.TotalSteps + 1,
					Path:            currentState.Path + string(moveChars[i]),
				}
				key := StateKey{
					YourRobotR:      newState.YourRobotPos.R,
					YourRobotC:      newState.YourRobotPos.C,
					Robot2R:         newState.Robot2Pos.R,
					Robot2C:         newState.Robot2Pos.C,
					Robot1R:         newState.Robot1Pos.R,
					Robot1C:         newState.Robot1Pos.C,
					TargetCharIndex: newState.TargetCharIndex,
				}
				if !visited[key] {
					visited[key] = true
					queue.PushBack(newState)
				}
			}
		}

		// --- 你的机器人也可以按 'A' 键 (激活下一个机器人) ---
		// 这个分支总是可以尝试，无论你的机器人当前在哪个键上。

		// 关键修正：你按 A 键时， Robot2 的行为取决于你当前在你的键盘上按下的那个键
		yourButtonOnApress := getButtonAtPos(currentState.YourRobotPos, directionalKeypad[:])

		nextRobot2Pos := currentState.Robot2Pos // 默认不变
		nextRobot1Pos := currentState.Robot1Pos // 默认不变
		nextTargetCharIndex := currentState.TargetCharIndex

		isThisAPressValid := true // 标记此次A键操作是否有效（没有导致机器人恐慌）

		// 调试：打印 A 键操作前的状态 (可选)
		// fmt.Printf("  Attempting A press. YourR=(%d,%d)%c (This is the button *you* pressed to activate R2). R2R=(%d,%d)%c, R1R=(%d,%d)%c\n",
		// 	currentState.YourRobotPos.R, currentState.YourRobotPos.C, yourButtonOnApress,
		// 	currentState.Robot2Pos.R, currentState.Robot2Pos.C, getButtonAtPos(currentState.Robot2Pos, directionalKeypad[:]),
		// 	currentState.Robot1Pos.R, currentState.Robot1Pos.C, getButtonAtPos(currentState.Robot1Pos, numericKeypad[:]))

		switch yourButtonOnApress { // 你的机器人当前指向的键决定了 Robot2 的行动
		case '^', 'v', '<', '>': // 如果你按的是方向键，然后按 A，这意味着你告诉 Robot2 移动
			var r2_dr, r2_dc int // 机器人2在它自己键盘上的移动方向
			switch yourButtonOnApress {
			case '^':
				r2_dr = -1
			case 'v':
				r2_dr = 1
			case '<':
				r2_dc = -1
			case '>':
				r2_dc = 1
			}

			tempNextRobot2Pos := Pos{R: currentState.Robot2Pos.R + r2_dr, C: currentState.Robot2Pos.C + r2_dc}
			if isValidMove(tempNextRobot2Pos, directionalKeypad[:]) { // 检查机器人2在它自己键盘上的新位置是否有效
				nextRobot2Pos = tempNextRobot2Pos
				// fmt.Printf("    You pressed '%c'+A. R2 moved from (%d,%d)%c to (%d,%d)%c\n",
				// 	yourButtonOnApress, currentState.Robot2Pos.R, currentState.Robot2Pos.C, getButtonAtPos(currentState.Robot2Pos, directionalKeypad[:]),
				// 	nextRobot2Pos.R, nextRobot2Pos.C, getButtonAtPos(nextRobot2Pos, directionalKeypad[:]))
			} else {
				isThisAPressValid = false // 机器人2移动到间隙，操作无效
				// fmt.Printf("    You pressed '%c'+A. R2 tried to move from (%d,%d)%c to INVALID pos (%d,%d) -> ABORT\n",
				// 	yourButtonOnApress, currentState.Robot2Pos.R, currentState.Robot2Pos.C, getButtonAtPos(currentState.Robot2Pos, directionalKeypad[:]),
				// 	tempNextRobot2Pos.R, tempNextRobot2Pos.C)
			}
			// Robot 1 的位置不变。

		case 'A': // 如果你按的是 A 键，然后按 A，意味着你告诉 Robot2 按下它当前指向的键
			robot2ButtonOnApress := getButtonAtPos(currentState.Robot2Pos, directionalKeypad[:]) // 获取机器人2当前指向的键

			// fmt.Printf("    You pressed 'A'+A. R2 is on '%c'. R2's action depends on this.\n", robot2ButtonOnApress)

			// 检查机器人2当前位置是否是有效按钮（不能在间隙上按）
			if !isValidMove(currentState.Robot2Pos, directionalKeypad[:]) {
				isThisAPressValid = false // 机器人2当前在间隙，无法按下
				// fmt.Printf("    R2 currently at invalid pos (%d,%d)%c -> ABORT\n",
				// 	currentState.Robot2Pos.R, currentState.Robot2Pos.C, getButtonAtPos(currentState.Robot2Pos, directionalKeypad[:]))
			} else {
				// Robot 2 实际的行动：
				switch robot2ButtonOnApress {
				case '^', 'v', '<', '>': // R2告诉R1移动 (R1在数字键盘上移动)
					var r1_dr, r1_dc int // 机器人1在数字键盘上的移动方向
					switch robot2ButtonOnApress {
					case '^':
						r1_dr = -1
					case 'v':
						r1_dr = 1
					case '<':
						r1_dc = -1
					case '>':
						r1_dc = 1
					}
					tempNextRobot1Pos := Pos{R: currentState.Robot1Pos.R + r1_dr, C: currentState.Robot1Pos.C + r1_dc}
					if isValidMove(tempNextRobot1Pos, numericKeypad[:]) { // 检查机器人1在数字键盘上的新位置是否有效
						nextRobot1Pos = tempNextRobot1Pos
						// fmt.Printf("      R2 sent '%c'. R1 moved from (%d,%d)%c to (%d,%d)%c\n",
						// 	robot2ButtonOnApress, currentState.Robot1Pos.R, currentState.Robot1Pos.C, getButtonAtPos(currentState.Robot1Pos, numericKeypad[:]),
						// 	nextRobot1Pos.R, nextRobot1Pos.C, getButtonAtPos(nextRobot1Pos, numericKeypad[:]))
					} else {
						isThisAPressValid = false // 机器人1移动到间隙，操作无效
						// fmt.Printf("      R2 sent '%c'. R1 tried to move from (%d,%d)%c to INVALID pos (%d,%d) -> ABORT\n",
						// 	robot2ButtonOnApress, currentState.Robot1Pos.R, currentState.Robot1Pos.C, getButtonAtPos(currentState.Robot1Pos, numericKeypad[:]),
						// 	tempNextRobot1Pos.R, tempNextRobot1Pos.C)
					}
					// 机器人1和机器人2的位置都不变。

				case 'A': // R2告诉R1按下当前键 (R1在数字键盘上按下当前键)
					robot1CurrentButton := getButtonAtPos(currentState.Robot1Pos, numericKeypad[:])
					if !isValidMove(currentState.Robot1Pos, numericKeypad[:]) {
						isThisAPressValid = false // 机器人1当前在间隙，无法按下
						// fmt.Printf("      R2 sent 'A'. R1 currently at invalid pos (%d,%d)%c -> ABORT\n",
						// 	currentState.Robot1Pos.R, currentState.Robot1Pos.C, getButtonAtPos(currentState.Robot1Pos, numericKeypad[:]))
					} else if nextTargetCharIndex < len(targetCode) {
						expectedTargetChar := KeypadButton(targetCode[nextTargetCharIndex])
						// fmt.Printf("      R2 sent 'A'. R1 at (%d,%d)%c attempts to press '%c', expecting '%c'\n",
						// 	 currentState.Robot1Pos.R, currentState.Robot1Pos.C, robot1CurrentButton, robot1CurrentButton, expectedTargetChar)
						if robot1CurrentButton == expectedTargetChar {
							nextTargetCharIndex++ // 成功按下目标数字
							// fmt.Printf("      *** TARGET CHAR MATCH! %c pressed. New TargetIdx: %d ***\n", expectedTargetChar, nextTargetCharIndex)
						} else {
							// fmt.Printf("      R1 pressed wrong button: '%c' != expected '%c'\n", robot1CurrentButton, expectedTargetChar)
						}
					}
					// 所有机器人的位置都不变。

				default: // Robot 2 当前指向 Empty，这也是 R2 的一个无效状态
					isThisAPressValid = false
					// fmt.Printf("      R2 is on Empty button (%d,%d)%c -> ABORT\n",
					// 	currentState.Robot2Pos.R, currentState.Robot2Pos.C, getButtonAtPos(currentState.Robot2Pos, directionalKeypad[:]))
				}
			}

		default: // 如果你的机器人当前指向的是一个 Empty 键，这也是你无法按下 'A' 来激活 R2 的无效状态
			isThisAPressValid = false
			// fmt.Printf("    You pressed an invalid button (%d,%d)%c to activate R2 -> ABORT\n",
			// 	currentState.YourRobotPos.R, currentState.YourRobotPos.C, yourButtonOnApress)
		}

		if isThisAPressValid {
			newState := State{
				YourRobotPos:    currentState.YourRobotPos, // 你的机器人位置不变 (它只是按A)
				Robot2Pos:       nextRobot2Pos,             // 机器人2位置根据你的操作可能改变
				Robot1Pos:       nextRobot1Pos,             // 机器人1位置根据R2的操作可能改变
				TargetCharIndex: nextTargetCharIndex,
				TotalSteps:      currentState.TotalSteps + 1, // 你的步数增加
				Path:            currentState.Path + "A",
			}
			key := StateKey{
				YourRobotR:      newState.YourRobotPos.R,
				YourRobotC:      newState.YourRobotPos.C,
				Robot2R:         newState.Robot2Pos.R,
				Robot2C:         newState.Robot2Pos.C,
				Robot1R:         newState.Robot1Pos.R,
				Robot1C:         newState.Robot1Pos.C,
				TargetCharIndex: newState.TargetCharIndex,
			}
			if !visited[key] {
				visited[key] = true
				queue.PushBack(newState)
				// 调试：打印新添加的 A-press 状态
				// fmt.Printf("  -> Added A-press state: Steps=%d, Path='%s', YR=(%d,%d)%c, R2=(%d,%d)%c, R1=(%d,%d)%c, TargetIdx=%d\n",
				// 	newState.TotalSteps, newState.Path,
				// 	newState.YourRobotPos.R, newState.YourRobotPos.C, getButtonAtPos(newState.YourRobotPos, directionalKeypad[:]),
				// 	newState.Robot2Pos.R, newState.Robot2Pos.C, getButtonAtPos(newState.Robot2Pos, directionalKeypad[:]),
				// 	newState.Robot1Pos.R, newState.Robot1Pos.C, getButtonAtPos(newState.Robot1Pos, numericKeypad[:]),
				// 	newState.TargetCharIndex)
			} else {
				// fmt.Printf("  -> A-press state already visited, skipping. Key: %+v\n", key) // 只有在需要时才打印这个
			}
		} else {
			// fmt.Printf("  -> Invalid A press, not adding new state. Path: %sA. (Final Valid: %t)\n",
			// 	currentState.Path, isThisAPressValid)
		}
	}

	// fmt.Printf("--- Path not found for code: %s ---\n", targetCode) // Debug info
	return -1 // 如果找不到路径
}

// CalculateComplexitySum 计算所有代码的总复杂度
func CalculateComplexitySum(codes []string) int {
	totalComplexity := 0
	for _, code := range codes {
		// 计算数字部分
		numericPartStr := ""
		for _, char := range code {
			if char >= '0' && char <= '9' {
				numericPartStr += string(char)
			}
		}
		numericPart := 0
		if numericPartStr != "" {
			var err error
			numericPart, err = strconv.Atoi(numericPartStr)
			if err != nil {
				fmt.Printf("Error converting numeric part '%s' to int: %v\n", numericPartStr, err)
				continue
			}
		}

		// 计算最短序列长度
		length := FindShortestSequenceLength(code) // 调用你实现的函数

		totalComplexity += length * numericPart
	}
	return totalComplexity
}

func main() {
	// 在 main 函数中处理输入文件
	data, err := os.ReadFile("input")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input file: %v\n", err)
		os.Exit(1)
	}

	// 假设每一行是一个代码
	lines := strings.Split(string(data), "\n")
	var codes []string
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			codes = append(codes, trimmedLine)
		}
	}

	sum := CalculateComplexitySum(codes)
	fmt.Println(sum)
}
