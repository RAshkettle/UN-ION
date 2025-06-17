package main

import (
	"math/rand"
	"time"

	stopwatch "github.com/RAshkettle/Stopwatch"
	"github.com/hajimehoshi/ebiten/v2"
)

type GameScene struct {
	sceneManager   *SceneManager
	gameboard      *Gameboard
	blockManager   *BlockManager
	gameLogic      *GameLogic
	inputHandler   *InputHandler
	renderer       *GameRenderer
	currentPiece   *TetrisPiece
	currentType    PieceType
	fallTimer      *stopwatch.Stopwatch
	CurrentScore   int
}

func (g *GameScene) Update() error {
	// Update the fall timer
	g.fallTimer.Update()
	
	// Handle input
	g.inputHandler.HandleInput(g.currentPiece, g.currentType)

	// Handle automatic falling (every 1 second)
	if g.fallTimer.IsDone() {
		if g.currentPiece != nil {
			if g.gameLogic.TryMovePiece(g.currentPiece, 0, 1) {
				// Piece fell successfully
			} else {
				// Piece can't fall further, place it
				g.gameLogic.PlacePiece(g.currentPiece)
				
				// Check for game over condition
				if g.gameLogic.IsGameOver() {
					// Transition to end scene with current score
					g.sceneManager.TransitionToEndScreen(g.CurrentScore)
					return nil
				}
				
				g.spawnNewPiece()
			}
		}
		g.fallTimer.Reset()
		g.fallTimer.Start()
	}

	return nil
}

func (g *GameScene) Draw(screen *ebiten.Image) {
	g.renderer.Render(screen, g.gameLogic.GetPlacedBlocks(), g.currentPiece)
}

func (g *GameScene) spawnNewPiece() {
	// Generate random piece type
	pieceTypes := []PieceType{IPiece, OPiece, TPiece, SPiece, ZPiece, JPiece, LPiece}
	g.currentType = pieceTypes[rand.Intn(len(pieceTypes))]
	
	g.currentPiece = g.gameLogic.SpawnNewPiece(g.currentType)
	
	// Check if the new piece can be placed at its spawn position
	if g.currentPiece != nil && !g.gameLogic.IsValidPosition(g.currentPiece, 0, 0) {
		// Game over - new piece can't be placed
		g.sceneManager.TransitionToEndScreen(g.CurrentScore)
	}
}

func (g *GameScene) Layout(outerWidth, outerHeight int) (int, int) {
	// Update gameboard scaling when layout changes
	g.gameboard.UpdateScale(outerWidth, outerHeight)
	return outerWidth, outerHeight
}

func NewGameScene(sm *SceneManager) *GameScene {
	// Create fall timer (1 second intervals)
	fallTimer := stopwatch.NewStopwatch(1 * time.Second)
	fallTimer.Start()
	
	// Create components
	gameboard := NewGameboard(192, 320) // 192px wide, 320px tall
	blockManager := NewBlockManager()
	gameLogic := NewGameLogic(gameboard, blockManager)
	inputHandler := NewInputHandler(gameLogic)
	renderer := NewGameRenderer(gameboard, blockManager)
	
	g := &GameScene{
		sceneManager:  sm,
		gameboard:     gameboard,
		blockManager:  blockManager,
		gameLogic:     gameLogic,
		inputHandler:  inputHandler,
		renderer:      renderer,
		fallTimer:     fallTimer,
		CurrentScore:  0,
	}

	// Spawn initial piece
	g.spawnNewPiece()

	return g
}
