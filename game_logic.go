package main

import "fmt"

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

// IsGameOver checks if any placed blocks have reached the top of the gameboard
func (gl *GameLogic) IsGameOver() bool {
	// Check if any placed blocks are at Y position 0 or negative (top of the screen)
	for _, block := range gl.placedBlocks {
		if block.Y <= 0 {
			return true
		}
	}
	return false
}

// CheckAndProcessReactions finds horizontal clusters and removes contiguous zero-sum subsequences
func (gl *GameLogic) CheckAndProcessReactions() {
	for {
		blocksToRemove := gl.findBlocksToRemove()

		if len(blocksToRemove) == 0 {
			break // No more reactions possible
		}

		// Remove the blocks
		gl.removeBlocks(blocksToRemove)

		// Make remaining blocks fall
		gl.processBlockFalling()
	}
}

// findBlocksToRemove finds all blocks that should be removed based on the rules
func (gl *GameLogic) findBlocksToRemove() []Block {
	var blocksToRemove []Block

	// Group blocks by row (Y coordinate)
	rowMap := make(map[int][]Block)
	for _, block := range gl.placedBlocks {
		rowMap[block.Y] = append(rowMap[block.Y], block)
	}

	// Process each row
	for _, rowBlocks := range rowMap {
		if len(rowBlocks) < 3 {
			continue // Need at least 3 blocks
		}

		// Sort blocks by X position
		for i := 0; i < len(rowBlocks); i++ {
			for j := i + 1; j < len(rowBlocks); j++ {
				if rowBlocks[i].X > rowBlocks[j].X {
					rowBlocks[i], rowBlocks[j] = rowBlocks[j], rowBlocks[i]
				}
			}
		}

		// Find contiguous clusters (broken by gaps or neutral blocks)
		clusters := gl.findClustersInRow(rowBlocks)

		// For each cluster, find zero-sum subsequences
		for _, cluster := range clusters {
			if len(cluster) >= 3 {
				zeroSumBlocks := gl.findZeroSumSubsequence(cluster)
				blocksToRemove = append(blocksToRemove, zeroSumBlocks...)
			}
		}
	}

	return blocksToRemove
}

// findClustersInRow splits a row into contiguous clusters (broken by gaps or neutral blocks)
func (gl *GameLogic) findClustersInRow(rowBlocks []Block) [][]Block {
	var clusters [][]Block
	var currentCluster []Block

	for i, block := range rowBlocks {
		// Check if this block breaks the cluster
		if block.BlockType == NeutralBlock {
			// Neutral blocks break clusters
			if len(currentCluster) >= 3 {
				clusters = append(clusters, currentCluster)
			}
			currentCluster = nil
		} else if i > 0 && block.X != rowBlocks[i-1].X+1 {
			// Gap in X coordinates breaks clusters
			if len(currentCluster) >= 3 {
				clusters = append(clusters, currentCluster)
			}
			currentCluster = []Block{block}
		} else {
			// Continue the cluster
			currentCluster = append(currentCluster, block)
		}
	}

	// Don't forget the last cluster
	if len(currentCluster) >= 3 {
		clusters = append(clusters, currentCluster)
	}

	return clusters
}

// findZeroSumSubsequence finds the longest contiguous subsequence that sums to zero
func (gl *GameLogic) findZeroSumSubsequence(cluster []Block) []Block {
	// Try all possible contiguous subsequences of length 3 or more
	for length := len(cluster); length >= 3; length-- {
		for start := 0; start <= len(cluster)-length; start++ {
			subsequence := cluster[start : start+length]

			// Calculate sum
			sum := 0
			for _, block := range subsequence {
				switch block.BlockType {
				case PositiveBlock:
					sum += 1
				case NegativeBlock:
					sum -= 1
				}
			}

			if sum == 0 {
				return subsequence
			}
		}
	}

	return nil
}

// removeBlocks removes specified blocks from the placed blocks array
func (gl *GameLogic) removeBlocks(blocksToRemove []Block) {
	if len(blocksToRemove) == 0 {
		return
	}

	// Create a map for fast lookup
	removeMap := make(map[string]bool)
	for _, block := range blocksToRemove {
		key := fmt.Sprintf("%d,%d", block.X, block.Y)
		removeMap[key] = true
	}

	// Filter out the blocks to remove
	var remainingBlocks []Block
	for _, block := range gl.placedBlocks {
		key := fmt.Sprintf("%d,%d", block.X, block.Y)
		if !removeMap[key] {
			remainingBlocks = append(remainingBlocks, block)
		}
	}

	gl.placedBlocks = remainingBlocks
}

// processBlockFalling makes remaining blocks fall individually
func (gl *GameLogic) processBlockFalling() {
	// Sort blocks by Y position (bottom to top)
	for i := 0; i < len(gl.placedBlocks); i++ {
		for j := i + 1; j < len(gl.placedBlocks); j++ {
			if gl.placedBlocks[i].Y < gl.placedBlocks[j].Y {
				gl.placedBlocks[i], gl.placedBlocks[j] = gl.placedBlocks[j], gl.placedBlocks[i]
			}
		}
	}

	// Make each block fall
	blockSize := gl.blockManager.GetScaledBlockSize(gl.gameboard.Width, gl.gameboard.Height)
	gameboardHeightInBlocks := int(float64(gl.gameboard.Height) / blockSize)

	for i := range gl.placedBlocks {
		block := &gl.placedBlocks[i]

		// Find the lowest valid Y position
		for newY := block.Y + 1; newY < gameboardHeightInBlocks; newY++ {
			// Check if position is occupied
			occupied := false
			for j := range gl.placedBlocks {
				if i != j && gl.placedBlocks[j].X == block.X && gl.placedBlocks[j].Y == newY {
					occupied = true
					break
				}
			}

			if occupied {
				break
			}

			block.Y = newY
		}
	}
}
