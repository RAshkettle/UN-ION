package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// ScorePopup represents a floating score text animation
type ScorePopup struct {
	X, Y       float64
	VX, VY     float64
	Life       float64
	MaxLife    float64
	Score      int
	Alpha      float64
	Scale      float64
}

// ScorePopupSystem manages all score popups
type ScorePopupSystem struct {
	popups []ScorePopup
}

// NewScorePopupSystem creates a new score popup system
func NewScorePopupSystem() *ScorePopupSystem {
	return &ScorePopupSystem{
		popups: make([]ScorePopup, 0),
	}
}

// AddScorePopup creates a new score popup at the specified location
func (sps *ScorePopupSystem) AddScorePopup(x, y float64, score int) {
	popup := ScorePopup{
		X:       x,
		Y:       y,
		VX:      0,
		VY:      -30.0, // Float upward
		Life:    1.5,   // 1.5 seconds duration
		MaxLife: 1.5,
		Score:   score,
		Alpha:   1.0,
		Scale:   1.0,
	}
	
	sps.popups = append(sps.popups, popup)
}

// Update updates all score popups
func (sps *ScorePopupSystem) Update(dt float64) {
	// Update existing popups
	for i := len(sps.popups) - 1; i >= 0; i-- {
		popup := &sps.popups[i]
		
		// Update position
		popup.X += popup.VX * dt
		popup.Y += popup.VY * dt
		
		// Slow down vertical movement over time
		popup.VY *= 0.95
		
		// Update life
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

// Draw renders all score popups
func (sps *ScorePopupSystem) Draw(screen *ebiten.Image) {
	for _, popup := range sps.popups {
		if popup.Alpha <= 0 {
			continue
		}
		
		// Create score text
		scoreText := fmt.Sprintf("+%d", popup.Score)
		
		// Calculate color with alpha
		alpha := uint8(popup.Alpha * 255)
		textColor := color.RGBA{255, 255, 0, alpha} // Bright yellow
		
		// Draw the text (we'll use a simple approach for now)
		// In a production game, you'd want to use a proper font with scaling
		if popup.Scale > 0.7 { // Only draw if large enough to be readable
			text.Draw(screen, scoreText, basicfont.Face7x13, int(popup.X), int(popup.Y), textColor)
		}
	}
}

// HasActivePopups returns true if there are any active popups
func (sps *ScorePopupSystem) HasActivePopups() bool {
	return len(sps.popups) > 0
}
