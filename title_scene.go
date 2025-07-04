package main

import (
	"bytes"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

type TitleScene struct {
	sceneManager *SceneManager
	titleFont    *text.GoTextFace
	subtitleFont *text.GoTextFace
	helpFont     *text.GoTextFace
	showHelp     bool
	prevHPressed bool // for just-pressed logic
}

func (t *TitleScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 15, 25, 255})
	t.drawTitleScreen(screen)
}

func (t *TitleScene) drawTitleScreen(screen *ebiten.Image) {
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()

	// Draw title
	titleText := "UN-ION"
	titleBounds, _ := text.Measure(titleText, t.titleFont, 0)
	titleX := (w - int(titleBounds)) / 2
	titleY := h/2 - 50

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(titleX), float64(titleY))
	op.ColorScale.ScaleWithColor(color.RGBA{220, 220, 255, 255})
	text.Draw(screen, titleText, t.titleFont, op)

	// Draw subtitle
	subtitleText := "Press any key to Start"
	subtitleBounds, _ := text.Measure(subtitleText, t.subtitleFont, 0)
	subtitleX := (w - int(subtitleBounds)) / 2
	subtitleY := titleY + 80

	op2 := &text.DrawOptions{}
	op2.GeoM.Translate(float64(subtitleX), float64(subtitleY))
	op2.ColorScale.ScaleWithColor(color.RGBA{180, 180, 200, 255})
	text.Draw(screen, subtitleText, t.subtitleFont, op2)

	// Draw help prompt
	helpPrompt := "Press H for Help"
	helpPromptBounds, _ := text.Measure(helpPrompt, t.subtitleFont, 0)
	helpPromptX := (w - int(helpPromptBounds)) / 2
	helpPromptY := subtitleY + 50 // More space from subtitle

	op3 := &text.DrawOptions{}
	op3.GeoM.Translate(float64(helpPromptX), float64(helpPromptY))
	op3.ColorScale.ScaleWithColor(color.RGBA{255, 255, 100, 255})
	text.Draw(screen, helpPrompt, t.subtitleFont, op3)

	controls := []string{
		"Quick Controls:",
		"WASD/Arrow Keys: Move piece",
		"Space: Rotate piece",
	}

	helpStartY := helpPromptY + 40
	for i, control := range controls {
		controlBounds, _ := text.Measure(control, t.helpFont, 0)
		controlX := (w - int(controlBounds)) / 2
		controlY := helpStartY + i*18

		op4 := &text.DrawOptions{}
		op4.GeoM.Translate(float64(controlX), float64(controlY))
		if i == 0 {
			op4.ColorScale.ScaleWithColor(color.RGBA{200, 200, 220, 255})
		} else {
			op4.ColorScale.ScaleWithColor(color.RGBA{150, 150, 170, 255})
		}
		text.Draw(screen, control, t.helpFont, op4)
	}
}

func (t *TitleScene) Update() error {
	hPressed := ebiten.IsKeyPressed(ebiten.KeyH)
	if hPressed && !t.prevHPressed {
		t.sceneManager.TransitionTo(SceneHelp)
	}
	t.prevHPressed = hPressed

	if ebiten.IsKeyPressed(ebiten.KeySpace) ||
		ebiten.IsKeyPressed(ebiten.KeyEnter) ||
		ebiten.IsKeyPressed(ebiten.KeyEscape) ||
		inpututil.IsKeyJustPressed(ebiten.KeyA) ||
		inpututil.IsKeyJustPressed(ebiten.KeyS) ||
		inpututil.IsKeyJustPressed(ebiten.KeyD) ||
		inpututil.IsKeyJustPressed(ebiten.KeyW) {
		t.sceneManager.TransitionTo(SceneGame)
		return nil
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		t.sceneManager.TransitionTo(SceneGame)
		return nil
	}

	return nil
}

func (t *TitleScene) Layout(outerWidth, outerHeight int) (int, int) {
	return outerWidth, outerHeight
}

func NewTitleScene(sm *SceneManager) *TitleScene {
	titleFontSource, _ := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	titleFont := &text.GoTextFace{
		Source: titleFontSource,
		Size:   48,
	}

	subtitleFontSource, _ := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	subtitleFont := &text.GoTextFace{
		Source: subtitleFontSource,
		Size:   24,
	}

	helpFontSource, _ := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	helpFont := &text.GoTextFace{
		Source: helpFontSource,
		Size:   12,
	}

	return &TitleScene{
		sceneManager: sm,
		titleFont:    titleFont,
		subtitleFont: subtitleFont,
		helpFont:     helpFont,
		showHelp:     false,
	}
}
