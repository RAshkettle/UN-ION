package main

import (
	_ "embed"
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed shaders/electrical_storm.kage
var electricalStormShader []byte

// Gameboard represents the main playing area for Tetris pieces
type Gameboard struct {
	Width       int
	Height      int
	X           int // X position on screen
	Y           int // Y position on screen
	shader      *ebiten.Shader
	startTime   time.Time
	baseWidth   int // Original design width (192)
	baseHeight  int // Original design height (320)
}

// NewGameboard creates a new gameboard with the specified dimensions and position
func NewGameboard(baseWidth, baseHeight int) *Gameboard {
	// Compile the electrical storm shader
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

// UpdateScale updates the gameboard size and position based on screen dimensions
func (gb *Gameboard) UpdateScale(screenWidth, screenHeight int) {
	// Calculate scale factor to maintain aspect ratio
	scaleX := float64(screenWidth) / 320.0  // 320 is our base screen width
	scaleY := float64(screenHeight) / 320.0 // 320 is our base screen height
	scale := min(scaleX, scaleY) // Use smaller scale to fit both dimensions
	
	// Apply scale to gameboard dimensions
	gb.Width = int(float64(gb.baseWidth) * scale)
	gb.Height = int(float64(gb.baseHeight) * scale)
	
	// Center the gameboard horizontally
	gb.X = (screenWidth - gb.Width) / 2
	gb.Y = 0 // Keep at top of screen
}

// Draw renders the gameboard on the screen with shader effect
func (gb *Gameboard) Draw(screen *ebiten.Image) {
	// Create a temporary image for the gameboard
	gameboardImage := ebiten.NewImage(gb.Width, gb.Height)
	
	// Calculate time for shader animation
	elapsed := time.Since(gb.startTime).Seconds()
	
	// Apply the electrical storm shader
	op := &ebiten.DrawTrianglesShaderOptions{}
	op.Uniforms = map[string]interface{}{
		"Time":       float32(elapsed),
		"Resolution": []float32{float32(gb.Width), float32(gb.Height)},
	}
	
	// Create vertices for a full-screen quad
	vertices := []ebiten.Vertex{
		{DstX: 0, DstY: 0, SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: float32(gb.Width), DstY: 0, SrcX: float32(gb.Width), SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: 0, DstY: float32(gb.Height), SrcX: 0, SrcY: float32(gb.Height), ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: float32(gb.Width), DstY: float32(gb.Height), SrcX: float32(gb.Width), SrcY: float32(gb.Height), ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
	}
	
	indices := []uint16{0, 1, 2, 1, 2, 3}
	
	gameboardImage.DrawTrianglesShader(vertices, indices, gb.shader, op)
	
	// Draw the gameboard at its position on the screen
	drawOp := &ebiten.DrawImageOptions{}
	drawOp.GeoM.Translate(float64(gb.X), float64(gb.Y))
	screen.DrawImage(gameboardImage, drawOp)
}

// GetBounds returns the gameboard boundaries
func (gb *Gameboard) GetBounds() (x, y, width, height int) {
	return gb.X, gb.Y, gb.Width, gb.Height
}

// Contains checks if a point is within the gameboard
func (gb *Gameboard) Contains(x, y int) bool {
	return x >= gb.X && x < gb.X+gb.Width && y >= gb.Y && y < gb.Y+gb.Height
}

// ToGridCoordinates converts screen coordinates to grid coordinates
// assuming each grid cell is blockSize pixels
func (gb *Gameboard) ToGridCoordinates(screenX, screenY int, blockSize float64) (gridX, gridY int) {
	relativeX := screenX - gb.X
	relativeY := screenY - gb.Y
	gridX = int(float64(relativeX) / blockSize)
	gridY = int(float64(relativeY) / blockSize)
	return
}

// ToScreenCoordinates converts grid coordinates to screen coordinates
func (gb *Gameboard) ToScreenCoordinates(gridX, gridY int, blockSize float64) (screenX, screenY int) {
	screenX = gb.X + int(float64(gridX)*blockSize)
	screenY = gb.Y + int(float64(gridY)*blockSize)
	return
}

// min returns the smaller of two values
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
