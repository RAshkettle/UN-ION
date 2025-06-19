package main

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type InputHandler struct {
	gameLogic        *GameLogic
	audioManager     *AudioManager
	leftRepeatTimer  time.Time
	rightRepeatTimer time.Time
	downRepeatTimer  time.Time
	initialDelay     time.Duration
	repeatDelay      time.Duration
}

func NewInputHandler(gameLogic *GameLogic, audioManager *AudioManager) *InputHandler {
	return &InputHandler{
		gameLogic:    gameLogic,
		audioManager: audioManager,
		initialDelay: 200 * time.Millisecond,
		repeatDelay:  100 * time.Millisecond,
	}
}

func (ih *InputHandler) HandleInput(currentPiece *TetrisPiece, currentType PieceType) bool {
	if currentPiece == nil {
		return false
	}

	now := time.Now()
	shouldPlace := false

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		ih.gameLogic.TryRotatePiece(currentPiece, currentType)
	}

	leftPressed := ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft)
	if leftPressed {
		if inpututil.IsKeyJustPressed(ebiten.KeyA) || inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
			if ih.gameLogic.TryMovePiece(currentPiece, -1, 0) {
				if ih.audioManager != nil {
					ih.audioManager.PlaySwooshSound()
				}
			}
			ih.leftRepeatTimer = now.Add(ih.initialDelay)
		} else if now.After(ih.leftRepeatTimer) {
			if ih.gameLogic.TryMovePiece(currentPiece, -1, 0) {
				if ih.audioManager != nil {
					ih.audioManager.PlaySwooshSound()
				}
			}
			ih.leftRepeatTimer = now.Add(ih.repeatDelay)
		}
	}

	rightPressed := ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight)
	if rightPressed {
		if inpututil.IsKeyJustPressed(ebiten.KeyD) || inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
			if ih.gameLogic.TryMovePiece(currentPiece, 1, 0) {
				if ih.audioManager != nil {
					ih.audioManager.PlaySwooshSound()
				}
			}
			ih.rightRepeatTimer = now.Add(ih.initialDelay)
		} else if now.After(ih.rightRepeatTimer) {
			if ih.gameLogic.TryMovePiece(currentPiece, 1, 0) {
				if ih.audioManager != nil {
					ih.audioManager.PlaySwooshSound()
				}
			}
			ih.rightRepeatTimer = now.Add(ih.repeatDelay)
		}
	}

	downPressed := ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown)
	if downPressed {
		if inpututil.IsKeyJustPressed(ebiten.KeyS) || inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			if ih.gameLogic.TryMovePiece(currentPiece, 0, 1) {
				if ih.audioManager != nil {
					ih.audioManager.PlaySwooshSound()
				}
			} else {
				shouldPlace = true
				dropHeight := 1
				ih.triggerHardDropShake(dropHeight)
			}
			ih.downRepeatTimer = now.Add(50 * time.Millisecond)
		} else if now.After(ih.downRepeatTimer) {
			if ih.gameLogic.TryMovePiece(currentPiece, 0, 1) {
				if ih.audioManager != nil {
					ih.audioManager.PlaySwooshSound()
				}
			} else {
				shouldPlace = true
				dropHeight := 1
				ih.triggerHardDropShake(dropHeight)
			}
			ih.downRepeatTimer = now.Add(50 * time.Millisecond)
		}
	}

	return shouldPlace
}

func (ih *InputHandler) triggerHardDropShake(dropHeight int) {
	if ih.gameLogic.hardDropCallback != nil {
		ih.gameLogic.hardDropCallback(dropHeight)
	}
}
