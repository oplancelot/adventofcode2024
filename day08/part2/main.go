package main

import (
	"fmt"
	"os"
	"strings"
)

// Position 表示网格中的一个位置
type Position struct {
	row, col int
}

// Antenna 表示一个天线及其频率和位置
type Antenna struct {
	frequency string
	position  Position
}

// readInput 读取输入文件并构建字符网格和天线列表
func readInput(filename string) ([][]string, []Antenna, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	grid := make([][]string, len(lines))
	var antennas []Antenna

	for i, line := range lines {
		grid[i] = strings.Split(line, "")

		for j, char := range grid[i] {
			if char != "." {
				antenna := Antenna{
					frequency: char,
					position:  Position{row: i, col: j},
				}
				antennas = append(antennas, antenna)
			}
		}
	}

	return grid, antennas, nil
}

// 计算最大公约数
func gcd(a, b int) int {
	a, b = abs(a), abs(b)
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// 将向量简化为最简形式
func simplifyVector(dr, dc int) (int, int) {
	if dr == 0 && dc == 0 {
		return 0, 0
	}
	if dr == 0 {
		return 0, dc/abs(dc)
	}
	if dc == 0 {
		return dr/abs(dr), 0
	}
	
	g := gcd(abs(dr), abs(dc))
	return dr/g, dc/g
}

// 绝对值函数
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// findAntinodes 计算所有反节点位置并返回唯一位置的数量
func findAntinodes(grid [][]string, antennas []Antenna) int {
	// 使用map跟踪唯一的反节点位置
	antinodes := make(map[Position]bool)
	
	// 按频率对天线进行分组
	frequencyGroups := make(map[string][]Antenna)
	for _, antenna := range antennas {
		frequencyGroups[antenna.frequency] = append(frequencyGroups[antenna.frequency], antenna)
	}
	
	// 对于每个频率组，找出所有成一直线的点
	for _, antennaGroup := range frequencyGroups {
		// 需要至少2个相同频率的天线才能形成反节点
		if len(antennaGroup) < 2 {
			continue
		}
		
		// 对于每对天线，找出它们之间和延长线上的所有点
		for i := 0; i < len(antennaGroup); i++ {
			for j := i + 1; j < len(antennaGroup); j++ {
				a1 := antennaGroup[i]
				a2 := antennaGroup[j]
				
				// 计算方向向量并简化
				dr := a2.position.row - a1.position.row
				dc := a2.position.col - a1.position.col
				
				// 简化方向向量
				dr, dc = simplifyVector(dr, dc)
				
				// 从a1开始，沿着方向向量移动，直到到达a2
				// 所有这些点都是反节点
				for k := 0; ; k++ {
					pos := Position{
						row: a1.position.row + k*dr,
						col: a1.position.col + k*dc,
					}
					
					if pos.row == a2.position.row && pos.col == a2.position.col {
						antinodes[pos] = true
						break
					}
					
					if isWithinBounds(grid, pos) {
						antinodes[pos] = true
					}
				}
				
				// 继续沿着方向向量移动，直到离开网格
				// 向a2方向移动
				for k := 1; ; k++ {
					pos := Position{
						row: a2.position.row + k*dr,
						col: a2.position.col + k*dc,
					}
					
					if !isWithinBounds(grid, pos) {
						break
					}
					
					antinodes[pos] = true
				}
				
				// 向a1反方向移动
				for k := 1; ; k++ {
					pos := Position{
						row: a1.position.row - k*dr,
						col: a1.position.col - k*dc,
					}
					
					if !isWithinBounds(grid, pos) {
						break
					}
					
					antinodes[pos] = true
				}
			}
		}
	}
	
	return len(antinodes)
}

// isWithinBounds 检查位置是否在网格边界内
func isWithinBounds(grid [][]string, pos Position) bool {
	return pos.row >= 0 && pos.row < len(grid) && pos.col >= 0 && pos.col < len(grid[0])
}

func main() {
	const inputFile = "input"
	grid, antennas, err := readInput(inputFile)
	if err != nil {
		fmt.Printf("读取输入文件失败 (%s): %v\n", inputFile, err)
		return
	}

	fmt.Printf("读取到 %d 个天线\n", len(antennas))

	// 打印一些天线信息以验证输入
	if len(antennas) > 0 {
		fmt.Printf("第一个天线: 频率=%s, 位置=(%d,%d)\n",
			antennas[0].frequency, antennas[0].position.row, antennas[0].position.col)
	}

	totalUniqueLocations := findAntinodes(grid, antennas)
	fmt.Printf("地图边界内包含反节点的唯一位置数量: %d\n", totalUniqueLocations)
}
