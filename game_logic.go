package main

// GameLogic handles game rules, collision detection, and piece management
type GameLogic struct {
	gameboard    *Gameboard
	blockManager *BlockManager
	placedBlocks []Block
}

// NewGameLogic creates a new game logic handler
func NewGameLogic(gameboard *Gameboard, blockManager *BlockManager) *GameLogic {
	return &GameLogic{
		gameboard:    gameboard,
		blockManager: blockManager,
		placedBlocks: make([]Block, 0),
	}
}

// IsValidPosition checks if a piece can be placed at the given position
func (gl *GameLogic) IsValidPosition(piece *TetrisPiece, offsetX, offsetY int) bool {
	// Use the actual scaled block size that's being used for rendering
	blockSize := gl.blockManager.GetScaledBlockSize(gl.gameboard.Width, gl.gameboard.Height)
	
	// Calculate grid dimensions based on actual gameboard size and block size
	gameboardWidthInBlocks := int(float64(gl.gameboard.Width) / blockSize)
	gameboardHeightInBlocks := int(float64(gl.gameboard.Height) / blockSize)
	
	for _, block := range piece.Blocks {
		newX := piece.X + block.X + offsetX
		newY := piece.Y + block.Y + offsetY
		
		// Check boundaries
		if newX < 0 || newX >= gameboardWidthInBlocks || newY >= gameboardHeightInBlocks {
			return false
		}
		
		// Check collision with placed blocks
		for _, placedBlock := range gl.placedBlocks {
			if placedBlock.X == newX && placedBlock.Y == newY {
				return false
			}
		}
	}
	
	return true
}

// PlacePiece adds the current piece to the placed blocks
func (gl *GameLogic) PlacePiece(piece *TetrisPiece) {
	if piece == nil {
		return
	}
	
	for _, block := range piece.Blocks {
		placedBlock := Block{
			X:         piece.X + block.X,
			Y:         piece.Y + block.Y,
			BlockType: block.BlockType,
		}
		gl.placedBlocks = append(gl.placedBlocks, placedBlock)
	}
}

// GetPlacedBlocks returns a copy of the placed blocks
func (gl *GameLogic) GetPlacedBlocks() []Block {
	return gl.placedBlocks
}

// SpawnNewPiece creates a new piece at the top center of the gameboard
func (gl *GameLogic) SpawnNewPiece(pieceType PieceType) *TetrisPiece {
	// Use the same block size calculation as collision detection
	blockSize := gl.blockManager.GetScaledBlockSize(gl.gameboard.Width, gl.gameboard.Height)
	gameboardWidthInBlocks := int(float64(gl.gameboard.Width) / blockSize)
	centerX := gameboardWidthInBlocks / 2
	
	return gl.blockManager.CreateTetrisPiece(pieceType, centerX, 0)
}

// TryRotatePiece attempts to rotate a piece, returns true if successful
func (gl *GameLogic) TryRotatePiece(piece *TetrisPiece, pieceType PieceType) bool {
	if piece == nil {
		return false
	}
	
	// Store original rotation in case we need to revert
	originalRotation := piece.Rotation
	originalBlocks := make([]Block, len(piece.Blocks))
	copy(originalBlocks, piece.Blocks)
	
	gl.blockManager.RotatePiece(piece, pieceType)
	
	// Check if rotation is valid
	if !gl.IsValidPosition(piece, 0, 0) {
		// Revert rotation
		piece.Rotation = originalRotation
		piece.Blocks = originalBlocks
		return false
	}
	
	return true
}

// TryMovePiece attempts to move a piece, returns true if successful
func (gl *GameLogic) TryMovePiece(piece *TetrisPiece, deltaX, deltaY int) bool {
	if piece == nil {
		return false
	}
	
	if gl.IsValidPosition(piece, deltaX, deltaY) {
		piece.X += deltaX
		piece.Y += deltaY
		return true
	}
	
	return false
}
