package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// GameRenderer handles all rendering operations
type GameRenderer struct {
	gameboard    *Gameboard
	blockManager *BlockManager
}

// NewGameRenderer creates a new game renderer
func NewGameRenderer(gameboard *Gameboard, blockManager *BlockManager) *GameRenderer {
	return &GameRenderer{
		gameboard:    gameboard,
		blockManager: blockManager,
	}
}

// Render draws the complete game state
func (gr *GameRenderer) Render(screen *ebiten.Image, placedBlocks []Block, currentPiece *TetrisPiece) {
	// Dark background
	screen.Fill(color.RGBA{15, 20, 30, 255})

	// Draw the gameboard with shader effect FIRST (background)
	gr.gameboard.Draw(screen)

	// Calculate block size for rendering
	blockSize := gr.blockManager.GetScaledBlockSize(gr.gameboard.Width, gr.gameboard.Height)

	// Create a temporary image for all blocks
	blocksImage := ebiten.NewImage(gr.gameboard.Width, gr.gameboard.Height)

	// Draw placed blocks first
	for _, block := range placedBlocks {
		worldX := float64(block.X) * blockSize
		worldY := float64(block.Y) * blockSize
		gr.blockManager.DrawBlock(blocksImage, block, worldX, worldY, blockSize)
	}

	// Draw current piece on top of placed blocks
	if currentPiece != nil {
		for _, block := range currentPiece.Blocks {
			worldX := float64(currentPiece.X+block.X) * blockSize
			worldY := float64(currentPiece.Y+block.Y) * blockSize
			gr.blockManager.DrawBlock(blocksImage, block, worldX, worldY, blockSize)
		}
	}

	// Draw all blocks on top of the gameboard
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(gr.gameboard.X), float64(gr.gameboard.Y))
	screen.DrawImage(blocksImage, op)
}

// RenderScore draws the current score on the left side of the gameboard
func (gr *GameRenderer) RenderScore(screen *ebiten.Image, currentScore int) {
	// Position score relative to gameboard but ensure it's visible
	margin := 10
	scoreX := max(margin, gr.gameboard.X - 80)
	scoreY := gr.gameboard.Y + 50
	
	// Draw "SCORE" label
	text.Draw(screen, "SCORE", basicfont.Face7x13, scoreX, scoreY, color.RGBA{200, 200, 255, 255})
	
	// Draw the actual score value with better formatting
	scoreText := fmt.Sprintf("%d", currentScore)
	text.Draw(screen, scoreText, basicfont.Face7x13, scoreX, scoreY+25, color.RGBA{255, 255, 255, 255})
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
