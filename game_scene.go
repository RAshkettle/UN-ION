package main

import (
	"math/rand"
	"time"

	stopwatch "github.com/RAshkettle/Stopwatch"
	"github.com/hajimehoshi/ebiten/v2"
)

type GameScene struct {
	sceneManager    *SceneManager
	gameboard       *Gameboard
	blockManager    *BlockManager
	gameLogic       *GameLogic
	inputHandler    *InputHandler
	renderer        *GameRenderer
	particleSystem  *ParticleSystem
	currentPiece    *TetrisPiece
	currentType     PieceType
	nextPiece       *TetrisPiece
	nextType        PieceType
	fallTimer       *stopwatch.Stopwatch
	CurrentScore    int
	lastUpdateTime  time.Time
}

func (g *GameScene) Update() error {
	// Calculate delta time
	now := time.Now()
	if g.lastUpdateTime.IsZero() {
		g.lastUpdateTime = now
	}
	dt := now.Sub(g.lastUpdateTime).Seconds()
	g.lastUpdateTime = now
	
	// Update particle system
	if g.particleSystem != nil {
		g.particleSystem.Update(dt)
	}

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

				// Process any chain reactions from placed blocks and add score
				reactionScore := g.gameLogic.CheckAndProcessReactions()
				g.CurrentScore += reactionScore

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
	g.renderer.RenderScore(screen, g.CurrentScore)
	g.renderNextPiecePreview(screen)
	
	// Render particles on top of everything
	if g.particleSystem != nil {
		g.particleSystem.Draw(screen)
	}
}

func (g *GameScene) renderNextPiecePreview(screen *ebiten.Image) {
	if g.nextPiece == nil {
		return
	}

	// Calculate preview position (to the right of the gameboard)
	screenWidth, _ := screen.Bounds().Dx(), screen.Bounds().Dy()

	// Position the preview to the right of the gameboard
	previewX := float64(g.gameboard.X + g.gameboard.Width + 20) // 20 pixels margin
	previewY := float64(g.gameboard.Y + 50)                     // 50 pixels from top of gameboard

	// Scale the preview blocks to be smaller
	blockSize := g.blockManager.GetScaledBlockSize(g.gameboard.Width, g.gameboard.Height)
	previewBlockSize := blockSize * 0.6 // Make preview blocks 60% of normal size

	// Only render if there's space on screen
	if previewX+previewBlockSize*4 < float64(screenWidth) {
		// Render each block of the next piece
		for _, block := range g.nextPiece.Blocks {
			worldX := previewX + float64(block.X)*previewBlockSize
			worldY := previewY + float64(block.Y)*previewBlockSize

			g.blockManager.DrawBlock(screen, block, worldX, worldY, previewBlockSize)
		}
	}
}

func (g *GameScene) spawnNewPiece() {
	// Use the next piece as current piece
	if g.nextPiece != nil {
		g.currentType = g.nextType
		// Copy the next piece and position it properly for gameplay
		g.currentPiece = g.copyPieceForGameplay(g.nextPiece, g.currentType)
	} else {
		// Fallback for first piece (shouldn't happen in normal flow)
		pieceTypes := []PieceType{IPiece, OPiece, TPiece, SPiece, ZPiece, JPiece, LPiece}
		g.currentType = pieceTypes[rand.Intn(len(pieceTypes))]
		g.currentPiece = g.gameLogic.SpawnNewPiece(g.currentType)
	}

	// Generate new next piece
	g.generateNextPiece()

	// Check if the current piece can be placed at its spawn position
	if g.currentPiece != nil && !g.gameLogic.IsValidPosition(g.currentPiece, 0, 0) {
		// Game over - new piece can't be placed
		g.sceneManager.TransitionToEndScreen(g.CurrentScore)
	}
}

func (g *GameScene) generateNextPiece() {
	// Generate random piece type for next piece
	pieceTypes := []PieceType{IPiece, OPiece, TPiece, SPiece, ZPiece, JPiece, LPiece}
	g.nextType = pieceTypes[rand.Intn(len(pieceTypes))]

	// Create the next piece at a preview position (we'll position it for display)
	g.nextPiece = g.blockManager.CreateTetrisPiece(g.nextType, 0, 0)
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
	particleSystem := NewParticleSystem()

	g := &GameScene{
		sceneManager:   sm,
		gameboard:      gameboard,
		blockManager:   blockManager,
		gameLogic:      gameLogic,
		inputHandler:   inputHandler,
		renderer:       renderer,
		particleSystem: particleSystem,
		fallTimer:      fallTimer,
		CurrentScore:   0,
		lastUpdateTime: time.Now(),
	}

	// Set up the explosion callback for particle effects
	gameLogic.SetExplosionCallback(func(worldX, worldY float64, blockType BlockType) {
		particleSystem.AddExplosion(worldX, worldY, blockType)
	})

	// Generate initial next piece
	g.generateNextPiece()

	// Spawn initial current piece
	g.spawnNewPiece()

	return g
}

// copyPieceForGameplay creates a copy of a piece and positions it for gameplay
func (g *GameScene) copyPieceForGameplay(piece *TetrisPiece, pieceType PieceType) *TetrisPiece {
	// Create a copy of the blocks
	blocksCopy := make([]Block, len(piece.Blocks))
	copy(blocksCopy, piece.Blocks)

	// Calculate spawn position (center of gameboard, top)
	blockSize := g.blockManager.GetScaledBlockSize(g.gameboard.Width, g.gameboard.Height)
	gameboardWidthInBlocks := int(float64(g.gameboard.Width) / blockSize)
	centerX := gameboardWidthInBlocks / 2

	return &TetrisPiece{
		Blocks:   blocksCopy,
		X:        centerX,
		Y:        0,
		Rotation: piece.Rotation,
	}
}
