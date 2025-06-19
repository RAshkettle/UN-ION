package main

import (
	"bytes"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

// GameRenderer handles all rendering operations
type GameRenderer struct {
	gameboard      *Gameboard
	blockManager   *BlockManager
	scoreFont      *text.GoTextFace
	scoreLabelFont *text.GoTextFace

	// Reusable text draw options to avoid per-frame allocations
	labelOp *text.DrawOptions
	scoreOp *text.DrawOptions
}

// NewGameRenderer creates a new game renderer
func NewGameRenderer(gameboard *Gameboard, blockManager *BlockManager) *GameRenderer {
	// Create fonts for score display
	scoreFontSource, _ := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	scoreFont := &text.GoTextFace{
		Source: scoreFontSource,
		Size:   32, // Large, bold-looking font for the score number
	}

	scoreLabelFontSource, _ := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	scoreLabelFont := &text.GoTextFace{
		Source: scoreLabelFontSource,
		Size:   18, // Slightly larger font for the "SCORE" label
	}

	return &GameRenderer{
		gameboard:      gameboard,
		blockManager:   blockManager,
		scoreFont:      scoreFont,
		scoreLabelFont: scoreLabelFont,

		// Initialize reusable text draw options
		labelOp: &text.DrawOptions{},
		scoreOp: &text.DrawOptions{},
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

// RenderWithDropShadow draws the complete game state including drop shadow
func (gr *GameRenderer) RenderWithDropShadow(screen *ebiten.Image, placedBlocks []Block, currentPiece *TetrisPiece, shadowPiece *TetrisPiece) {
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

	// Draw drop shadow before current piece (if different positions)
	if shadowPiece != nil && currentPiece != nil && shadowPiece.Y > currentPiece.Y {
		for _, block := range shadowPiece.Blocks {
			worldX := float64(shadowPiece.X+block.X) * blockSize
			worldY := float64(shadowPiece.Y+block.Y) * blockSize

			// Create a shadow block with reduced opacity effect
			shadowBlock := Block{
				X:         block.X,
				Y:         block.Y,
				BlockType: block.BlockType,
			}

			// Draw with special shadow rendering (we'll make it grayed out)
			gr.blockManager.DrawBlock(blocksImage, shadowBlock, worldX, worldY, blockSize)
		}
	}

	// Draw current piece on top of shadow and placed blocks
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

// RenderScore draws the current score at the top right, above the On-Deck preview
func (gr *GameRenderer) RenderScore(screen *ebiten.Image, currentScore int) {
	// Position score above the next piece preview area (top right)
	margin := 10
	scoreX := gr.gameboard.X + gr.gameboard.Width + 20 // Same X as preview area
	scoreY := max(margin, gr.gameboard.Y-15)           // Position above the gameboard, with more space for larger font

	// Draw "SCORE" label
	gr.labelOp.GeoM.Reset()
	gr.labelOp.GeoM.Translate(float64(scoreX), float64(scoreY))
	gr.labelOp.ColorScale.Reset()
	gr.labelOp.ColorScale.ScaleWithColor(color.RGBA{200, 200, 255, 255})
	text.Draw(screen, "SCORE", gr.scoreLabelFont, gr.labelOp)

	// Draw the actual score value with large, bold font
	scoreText := fmt.Sprintf("%d", currentScore)
	gr.scoreOp.GeoM.Reset()
	gr.scoreOp.GeoM.Translate(float64(scoreX), float64(scoreY+25)) // More space for larger font
	gr.scoreOp.ColorScale.Reset()
	gr.scoreOp.ColorScale.ScaleWithColor(color.RGBA{255, 255, 100, 255}) // Bright yellow for emphasis
	text.Draw(screen, scoreText, gr.scoreFont, gr.scoreOp)
}

// RenderDropShadow draws a translucent preview of where the piece will land
func (gr *GameRenderer) RenderDropShadow(screen *ebiten.Image, shadowPiece *TetrisPiece) {
	if shadowPiece == nil {
		return
	}

	// Calculate block size for rendering
	blockSize := gr.blockManager.GetScaledBlockSize(gr.gameboard.Width, gr.gameboard.Height)

	// Create a temporary image for the shadow blocks
	shadowImage := ebiten.NewImage(gr.gameboard.Width, gr.gameboard.Height)

	// Draw shadow blocks with reduced opacity
	for _, block := range shadowPiece.Blocks {
		worldX := float64(shadowPiece.X+block.X) * blockSize
		worldY := float64(shadowPiece.Y+block.Y) * blockSize

		// Draw the block normally first
		gr.blockManager.DrawBlock(shadowImage, block, worldX, worldY, blockSize)
	}

	// Apply the shadow image to screen with reduced opacity
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(gr.gameboard.X), float64(gr.gameboard.Y))
	op.ColorScale.ScaleAlpha(0.3) // Make it 30% transparent
	screen.DrawImage(shadowImage, op)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
