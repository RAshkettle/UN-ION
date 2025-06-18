package main

import (
	"math"
	"math/rand"
)

// ScreenShake handles camera shake effects
type ScreenShake struct {
	intensity     float64
	duration      float64
	currentTime   float64
	offsetX       float64
	offsetY       float64
	isShaking     bool
}

// NewScreenShake creates a new screen shake manager
func NewScreenShake() *ScreenShake {
	return &ScreenShake{}
}

// StartShake begins a screen shake effect
func (ss *ScreenShake) StartShake(intensity, duration float64) {
	ss.intensity = intensity
	ss.duration = duration
	ss.currentTime = 0
	ss.isShaking = true
}

// Update updates the screen shake effect
func (ss *ScreenShake) Update(deltaTime float64) {
	if !ss.isShaking {
		ss.offsetX = 0
		ss.offsetY = 0
		return
	}
	
	ss.currentTime += deltaTime
	
	if ss.currentTime >= ss.duration {
		ss.isShaking = false
		ss.offsetX = 0
		ss.offsetY = 0
		return
	}
	
	// Calculate shake intensity that decreases over time
	progress := ss.currentTime / ss.duration
	currentIntensity := ss.intensity * (1.0 - progress)
	
	// Generate random shake offset
	angle := rand.Float64() * 2 * math.Pi
	ss.offsetX = math.Cos(angle) * currentIntensity
	ss.offsetY = math.Sin(angle) * currentIntensity
}

// GetOffset returns the current shake offset
func (ss *ScreenShake) GetOffset() (float64, float64) {
	return ss.offsetX, ss.offsetY
}

// IsShaking returns true if currently shaking
func (ss *ScreenShake) IsShaking() bool {
	return ss.isShaking
}
