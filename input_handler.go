package main

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// InputHandler manages all game input
type InputHandler struct {
	gameLogic *GameLogic
	// Key repeat timers
	leftRepeatTimer  time.Time
	rightRepeatTimer time.Time
	downRepeatTimer  time.Time
	// Key repeat intervals
	initialDelay time.Duration
	repeatDelay  time.Duration
}

// NewInputHandler creates a new input handler
func NewInputHandler(gameLogic *GameLogic) *InputHandler {
	return &InputHandler{
		gameLogic:    gameLogic,
		initialDelay: 200 * time.Millisecond, // Initial delay before repeat starts
		repeatDelay:  100 * time.Millisecond, // Delay between repeats
	}
}

// HandleInput processes all input for the current frame
func (ih *InputHandler) HandleInput(currentPiece *TetrisPiece, currentType PieceType) {
	if currentPiece == nil {
		return
	}

	now := time.Now()

	// Handle piece rotation (only on key press, no repeat)
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		ih.gameLogic.TryRotatePiece(currentPiece, currentType)
	}

	// Handle left movement with key repeat
	leftPressed := ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft)
	if leftPressed {
		if inpututil.IsKeyJustPressed(ebiten.KeyA) || inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
			// First press - move immediately
			ih.gameLogic.TryMovePiece(currentPiece, -1, 0)
			ih.leftRepeatTimer = now.Add(ih.initialDelay)
		} else if now.After(ih.leftRepeatTimer) {
			// Key held - repeat movement
			ih.gameLogic.TryMovePiece(currentPiece, -1, 0)
			ih.leftRepeatTimer = now.Add(ih.repeatDelay)
		}
	}

	// Handle right movement with key repeat
	rightPressed := ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight)
	if rightPressed {
		if inpututil.IsKeyJustPressed(ebiten.KeyD) || inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
			// First press - move immediately
			ih.gameLogic.TryMovePiece(currentPiece, 1, 0)
			ih.rightRepeatTimer = now.Add(ih.initialDelay)
		} else if now.After(ih.rightRepeatTimer) {
			// Key held - repeat movement
			ih.gameLogic.TryMovePiece(currentPiece, 1, 0)
			ih.rightRepeatTimer = now.Add(ih.repeatDelay)
		}
	}

	// Handle down movement with key repeat (faster for quick drop)
	downPressed := ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown)
	if downPressed {
		if inpututil.IsKeyJustPressed(ebiten.KeyS) || inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			// First press - move immediately
			ih.gameLogic.TryMovePiece(currentPiece, 0, 1)
			ih.downRepeatTimer = now.Add(50 * time.Millisecond) // Faster initial delay for down
		} else if now.After(ih.downRepeatTimer) {
			// Key held - repeat movement
			ih.gameLogic.TryMovePiece(currentPiece, 0, 1)
			ih.downRepeatTimer = now.Add(50 * time.Millisecond) // Faster repeat for down
		}
	}
}
