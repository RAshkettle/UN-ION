package main

import (
	_ "embed"
	"fmt"
	"image/color"
	"math/rand"
	"time"

	stopwatch "github.com/RAshkettle/Stopwatch"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//go:embed shaders/electrical_storm.kage
var electricalStormShader []byte

// min returns the smaller of two values
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// Gameboard represents the main playing area for Tetris pieces
type Gameboard struct {
	Width       int
	Height      int
	X           int // X position on screen
	Y           int // Y position on screen
	shader      *ebiten.Shader
	startTime   time.Time
	baseWidth   int // Original design width (192)
	baseHeight  int // Original design height (320)
}

// NewGameboard creates a new gameboard with the specified dimensions and position
func NewGameboard(baseWidth, baseHeight int) *Gameboard {
	// Compile the electrical storm shader
	shader, err := ebiten.NewShader(electricalStormShader)
	if err != nil {
		panic(fmt.Sprintf("Failed to compile electrical storm shader: %v", err))
	}
	
	return &Gameboard{
		baseWidth:  baseWidth,
		baseHeight: baseHeight,
		Width:      baseWidth,
		Height:     baseHeight,
		shader:     shader,
		startTime:  time.Now(),
	}
}

// UpdateScale updates the gameboard size and position based on screen dimensions
func (gb *Gameboard) UpdateScale(screenWidth, screenHeight int) {
	// Calculate scale factor to maintain aspect ratio
	scaleX := float64(screenWidth) / 320.0  // 320 is our base screen width
	scaleY := float64(screenHeight) / 320.0 // 320 is our base screen height
	scale := min(scaleX, scaleY) // Use smaller scale to fit both dimensions
	
	// Apply scale to gameboard dimensions
	gb.Width = int(float64(gb.baseWidth) * scale)
	gb.Height = int(float64(gb.baseHeight) * scale)
	
	// Center the gameboard horizontally
	gb.X = (screenWidth - gb.Width) / 2
	gb.Y = 0 // Keep at top of screen
}

// Draw renders the gameboard on the screen with shader effect
func (gb *Gameboard) Draw(screen *ebiten.Image) {
	// Create a temporary image for the gameboard
	gameboardImage := ebiten.NewImage(gb.Width, gb.Height)
	
	// Calculate time for shader animation
	elapsed := time.Since(gb.startTime).Seconds()
	
	// Apply the electrical storm shader
	op := &ebiten.DrawTrianglesShaderOptions{}
	op.Uniforms = map[string]interface{}{
		"Time":       float32(elapsed),
		"Resolution": []float32{float32(gb.Width), float32(gb.Height)},
	}
	
	// Create vertices for a full-screen quad
	vertices := []ebiten.Vertex{
		{DstX: 0, DstY: 0, SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: float32(gb.Width), DstY: 0, SrcX: float32(gb.Width), SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: 0, DstY: float32(gb.Height), SrcX: 0, SrcY: float32(gb.Height), ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: float32(gb.Width), DstY: float32(gb.Height), SrcX: float32(gb.Width), SrcY: float32(gb.Height), ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
	}
	
	indices := []uint16{0, 1, 2, 1, 2, 3}
	
	gameboardImage.DrawTrianglesShader(vertices, indices, gb.shader, op)
	
	// Draw the gameboard at its position on the screen
	drawOp := &ebiten.DrawImageOptions{}
	drawOp.GeoM.Translate(float64(gb.X), float64(gb.Y))
	screen.DrawImage(gameboardImage, drawOp)
}

type GameScene struct {
	sceneManager   *SceneManager
	gameboard      *Gameboard
	blockManager   *BlockManager
	currentPiece   *TetrisPiece
	currentType    PieceType
	lastSpawnTime  time.Time
	spawnInterval  time.Duration
	fallTimer      *stopwatch.Stopwatch
	placedBlocks   []Block // Blocks that have been placed on the board
}

func (g *GameScene) Update() error {
	// Update the fall timer
	g.fallTimer.Update()
	
	// Handle piece rotation
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if g.currentPiece != nil {
			// Store original rotation in case we need to revert
			originalRotation := g.currentPiece.Rotation
			originalBlocks := make([]Block, len(g.currentPiece.Blocks))
			copy(originalBlocks, g.currentPiece.Blocks)
			
			g.blockManager.RotatePiece(g.currentPiece, g.currentType)
			
			// Check if rotation is valid
			if !g.isValidPosition(g.currentPiece, 0, 0) {
				// Revert rotation
				g.currentPiece.Rotation = originalRotation
				g.currentPiece.Blocks = originalBlocks
			}
		}
	}

	// Handle piece movement (16 pixels = 1 block)
	if inpututil.IsKeyJustPressed(ebiten.KeyA) || inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		if g.currentPiece != nil && g.isValidPosition(g.currentPiece, -1, 0) {
			g.currentPiece.X--
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) || inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		if g.currentPiece != nil && g.isValidPosition(g.currentPiece, 1, 0) {
			g.currentPiece.X++
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) || inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		if g.currentPiece != nil && g.isValidPosition(g.currentPiece, 0, 1) {
			g.currentPiece.Y++
		}
	}

	// Handle automatic falling (every 1 second)
	if g.fallTimer.IsDone() {
		if g.currentPiece != nil {
			if g.isValidPosition(g.currentPiece, 0, 1) {
				g.currentPiece.Y++ // Fall 1 block (16 pixels when scaled)
			} else {
				// Piece can't fall further, place it
				g.placePiece()
			}
		}
		g.fallTimer.Reset()
		g.fallTimer.Start()
	}

	return nil
}

func (g *GameScene) Draw(screen *ebiten.Image) {
	// Dark background
	screen.Fill(color.RGBA{15, 20, 30, 255})

	// Draw the gameboard with shader effect FIRST (background)
	g.gameboard.Draw(screen)
	
	// Calculate block size for rendering
	blockSize := g.blockManager.GetScaledBlockSize(g.gameboard.Width, g.gameboard.Height)
	
	// Create a temporary image for all blocks
	blocksImage := ebiten.NewImage(g.gameboard.Width, g.gameboard.Height)
	
	// Draw placed blocks first
	for _, block := range g.placedBlocks {
		worldX := float64(block.X) * blockSize
		worldY := float64(block.Y) * blockSize
		g.blockManager.DrawBlock(blocksImage, block, worldX, worldY, blockSize)
	}
	
	// Draw current piece on top of placed blocks
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

func (g *GameScene) spawnNewPiece() {
	// Generate random piece type
	pieceTypes := []PieceType{IPiece, OPiece, TPiece, SPiece, ZPiece, JPiece, LPiece}
	g.currentType = pieceTypes[rand.Intn(len(pieceTypes))]
	
	// Spawn piece at top center of gameboard
	// Use the same block size calculation as collision detection
	blockSize := g.blockManager.GetScaledBlockSize(g.gameboard.Width, g.gameboard.Height)
	gameboardWidthInBlocks := int(float64(g.gameboard.Width) / blockSize)
	centerX := gameboardWidthInBlocks / 2
	
	g.currentPiece = g.blockManager.CreateTetrisPiece(g.currentType, centerX, 0)
}

// isValidPosition checks if a piece can be placed at the given position
func (g *GameScene) isValidPosition(piece *TetrisPiece, offsetX, offsetY int) bool {
	// Use the actual scaled block size that's being used for rendering
	blockSize := g.blockManager.GetScaledBlockSize(g.gameboard.Width, g.gameboard.Height)
	
	// Calculate grid dimensions based on actual gameboard size and block size
	gameboardWidthInBlocks := int(float64(g.gameboard.Width) / blockSize)
	gameboardHeightInBlocks := int(float64(g.gameboard.Height) / blockSize)
	
	for _, block := range piece.Blocks {
		newX := piece.X + block.X + offsetX
		newY := piece.Y + block.Y + offsetY
		
		// Check boundaries
		if newX < 0 || newX >= gameboardWidthInBlocks || newY >= gameboardHeightInBlocks {
			return false
		}
		
		// Check collision with placed blocks
		for _, placedBlock := range g.placedBlocks {
			if placedBlock.X == newX && placedBlock.Y == newY {
				return false
			}
		}
	}
	
	return true
}

// placePiece adds the current piece to the placed blocks
func (g *GameScene) placePiece() {
	if g.currentPiece == nil {
		return
	}
	
	for _, block := range g.currentPiece.Blocks {
		placedBlock := Block{
			X:         g.currentPiece.X + block.X,
			Y:         g.currentPiece.Y + block.Y,
			BlockType: block.BlockType,
		}
		g.placedBlocks = append(g.placedBlocks, placedBlock)
	}
	
	// Spawn new piece
	g.spawnNewPiece()
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
	
	g := &GameScene{
		sceneManager:  sm,
		gameboard:     NewGameboard(192, 320), // 192px wide, 320px tall
		blockManager:  NewBlockManager(),
		fallTimer:     fallTimer,
		placedBlocks:  make([]Block, 0),
		lastSpawnTime: time.Now(),
		spawnInterval: time.Hour, // Long interval since we spawn on demand now
	}

	// Spawn initial piece
	g.spawnNewPiece()

	return g
}
