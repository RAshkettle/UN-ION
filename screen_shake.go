package main

import (
	"math"
	"math/rand"
)


type ScreenShake struct {
	intensity   float64
	duration    float64
	currentTime float64
	offsetX     float64
	offsetY     float64
	isShaking   bool
}


func NewScreenShake() *ScreenShake {
	return &ScreenShake{}
}

func (ss *ScreenShake) StartShake(intensity, duration float64) {
	ss.intensity = intensity
	ss.duration = duration
	ss.currentTime = 0
	ss.isShaking = true
}


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


func (ss *ScreenShake) GetOffset() (float64, float64) {
	return ss.offsetX, ss.offsetY
}


func (ss *ScreenShake) IsShaking() bool {
	return ss.isShaking
}
