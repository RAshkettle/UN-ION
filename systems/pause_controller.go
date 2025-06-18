package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// PauseController handles pause functionality
type PauseController struct {
	gameState    *GameState
	audioManager *AudioManager
}

// NewPauseController creates a new pause controller
func NewPauseController(gameState *GameState, audioManager *AudioManager) *PauseController {
	return &PauseController{
		gameState:    gameState,
		audioManager: audioManager,
	}
}

// Update handles pause input
func (pc *PauseController) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		pc.gameState.TogglePause()

		if pc.gameState.IsPaused {
			pc.audioManager.PauseBackgroundMusic()
		} else {
			pc.audioManager.ResumeBackgroundMusic()
		}
	}
}

// Draw renders the pause overlay if paused
func (pc *PauseController) Draw(screen *ebiten.Image) {
	if !pc.gameState.IsPaused {
		return
	}

	// Draw semi-transparent overlay
	overlay := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
	overlay.Fill(color.RGBA{0, 0, 0, 128}) // 50% transparent black
	screen.DrawImage(overlay, nil)

	// Draw "PAUSED" text in the center
	centerX := screen.Bounds().Dx() / 2
	centerY := screen.Bounds().Dy() / 2

	// Use basic font for text rendering
	fontFace := basicfont.Face7x13

	// Draw "PAUSED" text
	pausedText := "PAUSED"
	pausedBounds := text.BoundString(fontFace, pausedText)
	pausedX := centerX - pausedBounds.Dx()/2
	pausedY := centerY - 10
	text.Draw(screen, pausedText, fontFace, pausedX, pausedY, color.RGBA{255, 255, 255, 255})

	// Draw "Press P to Resume" below
	resumeText := "Press P to Resume"
	resumeBounds := text.BoundString(fontFace, resumeText)
	resumeX := centerX - resumeBounds.Dx()/2
	resumeY := centerY + 20
	text.Draw(screen, resumeText, fontFace, resumeX, resumeY, color.RGBA{200, 200, 200, 255})
}
