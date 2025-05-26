package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

// Robot 结构体表示一个机器人的位置和速度
type Robot struct {
	Px int // 初始X坐标
	Py int // 初始Y坐标
	Vx int // X轴速度
	Vy int // Y轴速度
}

// ParseInput 函数解析多行字符串输入，将其转换为 Robot 结构体切片。
func ParseInput(input string) ([]Robot, error) {
	var robots []Robot
	// TrimSpace 移除输入字符串首尾的空白，Split 通过换行符分割成行
	lines := strings.Split(strings.TrimSpace(input), "\n")

	for _, line := range lines {
		// 分割 "p=x,y v=x,y" 这样的行
		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			// 简单的检查，如果行不包含 "v="，则跳过（处理可能存在的末尾截断行）
			if !strings.Contains(line, "v=") {
				continue
			}
			return nil, fmt.Errorf("格式错误的行: %s", line)
		}

		// 提取位置和速度字符串，并移除前缀 "p=" 和 "v="
		pStr := strings.TrimPrefix(parts[0], "p=")
		vStr := strings.TrimPrefix(parts[1], "v=")

		// 分割坐标 "x,y"
		pCoords := strings.Split(pStr, ",")
		vCoords := strings.Split(vStr, ",")

		if len(pCoords) != 2 || len(vCoords) != 2 {
			return nil, fmt.Errorf("行中坐标格式错误: %s", line)
		}

		// 将字符串坐标转换为整数
		px, err := strconv.Atoi(pCoords[0])
		if err != nil {
			return nil, err
		}
		py, err := strconv.Atoi(pCoords[1])
		if err != nil {
			return nil, err
		}
		vx, err := strconv.Atoi(vCoords[0])
		if err != nil {
			return nil, err
		}
		vy, err := strconv.Atoi(vCoords[1])
		if err != nil {
			return nil, err
		}

		// 将解析出的机器人添加到切片中
		robots = append(robots, Robot{Px: px, Py: py, Vx: vx, Vy: vy})
	}
	return robots, nil
}

// mod 函数用于正确处理 Go 语言中负数的取模运算。
// Go 的 % 运算符在被除数为负时可能返回负结果。
func mod(a, n int) int {
	return (a%n + n) % n
}

// CalculateSafetyFactor 函数模拟机器人移动并计算安全系数。
func CalculateSafetyFactor(robots []Robot, width, height, simulationTime int) int {
	// 使用 map 统计每个最终位置的机器人数量
	finalPositions := make(map[[2]int]int) // [2]int 用作 x, y 坐标的键

	for _, robot := range robots {
		// 计算机器人经过 simulationTime 秒后的最终位置，并应用环绕效果
		finalX := mod(robot.Px+robot.Vx*simulationTime, width)
		finalY := mod(robot.Py+robot.Vy*simulationTime, height)
		// 增加该位置的机器人计数
		finalPositions[[2]int{finalX, finalY}]++
	}

	// 计算空间中线的位置
	midX := width / 2
	midY := height / 2

	// 初始化四个象限的机器人计数：[左上, 右上, 左下, 右下]
	quadrantCounts := [4]int{0, 0, 0, 0}

	for pos, count := range finalPositions {
		x, y := pos[0], pos[1]

		// 机器人如果正好在中间线上，则不计入任何象限
		if x == midX || y == midY {
			continue
		}

		// 判断象限并增加计数
		if x < midX && y < midY { // 左上象限
			quadrantCounts[0] += count
		} else if x > midX && y < midY { // 右上象限
			quadrantCounts[1] += count
		} else if x < midX && y > midY { // 左下象限
			quadrantCounts[2] += count
		} else if x > midX && y > midY { // 右下象限
			quadrantCounts[3] += count
		}
	}

	// 计算安全系数（四个象限计数的乘积）
	safetyFactor := 1
	for _, count := range quadrantCounts {
		safetyFactor *= count
	}

	return safetyFactor
}

func main() {
	// 定义主问题所需的空间尺寸和模拟时间
	const WIDTH = 101
	const HEIGHT = 103
	const SIMULATION_TIME = 100

	// 从 input.txt 文件中读取输入数据
	inputBytes, err := ioutil.ReadFile("input")
	if err != nil {
		// 如果读取文件失败，则记录错误并终止程序
		log.Fatalf("无法读取 input.txt 文件: %v", err)
	}
	inputData := string(inputBytes) // 将字节切片转换为字符串

	// 解析机器人数据
	robots, err := ParseInput(inputData)
	if err != nil {
		// 如果解析输入失败，则记录错误并终止程序
		log.Fatalf("解析机器人输入失败: %v", err)
	}

	// 计算安全系数
	safetyFactor := CalculateSafetyFactor(robots, WIDTH, HEIGHT, SIMULATION_TIME)

	// 打印最终的安全系数
	fmt.Println("安全系数是:", safetyFactor)
}
