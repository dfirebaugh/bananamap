package main

import (
	"image/color"
	"log"

	"bananamap/level"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	tileSize  = 32
	ngridLine = 128
)

var (
	canvas *ebiten.Image

	startX, startY   int
	offsetX, offsetY float64

	isCollidable = false

	gridLines    []Line
	lvl          *level.Level
	currentLayer = 0
)

func initCanvas() {
	canvas = ebiten.NewImage(canvasWidth, canvasHeight)
	go func() {
		gridLines = initGrid()
	}()
	var err error
	loadedSpriteSheet, _, err = ebitenutil.NewImageFromFile("resources/images/tiles.png")
	if err != nil {
		log.Fatal(err)
	}

	lvl = level.Create(level.Source{
		Img:      loadedSpriteSheet,
		TileSize: 16,
	},
		canvasWidth, canvasHeight, tileSize, 1)
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

func updateGrid() {
	updateGridPan()
	canvasInputs()
}

func updateGridPan() {
	mouseX, mouseY := ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
		offsetX = offsetX + (float64(mouseX) - float64(startX))
		offsetY = offsetY + (float64(mouseY) - float64(startY))
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		offsetX = 0
		offsetY = 0
	}
	startX, startY = ebiten.CursorPosition()
}

func drawGrid(lines []Line) {
	for _, g := range lines {
		g.Draw(canvas)
	}
}

func drawCanvas(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	canvas.Fill(color.NRGBA{0x00, 0x40, 0x80, 0xff})
	lvl.Draw(canvas)
	drawGrid(gridLines)
	op.GeoM.Translate(offsetX, offsetY)
	screen.DrawImage(canvas, op)
}

func getTileIndex(mouseX, mouseY int) (int, int) {
	mapX := mouseX - int(offsetX)
	mapY := mouseY - int(offsetY)
	return mapX / tileSize, mapY / tileSize
}

func canvasInputs() {
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		selection[0].coords.X--
		spriteIndicatorTranslate(selection[0])
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		selection[0].coords.X++
		spriteIndicatorTranslate(selection[0])
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		selection[0].coords.Y--
		spriteIndicatorTranslate(selection[0])
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		selection[0].coords.Y++
		spriteIndicatorTranslate(selection[0])
	}

	// toggle
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mouseX, mouseY := ebiten.CursorPosition()

		// if the click is in the spritesheet / toolbox area
		if mouseY > screenHeight-spriteSheetHeight {
			spriteSheetClick()
			return
		}

		tileX, tileY := getTileIndex(mouseX, mouseY)

		coords := level.Coordinates{
			X: tileX,
			Y: tileY,
		}

		index := level.CoordsToIndex(coords, canvasWidth/tileSize, tileSize)
		if index >= (canvasWidth/tileSize)*(canvasHeight/tileSize) {
			return
		}
		if index < 0 {
			return
		}

		for _, s := range selection {
			lvl.UpdateTile(level.Coordinates{
				X: coords.X + s.coords.X - selection[0].coords.X,
				Y: coords.Y + s.coords.Y - selection[0].coords.Y,
			}, 0, s.coords, spriteSize, loadedSpriteSheet)
		}
		lvl.Export()
		return
	}

	mouseX, mouseY := ebiten.CursorPosition()

	tileX, tileY := getTileIndex(mouseX, mouseY)

	coords := level.Coordinates{
		X: tileX,
		Y: tileY,
	}

	// stream delete
	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		// if the click is in the spritesheet / toolbox area
		if mouseY > screenHeight-spriteSheetHeight {
			spriteSheetClick()
			return
		}
		if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			return
		}
	}

	if !ebiten.IsKeyPressed(ebiten.KeyShift) {
		return
	}

	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		return
	}

	// if the click is in the spritesheet / toolbox area
	if mouseY > screenHeight-spriteSheetHeight {
		spriteSheetClick()
		return
	}

	// stream draw
	lvl.UpdateTile(coords, 0, selection[0].coords, spriteSize, loadedSpriteSheet)
}
