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

	startX, startY   int
	offsetX, offsetY float64

	gridLines []Line
)

func initCanvas() {
	canvas = ebiten.NewImage(canvasWidth, canvasHeight)
	lvl = ebiten.NewImage(spriteSheetWidth, spriteSheetHeight)
	go func() {
		gridLines = initGrid()
	}()
}

func updateGrid() {
	mouseX, mouseY := ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
		offsetX = offsetX + (float64(mouseX) - float64(startX))
		offsetY = offsetY + (float64(mouseY) - float64(startY))
	}
	startX, startY = ebiten.CursorPosition()
}

func initGrid() []Line {
	var lines []Line
	for n := 0.0; n <= ngridLine; n++ {
		horizontal := NewLine(canvasWidth, 1, 0, n*tileSize)
		vertical := NewLine(1, canvasHeight, n*tileSize, 0)

		lines = append(lines, *horizontal)
		lines = append(lines, *vertical)
	}
	return lines
}

func drawGrid(lines []Line) {
	for _, g := range lines {
		g.Draw(canvas)
	}
}

func drawCanvas(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	canvas.Fill(color.NRGBA{0x00, 0x40, 0x80, 0xff})
	lvl.Fill(color.NRGBA{0x40, 0x40, 0x80, 0xff})
	canvas.DrawImage(lvl, &ebiten.DrawImageOptions{})
	drawGrid(gridLines)
	op.GeoM.Translate(offsetX, offsetY)
	screen.DrawImage(canvas, op)
}
