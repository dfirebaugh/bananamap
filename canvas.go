package main

import (
	"image/color"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

const (
	tileSize  = 32
	ngridLine = 128
)

var (
	canvas *ebiten.Image
	lvl    *ebiten.Image

	startX, startY int

	offsetX, offsetY float64
)

func initCanvas() {
	canvas = ebiten.NewImage(canvasWidth, canvasHeight)
	lvl = ebiten.NewImage(spriteSheetWidth, spriteSheetHeight)
}

func updateGrid(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
		mouseX, mouseY := ebiten.CursorPosition()
		offsetX = float64(mouseX) - float64(startX)
		offsetY = float64(mouseY) - float64(startY)
	} else {
		startX, startY = ebiten.CursorPosition()
	}

	for n := 0.0; n <= ngridLine; n++ {
		horizontal := NewLine(canvasWidth, 1, 0, n*tileSize)
		horizontal.Draw(canvas)
		vertical := NewLine(1, canvasHeight, n*tileSize, 0)
		vertical.Draw(canvas)
	}

	op.GeoM.Translate(offsetX, offsetY)
	screen.DrawImage(canvas, op)
}

func drawCanvas(screen *ebiten.Image) {
	canvas.Fill(color.NRGBA{0x00, 0x40, 0x80, 0xff})
	lvl.Fill(color.NRGBA{0x40, 0x40, 0x80, 0xff})
	canvas.DrawImage(lvl, &ebiten.DrawImageOptions{})
	updateGrid(screen)
}
