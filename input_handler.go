package main

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// InputHandler manages all game input
type InputHandler struct {
	gameLogic    *GameLogic
	audioManager *AudioManager
	// Key repeat timers
	leftRepeatTimer  time.Time
	rightRepeatTimer time.Time
	downRepeatTimer  time.Time
	// Key repeat intervals
	initialDelay time.Duration
	repeatDelay  time.Duration
}

// NewInputHandler creates a new input handler
func NewInputHandler(gameLogic *GameLogic, audioManager *AudioManager) *InputHandler {
	return &InputHandler{
		gameLogic:    gameLogic,
		audioManager: audioManager,
		initialDelay: 200 * time.Millisecond, // Initial delay before repeat starts
		repeatDelay:  100 * time.Millisecond, // Delay between repeats
	}
}

// HandleInput processes all input for the current frame
// Returns true if the piece should be placed immediately due to manual drop
func (ih *InputHandler) HandleInput(currentPiece *TetrisPiece, currentType PieceType) bool {
	if currentPiece == nil {
		return false
	}

	now := time.Now()
	shouldPlace := false

	// Handle piece rotation (only on key press, no repeat)
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		ih.gameLogic.TryRotatePiece(currentPiece, currentType)
	}

	// Handle left movement with key repeat
	leftPressed := ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft)
	if leftPressed {
		if inpututil.IsKeyJustPressed(ebiten.KeyA) || inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
			// First press - move immediately
			if ih.gameLogic.TryMovePiece(currentPiece, -1, 0) {
				if ih.audioManager != nil {
					ih.audioManager.PlaySwooshSound()
				}
			}
			ih.leftRepeatTimer = now.Add(ih.initialDelay)
		} else if now.After(ih.leftRepeatTimer) {
			// Key held - repeat movement
			if ih.gameLogic.TryMovePiece(currentPiece, -1, 0) {
				if ih.audioManager != nil {
					ih.audioManager.PlaySwooshSound()
				}
			}
			ih.leftRepeatTimer = now.Add(ih.repeatDelay)
		}
	}

	// Handle right movement with key repeat
	rightPressed := ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight)
	if rightPressed {
		if inpututil.IsKeyJustPressed(ebiten.KeyD) || inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
			// First press - move immediately
			if ih.gameLogic.TryMovePiece(currentPiece, 1, 0) {
				if ih.audioManager != nil {
					ih.audioManager.PlaySwooshSound()
				}
			}
			ih.rightRepeatTimer = now.Add(ih.initialDelay)
		} else if now.After(ih.rightRepeatTimer) {
			// Key held - repeat movement
			if ih.gameLogic.TryMovePiece(currentPiece, 1, 0) {
				if ih.audioManager != nil {
					ih.audioManager.PlaySwooshSound()
				}
			}
			ih.rightRepeatTimer = now.Add(ih.repeatDelay)
		}
	}

	// Handle down movement with key repeat (faster for quick drop)
	downPressed := ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown)
	if downPressed {
		if inpututil.IsKeyJustPressed(ebiten.KeyS) || inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			// First press - try to move down, place if can't
			if ih.gameLogic.TryMovePiece(currentPiece, 0, 1) {
				// Successful move - play swoosh
				if ih.audioManager != nil {
					ih.audioManager.PlaySwooshSound()
				}
			} else {
				// Hit bottom - place piece
				shouldPlace = true
				// Calculate drop height for screen shake (minimum 1 to always have some effect)
				dropHeight := 1
				ih.triggerHardDropShake(dropHeight)
			}
			ih.downRepeatTimer = now.Add(50 * time.Millisecond) // Faster initial delay for down
		} else if now.After(ih.downRepeatTimer) {
			// Key held - try to move down, place if can't
			if ih.gameLogic.TryMovePiece(currentPiece, 0, 1) {
				// Successful move - play swoosh
				if ih.audioManager != nil {
					ih.audioManager.PlaySwooshSound()
				}
			} else {
				// Hit bottom - place piece
				shouldPlace = true
				// Calculate drop height for screen shake (minimum 1 to always have some effect)
				dropHeight := 1
				ih.triggerHardDropShake(dropHeight)
			}
			ih.downRepeatTimer = now.Add(50 * time.Millisecond) // Faster repeat for down
		}
	}

	return shouldPlace
}

// triggerHardDropShake triggers screen shake for hard drops
func (ih *InputHandler) triggerHardDropShake(dropHeight int) {
	if ih.gameLogic.hardDropCallback != nil {
		ih.gameLogic.hardDropCallback(dropHeight)
	}
}
