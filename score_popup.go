package main

import (
	"bytes"
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

type ScorePopup struct {
	X, Y    float64
	VX, VY  float64
	Life    float64
	MaxLife float64
	Score   int
	Alpha   float64
	Scale   float64
}

type ScorePopupSystem struct {
	popups []ScorePopup
	font   *text.GoTextFace
}

func NewScorePopupSystem() *ScorePopupSystem {
	fontSource, _ := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	font := &text.GoTextFace{
		Source: fontSource,
		Size:   20,
	}
	return &ScorePopupSystem{
		popups: make([]ScorePopup, 0),
		font:   font,
	}
}

func (sps *ScorePopupSystem) AddScorePopup(x, y float64, score int) {
	popup := ScorePopup{
		X:       x,
		Y:       y,
		VX:      0,
		VY:      -30.0,
		Life:    1.5,
		MaxLife: 1.5,
		Score:   score,
		Alpha:   1.0,
		Scale:   1.0,
	}

	sps.popups = append(sps.popups, popup)
}

func (sps *ScorePopupSystem) Update(dt float64) {
	for i := len(sps.popups) - 1; i >= 0; i-- {
		popup := &sps.popups[i]

		popup.X += popup.VX * dt
		popup.Y += popup.VY * dt

		dampingFactor := math.Pow(0.95, dt*60.0)
		popup.VY *= dampingFactor

		popup.Life -= dt

		// Update alpha and scale based on remaining life
		lifeRatio := popup.Life / popup.MaxLife

		if lifeRatio > 0.8 {
			// Growing phase
			popup.Scale = 0.5 + (1.0-lifeRatio)*5.0*0.5 // Grow from 0.5 to 1.0
			popup.Alpha = 1.0
		} else if lifeRatio > 0.2 {
			// Stable phase
			popup.Scale = 1.0
			popup.Alpha = 1.0
		} else {
			// Fading phase
			popup.Scale = 1.0
			popup.Alpha = lifeRatio / 0.2 // Fade from 1.0 to 0.0
		}

		// Remove dead popups
		if popup.Life <= 0 {
			// Remove popup by swapping with last and shrinking slice
			sps.popups[i] = sps.popups[len(sps.popups)-1]
			sps.popups = sps.popups[:len(sps.popups)-1]
		}
	}
}

func (sps *ScorePopupSystem) Draw(screen *ebiten.Image) {
	for _, popup := range sps.popups {
		if popup.Alpha <= 0 {
			continue
		}

		scoreText := fmt.Sprintf("+%d", popup.Score)

		alpha := uint8(popup.Alpha * 255)
		textColor := color.RGBA{255, 255, 0, alpha}

		if popup.Scale > 0.7 {
			op := &text.DrawOptions{}
			op.GeoM.Translate(popup.X, popup.Y)
			op.ColorScale.ScaleWithColor(textColor)
			text.Draw(screen, scoreText, sps.font, op)
		}
	}
}

func (sps *ScorePopupSystem) HasActivePopups() bool {
	return len(sps.popups) > 0
}
