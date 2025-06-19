package main

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
