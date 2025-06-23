package main

import (
	"image/color"
	"math/rand"
	"time"

	stopwatch "github.com/RAshkettle/Stopwatch"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	GameboardWidth  = 192
	GameboardHeight = 320
	FallInterval    = 1
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
	pauseController *PauseController
	gameState       *GameState
	currentPiece    *TetrisPiece
	currentType     PieceType
	nextPiece       *TetrisPiece
	nextType        PieceType
	fallTimer       *stopwatch.Stopwatch
	CurrentScore    int
	lastUpdateTime  time.Time

	tempImage     *ebiten.Image
	particleImage *ebiten.Image
	blocksImage   *ebiten.Image

	shakeOp    *ebiten.DrawImageOptions
	particleOp *ebiten.DrawImageOptions
	blocksOp   *ebiten.DrawImageOptions
}

func (g *GameScene) Update() error {
	g.pauseController.Update()

	now := time.Now()
	if g.lastUpdateTime.IsZero() {
		g.lastUpdateTime = now
	}
	dt := now.Sub(g.lastUpdateTime).Seconds()
	g.lastUpdateTime = now

	g.updateVisualEffects(dt)

	if g.gameState.IsPaused {
		return nil
	}

	g.fallTimer.Update()

	g.updateWobblingBlocks(dt)

	shouldPlacePiece := g.inputHandler.HandleInput(g.currentPiece, g.currentType)

	if shouldPlacePiece && g.currentPiece != nil {
		g.placePieceAndCheckReactions()
		return nil
	}

	if g.fallTimer.IsDone() {
		if g.currentPiece != nil {
			if g.gameLogic.TryMovePiece(g.currentPiece, 0, 1) {
			} else {
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
	shakeX, shakeY := g.screenShake.GetOffset()

	screenW, screenH := screen.Bounds().Dx(), screen.Bounds().Dy()
	if g.tempImage == nil || g.tempImage.Bounds().Dx() != screenW || g.tempImage.Bounds().Dy() != screenH {
		g.tempImage = ebiten.NewImage(screenW, screenH)
		g.particleImage = ebiten.NewImage(screenW, screenH)
	}

	g.tempImage.Clear()

	var shadowPiece *TetrisPiece
	if g.currentPiece != nil {
		shadowPiece = g.gameLogic.CalculateDropPosition(g.currentPiece)
	}

	g.renderGameWithShadow(g.tempImage, shadowPiece)
	g.renderer.RenderScore(g.tempImage, g.CurrentScore)
	g.renderNextPiecePreview(g.tempImage)

	if g.scorePopups != nil {
		g.scorePopups.Draw(g.tempImage)
	}

	g.shakeOp.GeoM.Reset()
	g.shakeOp.GeoM.Translate(shakeX, shakeY)
	screen.DrawImage(g.tempImage, g.shakeOp)

	if g.particleSystem != nil {
		g.particleOp.GeoM.Reset()
		g.particleOp.GeoM.Translate(shakeX, shakeY)

		g.particleImage.Clear()
		g.particleSystem.Draw(g.particleImage)
		screen.DrawImage(g.particleImage, g.particleOp)
	}

	g.pauseController.Draw(screen)
}

func (g *GameScene) renderGameWithShadow(screen *ebiten.Image, shadowPiece *TetrisPiece) {
	screen.Fill(color.RGBA{15, 20, 30, 255})

	g.gameboard.Draw(screen)

	blockSize := g.blockManager.GetScaledBlockSize(g.gameboard.Width, g.gameboard.Height)

	if g.blocksImage == nil || g.blocksImage.Bounds().Dx() != g.gameboard.Width || g.blocksImage.Bounds().Dy() != g.gameboard.Height {
		g.blocksImage = ebiten.NewImage(g.gameboard.Width, g.gameboard.Height)
	}

	g.blocksImage.Clear()

	for _, block := range g.gameLogic.GetPlacedBlocks() {
		renderX, renderY, rotation, scale := g.gameLogic.GetBlockRenderTransform(&block)
		worldX := renderX * blockSize
		worldY := renderY * blockSize

		if block.IsArcing || rotation != 0.0 || scale != 1.0 {
			g.blockManager.DrawBlockTransformed(g.blocksImage, block, worldX, worldY, rotation, scale, blockSize)
		} else {
			g.blockManager.DrawBlock(g.blocksImage, block, worldX, worldY, blockSize)
		}
	}

	if shadowPiece != nil && g.currentPiece != nil && shadowPiece.Y > g.currentPiece.Y {
		for _, block := range shadowPiece.Blocks {
			worldX := float64(shadowPiece.X+block.X) * blockSize
			worldY := float64(shadowPiece.Y+block.Y) * blockSize
			g.blockManager.DrawShadowBlock(g.blocksImage, block, worldX, worldY, blockSize)
		}
	}
	if g.currentPiece != nil {
		for _, block := range g.currentPiece.Blocks {
			worldX := float64(g.currentPiece.X+block.X) * blockSize
			worldY := float64(g.currentPiece.Y+block.Y) * blockSize
			g.blockManager.DrawBlock(g.blocksImage, block, worldX, worldY, blockSize)
		}
	}

	warnings := g.gameLogic.GetStormWarnings()
	for _, warning := range warnings {
		worldX := float64(warning.Column) * blockSize
		worldY := float64(warning.HighestBlockY) * blockSize
		g.blockManager.DrawWarningSprite(g.blocksImage, worldX, worldY, warning.WarningTime, blockSize, warning.Column, g.gameboard.Width)
	}

	g.blocksOp.GeoM.Reset()
	g.blocksOp.GeoM.Translate(float64(g.gameboard.X), float64(g.gameboard.Y))
	screen.DrawImage(g.blocksImage, g.blocksOp)
}

func (g *GameScene) renderNextPiecePreview(screen *ebiten.Image) {
	if g.nextPiece == nil {
		return
	}

	screenWidth, _ := screen.Bounds().Dx(), screen.Bounds().Dy()

	previewX := float64(g.gameboard.X + g.gameboard.Width + 20)
	previewY := float64(g.gameboard.Y + 100)

	blockSize := g.blockManager.GetScaledBlockSize(g.gameboard.Width, g.gameboard.Height)
	previewBlockSize := blockSize * 0.6

	if previewX+previewBlockSize*4 < float64(screenWidth) {
		for _, block := range g.nextPiece.Blocks {
			worldX := previewX + float64(block.X)*previewBlockSize
			worldY := previewY + float64(block.Y)*previewBlockSize

			g.blockManager.DrawBlock(screen, block, worldX, worldY, previewBlockSize)
		}
	}
}

func (g *GameScene) spawnNewPiece() {
	if g.nextPiece != nil {
		g.currentType = g.nextType
		g.currentPiece = g.copyPieceForGameplay(g.nextPiece, g.currentType)
	} else {
		pieceTypes := []PieceType{IPiece, OPiece, TPiece, SPiece, ZPiece, JPiece, LPiece}
		g.currentType = pieceTypes[rand.Intn(len(pieceTypes))]
		g.currentPiece = g.gameLogic.SpawnNewPiece(g.currentType)
	}

	g.generateNextPiece()

	if g.currentPiece != nil && !g.gameLogic.IsValidPositionIgnoreNeutral(g.currentPiece, 0, 0) {
		g.sceneManager.TransitionToEndScreen(g.CurrentScore)
	}
}

func (g *GameScene) generateNextPiece() {
	pieceTypes := []PieceType{IPiece, OPiece, TPiece, SPiece, ZPiece, JPiece, LPiece}
	g.nextType = pieceTypes[rand.Intn(len(pieceTypes))]

	g.nextPiece = g.blockManager.CreateTetrisPiece(g.nextType, 0, 0)
}

func (g *GameScene) Layout(outerWidth, outerHeight int) (int, int) {
	g.gameboard.UpdateScale(outerWidth, outerHeight)
	return outerWidth, outerHeight
}

func NewGameScene(sm *SceneManager) *GameScene {
	fallTimer := stopwatch.NewStopwatch(FallInterval * time.Second)
	fallTimer.Start()

	gameboard := NewGameboard(GameboardWidth, GameboardHeight)
	blockManager := NewBlockManager()
	gameLogic := NewGameLogic(gameboard, blockManager)
	audioManager := NewAudioManager()
	inputHandler := NewInputHandler(gameLogic, audioManager)
	renderer := NewGameRenderer(gameboard, blockManager)
	particleSystem := NewParticleSystem()
	screenShake := NewScreenShake()
	scorePopups := NewScorePopupSystem()
	gameState := NewGameState()
	pauseController := NewPauseController(gameState, audioManager)

	err := audioManager.Initialize()
	if err != nil {
		println("Warning: Could not initialize audio:", err.Error())
	}

	g := &GameScene{
		sceneManager:    sm,
		gameboard:       gameboard,
		blockManager:    blockManager,
		gameLogic:       gameLogic,
		inputHandler:    inputHandler,
		renderer:        renderer,
		particleSystem:  particleSystem,
		audioManager:    audioManager,
		screenShake:     screenShake,
		scorePopups:     scorePopups,
		pauseController: pauseController,
		gameState:       gameState,
		fallTimer:       fallTimer,
		CurrentScore:    0,
		lastUpdateTime:  time.Now(),

		shakeOp:    &ebiten.DrawImageOptions{},
		particleOp: &ebiten.DrawImageOptions{},
		blocksOp:   &ebiten.DrawImageOptions{},
	}

	gameLogic.SetExplosionCallback(func(worldX, worldY float64, blockType BlockType) {
		particleSystem.AddExplosion(worldX, worldY, blockType)
	})

	gameLogic.SetAudioCallback(func(blocksRemoved int) {
		audioManager.PlayBlockBreakMultiple(blocksRemoved)

		intensity := float64(blocksRemoved) * 2.0
		duration := 0.2 + float64(blocksRemoved)*0.05
		screenShake.StartShake(intensity, duration)
	})

	gameLogic.SetDustCallback(func(worldX, worldY float64) {
		particleSystem.AddDustCloud(worldX, worldY)
	})

	gameLogic.SetHardDropCallback(func(dropHeight int) {
		intensity := 1.0 + float64(dropHeight)*0.5
		duration := 0.1
		screenShake.StartShake(intensity, duration)
	})

	g.generateNextPiece()

	g.spawnNewPiece()

	audioManager.StartBackgroundMusic()

	return g
}

func (g *GameScene) copyPieceForGameplay(piece *TetrisPiece, pieceType PieceType) *TetrisPiece {
	blocksCopy := make([]Block, len(piece.Blocks))
	copy(blocksCopy, piece.Blocks)

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

func (g *GameScene) updateWobblingBlocks(dt float64) {
	anyBlocksFinishedArcing := g.gameLogic.UpdateArcingBlocks(dt)

	anyBlocksLanded := g.gameLogic.UpdateFallingBlocks(dt)

	anyBlocksFinished := g.gameLogic.UpdateWobblingBlocks(dt)

	g.gameLogic.UpdateElectricalStorms(dt)
	g.gameLogic.ClearInvalidStorms()

	newNeutralBlocks := g.gameLogic.UpdateStormTimers(dt)
	for _, neutralBlock := range newNeutralBlocks {
		g.gameLogic.AddNeutralBlock(neutralBlock)

		blockSize := g.blockManager.GetScaledBlockSize(g.gameboard.Width, g.gameboard.Height)
		worldX := float64(g.gameboard.X) + float64(neutralBlock.X)*blockSize + blockSize/2
		worldY := float64(g.gameboard.Y) + float64(neutralBlock.Y)*blockSize + blockSize/2
		g.particleSystem.AddDustCloud(worldX, worldY)
	}

	if len(newNeutralBlocks) > 0 {
		g.gameLogic.processBlockFalling()
	}

	if anyBlocksLanded || anyBlocksFinishedArcing {
		reactionScore := g.gameLogic.CheckForNewReactions()
		g.gameLogic.CheckForElectricalStorms()

		if reactionScore > 0 {
			g.CurrentScore += reactionScore
		}
	}

	if anyBlocksFinished {
		removedCount := g.gameLogic.RemoveFinishedWobblingBlocks()

		if removedCount > 0 {
			g.gameLogic.processBlockFalling()

			reactionScore := g.gameLogic.CheckForNewReactions()
			g.gameLogic.CheckForElectricalStorms()

			if reactionScore > 0 {
				g.CurrentScore += reactionScore

				popupX := float64(g.gameboard.X + g.gameboard.Width/2)
				popupY := float64(g.gameboard.Y + g.gameboard.Height/3)
				g.scorePopups.AddScorePopup(popupX, popupY, reactionScore)
			}
		}
	}
}

func (g *GameScene) placePieceAndCheckReactions() {
	if g.currentPiece == nil {
		return
	}

	g.gameLogic.PlacePiece(g.currentPiece)

	reactionScore := g.gameLogic.CheckForNewReactions()

	g.gameLogic.CheckForElectricalStorms()

	if reactionScore > 0 {
		g.CurrentScore += reactionScore

		popupX := float64(g.gameboard.X + g.gameboard.Width/2)
		popupY := float64(g.gameboard.Y + g.gameboard.Height/3)
		g.scorePopups.AddScorePopup(popupX, popupY, reactionScore)
	}

	if g.gameLogic.IsGameOver() {
		g.sceneManager.TransitionToEndScreen(g.CurrentScore)
		return
	}

	g.spawnNewPiece()
}
