package main

import (
	"bytes"
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

// Gameboard represents the main playing area for Tetris pieces
type Gameboard struct {
	Width  int
	Height int
	X      int // X position on screen
	Y      int // Y position on screen
}

// NewGameboard creates a new gameboard with the specified dimensions and position
func NewGameboard(x, y, width, height int) *Gameboard {
	return &Gameboard{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

// Draw renders the gameboard on the screen
func (gb *Gameboard) Draw(screen *ebiten.Image) {
	// Fill the gameboard area with light gray
	gameboardImage := ebiten.NewImage(gb.Width, gb.Height)
	gameboardImage.Fill(color.RGBA{200, 200, 200, 255}) // Light gray

	// Draw the gameboard at its position
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(gb.X), float64(gb.Y))
	screen.DrawImage(gameboardImage, op)
}

// GetBounds returns the gameboard boundaries
func (gb *Gameboard) GetBounds() (x, y, width, height int) {
	return gb.X, gb.Y, gb.Width, gb.Height
}

// Contains checks if a point is within the gameboard
func (gb *Gameboard) Contains(x, y int) bool {
	return x >= gb.X && x < gb.X+gb.Width && y >= gb.Y && y < gb.Y+gb.Height
}

// ToGridCoordinates converts screen coordinates to grid coordinates
// assuming each grid cell is blockSize pixels
func (gb *Gameboard) ToGridCoordinates(screenX, screenY int, blockSize float64) (gridX, gridY int) {
	relativeX := screenX - gb.X
	relativeY := screenY - gb.Y
	gridX = int(float64(relativeX) / blockSize)
	gridY = int(float64(relativeY) / blockSize)
	return
}

// ToScreenCoordinates converts grid coordinates to screen coordinates
func (gb *Gameboard) ToScreenCoordinates(gridX, gridY int, blockSize float64) (screenX, screenY int) {
	screenX = gb.X + int(float64(gridX)*blockSize)
	screenY = gb.Y + int(float64(gridY)*blockSize)
	return
}

type GameScene struct {
	sceneManager   *SceneManager
	blockManager   *BlockManager
	gameboard      *Gameboard
	currentPiece   *TetrisPiece
	currentType    PieceType
	lastSpawnTime  time.Time
	spawnInterval  time.Duration
	screenWidth    int
	screenHeight   int
	debugFont      *text.GoTextFace
}

func (g *GameScene) Update() error {
	// Handle piece rotation
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if g.currentPiece != nil {
			g.blockManager.RotatePiece(g.currentPiece, g.currentType)
		}
	}

	// Handle piece movement
	if inpututil.IsKeyJustPressed(ebiten.KeyA) || inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		if g.currentPiece != nil {
			g.currentPiece.X--
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) || inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		if g.currentPiece != nil {
			g.currentPiece.X++
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) || inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		if g.currentPiece != nil {
			g.currentPiece.Y++
		}
	}

	// Spawn new piece periodically for demonstration
	if time.Since(g.lastSpawnTime) > g.spawnInterval {
		g.spawnNewPiece()
		g.lastSpawnTime = time.Now()
	}

	return nil
}

func (g *GameScene) Draw(screen *ebiten.Image) {
	// Dark background
	screen.Fill(color.RGBA{15, 20, 30, 255})

	// Draw gameboard
	g.gameboard.Draw(screen)

	// Draw current piece if it exists
	if g.currentPiece != nil {
		g.blockManager.DrawTetrisPiece(screen, g.currentPiece, g.screenWidth, g.screenHeight)
	}

	// Draw debug information
	g.drawDebugInfo(screen)
}

func (g *GameScene) drawDebugInfo(screen *ebiten.Image) {
	if g.currentPiece == nil {
		return
	}

	// Draw piece type
	pieceNames := map[PieceType]string{
		IPiece: "I-Piece", OPiece: "O-Piece", TPiece: "T-Piece",
		SPiece: "S-Piece", ZPiece: "Z-Piece", JPiece: "J-Piece", LPiece: "L-Piece",
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(10, 20)
	op.ColorScale.ScaleWithColor(color.RGBA{255, 255, 255, 255})
	text.Draw(screen, fmt.Sprintf("Current: %s", pieceNames[g.currentType]), g.debugFont, op)

	// Draw block types
	y := 40
	for i, block := range g.currentPiece.Blocks {
		blockTypeNames := map[BlockType]string{
			PositiveBlock: "Positive", NegativeBlock: "Negative", NeutralBlock: "Neutral",
		}
		op2 := &text.DrawOptions{}
		op2.GeoM.Translate(10, float64(y))
		op2.ColorScale.ScaleWithColor(color.RGBA{200, 200, 200, 255})
		text.Draw(screen, fmt.Sprintf("Block %d: %s", i+1, blockTypeNames[block.BlockType]), g.debugFont, op2)
		y += 20
	}

	// Draw controls
	controls := []string{
		"WASD/Arrows: Move",
		"Space: Rotate",
		fmt.Sprintf("Gameboard: %dx%d at (%d,%d)", g.gameboard.Width, g.gameboard.Height, g.gameboard.X, g.gameboard.Y),
	}

	for i, control := range controls {
		op3 := &text.DrawOptions{}
		op3.GeoM.Translate(10, float64(g.screenHeight-75+i*15))
		op3.ColorScale.ScaleWithColor(color.RGBA{150, 150, 150, 255})
		text.Draw(screen, control, g.debugFont, op3)
	}
}

func (g *GameScene) Layout(outerWidth, outerHeight int) (int, int) {
	g.screenWidth = outerWidth
	g.screenHeight = outerHeight
	return outerWidth, outerHeight
}

func (g *GameScene) spawnNewPiece() {
	// Generate random piece type
	pieceTypes := []PieceType{IPiece, OPiece, TPiece, SPiece, ZPiece, JPiece, LPiece}
	g.currentType = pieceTypes[rand.Intn(len(pieceTypes))]
	
	// Spawn piece at top center of gameboard
	// Calculate center of gameboard in grid coordinates (assuming 16px blocks)
	blockSize := 16.0
	gameboardWidthInBlocks := int(float64(g.gameboard.Width) / blockSize)
	centerX := gameboardWidthInBlocks / 2
	
	g.currentPiece = g.blockManager.CreateTetrisPiece(g.currentType, centerX, 0)
}

func NewGameScene(sm *SceneManager) *GameScene {
	// Create debug font
	debugFontSource, _ := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	debugFont := &text.GoTextFace{
		Source: debugFontSource,
		Size:   12,
	}

	g := &GameScene{
		sceneManager:  sm,
		blockManager:  NewBlockManager(),
		lastSpawnTime: time.Now(),
		spawnInterval: 3 * time.Second, // Spawn new piece every 3 seconds for demo
		screenWidth:   320,
		screenHeight:  320,
		debugFont:     debugFont,
		gameboard:     NewGameboard(64, 0, 192, 320), // 192px wide, full height, centered horizontally
	}

	// Spawn initial piece
	g.spawnNewPiece()

	return g
}
