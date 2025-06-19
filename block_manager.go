package main

import (
	"math"
	"math/rand"
	"union/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	WobbleDuration   = 0.8
	WobbleIntensity  = 2.0
	WobbleFrequency  = 8.0
	StormDuration    = 1.2
	StormIntensity   = 3.0
	StormFrequency   = 12.0
	SparkFrequency   = 15.0
	FallSpeed        = 4.0
	ArcSpeed         = 4.0
	ArcHeight        = 3.0
	MinArcScale      = 0.1
	MaxRotation      = 360.0
	WarningDuration  = 1.0
	WarningIntensity = 3.0
	WarningFrequency = 20.0
)

type BlockManager struct {
	blockSize float64
}

func NewBlockManager() *BlockManager {
	return &BlockManager{
		blockSize: 16.0,
	}
}

func (bm *BlockManager) GetScaledBlockSize(gameboardWidth, gameboardHeight int) float64 {
	baseBlocksWide := 12.0
	baseBlocksTall := 20.0
	blockSizeFromWidth := float64(gameboardWidth) / baseBlocksWide
	blockSizeFromHeight := float64(gameboardHeight) / baseBlocksTall
	if blockSizeFromWidth < blockSizeFromHeight {
		return blockSizeFromWidth
	}
	return blockSizeFromHeight
}

func (bm *BlockManager) GenerateRandomBlockType() BlockType {
	r := rand.Float64()
	if r < 0.4 {
		return PositiveBlock
	} else if r < 0.8 {
		return NegativeBlock
	}
	return NeutralBlock
}

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

func (bm *BlockManager) DrawBlock(screen *ebiten.Image, block Block, worldX, worldY float64, blockSize float64) {
	sprite := bm.GetBlockSprite(block.BlockType)

	op := &ebiten.DrawImageOptions{}

	scaleX := blockSize / float64(sprite.Bounds().Dx())
	scaleY := blockSize / float64(sprite.Bounds().Dy())
	op.GeoM.Scale(scaleX, scaleY)

	if block.IsInStorm {
		stormX := math.Sin(block.StormPhase) * StormIntensity
		stormY := math.Cos(block.StormPhase*1.7) * StormIntensity * 0.3

		sparkOffset := math.Sin(block.SparkPhase) * 1.0
		stormX += sparkOffset
		stormY += math.Cos(block.SparkPhase*2.3) * 0.5

		op.GeoM.Translate(worldX+stormX, worldY+stormY)

		flickerIntensity := 0.2 + 0.1*math.Sin(block.SparkPhase*2)

		switch block.BlockType {
		case PositiveBlock:
			op.ColorScale.Scale(float32(1.0+flickerIntensity), float32(1.0+flickerIntensity*0.5), float32(1.0-flickerIntensity*0.3), 1.0)
		case NegativeBlock:
			op.ColorScale.Scale(float32(1.0-flickerIntensity*0.3), float32(1.0+flickerIntensity*0.5), float32(1.0+flickerIntensity), 1.0)
		}

	} else if block.IsWobbling {
		wobbleX := math.Sin(block.WobblePhase) * WobbleIntensity
		wobbleY := math.Cos(block.WobblePhase*1.3) * WobbleIntensity * 0.5

		op.GeoM.Translate(worldX+wobbleX, worldY+wobbleY)

		wobbleProgress := block.WobbleTime / WobbleDuration
		alpha := 1.0 - wobbleProgress*0.3
		op.ColorScale.Scale(1, 1, 1, float32(alpha))
	} else {
		op.GeoM.Translate(worldX, worldY)
	}

	screen.DrawImage(sprite, op)

	if block.IsWobbling && block.ShowPowSprite {
		bm.DrawPowSprite(screen, worldX, worldY, block.WobblePhase, blockSize)
	}
}

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
		sprite = assets.NeutralChargeSprite
	}

	if sprite != nil {
		op := &ebiten.DrawImageOptions{}
		spriteWidth := float64(sprite.Bounds().Dx())
		spriteHeight := float64(sprite.Bounds().Dy())
		scaleX := blockSize / spriteWidth
		scaleY := blockSize / spriteHeight
		op.GeoM.Scale(scaleX, scaleY)
		op.GeoM.Translate(worldX, worldY)

		op.ColorScale.ScaleAlpha(0.4)
		op.ColorScale.Scale(0.5, 0.5, 0.5, 1.0)

		screen.DrawImage(sprite, op)
	}
}

func (bm *BlockManager) DrawBlockTransformed(screen *ebiten.Image, block Block, worldX, worldY, rotation, scale, blockSize float64) {
	sprite := bm.GetBlockSprite(block.BlockType)

	op := &ebiten.DrawImageOptions{}

	scaleX := (blockSize * scale) / float64(sprite.Bounds().Dx())
	scaleY := (blockSize * scale) / float64(sprite.Bounds().Dy())

	centerX := float64(sprite.Bounds().Dx()) / 2
	centerY := float64(sprite.Bounds().Dy()) / 2

	op.GeoM.Translate(-centerX, -centerY)
	op.GeoM.Rotate(rotation * math.Pi / 180.0)
	op.GeoM.Scale(scaleX, scaleY)

	if block.IsInStorm {
		stormX := math.Sin(block.StormPhase) * StormIntensity
		stormY := math.Cos(block.StormPhase*1.7) * StormIntensity * 0.3

		sparkOffset := math.Sin(block.SparkPhase) * 1.0
		stormX += sparkOffset
		stormY += math.Cos(block.SparkPhase*2.3) * 0.5

		op.GeoM.Translate(worldX+stormX+(blockSize*scale)/2, worldY+stormY+(blockSize*scale)/2)

		flickerIntensity := 0.2 + 0.1*math.Sin(block.SparkPhase*2)
		switch block.BlockType {
		case PositiveBlock:
			op.ColorScale.Scale(float32(1.0+flickerIntensity), float32(1.0+flickerIntensity*0.5), float32(1.0-flickerIntensity*0.3), 1.0)
		case NegativeBlock:
			op.ColorScale.Scale(float32(1.0-flickerIntensity*0.3), float32(1.0+flickerIntensity*0.5), float32(1.0+flickerIntensity), 1.0)
		}
	} else if block.IsWobbling {
		wobbleX := math.Sin(block.WobblePhase) * WobbleIntensity
		wobbleY := math.Cos(block.WobblePhase*1.3) * WobbleIntensity * 0.5

		op.GeoM.Translate(worldX+wobbleX+(blockSize*scale)/2, worldY+wobbleY+(blockSize*scale)/2)

		wobbleProgress := block.WobbleTime / WobbleDuration
		alpha := 1.0 - wobbleProgress*0.3
		op.ColorScale.Scale(1, 1, 1, float32(alpha))
	} else {
		op.GeoM.Translate(worldX+(blockSize*scale)/2, worldY+(blockSize*scale)/2)
	}

	screen.DrawImage(sprite, op)

	if block.IsWobbling && block.ShowPowSprite {
		bm.DrawPowSprite(screen, worldX, worldY, block.WobblePhase, blockSize*scale)
	}
}

func (bm *BlockManager) DrawWarningSprite(screen *ebiten.Image, worldX, worldY, warningTime, blockSize float64, column int, gameboardWidth int) {
	sprite := assets.ZapSprite
	if sprite == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}

	scaleX := blockSize / float64(sprite.Bounds().Dx())
	scaleY := blockSize / float64(sprite.Bounds().Dy())
	op.GeoM.Scale(scaleX, scaleY)

	shakePhase := warningTime * WarningFrequency * 2 * math.Pi
	shakeX := math.Sin(shakePhase) * WarningIntensity
	shakeY := math.Cos(shakePhase*1.3) * WarningIntensity * 0.7

	wobblePhase := warningTime * WarningFrequency * 1.5 * math.Pi
	wobbleX := math.Sin(wobblePhase) * WarningIntensity * 0.5
	wobbleY := math.Cos(wobblePhase*0.8) * WarningIntensity * 0.3

	gameboardWidthInBlocks := int(float64(gameboardWidth) / blockSize)
	boardCenter := gameboardWidthInBlocks / 2

	var offsetX float64
	if column >= boardCenter {
		offsetX = -blockSize * 0.3
	} else {
		offsetX = blockSize * 0.7
	}
	offsetY := blockSize * 0.1

	op.GeoM.Translate(worldX+offsetX+shakeX+wobbleX, worldY+offsetY+shakeY+wobbleY)

	op.ColorScale.Scale(1.2, 1.1, 0.9, 1.0)

	screen.DrawImage(sprite, op)
}

func (bm *BlockManager) DrawPowSprite(screen *ebiten.Image, worldX, worldY, wobblePhase, blockSize float64) {
	sprite := assets.PowSprite
	if sprite == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}

	scaleX := blockSize / float64(sprite.Bounds().Dx())
	scaleY := blockSize / float64(sprite.Bounds().Dy())

	wobbleX := math.Sin(wobblePhase) * WobbleIntensity
	wobbleY := math.Cos(wobblePhase*1.3) * WobbleIntensity * 0.5

	spriteWidth := float64(sprite.Bounds().Dx()) * scaleX
	spriteHeight := float64(sprite.Bounds().Dy()) * scaleY
	centerOffsetX := (blockSize - spriteWidth) / 2
	centerOffsetY := (blockSize - spriteHeight) / 2

	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(worldX+centerOffsetX+wobbleX, worldY+centerOffsetY+wobbleY)

	op.ColorScale.Scale(1.0, 1.0, 1.0, 0.9)

	screen.DrawImage(sprite, op)
}
