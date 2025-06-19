package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type PieceType int

const (
	IPiece PieceType = iota
	OPiece
	TPiece
	SPiece
	ZPiece
	JPiece
	LPiece
)

type TetrisPiece struct {
	Blocks   []Block
	X, Y     int
	Rotation int
}

func GetPieceBlocks(pieceType PieceType, rotation int, genBlockType func() BlockType) []Block {
	switch pieceType {
	case IPiece:
		return getIPieceBlocks(rotation, genBlockType)
	case OPiece:
		return getOPieceBlocks(rotation, genBlockType)
	case TPiece:
		return getTPieceBlocks(rotation, genBlockType)
	case SPiece:
		return getSPieceBlocks(rotation, genBlockType)
	case ZPiece:
		return getZPieceBlocks(rotation, genBlockType)
	case JPiece:
		return getJPieceBlocks(rotation, genBlockType)
	case LPiece:
		return getLPieceBlocks(rotation, genBlockType)
	default:
		return nil
	}
}

func getIPieceBlocks(rotation int, genBlockType func() BlockType) []Block {
	var blocks []Block
	if rotation%2 == 0 {
		for i := 0; i < 4; i++ {
			blocks = append(blocks, Block{X: i, Y: 0, BlockType: genBlockType()})
		}
	} else {
		for i := 0; i < 4; i++ {
			blocks = append(blocks, Block{X: 0, Y: i, BlockType: genBlockType()})
		}
	}
	return blocks
}

func getOPieceBlocks(rotation int, genBlockType func() BlockType) []Block {
	switch rotation % 4 {
	case 0:
		return []Block{
			{X: 0, Y: 0, BlockType: genBlockType()},
			{X: 1, Y: 0, BlockType: genBlockType()},
			{X: 0, Y: 1, BlockType: genBlockType()},
			{X: 1, Y: 1, BlockType: genBlockType()},
		}
	case 1:
		return []Block{
			{X: 1, Y: 0, BlockType: genBlockType()},
			{X: 1, Y: 1, BlockType: genBlockType()},
			{X: 0, Y: 0, BlockType: genBlockType()},
			{X: 0, Y: 1, BlockType: genBlockType()},
		}
	case 2:
		return []Block{
			{X: 1, Y: 1, BlockType: genBlockType()},
			{X: 0, Y: 1, BlockType: genBlockType()},
			{X: 1, Y: 0, BlockType: genBlockType()},
			{X: 0, Y: 0, BlockType: genBlockType()},
		}
	case 3:
		return []Block{
			{X: 0, Y: 1, BlockType: genBlockType()},
			{X: 0, Y: 0, BlockType: genBlockType()},
			{X: 1, Y: 1, BlockType: genBlockType()},
			{X: 1, Y: 0, BlockType: genBlockType()},
		}
	}
	return nil
}

func getTPieceBlocks(rotation int, genBlockType func() BlockType) []Block {
	switch rotation % 4 {
	case 0:
		return []Block{
			{X: 1, Y: 0, BlockType: genBlockType()},
			{X: 0, Y: 1, BlockType: genBlockType()},
			{X: 1, Y: 1, BlockType: genBlockType()},
			{X: 2, Y: 1, BlockType: genBlockType()},
		}
	case 1:
		return []Block{
			{X: 0, Y: 0, BlockType: genBlockType()},
			{X: 0, Y: 1, BlockType: genBlockType()},
			{X: 1, Y: 1, BlockType: genBlockType()},
			{X: 0, Y: 2, BlockType: genBlockType()},
		}
	case 2:
		return []Block{
			{X: 0, Y: 0, BlockType: genBlockType()},
			{X: 1, Y: 0, BlockType: genBlockType()},
			{X: 2, Y: 0, BlockType: genBlockType()},
			{X: 1, Y: 1, BlockType: genBlockType()},
		}
	case 3:
		return []Block{
			{X: 1, Y: 0, BlockType: genBlockType()},
			{X: 0, Y: 1, BlockType: genBlockType()},
			{X: 1, Y: 1, BlockType: genBlockType()},
			{X: 1, Y: 2, BlockType: genBlockType()},
		}
	}
	return nil
}

func getSPieceBlocks(rotation int, genBlockType func() BlockType) []Block {
	if rotation%2 == 0 {
		return []Block{
			{X: 1, Y: 0, BlockType: genBlockType()},
			{X: 2, Y: 0, BlockType: genBlockType()},
			{X: 0, Y: 1, BlockType: genBlockType()},
			{X: 1, Y: 1, BlockType: genBlockType()},
		}
	} else {
		return []Block{
			{X: 0, Y: 0, BlockType: genBlockType()},
			{X: 0, Y: 1, BlockType: genBlockType()},
			{X: 1, Y: 1, BlockType: genBlockType()},
			{X: 1, Y: 2, BlockType: genBlockType()},
		}
	}
}

func getZPieceBlocks(rotation int, genBlockType func() BlockType) []Block {
	if rotation%2 == 0 {
		return []Block{
			{X: 0, Y: 0, BlockType: genBlockType()},
			{X: 1, Y: 0, BlockType: genBlockType()},
			{X: 1, Y: 1, BlockType: genBlockType()},
			{X: 2, Y: 1, BlockType: genBlockType()},
		}
	} else {
		return []Block{
			{X: 1, Y: 0, BlockType: genBlockType()},
			{X: 0, Y: 1, BlockType: genBlockType()},
			{X: 1, Y: 1, BlockType: genBlockType()},
			{X: 0, Y: 2, BlockType: genBlockType()},
		}
	}
}

func getJPieceBlocks(rotation int, genBlockType func() BlockType) []Block {
	switch rotation % 4 {
	case 0:
		return []Block{
			{X: 0, Y: 0, BlockType: genBlockType()},
			{X: 0, Y: 1, BlockType: genBlockType()},
			{X: 1, Y: 1, BlockType: genBlockType()},
			{X: 2, Y: 1, BlockType: genBlockType()},
		}
	case 1:
		return []Block{
			{X: 1, Y: 0, BlockType: genBlockType()},
			{X: 1, Y: 1, BlockType: genBlockType()},
			{X: 1, Y: 2, BlockType: genBlockType()},
			{X: 0, Y: 2, BlockType: genBlockType()},
		}
	case 2:
		return []Block{
			{X: 0, Y: 0, BlockType: genBlockType()},
			{X: 1, Y: 0, BlockType: genBlockType()},
			{X: 2, Y: 0, BlockType: genBlockType()},
			{X: 2, Y: 1, BlockType: genBlockType()},
		}
	case 3:
		return []Block{
			{X: 0, Y: 0, BlockType: genBlockType()},
			{X: 1, Y: 0, BlockType: genBlockType()},
			{X: 0, Y: 1, BlockType: genBlockType()},
			{X: 0, Y: 2, BlockType: genBlockType()},
		}
	}
	return nil
}

func getLPieceBlocks(rotation int, genBlockType func() BlockType) []Block {
	switch rotation % 4 {
	case 0:
		return []Block{
			{X: 2, Y: 0, BlockType: genBlockType()},
			{X: 0, Y: 1, BlockType: genBlockType()},
			{X: 1, Y: 1, BlockType: genBlockType()},
			{X: 2, Y: 1, BlockType: genBlockType()},
		}
	case 1:
		return []Block{
			{X: 0, Y: 0, BlockType: genBlockType()},
			{X: 0, Y: 1, BlockType: genBlockType()},
			{X: 0, Y: 2, BlockType: genBlockType()},
			{X: 1, Y: 2, BlockType: genBlockType()},
		}
	case 2:
		return []Block{
			{X: 0, Y: 0, BlockType: genBlockType()},
			{X: 1, Y: 0, BlockType: genBlockType()},
			{X: 2, Y: 0, BlockType: genBlockType()},
			{X: 0, Y: 1, BlockType: genBlockType()},
		}
	case 3:
		return []Block{
			{X: 0, Y: 0, BlockType: genBlockType()},
			{X: 1, Y: 0, BlockType: genBlockType()},
			{X: 1, Y: 1, BlockType: genBlockType()},
			{X: 1, Y: 2, BlockType: genBlockType()},
		}
	}
	return nil
}

func (bm *BlockManager) DrawTetrisPiece(screen *ebiten.Image, piece *TetrisPiece, screenWidth, screenHeight int) {
	blockSize := bm.GetScaledBlockSize(screenWidth, screenHeight)

	for _, block := range piece.Blocks {
		worldX := float64(piece.X+block.X) * blockSize
		worldY := float64(piece.Y+block.Y) * blockSize

		bm.DrawBlock(screen, block, worldX, worldY, blockSize)
	}
}

func (bm *BlockManager) CreateTetrisPiece(pieceType PieceType, x, y int) *TetrisPiece {
	blocks := GetPieceBlocks(pieceType, 0, bm.GenerateRandomBlockType)

	return &TetrisPiece{
		Blocks:   blocks,
		X:        x,
		Y:        y,
		Rotation: 0,
	}
}
func (bm *BlockManager) RotatePiece(piece *TetrisPiece, pieceType PieceType) {

	if pieceType == OPiece {
		newRotation := (piece.Rotation + 1) % 4
		newPositions := bm.GetPiecePositions(pieceType, newRotation)

		if len(newPositions) == len(piece.Blocks) {
			for i, pos := range newPositions {
				piece.Blocks[i].X = pos.X
				piece.Blocks[i].Y = pos.Y
			}
			piece.Rotation = newRotation
		}
		return
	}
	positionToType := make(map[string]BlockType)
	for _, block := range piece.Blocks {
		worldX := piece.X + block.X
		worldY := piece.Y + block.Y
		key := fmt.Sprintf("%d,%d", worldX, worldY)
		positionToType[key] = block.BlockType
	}

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

	for i := range piece.Blocks {
		block := &piece.Blocks[i]

		relX := float64(block.X) - centerX
		relY := float64(block.Y) - centerY

		newRelX := relY
		newRelY := -relX

		block.X = int(newRelX + centerX + 0.5)
		block.Y = int(newRelY + centerY + 0.5)
	}

	piece.Rotation = (piece.Rotation + 1) % 4
}

func (bm *BlockManager) GetPiecePositions(pieceType PieceType, rotation int) []struct{ X, Y int } {
	switch pieceType {
	case IPiece:
		return getIPiecePositions(rotation)
	case OPiece:
		return getOPiecePositions(rotation)
	case TPiece:
		return getTPiecePositions(rotation)
	case SPiece:
		return getSPiecePositions(rotation)
	case ZPiece:
		return getZPiecePositions(rotation)
	case JPiece:
		return getJPiecePositions(rotation)
	case LPiece:
		return getLPiecePositions(rotation)
	default:
		return nil
	}
}

func getIPiecePositions(rotation int) []struct{ X, Y int } {
	if rotation%2 == 0 {
		var positions []struct{ X, Y int }
		for i := range 4 {
			positions = append(positions, struct{ X, Y int }{X: i, Y: 0})
		}
		return positions
	} else {
		var positions []struct{ X, Y int }
		for i := range 4 {
			positions = append(positions, struct{ X, Y int }{X: 0, Y: i})
		}
		return positions
	}
}

func getOPiecePositions(rotation int) []struct{ X, Y int } {
	switch rotation % 4 {
	case 0:
		return []struct{ X, Y int }{
			{X: 0, Y: 0},
			{X: 1, Y: 0},
			{X: 0, Y: 1},
			{X: 1, Y: 1},
		}
	case 1:
		return []struct{ X, Y int }{
			{X: 1, Y: 0},
			{X: 1, Y: 1},
			{X: 0, Y: 0},
			{X: 0, Y: 1},
		}
	case 2:
		return []struct{ X, Y int }{
			{X: 1, Y: 1},
			{X: 0, Y: 1},
			{X: 1, Y: 0},
			{X: 0, Y: 0},
		}
	case 3:
		return []struct{ X, Y int }{
			{X: 0, Y: 1},
			{X: 0, Y: 0},
			{X: 1, Y: 1},
			{X: 1, Y: 0},
		}
	}
	return nil
}

func getTPiecePositions(rotation int) []struct{ X, Y int } {
	switch rotation % 4 {
	case 0:
		return []struct{ X, Y int }{
			{X: 1, Y: 0},
			{X: 0, Y: 1},
			{X: 1, Y: 1},
			{X: 2, Y: 1},
		}
	case 1:
		return []struct{ X, Y int }{
			{X: 0, Y: 0},
			{X: 0, Y: 1},
			{X: 1, Y: 1},
			{X: 0, Y: 2},
		}
	case 2:
		return []struct{ X, Y int }{
			{X: 0, Y: 0},
			{X: 1, Y: 0},
			{X: 2, Y: 0},
			{X: 1, Y: 1},
		}
	case 3:
		return []struct{ X, Y int }{
			{X: 1, Y: 0},
			{X: 0, Y: 1},
			{X: 1, Y: 1},
			{X: 1, Y: 2},
		}
	}
	return nil
}

func getSPiecePositions(rotation int) []struct{ X, Y int } {
	if rotation%2 == 0 {
		return []struct{ X, Y int }{
			{X: 1, Y: 0},
			{X: 2, Y: 0},
			{X: 0, Y: 1},
			{X: 1, Y: 1},
		}
	} else {
		return []struct{ X, Y int }{
			{X: 0, Y: 0},
			{X: 0, Y: 1},
			{X: 1, Y: 1},
			{X: 1, Y: 2},
		}
	}
}

func getZPiecePositions(rotation int) []struct{ X, Y int } {
	if rotation%2 == 0 {
		return []struct{ X, Y int }{
			{X: 0, Y: 0},
			{X: 1, Y: 0},
			{X: 1, Y: 1},
			{X: 2, Y: 1},
		}
	} else {
		return []struct{ X, Y int }{
			{X: 1, Y: 0},
			{X: 0, Y: 1},
			{X: 1, Y: 1},
			{X: 0, Y: 2},
		}
	}
}

func getJPiecePositions(rotation int) []struct{ X, Y int } {
	switch rotation % 4 {
	case 0:
		return []struct{ X, Y int }{
			{X: 0, Y: 0},
			{X: 0, Y: 1},
			{X: 1, Y: 1},
			{X: 2, Y: 1},
		}
	case 1:
		return []struct{ X, Y int }{
			{X: 1, Y: 0},
			{X: 1, Y: 1},
			{X: 1, Y: 2},
			{X: 0, Y: 2},
		}
	case 2:
		return []struct{ X, Y int }{
			{X: 0, Y: 0},
			{X: 1, Y: 0},
			{X: 2, Y: 0},
			{X: 2, Y: 1},
		}
	case 3:
		return []struct{ X, Y int }{
			{X: 0, Y: 0},
			{X: 1, Y: 0},
			{X: 0, Y: 1},
			{X: 0, Y: 2},
		}
	}
	return nil
}

func getLPiecePositions(rotation int) []struct{ X, Y int } {
	switch rotation % 4 {
	case 0:
		return []struct{ X, Y int }{
			{X: 2, Y: 0},
			{X: 0, Y: 1},
			{X: 1, Y: 1},
			{X: 2, Y: 1},
		}
	case 1:
		return []struct{ X, Y int }{
			{X: 0, Y: 0},
			{X: 0, Y: 1},
			{X: 0, Y: 2},
			{X: 1, Y: 2},
		}
	case 2:
		return []struct{ X, Y int }{
			{X: 0, Y: 0},
			{X: 1, Y: 0},
			{X: 2, Y: 0},
			{X: 0, Y: 1},
		}
	case 3:
		return []struct{ X, Y int }{
			{X: 0, Y: 0},
			{X: 1, Y: 0},
			{X: 1, Y: 1},
			{X: 1, Y: 2},
		}
	}
	return nil
}
