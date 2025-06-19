package main

import (
	"bytes"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

var pauseTitleFontSource, _ = text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
var pauseTitleFont = &text.GoTextFace{
	Source: pauseTitleFontSource,
	Size:   48,
}
var pauseSubtitleFontSource, _ = text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
var pauseSubtitleFont = &text.GoTextFace{
	Source: pauseSubtitleFontSource,
	Size:   24,
}


type PauseController struct {
	gameState    *GameState
	audioManager *AudioManager
}

func NewPauseController(gameState *GameState, audioManager *AudioManager) *PauseController {
	return &PauseController{
		gameState:    gameState,
		audioManager: audioManager,
	}
}


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


func (pc *PauseController) Draw(screen *ebiten.Image) {
	if !pc.gameState.IsPaused {
		return
	}

	// Draw semi-transparent overlay
	overlay := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
	overlay.Fill(color.RGBA{0, 0, 0, 128}) 
	screen.DrawImage(overlay, nil)


	centerX := screen.Bounds().Dx() / 2
	centerY := screen.Bounds().Dy() / 2

	pausedText := "PAUSED"
	pausedAdvance, _ := text.Measure(pausedText, pauseTitleFont, 0)
	pausedX := centerX - int(pausedAdvance)/2
	pausedY := centerY - 10
	pausedOp := &text.DrawOptions{}
	pausedOp.GeoM.Translate(float64(pausedX), float64(pausedY))
	pausedOp.ColorScale.ScaleWithColor(color.RGBA{220, 220, 255, 255})
	text.Draw(screen, pausedText, pauseTitleFont, pausedOp)


	resumeText := "Press P to Resume"
	resumeAdvance, _ := text.Measure(resumeText, pauseSubtitleFont, 0)
	resumeX := centerX - int(resumeAdvance)/2
	resumeY := centerY + 20
	resumeOp := &text.DrawOptions{}
	resumeOp.GeoM.Translate(float64(resumeX), float64(resumeY))
	resumeOp.ColorScale.ScaleWithColor(color.RGBA{200, 200, 200, 255})
	text.Draw(screen, resumeText, pauseSubtitleFont, resumeOp)
}
