package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

// Gate 结构体保持不变
type Gate struct{ input1, input2, op, output string }

// 这是 Rust 代码逻辑的 Go 语言实现
func solvePart2_staticAnalysis(input string) string {
	// 1. 解析电路，我们只需要门电路的连接信息
	gates, _, _ := parseGates(input)

	// 2. 创建一个“连接缓存”，用于快速查找“某个导线是否是某个类型门的输入”
	// connCacheKey 是 map 的键类型
	type connCacheKey struct {
		wire string
		op   string
	}
	connectionCache := make(map[connCacheKey]bool)
	for _, g := range gates {
		connectionCache[connCacheKey{g.input1, g.op}] = true
		connectionCache[connCacheKey{g.input2, g.op}] = true
	}

	// 3. 遍历所有门，根据启发式规则过滤出“可疑”的导线
	var suspiciousWires []string
	for _, g := range gates {
		isSuspicious := false
		switch g.op {
		case "AND":
			// 规则: 一个 AND 门的输出没有被用作任何 OR 门的输入
			// (并且排除了一个 x00 边界情况)
			if g.input1 != "x00" && g.input2 != "x00" && !connectionCache[connCacheKey{g.output, "OR"}] {
				isSuspicious = true
			}
		case "XOR":
			// 规则 1: 一个由 x/y 输入构成的 XOR 门，其输出没有被用作另一个 XOR 门的输入
			cond1 := (strings.HasPrefix(g.input1, "x") || strings.HasPrefix(g.input2, "x")) &&
				(g.input1 != "x00" && g.input2 != "x00" && !connectionCache[connCacheKey{g.output, "XOR"}])

			// 规则 2: 一个由非 x/y 输入构成的 XOR 门，其输出不是 z##
			cond2 := !strings.HasPrefix(g.output, "z") && !strings.HasPrefix(g.input1, "x") && !strings.HasPrefix(g.input2, "x")

			if cond1 || cond2 {
				isSuspicious = true
			}
		case "OR":
			// 规则: 一个 OR 门的输出直接连接到了 z## (最终的和)
			// (并且排除了 z45 边界情况)
			if strings.HasPrefix(g.output, "z") && g.output != "z45" {
				isSuspicious = true
			}
		}

		if isSuspicious {
			suspiciousWires = append(suspiciousWires, g.output)
		}
	}

	// 4. 排序并格式化输出
	sort.Strings(suspiciousWires)
	return strings.Join(suspiciousWires, ",")
}

// 解析器: 只解析门电路部分，因为我们不再需要初始值
func parseGates(input string) ([]Gate, []string, []string) {
	parts := strings.Split(strings.TrimSpace(input), "\n\n")
	var initialX, initialY []string // 虽然不用，但为了复用先保留
	if len(parts) == 0 {
		return nil, nil, nil
	}

	// 解析初始导线 (虽然此方法不用，但保留以防万一)
	scanner := bufio.NewScanner(strings.NewReader(parts[0]))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		wireName := strings.Split(line, ": ")[0]
		if strings.HasPrefix(wireName, "x") {
			initialX = append(initialX, wireName)
		} else if strings.HasPrefix(wireName, "y") {
			initialY = append(initialY, wireName)
		}
	}

	var gates []Gate
	gateRegex := regexp.MustCompile(`(\w+) (AND|OR|XOR) (\w+) -> (\w+)`)
	gateSection := parts[0]
	if len(parts) > 1 {
		gateSection = parts[1]
	}

	scanner = bufio.NewScanner(strings.NewReader(gateSection))
	for scanner.Scan() {
		matches := gateRegex.FindStringSubmatch(scanner.Text())
		if len(matches) == 5 {
			gates = append(gates, Gate{input1: matches[1], input2: matches[3], op: matches[2], output: matches[4]})
		}
	}
	return gates, initialX, initialY
}

func main() {
	data, err := os.ReadFile("input")
	if err != nil {
		log.Fatalf("读取输入文件失败: %v", err)
	}
	// 调用新的基于静态分析的求解器
	result := solvePart2_staticAnalysis(string(data))
	fmt.Printf("排序并用逗号连接的导线是: %s\n", result)
}
