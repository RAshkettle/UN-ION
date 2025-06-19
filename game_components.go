package main

import (
	"time"

	stopwatch "github.com/RAshkettle/Stopwatch"
)

type GameComponents struct {
	Gameboard       *Gameboard
	BlockManager    *BlockManager
	GameLogic       *GameLogic
	Renderer        *GameRenderer
	InputHandler    *InputHandler
	AudioManager    *AudioManager
	ParticleSystem  *ParticleSystem
	ScreenShake     *ScreenShake
	ScorePopups     *ScorePopupSystem
	PauseController *PauseController
	EventSystem     *EventSystem
	GameState       *GameState
	CurrentPiece    *TetrisPiece
	CurrentType     PieceType
	NextPiece       *TetrisPiece
	NextType        PieceType
	FallTimer       *stopwatch.Stopwatch
}

func NewGameComponents() *GameComponents {
	gameState := NewGameState()
	eventSystem := NewEventSystem()
	gameboard := NewGameboard(192, 320)
	blockManager := NewBlockManager()
	gameLogic := NewGameLogic(gameboard, blockManager)
	renderer := NewGameRenderer(gameboard, blockManager)
	audioManager := NewAudioManager()
	particleSystem := NewParticleSystem()
	screenShake := NewScreenShake()
	scorePopups := NewScorePopupSystem()
	pauseController := NewPauseController(gameState, audioManager)
	inputHandler := NewInputHandler(gameLogic, audioManager)
	fallTimer := stopwatch.NewStopwatch(1 * time.Second)
	fallTimer.Start()

	components := &GameComponents{
		Gameboard:       gameboard,
		BlockManager:    blockManager,
		GameLogic:       gameLogic,
		Renderer:        renderer,
		InputHandler:    inputHandler,
		AudioManager:    audioManager,
		ParticleSystem:  particleSystem,
		ScreenShake:     screenShake,
		ScorePopups:     scorePopups,
		PauseController: pauseController,
		EventSystem:     eventSystem,
		GameState:       gameState,
		FallTimer:       fallTimer,
	}

	err := audioManager.Initialize()
	if err != nil {
		println("Warning: Could not initialize audio:", err.Error())
	}

	components.setupEventListeners()

	return components
}

func (gc *GameComponents) setupEventListeners() {
	gc.EventSystem.Subscribe(EventBlocksRemoved, func(event GameEvent) {
		data := event.Data.(BlocksRemovedData)
		gc.AudioManager.PlayBlockBreakMultiple(data.Count)
		intensity := float64(data.Count) * 2.0
		duration := 0.2 + float64(data.Count)*0.05
		gc.ScreenShake.StartShake(intensity, duration)
		for _, pos := range data.Positions {
			gc.ParticleSystem.AddExplosion(pos.X, pos.Y, NeutralBlock)
		}
	})

	gc.EventSystem.Subscribe(EventPiecePlaced, func(event GameEvent) {
		data := event.Data.(PiecePlacedData)
		gc.ParticleSystem.AddDustCloud(data.Position.X, data.Position.Y)
	})

	gc.EventSystem.Subscribe(EventHardDrop, func(event GameEvent) {
		data := event.Data.(HardDropData)
		intensity := 1.0 + float64(data.DropHeight)*0.5
		duration := 0.1
		gc.ScreenShake.StartShake(intensity, duration)
	})
}
