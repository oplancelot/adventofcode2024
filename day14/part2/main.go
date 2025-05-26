package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

// Num 是用于坐标、速度和时间步的整数类型。
type Num int32

// --- 数据结构 ---

// Robot 存储机器人的初始位置和速度。
type Robot struct {
	Px Num
	Py Num
	Vx Num
	Vy Num
}

// --- 辅助函数 ---

// parseNumberAndAdvance 从字节切片中解析一个数字，并将索引向前推进。
func parseNumberAndAdvance(dataBytes []byte, originalIdx *int, endByte byte) (Num, error) { // 参数名改为 dataBytes
	var number Num = 0
	idx := *originalIdx
	var sign Num = 1

	if idx >= len(dataBytes) { // 使用 dataBytes
		return 0, fmt.Errorf("unexpected end of bytes at index %d", idx)
	}

	if dataBytes[idx] == '-' { // 使用 dataBytes
		sign = -1
		idx++
	}

	startNumIdx := idx
	for idx < len(dataBytes) { // 使用 dataBytes
		value := dataBytes[idx] // 使用 dataBytes
		if value == endByte {
			break
		}
		if value < '0' || value > '9' {
			return 0, fmt.Errorf("invalid character '%c' at index %d", value, idx)
		}
		number = number*10 + Num(value-'0')
		idx++
	}

	if idx == startNumIdx && sign == 1 {
		return 0, fmt.Errorf("no number found starting at index %d", startNumIdx)
	}

	*originalIdx = idx + 1
	return number * sign, nil
}

// mod 函数实现 Go 语言中正确的取模运算，确保结果在 [0, n-1) 范围内。
func mod(a, n Num) Num {
	return (a%n + n) % n
}

// --- 数学工具：用于中国剩余定理 ---

// extendedGCD 实现了扩展欧几里得算法，计算 ax + by = gcd(a,b)。
func extendedGCD(a, b Num) (Num, Num, Num) {
	if a == 0 {
		return b, 0, 1
	}
	gcd, x1, y1 := extendedGCD(b%a, a)
	x := y1 - (b/a)*x1
	y := x1
	return gcd, x, y
}

// modInverse 计算 a 的模 m 逆元，即 (a * x) % m = 1。
func modInverse(a, m Num) (Num, error) {
	gcd, x, _ := extendedGCD(a, m)
	if gcd != 1 {
		return 0, fmt.Errorf("%d has no modular inverse modulo %d", a, m)
	}
	return (x%m + m) % m, nil
}

// --- 主要逻辑 ---

// main 函数是程序的入口点，负责解决 Part 2 问题。
func main() {
	// 定义空间尺寸和机器人总数
	const WIDTH Num = 101
	const HEIGHT Num = 103
	const N_ROBOTS = 500

	// 1. 读取并解析输入数据
	inputBytes, err := ioutil.ReadFile("input") // inputBytes 在这里被定义为 []byte
	if err != nil {
		log.Fatalf("无法读取 input 文件: %v", err)
	}
	// inputData := string(inputBytes) // 字符串形式，这里不需要了，因为直接用 inputBytes

	var robots []Robot
	idx := 0
	bytesLen := len(inputBytes) // 使用 inputBytes 的长度

	for {
		if idx >= bytesLen {
			break
		}

		idx += 2 // 跳过 "p="
		// *** 关键修正: 将 'bytes' 替换为 'inputBytes' ***
		pX, err := parseNumberAndAdvance(inputBytes, &idx, ',')
		if err != nil {
			log.Fatalf("解析 pX 失败: %v", err)
		}
		pY, err := parseNumberAndAdvance(inputBytes, &idx, ' ')
		if err != nil {
			log.Fatalf("解析 pY 失败: %v", err)
		}

		idx += 2 // 跳过 "v="
		// *** 关键修正: 将 'bytes' 替换为 'inputBytes' ***
		vX, err := parseNumberAndAdvance(inputBytes, &idx, ',')
		if err != nil {
			log.Fatalf("解析 vX 失败: %v", err)
		}
		// *** 关键修正: 将 'bytes' 替换为 'inputBytes' ***
		vY, err := parseNumberAndAdvance(inputBytes, &idx, '\n')
		if err != nil {
			if idx == bytesLen {
				// 兼容文件末尾
			} else {
				log.Fatalf("解析 vY 失败: %v", err)
			}
		}

		robots = append(robots, Robot{Px: pX, Py: pY, Vx: vX, Vy: vY})
	}
	if len(robots) != N_ROBOTS {
		log.Fatalf("解析到的机器人数量不符，期望 %d，得到 %d", N_ROBOTS, len(robots))
	}

	fmt.Println("正在计算 X 轴和 Y 轴最聚集的时间点...")

	// 2. 计算 X 轴最聚集的时间点 (mostClusteredXIteration)
	var mostClusteredXIteration Num = 0
	biggestXCluster := 0

	xPositions := make([]Num, N_ROBOTS)
	for i := 0; i < N_ROBOTS; i++ {
		xPositions[i] = robots[i].Px
	}

	for t := Num(0); t < WIDTH; t++ {
		counts := make([]int, WIDTH)
		for _, pos := range xPositions {
			counts[pos]++
		}

		maxCount := 0
		for _, count := range counts {
			if count > maxCount {
				maxCount = count
			}
		}

		if maxCount > biggestXCluster {
			biggestXCluster = maxCount
			mostClusteredXIteration = t
		}

		for j := 0; j < N_ROBOTS; j++ {
			xPositions[j] = mod(xPositions[j]+robots[j].Vx, WIDTH)
		}
	}

	// 3. 计算 Y 轴最聚集的时间点 (mostClusteredYIteration)
	var mostClusteredYIteration Num = 0
	biggestYCluster := 0

	yPositions := make([]Num, N_ROBOTS)
	for i := 0; i < N_ROBOTS; i++ {
		yPositions[i] = robots[i].Py
	}

	for t := Num(0); t < HEIGHT; t++ {
		counts := make([]int, HEIGHT)
		for _, pos := range yPositions {
			counts[pos]++
		}

		maxCount := 0
		for _, count := range counts {
			if count > maxCount {
				maxCount = count
			}
		}

		if maxCount > biggestYCluster {
			biggestYCluster = maxCount
			mostClusteredYIteration = t
		}

		for j := 0; j < N_ROBOTS; j++ {
			yPositions[j] = mod(yPositions[j]+robots[j].Vy, HEIGHT)
		}
	}

	// 4. 使用中国剩余定理合成最终时间
	invMod, err := modInverse(WIDTH, HEIGHT)
	if err != nil {
		log.Fatalf("无法计算模逆元: %v", err)
	}

	diff := mostClusteredYIteration - mostClusteredXIteration
	k := mod(diff*invMod, HEIGHT)

	finalT := mostClusteredXIteration + k*WIDTH

	fmt.Printf("X 轴最聚集时间点: %d (最多 %d 个机器人)\n", mostClusteredXIteration, biggestXCluster)
	fmt.Printf("Y 轴最聚集时间点: %d (最多 %d 个机器人)\n", mostClusteredYIteration, biggestYCluster)
	fmt.Printf("Part 2 圣诞树图案时间 (由中国剩余定理计算): %d\n", finalT)
}
