package main

import (
	"time"

	stopwatch "github.com/RAshkettle/Stopwatch"
)

// GameComponents holds all game components in one place
type GameComponents struct {
	// Core game components
	Gameboard    *Gameboard
	BlockManager *BlockManager
	GameLogic    *GameLogic
	Renderer     *GameRenderer

	// System components
	InputHandler    *InputHandler
	AudioManager    *AudioManager
	ParticleSystem  *ParticleSystem
	ScreenShake     *ScreenShake
	ScorePopups     *ScorePopupSystem
	PauseController *PauseController
	EventSystem     *EventSystem
	GameState       *GameState

	// Game pieces
	CurrentPiece *TetrisPiece
	CurrentType  PieceType
	NextPiece    *TetrisPiece
	NextType     PieceType
	FallTimer    *stopwatch.Stopwatch
}

// NewGameComponents creates and initializes all game components
func NewGameComponents() *GameComponents {
	// Create core game state
	gameState := NewGameState()
	eventSystem := NewEventSystem()

	// Create core game components
	gameboard := NewGameboard(192, 320) // 192px wide, 320px tall
	blockManager := NewBlockManager()
	gameLogic := NewGameLogic(gameboard, blockManager)
	renderer := NewGameRenderer(gameboard, blockManager)

	// Create system components
	audioManager := NewAudioManager()
	particleSystem := NewParticleSystem()
	screenShake := NewScreenShake()
	scorePopups := NewScorePopupSystem()
	pauseController := NewPauseController(gameState, audioManager)
	inputHandler := NewInputHandler(gameLogic, audioManager)

	// Create fall timer (1 second intervals)
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

	// Initialize audio
	err := audioManager.Initialize()
	if err != nil {
		println("Warning: Could not initialize audio:", err.Error())
	}

	// Set up event listeners
	components.setupEventListeners()

	return components
}

// setupEventListeners configures all event listeners
func (gc *GameComponents) setupEventListeners() {
	// Blocks removed event - triggers particles, audio, and screen shake
	gc.EventSystem.Subscribe(EventBlocksRemoved, func(event GameEvent) {
		data := event.Data.(BlocksRemovedData)
		
		// Trigger audio
		gc.AudioManager.PlayBlockBreakMultiple(data.Count)
		
		// Trigger screen shake
		intensity := float64(data.Count) * 2.0
		duration := 0.2 + float64(data.Count)*0.05
		gc.ScreenShake.StartShake(intensity, duration)
		
		// Trigger particles for each position
		for _, pos := range data.Positions {
			gc.ParticleSystem.AddExplosion(pos.X, pos.Y, NeutralBlock) // Default block type
		}
	})

	// Piece placed event - triggers dust clouds
	gc.EventSystem.Subscribe(EventPiecePlaced, func(event GameEvent) {
		data := event.Data.(PiecePlacedData)
		gc.ParticleSystem.AddDustCloud(data.Position.X, data.Position.Y)
	})

	// Hard drop event - triggers screen shake
	gc.EventSystem.Subscribe(EventHardDrop, func(event GameEvent) {
		data := event.Data.(HardDropData)
		intensity := 1.0 + float64(data.DropHeight)*0.5
		duration := 0.1
		gc.ScreenShake.StartShake(intensity, duration)
	})
}
