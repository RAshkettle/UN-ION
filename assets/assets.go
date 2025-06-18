package assets

import (
	"bytes"
	"embed"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed *
var assets embed.FS

var (
	PositiveChargeSprite = loadImage("images/PositiveCharge.png")
	NegativeChargeSprite = loadImage("images/NegativeCharge.png")
	NeutralChargeSprite  = loadImage("images/NeutralCharge.png")

	BlockBreakSound = loadAudio("audio/coin.ogg")
	SwooshSound = loadAudio("audio/swoosh.ogg")
	BackgroundMusic = loadAudio("audio/background_music.mp3")
)

func loadImage(filePath string) *ebiten.Image {
	data, err := assets.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	ebitenImg := ebiten.NewImageFromImage(img)
	return ebitenImg
}

func loadAudio(filePath string) []byte {
	data, err := assets.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return data
}