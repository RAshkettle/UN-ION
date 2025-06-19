package main

import "time"

type GameState struct {
	IsPaused     bool
	Score        int
	Level        int
	LinesCleared int
	LastUpdate   time.Time
}

func NewGameState() *GameState {
	return &GameState{
		IsPaused:     false,
		Score:        0,
		Level:        1,
		LinesCleared: 0,
		LastUpdate:   time.Now(),
	}
}

func (gs *GameState) TogglePause() {
	gs.IsPaused = !gs.IsPaused
}

func (gs *GameState) AddScore(points int) {
	gs.Score += points
}

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
