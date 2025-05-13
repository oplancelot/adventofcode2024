package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"os"
	"strings"

	"github.com/fogleman/gg"
)

type Position struct {
	row, col int
}

type State struct {
	pos Position
	dir string
}

const (
	cellSize   = 20
	frameDelay = 5 // 单位是 10ms，5 表示每帧 50ms
)

func readInput(filename string) ([][]string, Position, string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, Position{}, "", err
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	grid := make([][]string, len(lines))

	var guardPos Position
	var guardDir string

	for i, line := range lines {
		line = strings.TrimSuffix(line, "\r")
		grid[i] = strings.Split(line, "")
		for j, char := range grid[i] {
			if char == "^" || char == ">" || char == "v" || char == "<" {
				guardPos = Position{i, j}
				guardDir = char
				grid[i][j] = "."
			}
		}
	}
	return grid, guardPos, guardDir, nil
}

func getNextPosition(pos Position, dir string) Position {
	switch dir {
	case "^":
		return Position{pos.row - 1, pos.col}
	case ">":
		return Position{pos.row, pos.col + 1}
	case "v":
		return Position{pos.row + 1, pos.col}
	case "<":
		return Position{pos.row, pos.col - 1}
	}
	return pos
}

func turnRight(dir string) string {
	switch dir {
	case "^":
		return ">"
	case ">":
		return "v"
	case "v":
		return "<"
	case "<":
		return "^"
	}
	return dir
}

func isInBounds(grid [][]string, pos Position) bool {
	return pos.row >= 0 && pos.row < len(grid) && pos.col >= 0 && pos.col < len(grid[0])
}

// 追踪路径并返回每一步状态
func tracePath(grid [][]string, startPos Position, startDir string) []State {
	var path []State
	pos := startPos
	dir := startDir
	path = append(path, State{pos, dir})

	visited := make(map[Position]bool)
	visited[pos] = true

	for {
		next := getNextPosition(pos, dir)
		if !isInBounds(grid, next) {
			break
		}
		if grid[next.row][next.col] == "#" {
			dir = turnRight(dir)
		} else {
			pos = next
			if visited[pos] {
				break // 避免循环（也可选用 step limit）
			}
			visited[pos] = true
		}
		path = append(path, State{pos, dir})
	}
	return path
}

// 渲染单帧
func renderFrame(grid [][]string, path []Position, guardPos Position) *gg.Context {
	rows := len(grid)
	cols := len(grid[0])
	dc := gg.NewContext(cols*cellSize, rows*cellSize)

	// 背景
	dc.SetColor(color.White)
	dc.Clear()

	// 画网格
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			x := float64(c * cellSize)
			y := float64(r * cellSize)

			switch grid[r][c] {
			case "#":
				dc.SetColor(color.Black)
			default:
				dc.SetColor(color.RGBA{230, 230, 230, 255})
			}
			dc.DrawRectangle(x, y, cellSize, cellSize)
			dc.Fill()
		}
	}

	// 画路径轨迹
	dc.SetColor(color.RGBA{0, 128, 255, 255})
	for _, p := range path {
		x := float64(p.col * cellSize)
		y := float64(p.row * cellSize)
		dc.DrawCircle(x+cellSize/2, y+cellSize/2, cellSize/4)
		dc.Fill()
	}

	// 画警卫当前位置
	dc.SetColor(color.RGBA{255, 0, 0, 255})
	x := float64(guardPos.col * cellSize)
	y := float64(guardPos.row * cellSize)
	dc.DrawCircle(x+cellSize/2, y+cellSize/2, float64(cellSize/2))
	dc.Fill()

	return dc
}

func saveGIF(frames []*gg.Context, filename string) error {
	var images []*image.Paletted
	var delays []int

	for _, dc := range frames {
		img := dc.Image()
		palettedImg := image.NewPaletted(img.Bounds(), []color.Color{
			color.White, color.Black,
			color.RGBA{230, 230, 230, 255},
			color.RGBA{0, 128, 255, 255},
			color.RGBA{255, 0, 0, 255},
		})
		// Draw to paletted
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
				palettedImg.Set(x, y, img.At(x, y))
			}
		}
		images = append(images, palettedImg)
		delays = append(delays, frameDelay)
	}

	anim := gif.GIF{
		Image: images,
		Delay: delays,
	}
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	return gif.EncodeAll(outFile, &anim)
}

func main() {
	grid, guardPos, guardDir, err := readInput("input")
	if err != nil {
		fmt.Println("Failed to read input:", err)
		return
	}

	pathStates := tracePath(grid, guardPos, guardDir)

	var frames []*gg.Context
	var pathTrace []Position
	for _, state := range pathStates {
		pathTrace = append(pathTrace, state.pos)
		frame := renderFrame(grid, pathTrace, state.pos)
		frames = append(frames, frame)
	}

	err = saveGIF(frames, "guard_path.gif")
	if err != nil {
		fmt.Println("Failed to save GIF:", err)
		return
	}

	fmt.Printf("GIF saved as guard_path.gif with %d frames.\n", len(frames))
}
