package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// InputHandler manages all game input
type InputHandler struct {
	gameLogic *GameLogic
}

// NewInputHandler creates a new input handler
func NewInputHandler(gameLogic *GameLogic) *InputHandler {
	return &InputHandler{
		gameLogic: gameLogic,
	}
}

// HandleInput processes all input for the current frame
func (ih *InputHandler) HandleInput(currentPiece *TetrisPiece, currentType PieceType) {
	if currentPiece == nil {
		return
	}

	// Handle piece rotation
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		ih.gameLogic.TryRotatePiece(currentPiece, currentType)
	}

	// Handle piece movement (16 pixels = 1 block)
	if inpututil.IsKeyJustPressed(ebiten.KeyA) || inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		ih.gameLogic.TryMovePiece(currentPiece, -1, 0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) || inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		ih.gameLogic.TryMovePiece(currentPiece, 1, 0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) || inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		ih.gameLogic.TryMovePiece(currentPiece, 0, 1)
	}
}
