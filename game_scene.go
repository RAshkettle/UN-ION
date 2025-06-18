package main

import (
	"image/color"
	"math/rand"
	"time"

	stopwatch "github.com/RAshkettle/Stopwatch"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// Game constants
const (
	GameboardWidth  = 192 // pixels
	GameboardHeight = 320 // pixels
	FallInterval    = 1   // seconds
)

type GameScene struct {
	sceneManager    *SceneManager
	gameboard       *Gameboard
	blockManager    *BlockManager
	gameLogic       *GameLogic
	inputHandler    *InputHandler
	renderer        *GameRenderer
	particleSystem  *ParticleSystem
	audioManager    *AudioManager
	screenShake     *ScreenShake
	scorePopups     *ScorePopupSystem
	currentPiece    *TetrisPiece
	currentType     PieceType
	nextPiece       *TetrisPiece
	nextType        PieceType
	fallTimer       *stopwatch.Stopwatch
	CurrentScore    int
	lastUpdateTime  time.Time
	isPaused        bool
}

func (g *GameScene) Update() error {
	// Handle pause input (always check, even when paused)
	g.handlePauseInput()
	
	// Calculate delta time
	now := time.Now()
	if g.lastUpdateTime.IsZero() {
		g.lastUpdateTime = now
	}
	dt := now.Sub(g.lastUpdateTime).Seconds()
	g.lastUpdateTime = now
	
	// Update visual effects even when paused
	g.updateVisualEffects(dt)
	
	// Skip game logic if paused
	if g.isPaused {
		return nil
	}

	// Update the fall timer
	g.fallTimer.Update()

	// Update wobbling blocks and handle chain reactions
	g.updateWobblingBlocks(dt)

	// Handle input
	shouldPlacePiece := g.inputHandler.HandleInput(g.currentPiece, g.currentType)
	
	// If input handler detected piece should be placed immediately
	if shouldPlacePiece && g.currentPiece != nil {
		g.placePieceAndCheckReactions()
		return nil
	}

	// Handle automatic falling (every 1 second)
	if g.fallTimer.IsDone() {
		if g.currentPiece != nil {
			if g.gameLogic.TryMovePiece(g.currentPiece, 0, 1) {
				// Piece fell successfully
			} else {
				// Piece can't fall further, place it
				g.placePieceAndCheckReactions()
				return nil
			}
		}
		g.fallTimer.Reset()
		g.fallTimer.Start()
	}

	return nil
}

func (g *GameScene) Draw(screen *ebiten.Image) {
	// Get screen shake offset
	shakeX, shakeY := g.screenShake.GetOffset()
	
	// Create a temporary image for shaken content
	tempImage := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
	
	// Calculate drop shadow for current piece
	var shadowPiece *TetrisPiece
	if g.currentPiece != nil {
		shadowPiece = g.gameLogic.CalculateDropPosition(g.currentPiece)
	}
	
	// Render game with drop shadow
	g.renderGameWithShadow(tempImage, shadowPiece)
	g.renderer.RenderScore(tempImage, g.CurrentScore)
	g.renderNextPiecePreview(tempImage)
	
	// Render score popups first
	if g.scorePopups != nil {
		g.scorePopups.Draw(tempImage)
	}
	
	// Apply shake offset when drawing to screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(shakeX, shakeY)
	screen.DrawImage(tempImage, op)
	
	// Render particles ABSOLUTELY LAST, directly to screen with shake offset
	if g.particleSystem != nil {
		particleOp := &ebiten.DrawImageOptions{}
		particleOp.GeoM.Translate(shakeX, shakeY)
		
		// Create a temporary image just for particles
		particleImage := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
		g.particleSystem.Draw(particleImage)
		screen.DrawImage(particleImage, particleOp)
	}
	
	// Draw pause overlay if paused
	g.drawPauseOverlay(screen)
}

// renderGameWithShadow renders the game state including drop shadow
func (g *GameScene) renderGameWithShadow(screen *ebiten.Image, shadowPiece *TetrisPiece) {
	// Dark background
	screen.Fill(color.RGBA{15, 20, 30, 255})

	// Draw the gameboard with shader effect FIRST (background)
	g.gameboard.Draw(screen)

	// Calculate block size for rendering
	blockSize := g.blockManager.GetScaledBlockSize(g.gameboard.Width, g.gameboard.Height)

	// Create a temporary image for all blocks
	blocksImage := ebiten.NewImage(g.gameboard.Width, g.gameboard.Height)

	// Draw placed blocks first
	for _, block := range g.gameLogic.GetPlacedBlocks() {
		worldX := float64(block.X) * blockSize
		worldY := float64(block.Y) * blockSize
		g.blockManager.DrawBlock(blocksImage, block, worldX, worldY, blockSize)
	}

	// Draw drop shadow (if different from current piece position)
	if shadowPiece != nil && g.currentPiece != nil && shadowPiece.Y > g.currentPiece.Y {
		for _, block := range shadowPiece.Blocks {
			worldX := float64(shadowPiece.X+block.X) * blockSize
			worldY := float64(shadowPiece.Y+block.Y) * blockSize
			g.blockManager.DrawShadowBlock(blocksImage, block, worldX, worldY, blockSize)
		}
	}

	// Draw current piece on top of shadow and placed blocks
	if g.currentPiece != nil {
		for _, block := range g.currentPiece.Blocks {
			worldX := float64(g.currentPiece.X+block.X) * blockSize
			worldY := float64(g.currentPiece.Y+block.Y) * blockSize
			g.blockManager.DrawBlock(blocksImage, block, worldX, worldY, blockSize)
		}
	}

	// Draw all blocks on top of the gameboard
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.gameboard.X), float64(g.gameboard.Y))
	screen.DrawImage(blocksImage, op)
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
	fallTimer := stopwatch.NewStopwatch(FallInterval * time.Second)
	fallTimer.Start()

	// Create components
	gameboard := NewGameboard(GameboardWidth, GameboardHeight)
	blockManager := NewBlockManager()
	gameLogic := NewGameLogic(gameboard, blockManager)
	audioManager := NewAudioManager()
	inputHandler := NewInputHandler(gameLogic, audioManager)
	renderer := NewGameRenderer(gameboard, blockManager)
	particleSystem := NewParticleSystem()
	screenShake := NewScreenShake()
	scorePopups := NewScorePopupSystem()
	
	// Initialize audio
	err := audioManager.Initialize()
	if err != nil {
		// Log error but continue without audio
		println("Warning: Could not initialize audio:", err.Error())
	}

	g := &GameScene{
		sceneManager:   sm,
		gameboard:      gameboard,
		blockManager:   blockManager,
		gameLogic:      gameLogic,
		inputHandler:   inputHandler,
		renderer:       renderer,
		particleSystem: particleSystem,
		audioManager:   audioManager,
		screenShake:    screenShake,
		scorePopups:    scorePopups,
		fallTimer:      fallTimer,
		CurrentScore:   0,
		lastUpdateTime: time.Now(),
	}

	// Set up the explosion callback for particle effects
	gameLogic.SetExplosionCallback(func(worldX, worldY float64, blockType BlockType) {
		particleSystem.AddExplosion(worldX, worldY, blockType)
	})

	// Set up the audio callback for block breaking sounds
	gameLogic.SetAudioCallback(func(blocksRemoved int) {
		audioManager.PlayBlockBreakMultiple(blocksRemoved)
		
		// Trigger screen shake based on number of blocks removed
		intensity := float64(blocksRemoved) * 2.0 // 2 pixels per block
		duration := 0.2 + float64(blocksRemoved)*0.05 // Longer shake for more blocks
		screenShake.StartShake(intensity, duration)
	})

	// Set up the dust callback for piece placement effects
	gameLogic.SetDustCallback(func(worldX, worldY float64) {
		particleSystem.AddDustCloud(worldX, worldY)
	})

	// Set up the hard drop callback for screen shake effects
	gameLogic.SetHardDropCallback(func(dropHeight int) {
		// Subtle screen shake for hard drops (much less intense than block explosions)
		intensity := 1.0 + float64(dropHeight)*0.5 // Subtle intensity
		duration := 0.1 // Short duration
		screenShake.StartShake(intensity, duration)
	})

	// Generate initial next piece
	g.generateNextPiece()

	// Spawn initial current piece
	g.spawnNewPiece()

	// Start background music at 10% volume
	audioManager.StartBackgroundMusic()

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

// handlePauseInput processes pause key input and manages pause state
func (g *GameScene) handlePauseInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.isPaused = !g.isPaused
		if g.isPaused {
			if g.audioManager != nil {
				g.audioManager.PauseBackgroundMusic()
			}
		} else {
			if g.audioManager != nil {
				g.audioManager.ResumeBackgroundMusic()
			}
		}
	}
}

// updateVisualEffects updates particles, screen shake, and score popups
func (g *GameScene) updateVisualEffects(dt float64) {
	if g.particleSystem != nil {
		g.particleSystem.Update(dt)
	}
	
	if g.screenShake != nil {
		g.screenShake.Update(dt)
	}
	
	if g.scorePopups != nil {
		g.scorePopups.Update(dt)
	}
}

// drawPauseOverlay renders the pause overlay with text
func (g *GameScene) drawPauseOverlay(screen *ebiten.Image) {
	if !g.isPaused {
		return
	}

	// Draw semi-transparent overlay
	overlay := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
	overlay.Fill(color.RGBA{0, 0, 0, 128}) // 50% transparent black
	screen.DrawImage(overlay, nil)

	// Draw "PAUSED" text in the center
	centerX := screen.Bounds().Dx() / 2
	centerY := screen.Bounds().Dy() / 2

	// Use basic font for text rendering
	fontFace := basicfont.Face7x13

	// Draw "PAUSED" text
	pausedText := "PAUSED"
	pausedBounds := text.BoundString(fontFace, pausedText)
	pausedX := centerX - pausedBounds.Dx()/2
	pausedY := centerY - 10
	text.Draw(screen, pausedText, fontFace, pausedX, pausedY, color.RGBA{255, 255, 255, 255})

	// Draw "Press P to Resume" below
	resumeText := "Press P to Resume"
	resumeBounds := text.BoundString(fontFace, resumeText)
	resumeX := centerX - resumeBounds.Dx()/2
	resumeY := centerY + 20
	text.Draw(screen, resumeText, fontFace, resumeX, resumeY, color.RGBA{200, 200, 200, 255})
}

// updateWobblingBlocks updates wobbling blocks and electrical storms, handles chain reactions
func (g *GameScene) updateWobblingBlocks(dt float64) {
	// Update wobbling animation
	anyBlocksFinished := g.gameLogic.UpdateWobblingBlocks(dt)
	
	// Update electrical storm animation (visual effects only, no removal)
	g.gameLogic.UpdateElectricalStorms(dt)
	
	// Handle finished wobbling blocks
	if anyBlocksFinished {
		removedCount := g.gameLogic.RemoveFinishedWobblingBlocks()
		
		if removedCount > 0 {
			// Make remaining blocks fall
			g.gameLogic.processBlockFalling()
			
			// Check for new chain reactions
			reactionScore := g.gameLogic.CheckForNewReactions()
			g.gameLogic.CheckForElectricalStorms() // Check for new storms (visual only)
			
			if reactionScore > 0 {
				g.CurrentScore += reactionScore
				
				// Add score popup for chain reaction
				popupX := float64(g.gameboard.X + g.gameboard.Width/2)
				popupY := float64(g.gameboard.Y + g.gameboard.Height/3)
				g.scorePopups.AddScorePopup(popupX, popupY, reactionScore)
			}
		}
	}
}

// placePieceAndCheckReactions places a piece and handles reactions/scoring
func (g *GameScene) placePieceAndCheckReactions() {
	if g.currentPiece == nil {
		return
	}

	// Place the piece
	g.gameLogic.PlacePiece(g.currentPiece)

	// Check for new horizontal reactions and start wobbling on blocks
	reactionScore := g.gameLogic.CheckForNewReactions()
	
	// Check for electrical storms (vertical sequences of 4+ same type) - visual effect only
	g.gameLogic.CheckForElectricalStorms()
	
	if reactionScore > 0 {
		g.CurrentScore += reactionScore
		
		// Add score popup at center of gameboard
		popupX := float64(g.gameboard.X + g.gameboard.Width/2)
		popupY := float64(g.gameboard.Y + g.gameboard.Height/3)
		g.scorePopups.AddScorePopup(popupX, popupY, reactionScore)
	}

	// Check for game over condition
	if g.gameLogic.IsGameOver() {
		// Transition to end scene with current score
		g.sceneManager.TransitionToEndScreen(g.CurrentScore)
		return
	}

	g.spawnNewPiece()
}
