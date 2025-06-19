package main

import (
	_ "embed"
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed shaders/electrical_storm.kage
var electricalStormShader []byte

type Gameboard struct {
	Width      int
	Height     int
	X          int
	Y          int
	shader     *ebiten.Shader
	startTime  time.Time
	baseWidth  int
	baseHeight int
}

func NewGameboard(baseWidth, baseHeight int) *Gameboard {
	shader, err := ebiten.NewShader(electricalStormShader)
	if err != nil {
		panic(fmt.Sprintf("Failed to compile electrical storm shader: %v", err))
	}

	return &Gameboard{
		baseWidth:  baseWidth,
		baseHeight: baseHeight,
		Width:      baseWidth,
		Height:     baseHeight,
		shader:     shader,
		startTime:  time.Now(),
	}
}

func (gb *Gameboard) UpdateScale(screenWidth, screenHeight int) {
	scaleX := float64(screenWidth) / 320.0
	scaleY := float64(screenHeight) / 320.0
	scale := min(scaleX, scaleY)

	gb.Width = int(float64(gb.baseWidth) * scale)
	gb.Height = int(float64(gb.baseHeight) * scale)

	gb.X = (screenWidth - gb.Width) / 2
	gb.Y = 0
}

func (gb *Gameboard) Draw(screen *ebiten.Image) {
	gameboardImage := ebiten.NewImage(gb.Width, gb.Height)

	elapsed := time.Since(gb.startTime).Seconds()

	op := &ebiten.DrawTrianglesShaderOptions{}
	op.Uniforms = map[string]interface{}{
		"Time":       float32(elapsed),
		"Resolution": []float32{float32(gb.Width), float32(gb.Height)},
	}

	vertices := []ebiten.Vertex{
		{DstX: 0, DstY: 0, SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: float32(gb.Width), DstY: 0, SrcX: float32(gb.Width), SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: 0, DstY: float32(gb.Height), SrcX: 0, SrcY: float32(gb.Height), ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: float32(gb.Width), DstY: float32(gb.Height), SrcX: float32(gb.Width), SrcY: float32(gb.Height), ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
	}

	indices := []uint16{0, 1, 2, 1, 2, 3}

	gameboardImage.DrawTrianglesShader(vertices, indices, gb.shader, op)

	drawOp := &ebiten.DrawImageOptions{}
	drawOp.GeoM.Translate(float64(gb.X), float64(gb.Y))
	screen.DrawImage(gameboardImage, drawOp)
}

func (gb *Gameboard) GetBounds() (x, y, width, height int) {
	return gb.X, gb.Y, gb.Width, gb.Height
}

func (gb *Gameboard) Contains(x, y int) bool {
	return x >= gb.X && x < gb.X+gb.Width && y >= gb.Y && y < gb.Y+gb.Height
}

func (gb *Gameboard) ToGridCoordinates(screenX, screenY int, blockSize float64) (gridX, gridY int) {
	relativeX := screenX - gb.X
	relativeY := screenY - gb.Y
	gridX = int(float64(relativeX) / blockSize)
	gridY = int(float64(relativeY) / blockSize)
	return
}

func (gb *Gameboard) ToScreenCoordinates(gridX, gridY int, blockSize float64) (screenX, screenY int) {
	screenX = gb.X + int(float64(gridX)*blockSize)
	screenY = gb.Y + int(float64(gridY)*blockSize)
	return
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
