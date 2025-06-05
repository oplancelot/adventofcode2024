package main

import (
	"strconv"
	"strings"
	"testing"
	// "fmt" // 用于调试
)

// 包级变量，供 fastSolverRecursiveHelperTesting 使用
// 在 solveProgramForSelfReplication 中会被正确设置和重置
var (
	testTargetProgramDigits []int
	testFoundDigits         [16]int
)

// fastSolverRecursiveHelperTesting 是核心递归逻辑的测试版本
// k 是当前数字的索引 (从15到0)
// a_k_plus_1 是前一个（更高位）迭代计算出的A值，代表 A_val / 8^(k+1)
func fastSolverRecursiveHelperTesting(k int, a_k_plus_1 int) bool {
	if k < 0 { // 所有数字都已找到
		return true
	}

	for dk := 0; dk <= 7; dk++ { // 尝试当前数字 d_k (0到7)
		// a_k_current_iter 是当前迭代的A值: A^(k) = d_k + 8*A^(k+1)
		// 这里使用参数 a_k_plus_1
		a_k_current_iter := dk + 8*a_k_plus_1

		// --- 开始模拟特定谜题程序的单次循环逻辑 ---
		// IP 0: 2,4 (bst A -> B_val = A%8)
		// 因为 a_k_current_iter = dk + 8*a_k_plus_1, 所以 a_k_current_iter % 8 等于 dk
		bVal := dk

		// IP 2: 1,3 (bxl B,3 -> B_val = B_val^3)
		bVal = bVal ^ 3

		// IP 4: 7,5 (cdv C -> C_val = a_k_current_iter / (2^B_val))
		denominatorC := 1
		// 小心处理2的幂的计算，特别是bVal为负或过大时
		if bVal >= 0 && bVal < 60 { // 限制bVal以避免 1 << uint(bVal) 溢出或行为异常
			denominatorC = 1 << uint(bVal)
		} else if bVal < 0 {
			// 2的负数次幂，对于整数除法，如果分子非零，结果为0。如果分母为0，会导致错误。
			// 稳健起见，如果bVal为负，使得2^bVal不是整数，分母视为“非常大”或导致C为0（除非A也为0）
			if a_k_current_iter == 0 {
				denominatorC = 1 // 0/1 = 0
			} else {
				denominatorC = (1 << 62) // 结果cVal将为0
			}
		} else { // bVal 过大
			if a_k_current_iter == 0 {
				denominatorC = 1
			} else {
				denominatorC = (1 << 62) // 结果cVal将为0
			}
		}

		cVal := 0
		if denominatorC != 0 {
			cVal = a_k_current_iter / denominatorC
		} else {
			// 理论上，如果 denominatorC 为 0，意味着 bVal 是一个导致 2^bVal=0 的特殊负值（不可能）
			// 或者 bVal < 0 使得 2^bVal 是小数，int(小数)=0。
			// 如果 a_k_current_iter 非零，这通常是一个无效路径。
			if a_k_current_iter != 0 {
				continue // 跳过此 dk
			}
			// 如果 a_k_current_iter 也为 0, cVal 为 0 是合理的
		}

		// IP 6: 0,3 (adv A -> A_val_for_next_iter_calc = a_k_current_iter/8)
		// 这个A值的变化是通过将 a_k_current_iter 传递给下一个递归调用来隐式处理的。

		// IP 8: 1,5 (bxl B,5 -> B_val = B_val^5)
		bVal = bVal ^ 5

		// IP 10: 4,4 (bxc B,C -> B_val = B_val^C_val)
		bVal = bVal ^ cVal

		// IP 12: 5,5 (out B -> output B_val%8)
		outputDigit := (bVal%8 + 8) % 8 // 确保模运算结果为正

		// --- 模拟结束 ---

		if outputDigit == testTargetProgramDigits[k] { // 检查输出是否与目标程序数字匹配
			testFoundDigits[k] = dk // 存储找到的数字
			// 递归到下一个数字 (k-1)，并传递当前的A值 (a_k_current_iter)
			// 作为下一个迭代的 A_k_plus_1 (即 A_k_current_iter 就是下一个迭代的 A_val_for_next_iter)
			if fastSolverRecursiveHelperTesting(k-1, a_k_current_iter) {
				return true
			}
		}
	}
	return false // 对于当前k和a_k_plus_1，没有找到合适的dk
}

// solveProgramForSelfReplication 是我们将要进行表格测试的函数
func solveProgramForSelfReplication(programToReplicateStr string) (finalA int64, success bool) {
	parts := strings.Split(programToReplicateStr, ",")
	if len(parts) != 16 {
		return 0, false // 程序长度必须为16
	}

	localTargetDigits := make([]int, 16)
	for i, p := range parts {
		val, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil || val < 0 || val > 7 {
			return 0, false // 无效的程序字符串数字
		}
		localTargetDigits[i] = val
	}

	testTargetProgramDigits = localTargetDigits // 设置供递归函数使用的包级变量
	testFoundDigits = [16]int{}                 // 重置包级变量

	if fastSolverRecursiveHelperTesting(15, 0) { // 初始调用：k=15, A^(16)=0
		resultA := int64(0)
		powerOf8 := int64(1)
		for i := 0; i < 16; i++ {
			resultA += int64(testFoundDigits[i]) * powerOf8
			if i < 15 {
				powerOf8 *= 8
			}
		}
		return resultA, true
	}
	return 0, false
}

// --- 表格驱动测试 ---
func TestSolveProgramSelfReplication_TableDriven(t *testing.T) {
	const actualPuzzleProgramString = "2,4,1,3,7,5,0,3,1,5,4,4,5,5,3,0" // 请确保这是您的实际程序字符串

	tests := []struct {
		name                  string
		programToReplicateStr string
		expectSuccess         bool
		expectedA             int64 // 仅当 expectSuccess 为 true 且此值非0时检查
	}{
		{
			name:                  "实际谜题程序字符串 (求解)",
			programToReplicateStr: actualPuzzleProgramString,
			expectSuccess:         true,
			expectedA:             0, // 初始设为0；在知道解之后，更新此值
		},
		{
			name:                  "无效程序长度 (过短)",
			programToReplicateStr: "1,2,3",
			expectSuccess:         false,
			expectedA:             0,
		},
		{
			name:                  "程序字符串包含非数字字符",
			programToReplicateStr: "2,4,1,3,7,5,X,3,1,5,4,4,5,5,3,0",
			expectSuccess:         false,
			expectedA:             0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualA, success := solveProgramForSelfReplication(tt.programToReplicateStr)

			if success != tt.expectSuccess {
				t.Errorf("solveProgramForSelfReplication() success = %v, want %v", success, tt.expectSuccess)
			}

			if success && tt.expectSuccess {
				if tt.expectedA != 0 {
					if actualA != tt.expectedA {
						t.Errorf("solveProgramForSelfReplication() found A = %d, want %d", actualA, tt.expectedA)
					} else {
						t.Logf("测试 '%s': 成功找到预期的 A = %d", tt.name, actualA)
					}
				} else {
					t.Logf("测试 '%s': 求解器按预期成功, 找到 A = %d (未断言特定A值)。Found digits: %v", tt.name, actualA, testFoundDigits)
				}
			} else if !success && tt.expectSuccess {
				t.Errorf("测试 '%s': 期望求解器成功，但它失败了。", tt.name)
			} else if success && !tt.expectSuccess {
				t.Errorf("测试 '%s': 期望求解器失败，但它成功了，A = %d。", tt.name, actualA)
			}
		})
	}
}

// (TestSimpleProgramExample_ReverseLogic 函数可以保持不变，因为它测试的是不同的逻辑)
func TestSimpleProgramExample_ReverseLogic(t *testing.T) {
	// ... (此函数与上一版本相同) ...
	tests := []struct {
		name            string
		programStr      string
		expectedA       int64
		expectFailParse bool
	}{
		{
			name:       "AoC Part 2 简单示例",
			programStr: "0,3,5,4,3,0",
			expectedA:  117440,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parts := strings.Split(tt.programStr, ",")
			s := make([]int, len(parts))
			validParse := true
			for i, p := range parts {
				val, err := strconv.Atoi(p)
				if err != nil || val < 0 || val > 7 {
					validParse = false
					break
				}
				s[i] = val
			}

			if tt.expectFailParse {
				if validParse {
					t.Errorf("期望程序字符串 '%s' 解析失败，但成功了", tt.programStr)
				}
				return
			}
			if !validParse {
				t.Fatalf("程序字符串 '%s' 解析失败，但测试用例未预期此情况", tt.programStr)
			}

			L := len(s)
			if L == 0 {
				if tt.expectedA == 0 {
					return
				}
				t.Fatalf("程序 '%s' 为空", tt.programStr)
			}
			// 对于简单示例，其反向逻辑依赖 s[L-1] == 0
			// if s[L-1] != 0 {
			// t.Logf("警告: 简单示例程序 '%s' 的最后一个数字不是0。反向逻辑可能不直接适用或有特定条件。", tt.programStr)
			// }

			A_val_internal := make([]int64, L+1)
			A_val_internal[L] = 0

			for k := L - 1; k >= 0; k-- {
				A_val_internal[k] = 8*A_val_internal[k+1] + int64(s[k])
			}
			calculatedA := 8 * A_val_internal[0]

			if calculatedA != tt.expectedA {
				t.Errorf("对于简单程序 '%s', 反向计算得到 A=%d, 期望 A=%d", tt.programStr, calculatedA, tt.expectedA)
			}
		})
	}
}
