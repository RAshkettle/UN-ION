package main

import (
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Un-ion")
	ebiten.SetWindowSize(1200, 800) 

	sceneManager := NewSceneManager()

	err := ebiten.RunGame(sceneManager)
	if err != nil {
		panic(err)
	}
}
