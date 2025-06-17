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

type GameScene struct {
	sceneManager   *SceneManager
	blockManager   *BlockManager
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
	}
	
	for i, control := range controls {
		op3 := &text.DrawOptions{}
		op3.GeoM.Translate(10, float64(g.screenHeight-60+i*15))
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
	
	// Spawn piece at top center
	centerX := 5 // Approximate center of a typical Tetris board
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
	}

	// Spawn initial piece
	g.spawnNewPiece()

	return g
}
