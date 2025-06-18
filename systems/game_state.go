package main

import "time"

// GameState represents the current state of the game
type GameState struct {
	IsPaused     bool
	Score        int
	Level        int
	LinesCleared int
	LastUpdate   time.Time
}

// NewGameState creates a new game state
func NewGameState() *GameState {
	return &GameState{
		IsPaused:     false,
		Score:        0,
		Level:        1,
		LinesCleared: 0,
		LastUpdate:   time.Now(),
	}
}

// TogglePause toggles the pause state
func (gs *GameState) TogglePause() {
	gs.IsPaused = !gs.IsPaused
}

// AddScore adds points to the current score
func (gs *GameState) AddScore(points int) {
	gs.Score += points
}

// GetDeltaTime calculates and returns the delta time since last update
func (gs *GameState) GetDeltaTime() float64 {
	now := time.Now()
	if gs.LastUpdate.IsZero() {
		gs.LastUpdate = now
		return 0
	}
	dt := now.Sub(gs.LastUpdate).Seconds()
	gs.LastUpdate = now
	return dt
}
