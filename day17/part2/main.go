package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"math" // 仅用于 math.Pow，但在这里可以用位移代替
)

var (
	targetProgramDigits []int
	foundDigits         [16]int // 存储找到的8进制位
)

// fastSolver recursively finds the base-8 digits of the initial A
// k is the current digit index (from 15 down to 0)
// a_k_plus_1 is the value of A / 8^(k+1), which is A for the next iteration if k was A_val_for_iter_k / 8
// More accurately, a_k_plus_1 is A after (15-k) divisions by 8 from A_initial_at_iter_k_plus_1
// Let's define a_k_plus_1 = A_val / 8, where A_val was for iteration k+1
// For solve(k, val_A_div_8_from_prev_iter):
// val_A_div_8_from_prev_iter is essentially A^(k+1) in our notation A^(j+1) = floor(A^(j)/8)
func fastSolver(k int, a_val_for_next_iter int) bool {
	if k < 0 { // All digits found
		return true
	}

	for dk := 0; dk <= 7; dk++ { // Try current digit d_k from 0 to 7
		// This A_k_current_iter is A^(k) = d_k + 8*A^(k+1)
		a_k_current_iter := dk + 8*a_val_for_next_iter

		// Simulate one loop of the specific program 2,4,1,3,7,5,0,3,1,5,4,4,5,5,3,0
		// IP 0: 2,4 (bst A -> B_val = A%8)
		// Since a_k_current_iter = dk + 8*a_val_for_next_iter, a_k_current_iter % 8 is dk
		bVal := dk // B after "bst A" (A%8)

		// IP 2: 1,3 (bxl B,3 -> B_val = B_val^3)
		bVal = bVal ^ 3

		// IP 4: 7,5 (cdv C -> C_val = A_k_current_iter / (2^B_val))
		denominatorC := 1
		if bVal >= 0 && bVal < 63 { // Avoid overflow for 1 << bVal if bVal is large or negative
			denominatorC = 1 << uint(bVal)
		} else if bVal < 0 { // 2 to a negative power results in 0 for int division unless A is also 0.
			denominatorC = 0 // This will cause division by zero if A != 0.
		} else { // bVal is too large, 2^bVal overflows int. Denominator effectively infinite.
			if a_k_current_iter != 0 { // If A is not zero, C becomes 0.
				denominatorC = math.MaxInt32 // Effectively makes C=0 unless A is also huge
			} else { // A is 0, C is 0.
				denominatorC = 1 // Avoid div by zero if A is 0
			}
		}

		cVal := 0
		if denominatorC != 0 {
			cVal = a_k_current_iter / denominatorC
		} else {
			// This case implies a problem state, likely not leading to a solution
			// Or indicates bVal was negative leading to fractional power.
			// The problem likely doesn't hit these edge cases for valid A.
			// If a_k_current_iter is non-zero, division by zero would error in real machine.
			// Let's assume this path doesn't lead to a solution.
			continue
		}

		// IP 6: 0,3 (adv A -> A_val_for_next_iter_calc = A_k_current_iter/8)
		// This is implicitly handled by passing `a_k_current_iter` to the next recursive call `solve(k-1, a_k_current_iter)`

		// IP 8: 1,5 (bxl B,5 -> B_val = B_val^5)
		bVal = bVal ^ 5 // bVal was B after (A%8)^3

		// IP 10: 4,4 (bxc B,C -> B_val = B_val^C_val)
		bVal = bVal ^ cVal

		// IP 12: 5,5 (out B -> output B_val%8)
		outputDigit := (bVal%8 + 8) % 8 // Ensure positive modulo

		if outputDigit == targetProgramDigits[k] {
			foundDigits[k] = dk
			if fastSolver(k-1, a_k_current_iter) { // Pass current A as A for next iter (it's already A_current_iter/8^0 effectively)
				return true
			}
		}
	}
	return false
}

// parseRegisterValue and input parsing (slightly simplified as we only need program string for solver)
func parseProgramStringFromInput(filePath string) (string, int, int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", 0, 0, fmt.Errorf("打开文件 %s 失败: %v", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var programStr string
	var initialA, initialB, initialC int // Read them but only B,C might be used if general
	lineNum := 0

	registerPrefixes := []string{"Register A:", "Register B:", "Register C:"}
	registerValues := []*int{&initialA, &initialB, &initialC}

	for i := 0; i < len(registerPrefixes); i++ {
		if !scanner.Scan() {
			return "", 0, 0, fmt.Errorf("解析文件失败: 读取第 %d 个寄存器行时遇到文件末尾", i+1)
		}
		lineNum++
		line := scanner.Text()
		// Simple parsing for prefix
		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 || !strings.HasPrefix(line, registerPrefixes[i]) {
			return "", 0, 0, fmt.Errorf("解析文件第 %d 行 (%s) 格式错误: %s", lineNum, registerPrefixes[i], line)
		}
		val, parseErr := strconv.Atoi(strings.TrimSpace(parts[1]))
		if parseErr != nil {
			return "", 0, 0, fmt.Errorf("解析文件第 %d 行 (%s) 数值错误: %v", lineNum, registerPrefixes[i], parseErr)
		}
		*registerValues[i] = val
	}

	foundProgramLine := false
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		} // Skip empty lines
		if strings.HasPrefix(line, "Program:") {
			programStr = strings.TrimSpace(strings.TrimPrefix(line, "Program:"))
			foundProgramLine = true
			break
		} else {
			return "", 0, 0, fmt.Errorf("解析文件第 %d 行失败: 期望空行或 'Program:', 实际为 '%s'", lineNum, line)
		}
	}
	if !foundProgramLine {
		return "", 0, 0, fmt.Errorf("解析文件失败: 未找到 'Program:' 行")
	}
	if err := scanner.Err(); err != nil {
		return "", 0, 0, fmt.Errorf("读取文件时发生错误: %v", err)
	}
	return programStr, initialB, initialC, nil
}

func main() {
	fmt.Println("Advent of Code Day 17 - Part Two (Fast Solver)")
	startTime := time.Now()

	programStrFromFile, _, _, err := parseProgramStringFromInput("input")
	if err != nil {
		fmt.Printf("无法解析输入文件: %v\n", err)
		return
	}

	fmt.Printf("目标程序字符串: %s\n", programStrFromFile)

	parts := strings.Split(programStrFromFile, ",")
	if len(parts) != 16 {
		fmt.Printf("错误: 预期程序长度为16, 实际为 %d\n", len(parts))
		return
	}
	targetProgramDigits = make([]int, 16)
	for i, p := range parts {
		val, convErr := strconv.Atoi(strings.TrimSpace(p))
		if convErr != nil || val < 0 || val > 7 {
			fmt.Printf("错误: 程序字符串包含无效数字 '%s'\n", p)
			return
		}
		targetProgramDigits[i] = val
	}

	// Initial call to the solver: solve for d_15 down to d_0.
	// The second argument to fastSolver is A^(k+1), so for k=15, A^(16) is 0.
	if fastSolver(15, 0) {
		resultA := int64(0)
		powerOf8 := int64(1)
		for i := 0; i < 16; i++ {
			resultA += int64(foundDigits[i]) * powerOf8
			if i < 15 { // Avoid overflow if powerOf8 becomes too large for the last multiplication
				if math.MaxInt64/8 < powerOf8 { // Check before multiplication
					if i < 15 && foundDigits[i+1] > 0 { // Only matters if higher digits exist
						fmt.Println("警告: powerOf8 可能溢出 int64")
						// This indicates A might be extremely large, check constraints.
						// For 8^15, int64 is fine. 8^15 approx 3.5e13. 8^20 > MaxInt64.
						// 8^16 is fine. 8^15 * 7 < MaxInt64.
					}
				}
				powerOf8 *= 8
			}
		}
		fmt.Printf("找到的最低正初始 A 值为: %d\n", resultA)
	} else {
		fmt.Println("未能找到解决方案。")
	}

	fmt.Printf("总耗时: %s\n", time.Since(startTime))
}
