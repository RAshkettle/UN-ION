package main

import (
	"bytes"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

type HelpScene struct {
	titleFont    *text.GoTextFace
	subtitleFont *text.GoTextFace
	helpFont     *text.GoTextFace
	sceneManager *SceneManager
	prevHPressed bool
}

func NewHelpScene(sm *SceneManager) *HelpScene {
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

	return &HelpScene{
		titleFont:    titleFont,
		subtitleFont: subtitleFont,
		helpFont:     helpFont,
		sceneManager: sm,
	}
}

func (h *HelpScene) Draw(screen *ebiten.Image) {
	w, hgt := screen.Bounds().Dx(), screen.Bounds().Dy()

	overlayImg := ebiten.NewImage(w, hgt)
	overlayImg.Fill(color.RGBA{5, 10, 20, 220})
	screen.DrawImage(overlayImg, &ebiten.DrawImageOptions{})

	titleText := "HOW TO PLAY UN-ION"
	titleBounds, _ := text.Measure(titleText, h.titleFont, 0)
	titleX := (w - int(titleBounds)) / 2
	startY := 60

	titleOp := &text.DrawOptions{}
	titleOp.GeoM.Translate(float64(titleX), float64(startY))
	titleOp.ColorScale.ScaleWithColor(color.RGBA{255, 255, 100, 255})
	text.Draw(screen, titleText, h.titleFont, titleOp)

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

		sectionOp := &text.DrawOptions{}
		sectionBounds, _ := text.Measure(section.title, h.subtitleFont, 0)
		sectionX := (w - int(sectionBounds)) / 2
		sectionOp.GeoM.Translate(float64(sectionX), float64(currentY))
		sectionOp.ColorScale.ScaleWithColor(color.RGBA{150, 255, 150, 255})
		text.Draw(screen, section.title, h.subtitleFont, sectionOp)
		currentY += lineHeight + 5

		for _, line := range section.lines {
			lineOp := &text.DrawOptions{}
			lineBounds, _ := text.Measure(line, h.helpFont, 0)
			lineX := (w - int(lineBounds)) / 2
			lineOp.GeoM.Translate(float64(lineX), float64(currentY))
			lineOp.ColorScale.ScaleWithColor(color.RGBA{200, 200, 255, 255})
			text.Draw(screen, line, h.helpFont, lineOp)
			currentY += lineHeight
		}
		currentY += sectionSpacing
	}

	footerText := "Press H to close help and return to title screen"
	footerBounds, _ := text.Measure(footerText, h.helpFont, 0)
	footerX := (w - int(footerBounds)) / 2
	footerY := hgt - 40

	footerOp := &text.DrawOptions{}
	footerOp.GeoM.Translate(float64(footerX), float64(footerY))
	footerOp.ColorScale.ScaleWithColor(color.RGBA{255, 255, 100, 255})
	text.Draw(screen, footerText, h.helpFont, footerOp)
}

func (h *HelpScene) Update() error {
	hPressed := ebiten.IsKeyPressed(ebiten.KeyH)
	if hPressed && !h.prevHPressed {
		h.sceneManager.TransitionTo(SceneTitleScreen)
	}
	h.prevHPressed = hPressed
	return nil
}

func (h *HelpScene) Layout(outerWidth, outerHeight int) (int, int) {
	return outerWidth, outerHeight
}
