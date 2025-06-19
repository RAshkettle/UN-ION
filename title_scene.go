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
}

func (t *TitleScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 15, 25, 255})

	if t.showHelp {
		t.drawHelpOverlay(screen)
	} else {
		t.drawTitleScreen(screen)
	}
}

func (t *TitleScene) drawTitleScreen(screen *ebiten.Image) {
	// Get screen dimensions
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

	// Draw help prompt - VERY VISIBLE
	helpPrompt := "Press H for Help"
	helpPromptBounds, _ := text.Measure(helpPrompt, t.subtitleFont, 0) // Use subtitle font (larger)
	helpPromptX := (w - int(helpPromptBounds)) / 2
	helpPromptY := subtitleY + 50 // More space from subtitle

	op3 := &text.DrawOptions{}
	op3.GeoM.Translate(float64(helpPromptX), float64(helpPromptY))
	op3.ColorScale.ScaleWithColor(color.RGBA{255, 255, 100, 255}) // Bright yellow - very visible
	text.Draw(screen, helpPrompt, t.subtitleFont, op3)            // Use larger font

	// Draw controls help - positioned right under the help prompt
	controls := []string{
		"Quick Controls:",
		"WASD/Arrow Keys: Move piece",
		"Space: Rotate piece",
	}

	helpStartY := helpPromptY + 40 // Adjusted for new spacing
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

func (t *TitleScene) drawHelpOverlay(screen *ebiten.Image) {
	// Get screen dimensions
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()

	// Draw semi-transparent overlay background
	overlayImg := ebiten.NewImage(w, h)
	overlayImg.Fill(color.RGBA{5, 10, 20, 220})
	screen.DrawImage(overlayImg, &ebiten.DrawImageOptions{})

	// Title
	titleText := "HOW TO PLAY UN-ION"
	titleBounds, _ := text.Measure(titleText, t.titleFont, 0)
	titleX := (w - int(titleBounds)) / 2
	startY := 60

	titleOp := &text.DrawOptions{}
	titleOp.GeoM.Translate(float64(titleX), float64(startY))
	titleOp.ColorScale.ScaleWithColor(color.RGBA{255, 255, 100, 255})
	text.Draw(screen, titleText, t.titleFont, titleOp)

	// Help content sections
	sections := []struct {
		title string
		lines []string
	}{
		{
			title: "OBJECTIVE:",
			lines: []string{
				"Create horizontal lines where charges sum to zero",
				"Minimum 4 blocks: equal + and - blocks (e.g., ++--)",
				"Neutral blocks (○) have zero charge value and disrupt your chains",
			},
		},
		{
			title: "SCORING:",
			lines: []string{
				"4 blocks = 10 points",
				"5 blocks = 20 points, 6 blocks = 40 points (doubles each block)",
				"Chain reactions add to your total score!",
			},
		},
		{
			title: "STORM MECHANIC:",
			lines: []string{
				"4+ vertical same-charge blocks create electrical storms",
				"Storms spawn neutral blocks that disrupt your plans",
				"Break them quickly!",
			},
		},
		{
			title: "CONTROLS:",
			lines: []string{
				"WASD or Arrow Keys: Move piece",
				"Space: Rotate piece",
				"P: Pause game",
				"H: Toggle this help (from title screen)",
			},
		},
	}

	currentY := startY + 80
	lineHeight := 20
	sectionSpacing := 30

	for _, section := range sections {
		// Section title
		sectionOp := &text.DrawOptions{}
		sectionBounds, _ := text.Measure(section.title, t.subtitleFont, 0)
		sectionX := (w - int(sectionBounds)) / 2
		sectionOp.GeoM.Translate(float64(sectionX), float64(currentY))
		sectionOp.ColorScale.ScaleWithColor(color.RGBA{150, 255, 150, 255})
		text.Draw(screen, section.title, t.subtitleFont, sectionOp)
		currentY += lineHeight + 5

		// Section lines
		for _, line := range section.lines {
			lineOp := &text.DrawOptions{}
			lineBounds, _ := text.Measure(line, t.helpFont, 0)
			lineX := (w - int(lineBounds)) / 2
			lineOp.GeoM.Translate(float64(lineX), float64(currentY))
			lineOp.ColorScale.ScaleWithColor(color.RGBA{200, 200, 255, 255})
			text.Draw(screen, line, t.helpFont, lineOp)
			currentY += lineHeight
		}
		currentY += sectionSpacing
	}

	// Footer
	footerText := "Press H to close help and return to title screen"
	footerBounds, _ := text.Measure(footerText, t.helpFont, 0)
	footerX := (w - int(footerBounds)) / 2
	footerY := h - 40

	footerOp := &text.DrawOptions{}
	footerOp.GeoM.Translate(float64(footerX), float64(footerY))
	footerOp.ColorScale.ScaleWithColor(color.RGBA{255, 255, 100, 255})
	text.Draw(screen, footerText, t.helpFont, footerOp)
}

func (t *TitleScene) Update() error {
	// Handle help toggle
	if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		t.showHelp = !t.showHelp
		return nil
	}

	// If help is showing, only allow H to close it
	if t.showHelp {
		return nil
	}

	// Check for key presses to start game (only when help is not showing)
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

	// Check for mouse clicks to start game (only when help is not showing)
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
	// Create fonts
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
		Size:   12, // Reduced from 16 to 12
	}

	return &TitleScene{
		sceneManager: sm,
		titleFont:    titleFont,
		subtitleFont: subtitleFont,
		helpFont:     helpFont,
		showHelp:     false,
	}
}
