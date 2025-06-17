package main

import (
	"math/rand"
	"union/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

// BlockType represents the three different charge types
type BlockType int

const (
	PositiveBlock BlockType = iota
	NegativeBlock
	NeutralBlock
)

// Block represents a single block in a Tetris piece
type Block struct {
	X, Y      int
	BlockType BlockType
}

// TetrisPiece represents a complete Tetris piece with multiple blocks
type TetrisPiece struct {
	Blocks   []Block
	X, Y     int // Position of the piece
	Rotation int // Current rotation state (0-3)
}

// PieceType represents the different Tetris piece shapes
type PieceType int

const (
	IPiece PieceType = iota // Straight line
	OPiece                  // Square
	TPiece                  // T-shape
	SPiece                  // S-shape
	ZPiece                  // Z-shape
	JPiece                  // J-shape
	LPiece                  // L-shape
)

// BlockManager handles creation and rendering of Tetris pieces
type BlockManager struct {
	blockSize float64
}

// NewBlockManager creates a new block manager
func NewBlockManager() *BlockManager {
	return &BlockManager{
		blockSize: 16.0, // Base size of 16x16 pixels
	}
}

// GetScaledBlockSize returns the block size scaled for the current gameboard
func (bm *BlockManager) GetScaledBlockSize(gameboardWidth, gameboardHeight int) float64 {
	// Calculate how many blocks should fit in the gameboard
	// We want 12 blocks wide (192/16) and 20 blocks tall (320/16)
	baseBlocksWide := 12.0
	baseBlocksTall := 20.0
	
	// Calculate block size based on gameboard dimensions
	blockSizeFromWidth := float64(gameboardWidth) / baseBlocksWide
	blockSizeFromHeight := float64(gameboardHeight) / baseBlocksTall
	
	// Use the smaller of the two to maintain aspect ratio
	if blockSizeFromWidth < blockSizeFromHeight {
		return blockSizeFromWidth
	}
	return blockSizeFromHeight
}

// GenerateRandomBlockType returns a random block type with specified probabilities
// Positive: 40%, Negative: 40%, Neutral: 20%
func (bm *BlockManager) GenerateRandomBlockType() BlockType {
	r := rand.Float64()
	if r < 0.4 {
		return PositiveBlock
	} else if r < 0.8 {
		return NegativeBlock
	}
	return NeutralBlock
}

// GetBlockSprite returns the appropriate sprite for a block type
func (bm *BlockManager) GetBlockSprite(blockType BlockType) *ebiten.Image {
	switch blockType {
	case PositiveBlock:
		return assets.PositiveChargeSprite
	case NegativeBlock:
		return assets.NegativeChargeSprite
	case NeutralBlock:
		return assets.NeutralChargeSprite
	default:
		return assets.NeutralChargeSprite
	}
}

// GetPieceBlocks returns the block positions for a given piece type and rotation
func (bm *BlockManager) GetPieceBlocks(pieceType PieceType, rotation int) []Block {
	var blocks []Block
	
	switch pieceType {
	case IPiece: // Straight line piece
		if rotation%2 == 0 {
			// Horizontal
			for i := 0; i < 4; i++ {
				blocks = append(blocks, Block{X: i, Y: 0, BlockType: bm.GenerateRandomBlockType()})
			}
		} else {
			// Vertical
			for i := 0; i < 4; i++ {
				blocks = append(blocks, Block{X: 0, Y: i, BlockType: bm.GenerateRandomBlockType()})
			}
		}
	
	case OPiece: // Square piece - now rotates to move charged blocks
		switch rotation % 4 {
		case 0: // Original position
			blocks = []Block{
				{X: 0, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 1, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 0, Y: 1, BlockType: bm.GenerateRandomBlockType()},
				{X: 1, Y: 1, BlockType: bm.GenerateRandomBlockType()},
			}
		case 1: // 90 degrees clockwise
			blocks = []Block{
				{X: 1, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 1, Y: 1, BlockType: bm.GenerateRandomBlockType()},
				{X: 0, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 0, Y: 1, BlockType: bm.GenerateRandomBlockType()},
			}
		case 2: // 180 degrees
			blocks = []Block{
				{X: 1, Y: 1, BlockType: bm.GenerateRandomBlockType()},
				{X: 0, Y: 1, BlockType: bm.GenerateRandomBlockType()},
				{X: 1, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 0, Y: 0, BlockType: bm.GenerateRandomBlockType()},
			}
		case 3: // 270 degrees clockwise
			blocks = []Block{
				{X: 0, Y: 1, BlockType: bm.GenerateRandomBlockType()},
				{X: 0, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 1, Y: 1, BlockType: bm.GenerateRandomBlockType()},
				{X: 1, Y: 0, BlockType: bm.GenerateRandomBlockType()},
			}
		}
	
	case TPiece: // T-shaped piece
		switch rotation % 4 {
		case 0: // T pointing up
			blocks = []Block{
				{X: 1, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 0, Y: 1, BlockType: bm.GenerateRandomBlockType()},
				{X: 1, Y: 1, BlockType: bm.GenerateRandomBlockType()},
				{X: 2, Y: 1, BlockType: bm.GenerateRandomBlockType()},
			}
		case 1: // T pointing right
			blocks = []Block{
				{X: 0, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 0, Y: 1, BlockType: bm.GenerateRandomBlockType()},
				{X: 1, Y: 1, BlockType: bm.GenerateRandomBlockType()},
				{X: 0, Y: 2, BlockType: bm.GenerateRandomBlockType()},
			}
		case 2: // T pointing down
			blocks = []Block{
				{X: 0, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 1, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 2, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 1, Y: 1, BlockType: bm.GenerateRandomBlockType()},
			}
		case 3: // T pointing left
			blocks = []Block{
				{X: 1, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 0, Y: 1, BlockType: bm.GenerateRandomBlockType()},
				{X: 1, Y: 1, BlockType: bm.GenerateRandomBlockType()},
				{X: 1, Y: 2, BlockType: bm.GenerateRandomBlockType()},
			}
		}
	
	case SPiece: // S-shaped piece
		if rotation%2 == 0 {
			blocks = []Block{
				{X: 1, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 2, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 0, Y: 1, BlockType: bm.GenerateRandomBlockType()},
				{X: 1, Y: 1, BlockType: bm.GenerateRandomBlockType()},
			}
		} else {
			blocks = []Block{
				{X: 0, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 0, Y: 1, BlockType: bm.GenerateRandomBlockType()},
				{X: 1, Y: 1, BlockType: bm.GenerateRandomBlockType()},
				{X: 1, Y: 2, BlockType: bm.GenerateRandomBlockType()},
			}
		}
	
	case ZPiece: // Z-shaped piece
		if rotation%2 == 0 {
			blocks = []Block{
				{X: 0, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 1, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 1, Y: 1, BlockType: bm.GenerateRandomBlockType()},
				{X: 2, Y: 1, BlockType: bm.GenerateRandomBlockType()},
			}
		} else {
			blocks = []Block{
				{X: 1, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 0, Y: 1, BlockType: bm.GenerateRandomBlockType()},
				{X: 1, Y: 1, BlockType: bm.GenerateRandomBlockType()},
				{X: 0, Y: 2, BlockType: bm.GenerateRandomBlockType()},
			}
		}
	
	case JPiece: // J-shaped piece
		switch rotation % 4 {
		case 0:
			blocks = []Block{
				{X: 0, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 0, Y: 1, BlockType: bm.GenerateRandomBlockType()},
				{X: 1, Y: 1, BlockType: bm.GenerateRandomBlockType()},
				{X: 2, Y: 1, BlockType: bm.GenerateRandomBlockType()},
			}
		case 1:
			blocks = []Block{
				{X: 1, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 1, Y: 1, BlockType: bm.GenerateRandomBlockType()},
				{X: 1, Y: 2, BlockType: bm.GenerateRandomBlockType()},
				{X: 0, Y: 2, BlockType: bm.GenerateRandomBlockType()},
			}
		case 2:
			blocks = []Block{
				{X: 0, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 1, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 2, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 2, Y: 1, BlockType: bm.GenerateRandomBlockType()},
			}
		case 3:
			blocks = []Block{
				{X: 0, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 1, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 0, Y: 1, BlockType: bm.GenerateRandomBlockType()},
				{X: 0, Y: 2, BlockType: bm.GenerateRandomBlockType()},
			}
		}
	
	case LPiece: // L-shaped piece
		switch rotation % 4 {
		case 0:
			blocks = []Block{
				{X: 2, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 0, Y: 1, BlockType: bm.GenerateRandomBlockType()},
				{X: 1, Y: 1, BlockType: bm.GenerateRandomBlockType()},
				{X: 2, Y: 1, BlockType: bm.GenerateRandomBlockType()},
			}
		case 1:
			blocks = []Block{
				{X: 0, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 0, Y: 1, BlockType: bm.GenerateRandomBlockType()},
				{X: 0, Y: 2, BlockType: bm.GenerateRandomBlockType()},
				{X: 1, Y: 2, BlockType: bm.GenerateRandomBlockType()},
			}
		case 2:
			blocks = []Block{
				{X: 0, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 1, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 2, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 0, Y: 1, BlockType: bm.GenerateRandomBlockType()},
			}
		case 3:
			blocks = []Block{
				{X: 0, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 1, Y: 0, BlockType: bm.GenerateRandomBlockType()},
				{X: 1, Y: 1, BlockType: bm.GenerateRandomBlockType()},
				{X: 1, Y: 2, BlockType: bm.GenerateRandomBlockType()},
			}
		}
	}
	
	return blocks
}

// GetPiecePositions returns only the block positions for a given piece type and rotation
// This is used for rotation to preserve block types while updating positions
func (bm *BlockManager) GetPiecePositions(pieceType PieceType, rotation int) []struct{ X, Y int } {
	var positions []struct{ X, Y int }
	
	switch pieceType {
	case IPiece: // Straight line piece
		if rotation%2 == 0 {
			// Horizontal
			for i := 0; i < 4; i++ {
				positions = append(positions, struct{ X, Y int }{X: i, Y: 0})
			}
		} else {
			// Vertical
			for i := 0; i < 4; i++ {
				positions = append(positions, struct{ X, Y int }{X: 0, Y: i})
			}
		}
	
	case OPiece: // Square piece - now rotates to move charged blocks
		switch rotation % 4 {
		case 0: // Original position
			positions = []struct{ X, Y int }{
				{X: 0, Y: 0},
				{X: 1, Y: 0},
				{X: 0, Y: 1},
				{X: 1, Y: 1},
			}
		case 1: // 90 degrees clockwise
			positions = []struct{ X, Y int }{
				{X: 1, Y: 0},
				{X: 1, Y: 1},
				{X: 0, Y: 0},
				{X: 0, Y: 1},
			}
		case 2: // 180 degrees
			positions = []struct{ X, Y int }{
				{X: 1, Y: 1},
				{X: 0, Y: 1},
				{X: 1, Y: 0},
				{X: 0, Y: 0},
			}
		case 3: // 270 degrees clockwise
			positions = []struct{ X, Y int }{
				{X: 0, Y: 1},
				{X: 0, Y: 0},
				{X: 1, Y: 1},
				{X: 1, Y: 0},
			}
		}
	
	case TPiece: // T-shaped piece
		switch rotation % 4 {
		case 0: // T pointing up
			positions = []struct{ X, Y int }{
				{X: 1, Y: 0},
				{X: 0, Y: 1},
				{X: 1, Y: 1},
				{X: 2, Y: 1},
			}
		case 1: // T pointing right
			positions = []struct{ X, Y int }{
				{X: 0, Y: 0},
				{X: 0, Y: 1},
				{X: 1, Y: 1},
				{X: 0, Y: 2},
			}
		case 2: // T pointing down
			positions = []struct{ X, Y int }{
				{X: 0, Y: 0},
				{X: 1, Y: 0},
				{X: 2, Y: 0},
				{X: 1, Y: 1},
			}
		case 3: // T pointing left
			positions = []struct{ X, Y int }{
				{X: 1, Y: 0},
				{X: 0, Y: 1},
				{X: 1, Y: 1},
				{X: 1, Y: 2},
			}
		}
	
	case SPiece: // S-shaped piece
		if rotation%2 == 0 {
			positions = []struct{ X, Y int }{
				{X: 1, Y: 0},
				{X: 2, Y: 0},
				{X: 0, Y: 1},
				{X: 1, Y: 1},
			}
		} else {
			positions = []struct{ X, Y int }{
				{X: 0, Y: 0},
				{X: 0, Y: 1},
				{X: 1, Y: 1},
				{X: 1, Y: 2},
			}
		}
	
	case ZPiece: // Z-shaped piece
		if rotation%2 == 0 {
			positions = []struct{ X, Y int }{
				{X: 0, Y: 0},
				{X: 1, Y: 0},
				{X: 1, Y: 1},
				{X: 2, Y: 1},
			}
		} else {
			positions = []struct{ X, Y int }{
				{X: 1, Y: 0},
				{X: 0, Y: 1},
				{X: 1, Y: 1},
				{X: 0, Y: 2},
			}
		}
	
	case JPiece: // J-shaped piece
		switch rotation % 4 {
		case 0:
			positions = []struct{ X, Y int }{
				{X: 0, Y: 0},
				{X: 0, Y: 1},
				{X: 1, Y: 1},
				{X: 2, Y: 1},
			}
		case 1:
			positions = []struct{ X, Y int }{
				{X: 1, Y: 0},
				{X: 1, Y: 1},
				{X: 1, Y: 2},
				{X: 0, Y: 2},
			}
		case 2:
			positions = []struct{ X, Y int }{
				{X: 0, Y: 0},
				{X: 1, Y: 0},
				{X: 2, Y: 0},
				{X: 2, Y: 1},
			}
		case 3:
			positions = []struct{ X, Y int }{
				{X: 0, Y: 0},
				{X: 1, Y: 0},
				{X: 0, Y: 1},
				{X: 0, Y: 2},
			}
		}
	
	case LPiece: // L-shaped piece
		switch rotation % 4 {
		case 0:
			positions = []struct{ X, Y int }{
				{X: 2, Y: 0},
				{X: 0, Y: 1},
				{X: 1, Y: 1},
				{X: 2, Y: 1},
			}
		case 1:
			positions = []struct{ X, Y int }{
				{X: 0, Y: 0},
				{X: 0, Y: 1},
				{X: 0, Y: 2},
				{X: 1, Y: 2},
			}
		case 2:
			positions = []struct{ X, Y int }{
				{X: 0, Y: 0},
				{X: 1, Y: 0},
				{X: 2, Y: 0},
				{X: 0, Y: 1},
			}
		case 3:
			positions = []struct{ X, Y int }{
				{X: 0, Y: 0},
				{X: 1, Y: 0},
				{X: 1, Y: 1},
				{X: 1, Y: 2},
			}
		}
	}
	
	return positions
}

// CreateTetrisPiece creates a new Tetris piece with random block types
func (bm *BlockManager) CreateTetrisPiece(pieceType PieceType, x, y int) *TetrisPiece {
	blocks := bm.GetPieceBlocks(pieceType, 0) // Start with rotation 0
	
	return &TetrisPiece{
		Blocks:   blocks,
		X:        x,
		Y:        y,
		Rotation: 0,
	}
}

// DrawBlock renders a single block at the specified position
func (bm *BlockManager) DrawBlock(screen *ebiten.Image, block Block, worldX, worldY float64, blockSize float64) {
	sprite := bm.GetBlockSprite(block.BlockType)
	
	op := &ebiten.DrawImageOptions{}
	
	// Scale the sprite to match the block size
	scaleX := blockSize / float64(sprite.Bounds().Dx())
	scaleY := blockSize / float64(sprite.Bounds().Dy())
	op.GeoM.Scale(scaleX, scaleY)
	
	// Position the block
	op.GeoM.Translate(worldX, worldY)
	
	screen.DrawImage(sprite, op)
}

// DrawTetrisPiece renders a complete Tetris piece
func (bm *BlockManager) DrawTetrisPiece(screen *ebiten.Image, piece *TetrisPiece, screenWidth, screenHeight int) {
	blockSize := bm.GetScaledBlockSize(screenWidth, screenHeight)
	
	for _, block := range piece.Blocks {
		worldX := float64(piece.X+block.X) * blockSize
		worldY := float64(piece.Y+block.Y) * blockSize
		
		bm.DrawBlock(screen, block, worldX, worldY, blockSize)
	}
}

// RotatePiece rotates a Tetris piece clockwise
func (bm *BlockManager) RotatePiece(piece *TetrisPiece, pieceType PieceType) {
	newRotation := (piece.Rotation + 1) % 4
	
	// Get the new block positions for this rotation
	newPositions := bm.GetPiecePositions(pieceType, newRotation)
	
	// Preserve the existing block types, only update positions
	if len(newPositions) == len(piece.Blocks) {
		for i, pos := range newPositions {
			piece.Blocks[i].X = pos.X
			piece.Blocks[i].Y = pos.Y
			// Keep piece.Blocks[i].BlockType unchanged
		}
		piece.Rotation = newRotation
	}
}

// TestBlockDistribution tests the probability distribution of block types
func (bm *BlockManager) TestBlockDistribution(numTests int) (positive, negative, neutral int) {
	for i := 0; i < numTests; i++ {
		blockType := bm.GenerateRandomBlockType()
		switch blockType {
		case PositiveBlock:
			positive++
		case NegativeBlock:
			negative++
		case NeutralBlock:
			neutral++
		}
	}
	return positive, negative, neutral
}