package main

import (
	"fmt"
	"math"
	"math/rand"
)

type Storm struct {
	Column      int
	Timer       float64
	NextDrop    float64
	IsActive    bool
	IsWarning   bool
	WarningTime float64
}

func (gl *GameLogic) findVerticalElectricalStorms() []Block {
	var stormBlocks []Block
	columnMap := make(map[int][]Block)
	for _, block := range gl.placedBlocks {
		if !block.IsWobbling && block.BlockType != NeutralBlock {
			columnMap[block.X] = append(columnMap[block.X], block)
		}
	}
	for _, columnBlocks := range columnMap {
		if len(columnBlocks) < 4 {
			continue
		}
		for i := 0; i < len(columnBlocks); i++ {
			for j := i + 1; j < len(columnBlocks); j++ {
				if columnBlocks[i].Y > columnBlocks[j].Y {
					columnBlocks[i], columnBlocks[j] = columnBlocks[j], columnBlocks[i]
				}
			}
		}
		stormSequences := gl.findVerticalStormSequences(columnBlocks)
		for _, sequence := range stormSequences {
			stormBlocks = append(stormBlocks, sequence...)
		}
	}
	return stormBlocks
}

func (gl *GameLogic) findVerticalStormSequences(columnBlocks []Block) [][]Block {
	var sequences [][]Block
	var currentSequence []Block
	var currentType BlockType = -1
	for i, block := range columnBlocks {
		if block.BlockType == currentType && (i == 0 || block.Y == columnBlocks[i-1].Y+1) {
			currentSequence = append(currentSequence, block)
		} else {
			if len(currentSequence) >= 4 {
				sequences = append(sequences, currentSequence)
			}
			currentSequence = []Block{block}
			currentType = block.BlockType
		}
	}
	if len(currentSequence) >= 4 {
		sequences = append(sequences, currentSequence)
	}
	return sequences
}

func (gl *GameLogic) StartElectricalStorm(stormBlocks []Block) {
	if len(stormBlocks) == 0 {
		return
	}
	stormMap := make(map[string]bool)
	for _, block := range stormBlocks {
		key := fmt.Sprintf("%d,%d", block.X, block.Y)
		stormMap[key] = true
	}
	for i := range gl.placedBlocks {
		block := &gl.placedBlocks[i]
		key := fmt.Sprintf("%d,%d", block.X, block.Y)
		if stormMap[key] {
			if !block.IsInStorm {
				block.IsInStorm = true
				block.StormTime = 0
				block.StormPhase = 0
				block.SparkPhase = 0
			}
		}
	}
}

func (gl *GameLogic) UpdateElectricalStorms(deltaTime float64) {
	for i := range gl.placedBlocks {
		block := &gl.placedBlocks[i]
		if block.IsInStorm {
			block.StormTime += deltaTime
			block.StormPhase += deltaTime * StormFrequency * 2 * math.Pi
			block.SparkPhase += deltaTime * SparkFrequency * 2 * math.Pi
		}
	}
}

func (gl *GameLogic) ClearInvalidStorms() {
	validStormBlocks := gl.findVerticalElectricalStorms()
	validStormMap := make(map[string]bool)
	for _, block := range validStormBlocks {
		key := fmt.Sprintf("%d,%d", block.X, block.Y)
		validStormMap[key] = true
	}
	// Track columns and block types with <4 blocks
	columnTypeCount := make(map[[2]int]int)
	for _, block := range gl.placedBlocks {
		if block.IsInStorm {
			columnTypeCount[[2]int{block.X, int(block.BlockType)}]++
		}
	}
	for i := range gl.placedBlocks {
		block := &gl.placedBlocks[i]
		if block.IsInStorm {
			key := fmt.Sprintf("%d,%d", block.X, block.Y)
			if !validStormMap[key] || columnTypeCount[[2]int{block.X, int(block.BlockType)}] < 4 {
				block.IsInStorm = false
				block.StormTime = 0
				block.StormPhase = 0
				block.SparkPhase = 0
			}
		}
	}
}

func (gl *GameLogic) CheckForElectricalStorms() int {
	stormBlocks := gl.findVerticalElectricalStorms()
	if len(stormBlocks) == 0 {
		return 0
	}
	gl.StartElectricalStorm(stormBlocks)
	gl.UpdateActiveStorms()
	return 0
}

func (gl *GameLogic) generateStormTimer() float64 {
	return 3.0 + rand.Float64()*2.0
}

func (gl *GameLogic) UpdateStormTimers(deltaTime float64) []Block {
	var newNeutralBlocks []Block
	for _, storm := range gl.activeStorms {
		if storm.IsActive {
			storm.Timer += deltaTime
			timeUntilSpawn := storm.NextDrop - storm.Timer
			if timeUntilSpawn <= WarningDuration && !storm.IsWarning {
				storm.IsWarning = true
				storm.WarningTime = 0
			}
			if storm.IsWarning {
				storm.WarningTime += deltaTime
			}
			if storm.Timer >= storm.NextDrop {
				storm.IsWarning = false
				storm.WarningTime = 0
				highestStormBlock := gl.FindHighestStormBlock(storm.Column)
				if highestStormBlock == nil {
					storm.Timer = 0
					storm.NextDrop = gl.generateStormTimer()
					continue
				}
				blockSize := gl.blockManager.GetScaledBlockSize(gl.gameboard.Width, gl.gameboard.Height)
				gameboardWidthInBlocks := int(float64(gl.gameboard.Width) / blockSize)
				targetColumn := rand.Intn(gameboardWidthInBlocks)
				neutralBlock := Block{
					X:         targetColumn,
					Y:         0,
					BlockType: NeutralBlock,
					IsArcing:  true,
					IsFalling: false,
				}
				gl.StartBlockArc(&neutralBlock,
					float64(highestStormBlock.X), float64(highestStormBlock.Y),
					float64(targetColumn), 0.0)
				positionFree := true
				for _, placedBlock := range gl.placedBlocks {
					if placedBlock.X == targetColumn && placedBlock.Y == 0 {
						positionFree = false
						break
					}
				}
				if positionFree {
					newNeutralBlocks = append(newNeutralBlocks, neutralBlock)
				}
				storm.Timer = 0
				storm.NextDrop = gl.generateStormTimer()
			}
		}
	}
	return newNeutralBlocks
}

func (gl *GameLogic) AddNeutralBlock(block Block) *Block {
	gl.placedBlocks = append(gl.placedBlocks, block)
	return &gl.placedBlocks[len(gl.placedBlocks)-1]
}

func (gl *GameLogic) UpdateActiveStorms() {
	stormColumns := make(map[int]bool)
	for _, block := range gl.placedBlocks {
		if block.IsInStorm {
			stormColumns[block.X] = true
		}
	}
	for column := range stormColumns {
		if _, exists := gl.activeStorms[column]; !exists {
			gl.activeStorms[column] = &Storm{
				Column:   column,
				Timer:    0,
				NextDrop: gl.generateStormTimer(),
				IsActive: true,
			}
		}
	}
	for column, storm := range gl.activeStorms {
		if !stormColumns[column] {
			storm.IsActive = false
			delete(gl.activeStorms, column)
		}
	}
}

func (gl *GameLogic) FindHighestStormBlock(column int) *Block {
	var highestBlock *Block
	highestY := 999
	for i := range gl.placedBlocks {
		block := &gl.placedBlocks[i]
		if block.X == column && block.IsInStorm && block.Y < highestY {
			highestY = block.Y
			highestBlock = block
		}
	}
	return highestBlock
}

func (gl *GameLogic) GetStormWarnings() []struct {
	Column        int
	WarningTime   float64
	HighestBlockY int
} {
	var warnings []struct {
		Column        int
		WarningTime   float64
		HighestBlockY int
	}
	for _, storm := range gl.activeStorms {
		if storm.IsWarning {
			highestBlock := gl.FindHighestStormBlock(storm.Column)
			if highestBlock != nil {
				warnings = append(warnings, struct {
					Column        int
					WarningTime   float64
					HighestBlockY int
				}{
					Column:        storm.Column,
					WarningTime:   storm.WarningTime,
					HighestBlockY: highestBlock.Y,
				})
			}
		}
	}
	return warnings
}
