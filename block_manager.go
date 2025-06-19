package main

import (
	"fmt"
	"math"
	"math/rand"
	"union/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	WobbleDuration   = 0.8   // Duration in seconds before block is destroyed
	WobbleIntensity  = 2.0   // Maximum wobble offset in pixels
	WobbleFrequency  = 8.0   // Wobbles per second
	StormDuration    = 1.2   // Duration in seconds before storm blocks are destroyed
	StormIntensity   = 3.0   // Maximum storm wobble offset in pixels
	StormFrequency   = 12.0  // Storm wobbles per second (faster than normal wobble)
	SparkFrequency   = 15.0  // Sparks per second
	FallSpeed        = 4.0   // Blocks per second falling speed
	ArcSpeed         = 4.0   // Arc animation speed (progress per second)
	ArcHeight        = 3.0   // Maximum arc height in blocks above start position
	MinArcScale      = 0.1   // Starting scale for arcing blocks
	MaxRotation      = 360.0 // Maximum rotation in degrees during arc
	WarningDuration  = 1.0   // Duration in seconds before spawn
	WarningIntensity = 3.0   // Maximum shake offset in pixels
	WarningFrequency = 20.0  // Shakes per second
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

// GetPiecePositions returns only the block positions for a given piece type and rotation
// This is used for rotation to preserve block type positions while updating positions
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
	blocks := GetPieceBlocks(pieceType, 0, bm.GenerateRandomBlockType)

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

	// Apply electrical storm effect (takes priority over normal wobble)
	if block.IsInStorm {
		// More intense wobble for electrical storms
		stormX := math.Sin(block.StormPhase) * StormIntensity
		stormY := math.Cos(block.StormPhase*1.7) * StormIntensity * 0.3 // Different frequency for Y

		// Add random sparking motion
		sparkOffset := math.Sin(block.SparkPhase) * 1.0
		stormX += sparkOffset
		stormY += math.Cos(block.SparkPhase*2.3) * 0.5

		// Position with storm offset
		op.GeoM.Translate(worldX+stormX, worldY+stormY)

		// Add electrical storm visual effects (no transparency fade since blocks don't get destroyed)
		// Flickering effect based on spark phase
		flickerIntensity := 0.2 + 0.1*math.Sin(block.SparkPhase*2)

		// Color modulation for electrical effect
		switch block.BlockType {
		case PositiveBlock:
			// Positive blocks get more red/yellow during storms
			op.ColorM.Scale(1.0+flickerIntensity, 1.0+flickerIntensity*0.5, 1.0-flickerIntensity*0.3, 1.0)
		case NegativeBlock:
			// Negative blocks get more blue/cyan during storms
			op.ColorM.Scale(1.0-flickerIntensity*0.3, 1.0+flickerIntensity*0.5, 1.0+flickerIntensity, 1.0)
		}

	} else if block.IsWobbling {
		// Normal wobble effect
		wobbleX := math.Sin(block.WobblePhase) * WobbleIntensity
		wobbleY := math.Cos(block.WobblePhase*1.3) * WobbleIntensity * 0.5

		// Position with wobble offset
		op.GeoM.Translate(worldX+wobbleX, worldY+wobbleY)

		// Add slight transparency to indicate impending destruction
		wobbleProgress := block.WobbleTime / WobbleDuration
		alpha := 1.0 - wobbleProgress*0.3 // Fade to 70% opacity
		op.ColorM.Scale(1, 1, 1, alpha)
	} else {
		// Position the block normally
		op.GeoM.Translate(worldX, worldY)
	}

	screen.DrawImage(sprite, op)

	// Draw POW sprite on top of wobbling blocks (but only when ShowPowSprite is true)
	if block.IsWobbling && block.ShowPowSprite {
		bm.DrawPowSprite(screen, worldX, worldY, block.WobblePhase, blockSize)
	}
}

// DrawShadowBlock renders a translucent shadow block at the specified world coordinates
func (bm *BlockManager) DrawShadowBlock(screen *ebiten.Image, block Block, worldX, worldY, blockSize float64) {
	var sprite *ebiten.Image

	switch block.BlockType {
	case PositiveBlock:
		sprite = assets.PositiveChargeSprite
	case NegativeBlock:
		sprite = assets.NegativeChargeSprite
	case NeutralBlock:
		sprite = assets.NeutralChargeSprite
	default:
		// Fallback for unknown block types
		sprite = assets.NeutralChargeSprite
	}

	if sprite != nil {
		op := &ebiten.DrawImageOptions{}
		// Scale the sprite to fit the block size
		spriteWidth := float64(sprite.Bounds().Dx())
		spriteHeight := float64(sprite.Bounds().Dy())
		scaleX := blockSize / spriteWidth
		scaleY := blockSize / spriteHeight
		op.GeoM.Scale(scaleX, scaleY)
		op.GeoM.Translate(worldX, worldY)

		// Make it translucent and darker for shadow effect
		op.ColorScale.ScaleAlpha(0.4)           // 40% opacity
		op.ColorScale.Scale(0.5, 0.5, 0.5, 1.0) // Darker

		screen.DrawImage(sprite, op)
	}
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

// RotatePiece rotates a Tetris piece clockwise while preserving block type positions
func (bm *BlockManager) RotatePiece(piece *TetrisPiece, pieceType PieceType) {
	// Special handling for O piece to preserve visual consistency
	if pieceType == OPiece {
		// For O piece, we want to rotate the charge pattern, not the physical positions
		newRotation := (piece.Rotation + 1) % 4
		newPositions := bm.GetPiecePositions(pieceType, newRotation)

		if len(newPositions) == len(piece.Blocks) {
			for i, pos := range newPositions {
				piece.Blocks[i].X = pos.X
				piece.Blocks[i].Y = pos.Y
				// Keep piece.Blocks[i].BlockType unchanged for O piece
			}
			piece.Rotation = newRotation
		}
		return
	}

	// For other pieces, we need to physically rotate while preserving block types at positions
	// Create a map of current world positions to block types
	positionToType := make(map[string]BlockType)
	for _, block := range piece.Blocks {
		worldX := piece.X + block.X
		worldY := piece.Y + block.Y
		key := fmt.Sprintf("%d,%d", worldX, worldY)
		positionToType[key] = block.BlockType
	}

	// Calculate center of rotation (bounding box center)
	minX, maxX := piece.Blocks[0].X, piece.Blocks[0].X
	minY, maxY := piece.Blocks[0].Y, piece.Blocks[0].Y
	for _, block := range piece.Blocks {
		if block.X < minX {
			minX = block.X
		}
		if block.X > maxX {
			maxX = block.X
		}
		if block.Y < minY {
			minY = block.Y
		}
		if block.Y > maxY {
			maxY = block.Y
		}
	}
	centerX := float64(minX+maxX) / 2.0
	centerY := float64(minY+maxY) / 2.0

	// Rotate each block around the center
	for i := range piece.Blocks {
		block := &piece.Blocks[i]

		// Translate to origin
		relX := float64(block.X) - centerX
		relY := float64(block.Y) - centerY

		// Rotate 90 degrees clockwise: (x,y) -> (y,-x)
		newRelX := relY
		newRelY := -relX

		// Translate back and round to nearest integer
		block.X = int(newRelX + centerX + 0.5)
		block.Y = int(newRelY + centerY + 0.5)
	}

	piece.Rotation = (piece.Rotation + 1) % 4
}

// DrawBlockTransformed renders a single block with rotation and scale at the specified position
func (bm *BlockManager) DrawBlockTransformed(screen *ebiten.Image, block Block, worldX, worldY, rotation, scale, blockSize float64) {
	sprite := bm.GetBlockSprite(block.BlockType)

	op := &ebiten.DrawImageOptions{}

	// Scale the sprite to match the block size and custom scale
	scaleX := (blockSize * scale) / float64(sprite.Bounds().Dx())
	scaleY := (blockSize * scale) / float64(sprite.Bounds().Dy())

	// Apply rotation around center of block
	centerX := float64(sprite.Bounds().Dx()) / 2
	centerY := float64(sprite.Bounds().Dy()) / 2

	// Translate to center, apply rotation and scale, then translate back
	op.GeoM.Translate(-centerX, -centerY)
	op.GeoM.Rotate(rotation * math.Pi / 180.0) // Convert degrees to radians
	op.GeoM.Scale(scaleX, scaleY)

	// Apply electrical storm effect (takes priority over normal wobble)
	if block.IsInStorm {
		// More intense wobble for electrical storms
		stormX := math.Sin(block.StormPhase) * StormIntensity
		stormY := math.Cos(block.StormPhase*1.7) * StormIntensity * 0.3

		// Add random sparking motion
		sparkOffset := math.Sin(block.SparkPhase) * 1.0
		stormX += sparkOffset
		stormY += math.Cos(block.SparkPhase*2.3) * 0.5

		// Position with storm offset and center the scaled/rotated sprite
		op.GeoM.Translate(worldX+stormX+(blockSize*scale)/2, worldY+stormY+(blockSize*scale)/2)

		// Add electrical storm visual effects
		flickerIntensity := 0.2 + 0.1*math.Sin(block.SparkPhase*2)
		if block.BlockType == PositiveBlock {
			op.ColorM.Scale(1.0+flickerIntensity, 1.0+flickerIntensity*0.5, 1.0-flickerIntensity*0.3, 1.0)
		} else if block.BlockType == NegativeBlock {
			op.ColorM.Scale(1.0-flickerIntensity*0.3, 1.0+flickerIntensity*0.5, 1.0+flickerIntensity, 1.0)
		}
	} else if block.IsWobbling {
		// Normal wobble effect
		wobbleX := math.Sin(block.WobblePhase) * WobbleIntensity
		wobbleY := math.Cos(block.WobblePhase*1.3) * WobbleIntensity * 0.5

		// Position with wobble offset and center the scaled/rotated sprite
		op.GeoM.Translate(worldX+wobbleX+(blockSize*scale)/2, worldY+wobbleY+(blockSize*scale)/2)

		// Add slight transparency to indicate impending destruction
		wobbleProgress := block.WobbleTime / WobbleDuration
		alpha := 1.0 - wobbleProgress*0.3
		op.ColorM.Scale(1, 1, 1, alpha)
	} else {
		// Position the block normally and center the scaled/rotated sprite
		op.GeoM.Translate(worldX+(blockSize*scale)/2, worldY+(blockSize*scale)/2)
	}

	screen.DrawImage(sprite, op)

	// Draw POW sprite on top of wobbling blocks (but only when ShowPowSprite is true)
	// Note: For transformed blocks, we draw the POW sprite at the original block position without transformation
	if block.IsWobbling && block.ShowPowSprite {
		bm.DrawPowSprite(screen, worldX, worldY, block.WobblePhase, blockSize*scale)
	}
}

// DrawWarningSprite renders a shaking ZAP sprite at the specified position
func (bm *BlockManager) DrawWarningSprite(screen *ebiten.Image, worldX, worldY, warningTime, blockSize float64, column int, gameboardWidth int) {
	sprite := assets.ZapSprite
	if sprite == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}

	// Scale the sprite to match the block size
	scaleX := blockSize / float64(sprite.Bounds().Dx())
	scaleY := blockSize / float64(sprite.Bounds().Dy())
	op.GeoM.Scale(scaleX, scaleY)

	// Add shake/wobble effect
	shakePhase := warningTime * WarningFrequency * 2 * math.Pi
	shakeX := math.Sin(shakePhase) * WarningIntensity
	shakeY := math.Cos(shakePhase*1.3) * WarningIntensity * 0.7

	// Add secondary wobble for more dynamic movement
	wobblePhase := warningTime * WarningFrequency * 1.5 * math.Pi
	wobbleX := math.Sin(wobblePhase) * WarningIntensity * 0.5
	wobbleY := math.Cos(wobblePhase*0.8) * WarningIntensity * 0.3

	// Calculate grid dimensions to determine board center
	gameboardWidthInBlocks := int(float64(gameboardWidth) / blockSize)
	boardCenter := gameboardWidthInBlocks / 2

	// Position sprite on the opposite side if storm is on right half
	var offsetX float64
	if column >= boardCenter {
		// Storm is on right half, position sprite to the left
		offsetX = -blockSize * 0.3 // Position towards top-left
	} else {
		// Storm is on left half, position sprite to the right
		offsetX = blockSize * 0.7 // Position towards top-right
	}
	offsetY := blockSize * 0.1

	op.GeoM.Translate(worldX+offsetX+shakeX+wobbleX, worldY+offsetY+shakeY+wobbleY)

	// Add slight color enhancement to make it more visible (no pulsing)
	op.ColorM.Scale(1.2, 1.1, 0.9, 1.0) // Slightly yellow/orange tint with full opacity

	screen.DrawImage(sprite, op)
}

// DrawPowSprite renders a wobbling POW sprite at the specified position
func (bm *BlockManager) DrawPowSprite(screen *ebiten.Image, worldX, worldY, wobblePhase, blockSize float64) {
	sprite := assets.PowSprite
	if sprite == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}

	// Scale the sprite to match the block size
	scaleX := blockSize / float64(sprite.Bounds().Dx())
	scaleY := blockSize / float64(sprite.Bounds().Dy())

	// Add wobble effect matching the block's wobble
	wobbleX := math.Sin(wobblePhase) * WobbleIntensity
	wobbleY := math.Cos(wobblePhase*1.3) * WobbleIntensity * 0.5

	// Center the sprite on the block by offsetting by half the scaled size
	spriteWidth := float64(sprite.Bounds().Dx()) * scaleX
	spriteHeight := float64(sprite.Bounds().Dy()) * scaleY
	centerOffsetX := (blockSize - spriteWidth) / 2
	centerOffsetY := (blockSize - spriteHeight) / 2

	// Apply scale and position with wobble offset and center alignment
	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(worldX+centerOffsetX+wobbleX, worldY+centerOffsetY+wobbleY)

	// Add slight transparency and color enhancement for visibility
	op.ColorM.Scale(1.0, 1.0, 1.0, 0.9) // Slight transparency so block is still visible underneath

	screen.DrawImage(sprite, op)
}
