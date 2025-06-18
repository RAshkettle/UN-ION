package main

import (
	"fmt"
	"math"
)

// ExplosionCallback is called when blocks are removed to trigger particle effects
type ExplosionCallback func(worldX, worldY float64, blockType BlockType)

// AudioCallback is called when blocks are removed to trigger audio effects
type AudioCallback func(blocksRemoved int)

// DustCallback is called when pieces are placed to trigger dust cloud effects
type DustCallback func(worldX, worldY float64)

// HardDropCallback is called when pieces are hard dropped to trigger screen shake
type HardDropCallback func(dropHeight int)

// GameLogic handles game rules, collision detection, and piece management
type GameLogic struct {
	gameboard         *Gameboard
	blockManager      *BlockManager
	placedBlocks      []Block
	explosionCallback ExplosionCallback
	audioCallback     AudioCallback
	dustCallback      DustCallback
	hardDropCallback  HardDropCallback
}

// NewGameLogic creates a new game logic handler
func NewGameLogic(gameboard *Gameboard, blockManager *BlockManager) *GameLogic {
	return &GameLogic{
		gameboard:    gameboard,
		blockManager: blockManager,
		placedBlocks: make([]Block, 0),
	}
}

// SetExplosionCallback sets the callback function for particle explosions
func (gl *GameLogic) SetExplosionCallback(callback ExplosionCallback) {
	gl.explosionCallback = callback
}

// SetAudioCallback sets the callback function for audio effects
func (gl *GameLogic) SetAudioCallback(callback AudioCallback) {
	gl.audioCallback = callback
}

// SetDustCallback sets the callback function for dust cloud effects
func (gl *GameLogic) SetDustCallback(callback DustCallback) {
	gl.dustCallback = callback
}

// SetHardDropCallback sets the callback function for hard drop screen shake effects
func (gl *GameLogic) SetHardDropCallback(callback HardDropCallback) {
	gl.hardDropCallback = callback
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
	
	// Trigger dust cloud effect at the center bottom of the piece
	if gl.dustCallback != nil {
		blockSize := gl.blockManager.GetScaledBlockSize(gl.gameboard.Width, gl.gameboard.Height)
		// Find the bottom-most Y position of the piece
		bottomY := piece.Y
		for _, block := range piece.Blocks {
			if piece.Y + block.Y > bottomY {
				bottomY = piece.Y + block.Y
			}
		}
		// Calculate world position at bottom center of piece
		worldX := float64(gl.gameboard.X) + float64(piece.X)*blockSize + blockSize/2
		worldY := float64(gl.gameboard.Y) + float64(bottomY+1)*blockSize // Just below the piece
		gl.dustCallback(worldX, worldY)
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

// CalculateDropPosition calculates where a piece would land if dropped straight down
func (gl *GameLogic) CalculateDropPosition(piece *TetrisPiece) *TetrisPiece {
	if piece == nil {
		return nil
	}

	// Create a copy of the piece
	shadowPiece := &TetrisPiece{
		X:        piece.X,
		Y:        piece.Y,
		Rotation: piece.Rotation,
		Blocks:   make([]Block, len(piece.Blocks)),
	}

	// Copy all blocks
	for i, block := range piece.Blocks {
		shadowPiece.Blocks[i] = Block{
			X:         block.X,
			Y:         block.Y,
			BlockType: block.BlockType,
		}
	}

	// Move the shadow piece down until it can't move anymore
	for gl.IsValidPosition(shadowPiece, 0, 1) {
		shadowPiece.Y++
	}

	return shadowPiece
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
// Returns the total score earned from all reactions
func (gl *GameLogic) CheckAndProcessReactions() int {
	totalScore := 0

	for {
		blocksToRemove := gl.findBlocksToRemove()

		if len(blocksToRemove) == 0 {
			break // No more reactions possible
		}

		// Calculate score for this reaction
		reactionScore := gl.calculateReactionScore(len(blocksToRemove))
		totalScore += reactionScore

		// Remove the blocks
		gl.removeBlocks(blocksToRemove)

		// Make remaining blocks fall
		gl.processBlockFalling()
	}

	return totalScore
}

// calculateReactionScore calculates score based on number of blocks removed
// 4 blocks = 10 points, each additional block multiplies by 2
func (gl *GameLogic) calculateReactionScore(blocksRemoved int) int {
	if blocksRemoved < 4 {
		return 0 // No score for less than 4 blocks
	}

	score := 10 // Base score for 4 blocks

	// For each block beyond 4, multiply by 2
	for i := 4; i < blocksRemoved; i++ {
		score *= 2
	}

	return score
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

	// Trigger audio callback for block breaking sound
	if gl.audioCallback != nil {
		gl.audioCallback(len(blocksToRemove))
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
		} else {
			// Trigger explosion callback for particle effects
			if gl.explosionCallback != nil {
				blockSize := gl.blockManager.GetScaledBlockSize(gl.gameboard.Width, gl.gameboard.Height)
				worldX := float64(gl.gameboard.X) + float64(block.X)*blockSize + blockSize/2
				worldY := float64(gl.gameboard.Y) + float64(block.Y)*blockSize + blockSize/2
				gl.explosionCallback(worldX, worldY, block.BlockType)
			}
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

// UpdateWobblingBlocks updates the wobble animation for all wobbling blocks
// Returns true if any blocks finished wobbling and should be removed
func (gl *GameLogic) UpdateWobblingBlocks(deltaTime float64) bool {
	anyBlocksFinished := false
	
	for i := range gl.placedBlocks {
		block := &gl.placedBlocks[i]
		if block.IsWobbling {
			// Update wobble time and phase
			block.WobbleTime += deltaTime
			block.WobblePhase += deltaTime * WobbleFrequency * 2 * math.Pi
			
			// Check if wobble duration is finished
			if block.WobbleTime >= WobbleDuration {
				anyBlocksFinished = true
			}
		}
	}
	
	return anyBlocksFinished
}

// RemoveFinishedWobblingBlocks removes blocks that have finished wobbling
func (gl *GameLogic) RemoveFinishedWobblingBlocks() int {
	var blocksToRemove []Block
	var remainingBlocks []Block
	
	// Separate blocks that finished wobbling from remaining blocks
	for _, block := range gl.placedBlocks {
		if block.IsWobbling && block.WobbleTime >= WobbleDuration {
			blocksToRemove = append(blocksToRemove, block)
		} else {
			remainingBlocks = append(remainingBlocks, block)
		}
	}
	
	if len(blocksToRemove) == 0 {
		return 0
	}
	
	// Trigger audio callback for block breaking sound
	if gl.audioCallback != nil {
		gl.audioCallback(len(blocksToRemove))
	}
	
	// Trigger explosion effects for removed blocks
	for _, block := range blocksToRemove {
		if gl.explosionCallback != nil {
			blockSize := gl.blockManager.GetScaledBlockSize(gl.gameboard.Width, gl.gameboard.Height)
			worldX := float64(gl.gameboard.X) + float64(block.X)*blockSize + blockSize/2
			worldY := float64(gl.gameboard.Y) + float64(block.Y)*blockSize + blockSize/2
			gl.explosionCallback(worldX, worldY, block.BlockType)
		}
	}
	
	// Update placed blocks
	gl.placedBlocks = remainingBlocks
	
	// Clean up any invalid storms after removing blocks
	gl.ClearInvalidStorms()
	
	return len(blocksToRemove)
}

// StartBlockWobbling marks blocks for wobbling instead of immediate removal
func (gl *GameLogic) StartBlockWobbling(blocksToWobble []Block) {
	if len(blocksToWobble) == 0 {
		return
	}
	
	// Create a map for fast lookup
	wobbleMap := make(map[string]bool)
	for _, block := range blocksToWobble {
		key := fmt.Sprintf("%d,%d", block.X, block.Y)
		wobbleMap[key] = true
	}
	
	// Mark matching blocks as wobbling
	for i := range gl.placedBlocks {
		block := &gl.placedBlocks[i]
		key := fmt.Sprintf("%d,%d", block.X, block.Y)
		if wobbleMap[key] && !block.IsWobbling {
			block.IsWobbling = true
			block.WobbleTime = 0
			block.WobblePhase = 0
		}
	}
}

// CheckForNewReactions finds blocks that should start wobbling (only non-wobbling blocks)
func (gl *GameLogic) CheckForNewReactions() int {
	blocksToWobble := gl.findNonWobblingBlocksToRemove()
	
	if len(blocksToWobble) == 0 {
		return 0
	}
	
	// Start wobbling on these blocks
	gl.StartBlockWobbling(blocksToWobble)
	
	// Calculate and return score
	return gl.calculateReactionScore(len(blocksToWobble))
}

// findNonWobblingBlocksToRemove finds blocks that should be removed, excluding already wobbling blocks
func (gl *GameLogic) findNonWobblingBlocksToRemove() []Block {
	var blocksToRemove []Block

	// Group non-wobbling blocks by row (Y coordinate)
	rowMap := make(map[int][]Block)
	for _, block := range gl.placedBlocks {
		if !block.IsWobbling {
			rowMap[block.Y] = append(rowMap[block.Y], block)
		}
	}

	// Process each row (same logic as findBlocksToRemove but only for non-wobbling blocks)
	for _, rowBlocks := range rowMap {
		if len(rowBlocks) < 3 {
			continue // Need at least 3 blocks
		}

		// Sort blocks by X coordinate for processing
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

// findVerticalElectricalStorms finds vertical sequences of 4+ positive or negative blocks
func (gl *GameLogic) findVerticalElectricalStorms() []Block {
	var stormBlocks []Block

	// Group ALL blocks by column (X coordinate) - including existing storm blocks
	columnMap := make(map[int][]Block)
	for _, block := range gl.placedBlocks {
		if !block.IsWobbling && block.BlockType != NeutralBlock {
			columnMap[block.X] = append(columnMap[block.X], block)
		}
	}

	// Process each column
	for _, columnBlocks := range columnMap {
		if len(columnBlocks) < 4 {
			continue // Need at least 4 blocks for electrical storm
		}

		// Sort blocks by Y position (top to bottom)
		for i := 0; i < len(columnBlocks); i++ {
			for j := i + 1; j < len(columnBlocks); j++ {
				if columnBlocks[i].Y > columnBlocks[j].Y {
					columnBlocks[i], columnBlocks[j] = columnBlocks[j], columnBlocks[i]
				}
			}
		}

		// Find contiguous vertical sequences of same type
		stormSequences := gl.findVerticalStormSequences(columnBlocks)

		// Add all storm sequences to the result
		for _, sequence := range stormSequences {
			stormBlocks = append(stormBlocks, sequence...)
		}
	}

	return stormBlocks
}

// findVerticalStormSequences finds contiguous vertical sequences of 4+ same-type blocks
func (gl *GameLogic) findVerticalStormSequences(columnBlocks []Block) [][]Block {
	var sequences [][]Block
	var currentSequence []Block
	var currentType BlockType = -1

	for i, block := range columnBlocks {
		// Check if this block continues the current sequence
		if block.BlockType == currentType && (i == 0 || block.Y == columnBlocks[i-1].Y+1) {
			// Continue the sequence
			currentSequence = append(currentSequence, block)
		} else {
			// Sequence broken - check if previous sequence was long enough for storm
			if len(currentSequence) >= 4 {
				sequences = append(sequences, currentSequence)
			}
			// Start new sequence
			currentSequence = []Block{block}
			currentType = block.BlockType
		}
	}

	// Don't forget the last sequence
	if len(currentSequence) >= 4 {
		sequences = append(sequences, currentSequence)
	}

	return sequences
}

// StartElectricalStorm marks blocks as being in an electrical storm
func (gl *GameLogic) StartElectricalStorm(stormBlocks []Block) {
	if len(stormBlocks) == 0 {
		return
	}
	
	// Create a map for fast lookup
	stormMap := make(map[string]bool)
	for _, block := range stormBlocks {
		key := fmt.Sprintf("%d,%d", block.X, block.Y)
		stormMap[key] = true
	}
	
	// Mark matching blocks as being in storm (or refresh existing storm blocks)
	for i := range gl.placedBlocks {
		block := &gl.placedBlocks[i]
		key := fmt.Sprintf("%d,%d", block.X, block.Y)
		if stormMap[key] {
			if !block.IsInStorm {
				// New block joining the storm
				block.IsInStorm = true
				block.StormTime = 0
				block.StormPhase = 0
				block.SparkPhase = 0
			}
			// Note: We don't reset storm time for existing storm blocks,
			// so they maintain their continuous animation
		}
	}
}

// UpdateElectricalStorms updates the storm animation for all storm blocks
// Storms are purely visual effects and don't destroy blocks
func (gl *GameLogic) UpdateElectricalStorms(deltaTime float64) {
	for i := range gl.placedBlocks {
		block := &gl.placedBlocks[i]
		if block.IsInStorm {
			// Update storm animation phases
			block.StormTime += deltaTime
			block.StormPhase += deltaTime * StormFrequency * 2 * math.Pi
			block.SparkPhase += deltaTime * SparkFrequency * 2 * math.Pi
			
			// Storm effects continue indefinitely (until block is removed by other means)
		}
	}
}

// ClearInvalidStorms removes storm status from blocks that are no longer part of valid storm sequences
func (gl *GameLogic) ClearInvalidStorms() {
	// Find all blocks that should currently be in storms
	validStormBlocks := gl.findVerticalElectricalStorms()
	
	// Create a map for fast lookup of valid storm positions
	validStormMap := make(map[string]bool)
	for _, block := range validStormBlocks {
		key := fmt.Sprintf("%d,%d", block.X, block.Y)
		validStormMap[key] = true
	}
	
	// Clear storm status from blocks that are no longer part of valid storms
	for i := range gl.placedBlocks {
		block := &gl.placedBlocks[i]
		if block.IsInStorm {
			key := fmt.Sprintf("%d,%d", block.X, block.Y)
			if !validStormMap[key] {
				// This block is no longer part of a valid storm
				block.IsInStorm = false
				block.StormTime = 0
				block.StormPhase = 0
				block.SparkPhase = 0
			}
		}
	}
}

// CheckForElectricalStorms finds vertical sequences and starts electrical storms
// Returns 0 since storms are visual effects only (no points for non-destructive effects)
func (gl *GameLogic) CheckForElectricalStorms() int {
	stormBlocks := gl.findVerticalElectricalStorms()
	
	if len(stormBlocks) == 0 {
		return 0
	}
	
	// Start electrical storm on these blocks (visual effect only)
	gl.StartElectricalStorm(stormBlocks)
	
	// No score for electrical storms since they don't destroy blocks
	return 0
}

