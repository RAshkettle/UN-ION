package main

import (
	"fmt"
	"math"
)

type ExplosionCallback func(worldX, worldY float64, blockType BlockType)
type AudioCallback func(blocksRemoved int)
type DustCallback func(worldX, worldY float64)
type HardDropCallback func(dropHeight int)

type GameLogic struct {
	gameboard         *Gameboard
	blockManager      *BlockManager
	placedBlocks      []Block
	explosionCallback ExplosionCallback
	audioCallback     AudioCallback
	dustCallback      DustCallback
	hardDropCallback  HardDropCallback
	activeStorms      map[int]*Storm
}

func NewGameLogic(gameboard *Gameboard, blockManager *BlockManager) *GameLogic {
	return &GameLogic{
		gameboard:    gameboard,
		blockManager: blockManager,
		placedBlocks: make([]Block, 0),
		activeStorms: make(map[int]*Storm),
	}
}

func (gl *GameLogic) SetExplosionCallback(callback ExplosionCallback) {
	gl.explosionCallback = callback
}

func (gl *GameLogic) SetAudioCallback(callback AudioCallback) {
	gl.audioCallback = callback
}

func (gl *GameLogic) SetDustCallback(callback DustCallback) {
	gl.dustCallback = callback
}

func (gl *GameLogic) SetHardDropCallback(callback HardDropCallback) {
	gl.hardDropCallback = callback
}

func (gl *GameLogic) IsValidPosition(piece *TetrisPiece, offsetX, offsetY int) bool {
	blockSize := gl.blockManager.GetScaledBlockSize(gl.gameboard.Width, gl.gameboard.Height)
	gameboardWidthInBlocks := int(float64(gl.gameboard.Width) / blockSize)
	gameboardHeightInBlocks := int(float64(gl.gameboard.Height) / blockSize)
	for _, block := range piece.Blocks {
		newX := piece.X + block.X + offsetX
		newY := piece.Y + block.Y + offsetY
		if newX < 0 || newX >= gameboardWidthInBlocks || newY >= gameboardHeightInBlocks {
			return false
		}
		for _, placedBlock := range gl.placedBlocks {
			if placedBlock.X == newX && placedBlock.Y == newY {
				return false
			}
		}
	}
	return true
}

func (gl *GameLogic) IsValidPositionIgnoreNeutral(piece *TetrisPiece, offsetX, offsetY int) bool {
	blockSize := gl.blockManager.GetScaledBlockSize(gl.gameboard.Width, gl.gameboard.Height)
	gameboardWidthInBlocks := int(float64(gl.gameboard.Width) / blockSize)
	gameboardHeightInBlocks := int(float64(gl.gameboard.Height) / blockSize)
	for _, block := range piece.Blocks {
		newX := piece.X + block.X + offsetX
		newY := piece.Y + block.Y + offsetY
		if newX < 0 || newX >= gameboardWidthInBlocks || newY >= gameboardHeightInBlocks {
			return false
		}
		for _, placedBlock := range gl.placedBlocks {
			if placedBlock.BlockType != NeutralBlock && placedBlock.X == newX && placedBlock.Y == newY {
				return false
			}
		}
	}
	return true
}

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
	if gl.dustCallback != nil {
		blockSize := gl.blockManager.GetScaledBlockSize(gl.gameboard.Width, gl.gameboard.Height)
		bottomY := piece.Y
		for _, block := range piece.Blocks {
			if piece.Y+block.Y > bottomY {
				bottomY = piece.Y + block.Y
			}
		}
		worldX := float64(gl.gameboard.X) + float64(piece.X)*blockSize + blockSize/2
		worldY := float64(gl.gameboard.Y) + float64(bottomY+1)*blockSize
		gl.dustCallback(worldX, worldY)
	}
}

func (gl *GameLogic) GetPlacedBlocks() []Block {
	return gl.placedBlocks
}

func (gl *GameLogic) SpawnNewPiece(pieceType PieceType) *TetrisPiece {
	blockSize := gl.blockManager.GetScaledBlockSize(gl.gameboard.Width, gl.gameboard.Height)
	gameboardWidthInBlocks := int(float64(gl.gameboard.Width) / blockSize)
	centerX := gameboardWidthInBlocks / 2
	return gl.blockManager.CreateTetrisPiece(pieceType, centerX, 0)
}

func (gl *GameLogic) TryRotatePiece(piece *TetrisPiece, pieceType PieceType) bool {
	if piece == nil {
		return false
	}
	originalRotation := piece.Rotation
	originalBlocks := make([]Block, len(piece.Blocks))
	copy(originalBlocks, piece.Blocks)
	gl.blockManager.RotatePiece(piece, pieceType)
	if !gl.IsValidPosition(piece, 0, 0) {
		piece.Rotation = originalRotation
		piece.Blocks = originalBlocks
		return false
	}
	return true
}

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

func (gl *GameLogic) CalculateDropPosition(piece *TetrisPiece) *TetrisPiece {
	if piece == nil {
		return nil
	}
	shadowPiece := &TetrisPiece{
		X:        piece.X,
		Y:        piece.Y,
		Rotation: piece.Rotation,
		Blocks:   make([]Block, len(piece.Blocks)),
	}
	for i, block := range piece.Blocks {
		shadowPiece.Blocks[i] = Block{
			X:         block.X,
			Y:         block.Y,
			BlockType: block.BlockType,
		}
	}
	for gl.IsValidPosition(shadowPiece, 0, 1) {
		shadowPiece.Y++
	}
	return shadowPiece
}

func (gl *GameLogic) IsGameOver() bool {
	for _, block := range gl.placedBlocks {
		if block.Y <= 0 && block.BlockType != NeutralBlock {
			return true
		}
	}
	return false
}

func (gl *GameLogic) CheckAndProcessReactions() int {
	totalScore := 0
	for {
		blocksToRemove := gl.findBlocksToRemove()
		if len(blocksToRemove) == 0 {
			break
		}
		reactionScore := gl.calculateReactionScore(len(blocksToRemove))
		totalScore += reactionScore
		gl.removeBlocks(blocksToRemove)
		gl.processBlockFalling()
	}
	return totalScore
}

func (gl *GameLogic) calculateReactionScore(blocksRemoved int) int {
	if blocksRemoved < 4 {
		return 0
	}
	score := 10
	for i := 4; i < blocksRemoved; i++ {
		score *= 2
	}
	return score
}

func (gl *GameLogic) findBlocksToRemove() []Block {
	var blocksToRemove []Block
	rowMap := make(map[int][]Block)
	for _, block := range gl.placedBlocks {
		rowMap[block.Y] = append(rowMap[block.Y], block)
	}
	for _, rowBlocks := range rowMap {
		if len(rowBlocks) < 3 {
			continue
		}
		for i := 0; i < len(rowBlocks); i++ {
			for j := i + 1; j < len(rowBlocks); j++ {
				if rowBlocks[i].X > rowBlocks[j].X {
					rowBlocks[i], rowBlocks[j] = rowBlocks[j], rowBlocks[i]
				}
			}
		}
		clusters := gl.findClustersInRow(rowBlocks)
		for _, cluster := range clusters {
			if len(cluster) >= 3 {
				zeroSumBlocks := gl.findZeroSumSubsequence(cluster)
				blocksToRemove = append(blocksToRemove, zeroSumBlocks...)
			}
		}
	}
	return blocksToRemove
}

func (gl *GameLogic) findClustersInRow(rowBlocks []Block) [][]Block {
	var clusters [][]Block
	var currentCluster []Block
	for i, block := range rowBlocks {
		if block.BlockType == NeutralBlock {
			if len(currentCluster) >= 3 {
				clusters = append(clusters, currentCluster)
			}
			currentCluster = nil
		} else if i > 0 && block.X != rowBlocks[i-1].X+1 {
			if len(currentCluster) >= 3 {
				clusters = append(clusters, currentCluster)
			}
			currentCluster = []Block{block}
		} else {
			currentCluster = append(currentCluster, block)
		}
	}
	if len(currentCluster) >= 3 {
		clusters = append(clusters, currentCluster)
	}
	return clusters
}

func (gl *GameLogic) findZeroSumSubsequence(cluster []Block) []Block {
	for length := len(cluster); length >= 3; length-- {
		for start := 0; start <= len(cluster)-length; start++ {
			subsequence := cluster[start : start+length]
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

func (gl *GameLogic) removeBlocks(blocksToRemove []Block) {
	if len(blocksToRemove) == 0 {
		return
	}
	if gl.audioCallback != nil {
		gl.audioCallback(len(blocksToRemove))
	}
	removeMap := make(map[string]bool)
	for _, block := range blocksToRemove {
		key := fmt.Sprintf("%d,%d", block.X, block.Y)
		removeMap[key] = true
	}
	var remainingBlocks []Block
	for _, block := range gl.placedBlocks {
		key := fmt.Sprintf("%d,%d", block.X, block.Y)
		if !removeMap[key] {
			remainingBlocks = append(remainingBlocks, block)
		} else {
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

func (gl *GameLogic) processBlockFalling() {
	for i := 0; i < len(gl.placedBlocks); i++ {
		for j := i + 1; j < len(gl.placedBlocks); j++ {
			if gl.placedBlocks[i].Y < gl.placedBlocks[j].Y {
				gl.placedBlocks[i], gl.placedBlocks[j] = gl.placedBlocks[j], gl.placedBlocks[i]
			}
		}
	}
	for i := range gl.placedBlocks {
		block := &gl.placedBlocks[i]
		if !block.IsFalling {
			gl.StartBlockFall(block)
		}
	}
}

func (gl *GameLogic) UpdateWobblingBlocks(deltaTime float64) bool {
	anyBlocksFinished := false
	for i := range gl.placedBlocks {
		block := &gl.placedBlocks[i]
		if block.IsWobbling {
			block.WobbleTime += deltaTime
			block.WobblePhase += deltaTime * WobbleFrequency * 2 * math.Pi
			if block.WobbleTime >= WobbleDuration {
				block.ShowPowSprite = false
				anyBlocksFinished = true
			}
		}
	}
	return anyBlocksFinished
}

func (gl *GameLogic) RemoveFinishedWobblingBlocks() int {
	var blocksToRemove []Block
	var remainingBlocks []Block
	for _, block := range gl.placedBlocks {
		if block.IsWobbling && block.WobbleTime >= WobbleDuration {
			block.ShowPowSprite = false
			blocksToRemove = append(blocksToRemove, block)
		} else {
			remainingBlocks = append(remainingBlocks, block)
		}
	}
	if len(blocksToRemove) == 0 {
		return 0
	}
	if gl.audioCallback != nil {
		gl.audioCallback(len(blocksToRemove))
	}
	for _, block := range blocksToRemove {
		if gl.explosionCallback != nil {
			blockSize := gl.blockManager.GetScaledBlockSize(gl.gameboard.Width, gl.gameboard.Height)
			worldX := float64(gl.gameboard.X) + float64(block.X)*blockSize + blockSize/2
			worldY := float64(gl.gameboard.Y) + float64(block.Y)*blockSize + blockSize/2
			gl.explosionCallback(worldX, worldY, block.BlockType)
		}
	}
	gl.placedBlocks = remainingBlocks
	return len(blocksToRemove)
}

func (gl *GameLogic) StartBlockWobbling(blocksToWobble []Block) {
	if len(blocksToWobble) == 0 {
		return
	}
	wobbleMap := make(map[string]bool)
	for _, block := range blocksToWobble {
		key := fmt.Sprintf("%d,%d", block.X, block.Y)
		wobbleMap[key] = true
	}
	for i := range gl.placedBlocks {
		block := &gl.placedBlocks[i]
		key := fmt.Sprintf("%d,%d", block.X, block.Y)
		if wobbleMap[key] && !block.IsWobbling {
			block.IsWobbling = true
			block.WobbleTime = 0
			block.WobblePhase = 0
			block.ShowPowSprite = true
		}
	}
}

func (gl *GameLogic) CheckForNewReactions() int {
	blocksToWobble := gl.findNonWobblingBlocksToRemove()
	if len(blocksToWobble) == 0 {
		return 0
	}
	gl.StartBlockWobbling(blocksToWobble)
	return gl.calculateReactionScore(len(blocksToWobble))
}

func (gl *GameLogic) findNonWobblingBlocksToRemove() []Block {
	var blocksToRemove []Block
	rowMap := make(map[int][]Block)
	for _, block := range gl.placedBlocks {
		if !block.IsWobbling {
			rowMap[block.Y] = append(rowMap[block.Y], block)
		}
	}
	for _, rowBlocks := range rowMap {
		if len(rowBlocks) < 3 {
			continue
		}
		for i := 0; i < len(rowBlocks); i++ {
			for j := i + 1; j < len(rowBlocks); j++ {
				if rowBlocks[i].X > rowBlocks[j].X {
					rowBlocks[i], rowBlocks[j] = rowBlocks[j], rowBlocks[i]
				}
			}
		}
		clusters := gl.findClustersInRow(rowBlocks)
		for _, cluster := range clusters {
			if len(cluster) >= 3 {
				zeroSumBlocks := gl.findZeroSumSubsequence(cluster)
				blocksToRemove = append(blocksToRemove, zeroSumBlocks...)
			}
		}
	}
	return blocksToRemove
}

// All storm-related methods have been removed from this file. See storm.go for their implementations.

func (gl *GameLogic) UpdateFallingBlocks(deltaTime float64) bool {
	anyBlocksLanded := false
	for i := range gl.placedBlocks {
		block := &gl.placedBlocks[i]
		if block.IsFalling {
			fallDistance := block.FallTargetY - block.FallStartY
			if fallDistance > 0 {
				block.FallProgress += deltaTime * FallSpeed / fallDistance
				if block.FallProgress >= 1.0 {
					block.FallProgress = 1.0
					block.Y = int(block.FallTargetY)
					block.IsFalling = false
					anyBlocksLanded = true
				}
			} else {
				block.IsFalling = false
				anyBlocksLanded = true
			}
		}
	}
	return anyBlocksLanded
}

func (gl *GameLogic) GetBlockRenderPosition(block *Block) (float64, float64) {
	if block.IsArcing {
		x, y, _, _ := gl.GetBlockArcPosition(block)
		return x, y
	} else if block.IsFalling {
		currentY := block.FallStartY + (block.FallTargetY-block.FallStartY)*block.FallProgress
		return float64(block.X), currentY
	}
	return float64(block.X), float64(block.Y)
}

func (gl *GameLogic) GetBlockRenderTransform(block *Block) (float64, float64, float64, float64) {
	if block.IsArcing {
		return gl.GetBlockArcPosition(block)
	} else if block.IsFalling {
		currentY := block.FallStartY + (block.FallTargetY-block.FallStartY)*block.FallProgress
		return float64(block.X), currentY, 0.0, 1.0
	}
	return float64(block.X), float64(block.Y), 0.0, 1.0
}

func (gl *GameLogic) StartBlockArc(block *Block, startX, startY, targetX, targetY float64) {
	block.IsArcing = true
	block.ArcStartX = startX
	block.ArcStartY = startY
	block.ArcTargetX = targetX
	block.ArcTargetY = targetY
	block.ArcProgress = 0.0
	block.ArcRotation = 0.0
	block.ArcScale = MinArcScale
}

func (gl *GameLogic) GetBlockArcPosition(block *Block) (float64, float64, float64, float64) {
	if !block.IsArcing {
		return float64(block.X), float64(block.Y), 0.0, 1.0
	}
	t := block.ArcProgress
	currentX := block.ArcStartX + (block.ArcTargetX-block.ArcStartX)*t
	linearY := block.ArcStartY + (block.ArcTargetY-block.ArcStartY)*t
	arcOffset := ArcHeight * 4 * t * (1 - t)
	currentY := linearY - arcOffset
	currentRotation := block.ArcRotation + MaxRotation*t
	currentScale := MinArcScale + (1.0-MinArcScale)*t
	return currentX, currentY, currentRotation, currentScale
}

func (gl *GameLogic) UpdateArcingBlocks(deltaTime float64) bool {
	anyBlocksFinishedArcing := false
	for i := range gl.placedBlocks {
		block := &gl.placedBlocks[i]
		if block.IsArcing {
			block.ArcProgress += deltaTime * ArcSpeed
			if block.ArcProgress >= 1.0 {
				block.ArcProgress = 1.0
				blockSize := gl.blockManager.GetScaledBlockSize(gl.gameboard.Width, gl.gameboard.Height)
				gameboardHeightInBlocks := int(float64(gl.gameboard.Height) / blockSize)
				targetColumn := int(block.ArcTargetX)
				finalY := gameboardHeightInBlocks - 1
				for _, placedBlock := range gl.placedBlocks {
					if &placedBlock != block && placedBlock.X == targetColumn && placedBlock.Y < finalY {
						finalY = placedBlock.Y - 1
					}
				}
				if finalY < 0 {
					finalY = 0
				}
				block.X = targetColumn
				block.Y = finalY
				block.IsArcing = false
				block.ArcScale = 1.0
				block.ArcRotation = 0.0
				anyBlocksFinishedArcing = true
			}
		}
	}
	return anyBlocksFinishedArcing
}

func (gl *GameLogic) StartBlockFall(block *Block) {
	blockSize := gl.blockManager.GetScaledBlockSize(gl.gameboard.Width, gl.gameboard.Height)
	gameboardHeightInBlocks := int(float64(gl.gameboard.Height) / blockSize)
	targetY := block.Y
	for newY := block.Y + 1; newY < gameboardHeightInBlocks; newY++ {
		occupied := false
		for _, placedBlock := range gl.placedBlocks {
			if &placedBlock == block {
				continue
			}
			checkY := placedBlock.Y
			if placedBlock.IsFalling {
				checkY = int(placedBlock.FallTargetY)
			}
			if placedBlock.X == block.X && checkY == newY {
				occupied = true
				break
			}
		}
		if occupied {
			break
		}
		targetY = newY
	}
	block.IsFalling = true
	block.FallStartY = float64(block.Y)
	block.FallTargetY = float64(targetY)
	block.FallProgress = 0
}
