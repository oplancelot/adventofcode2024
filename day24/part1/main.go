package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Gate 代表一个逻辑门及其连接和操作。
type Gate struct {
	input1, input2 string
	op             string
	output         string
}

// solve 模拟电路并返回十进制输出。
func solve(input string) int64 {
	wires, gates := parseInput(input)

	// 模拟电路，直到所有门都触发。
	for len(gates) > 0 {
		remainingGates := []Gate{}
		progressMade := false
		for _, gate := range gates {
			val1, ok1 := wires[gate.input1]
			val2, ok2 := wires[gate.input2]

			if ok1 && ok2 {
				var result int
				switch gate.op {
				case "AND":
					result = val1 & val2
				case "OR":
					result = val1 | val2
				case "XOR":
					result = val1 ^ val2
				}
				wires[gate.output] = result
				progressMade = true
			} else {
				remainingGates = append(remainingGates, gate)
			}
		}
		gates = remainingGates
		// 这个检查可以防止无限循环，尽管题目保证了这是一个DAG。
		if !progressMade && len(gates) > 0 {
			log.Fatalf("模拟停滞，仍有 %d 个门未计算", len(gates))
		}
	}

	// 从 'z' 导线计算最终数字。
	return calculateOutput(wires)
}

// parseInput 解析原始输入字符串，得到初始导线值和门列表。
func parseInput(input string) (map[string]int, []Gate) {
	parts := strings.Split(strings.TrimSpace(input), "\n\n")
	if len(parts) != 2 {
		log.Fatalf("输入格式无效：应由一个空行分隔成两部分")
	}

	wires := make(map[string]int)
	scanner := bufio.NewScanner(strings.NewReader(parts[0]))
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ": ")
		val, _ := strconv.Atoi(line[1])
		wires[line[0]] = val
	}

	var gates []Gate
	gateRegex := regexp.MustCompile(`(\w+) (AND|OR|XOR) (\w+) -> (\w+)`)
	scanner = bufio.NewScanner(strings.NewReader(parts[1]))
	for scanner.Scan() {
		matches := gateRegex.FindStringSubmatch(scanner.Text())
		if len(matches) == 5 {
			gates = append(gates, Gate{
				input1: matches[1],
				input2: matches[3],
				op:     matches[2],
				output: matches[4],
			})
		}
	}
	return wires, gates
}

// calculateOutput 从 'z' 导线计算十进制值。
func calculateOutput(wires map[string]int) int64 {
	type zWire struct {
		name  string
		index int
		value int
	}

	var zWires []zWire
	for name, value := range wires {
		if strings.HasPrefix(name, "z") {
			index, err := strconv.Atoi(name[1:])
			if err == nil {
				zWires = append(zWires, zWire{name: name, index: index, value: value})
			}
		}
	}

	// 按索引排序 (z00, z01, z02, ...)
	sort.Slice(zWires, func(i, j int) bool {
		return zWires[i].index < zWires[j].index
	})

	// 构建二进制字符串（从最高索引开始，以保证正确的比特顺序）
	var binaryString strings.Builder
	for i := len(zWires) - 1; i >= 0; i-- {
		binaryString.WriteString(strconv.Itoa(zWires[i].value))
	}

	binaryStr := binaryString.String()
	if binaryStr == "" {
		return 0
	}

	// 将二进制字符串转换为十进制
	result, err := strconv.ParseInt(binaryStr, 2, 64)
	if err != nil {
		log.Fatalf("解析二进制字符串 '%s' 失败: %v", binaryStr, err)
	}

	return result
}

func main() {
	// 你的谜题输入文件名是 "input"
	data, err := os.ReadFile("input")
	if err != nil {
		log.Fatalf("读取输入文件失败: %v", err)
	}
	result := solve(string(data))
	fmt.Printf("最终输出的十进制数是: %d\n", result)
}
