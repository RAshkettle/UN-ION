package main

import (
	"bytes"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

type GameRenderer struct {
	gameboard      *Gameboard
	blockManager   *BlockManager
	scoreFont      *text.GoTextFace
	scoreLabelFont *text.GoTextFace
	labelOp        *text.DrawOptions
	scoreOp        *text.DrawOptions
}

func NewGameRenderer(gameboard *Gameboard, blockManager *BlockManager) *GameRenderer {
	scoreFontSource, _ := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	scoreFont := &text.GoTextFace{
		Source: scoreFontSource,
		Size:   32,
	}

	scoreLabelFontSource, _ := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	scoreLabelFont := &text.GoTextFace{
		Source: scoreLabelFontSource,
		Size:   18,
	}

	return &GameRenderer{
		gameboard:      gameboard,
		blockManager:   blockManager,
		scoreFont:      scoreFont,
		scoreLabelFont: scoreLabelFont,
		labelOp:        &text.DrawOptions{},
		scoreOp:        &text.DrawOptions{},
	}
}

func (gr *GameRenderer) Render(screen *ebiten.Image, placedBlocks []Block, currentPiece *TetrisPiece) {
	screen.Fill(color.RGBA{15, 20, 30, 255})
	gr.gameboard.Draw(screen)
	blockSize := gr.blockManager.GetScaledBlockSize(gr.gameboard.Width, gr.gameboard.Height)
	blocksImage := ebiten.NewImage(gr.gameboard.Width, gr.gameboard.Height)
	for _, block := range placedBlocks {
		worldX := float64(block.X) * blockSize
		worldY := float64(block.Y) * blockSize
		gr.blockManager.DrawBlock(blocksImage, block, worldX, worldY, blockSize)
	}
	if currentPiece != nil {
		for _, block := range currentPiece.Blocks {
			worldX := float64(currentPiece.X+block.X) * blockSize
			worldY := float64(currentPiece.Y+block.Y) * blockSize
			gr.blockManager.DrawBlock(blocksImage, block, worldX, worldY, blockSize)
		}
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(gr.gameboard.X), float64(gr.gameboard.Y))
	screen.DrawImage(blocksImage, op)
}

func (gr *GameRenderer) RenderWithDropShadow(screen *ebiten.Image, placedBlocks []Block, currentPiece *TetrisPiece, shadowPiece *TetrisPiece) {
	screen.Fill(color.RGBA{15, 20, 30, 255})
	gr.gameboard.Draw(screen)
	blockSize := gr.blockManager.GetScaledBlockSize(gr.gameboard.Width, gr.gameboard.Height)
	blocksImage := ebiten.NewImage(gr.gameboard.Width, gr.gameboard.Height)
	for _, block := range placedBlocks {
		worldX := float64(block.X) * blockSize
		worldY := float64(block.Y) * blockSize
		gr.blockManager.DrawBlock(blocksImage, block, worldX, worldY, blockSize)
	}
	if shadowPiece != nil && currentPiece != nil && shadowPiece.Y > currentPiece.Y {
		for _, block := range shadowPiece.Blocks {
			worldX := float64(shadowPiece.X+block.X) * blockSize
			worldY := float64(shadowPiece.Y+block.Y) * blockSize
			shadowBlock := Block{
				X:         block.X,
				Y:         block.Y,
				BlockType: block.BlockType,
			}
			gr.blockManager.DrawBlock(blocksImage, shadowBlock, worldX, worldY, blockSize)
		}
	}
	if currentPiece != nil {
		for _, block := range currentPiece.Blocks {
			worldX := float64(currentPiece.X+block.X) * blockSize
			worldY := float64(currentPiece.Y+block.Y) * blockSize
			gr.blockManager.DrawBlock(blocksImage, block, worldX, worldY, blockSize)
		}
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(gr.gameboard.X), float64(gr.gameboard.Y))
	screen.DrawImage(blocksImage, op)
}

func (gr *GameRenderer) RenderScore(screen *ebiten.Image, currentScore int) {
	margin := 10
	scoreX := gr.gameboard.X + gr.gameboard.Width + 20
	scoreY := max(margin, gr.gameboard.Y-15)
	gr.labelOp.GeoM.Reset()
	gr.labelOp.GeoM.Translate(float64(scoreX), float64(scoreY))
	gr.labelOp.ColorScale.Reset()
	gr.labelOp.ColorScale.ScaleWithColor(color.RGBA{200, 200, 255, 255})
	text.Draw(screen, "SCORE", gr.scoreLabelFont, gr.labelOp)
	scoreText := fmt.Sprintf("%d", currentScore)
	gr.scoreOp.GeoM.Reset()
	gr.scoreOp.GeoM.Translate(float64(scoreX), float64(scoreY+25))
	gr.scoreOp.ColorScale.Reset()
	gr.scoreOp.ColorScale.ScaleWithColor(color.RGBA{255, 255, 100, 255})
	text.Draw(screen, scoreText, gr.scoreFont, gr.scoreOp)
}

func (gr *GameRenderer) RenderDropShadow(screen *ebiten.Image, shadowPiece *TetrisPiece) {
	if shadowPiece == nil {
		return
	}
	blockSize := gr.blockManager.GetScaledBlockSize(gr.gameboard.Width, gr.gameboard.Height)
	shadowImage := ebiten.NewImage(gr.gameboard.Width, gr.gameboard.Height)
	for _, block := range shadowPiece.Blocks {
		worldX := float64(shadowPiece.X+block.X) * blockSize
		worldY := float64(shadowPiece.Y+block.Y) * blockSize
		gr.blockManager.DrawBlock(shadowImage, block, worldX, worldY, blockSize)
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(gr.gameboard.X), float64(gr.gameboard.Y))
	op.ColorScale.ScaleAlpha(0.3)
	screen.DrawImage(shadowImage, op)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
