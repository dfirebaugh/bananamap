package main

import (
	"image"
	"image/color"

	ebiten "github.com/hajimehoshi/ebiten/v2"
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

	gridLines  []Line
	worldTiles map[coordinates]*ebiten.Image
)

type coordinates struct {
	x int
	y int
}

func initCanvas() {
	canvas = ebiten.NewImage(canvasWidth, canvasHeight)
	go func() {
		gridLines = initGrid()
	}()

	worldTiles = make(map[coordinates]*ebiten.Image)
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
	canvasClick()
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
	drawMap(canvas)
	drawGrid(gridLines)
	op.GeoM.Translate(offsetX, offsetY)
	screen.DrawImage(canvas, op)
}

func drawMap(canvas *ebiten.Image) {
	for coords, tile := range worldTiles {
		op := &ebiten.DrawImageOptions{}
		// tile.Fill(color.White)
		op.GeoM.Translate(float64(coords.x), float64(coords.y))
		canvas.DrawImage(tile, op)
	}
}

func getTileIndex(mouseX, mouseY int) (int, int) {
	mapX := mouseX - int(offsetX)
	mapY := mouseY - int(offsetY)
	return mapX / tileSize, mapY / tileSize
}

func canvasClick() {

	// toggle
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mouseX, mouseY := ebiten.CursorPosition()

		// if the click is in the spritesheet / toolbox area
		if mouseY > screenHeight-spriteSheetHeight {
			spriteSheetClick()
			return
		}

		tileX, tileY := getTileIndex(mouseX, mouseY)

		coords := coordinates{
			x: tileX * tileSize,
			y: tileY * tileSize,
		}

		if worldTiles[coords] != nil {
			delete(worldTiles, coords)
			return
		}

		// worldTiles[coords] = ebiten.NewImage(tileSize, tileSize)
		worldTiles[coords] = ebiten.NewImageFromImage(spriteSheet.SubImage(image.Rectangle{
			image.Point{X: selectedSpriteCoords.x, Y: selectedSpriteCoords.y},
			image.Point{X: selectedSpriteCoords.x + 16, Y: selectedSpriteCoords.y + 16},
		}))
		return
	}

	mouseX, mouseY := ebiten.CursorPosition()

	tileX, tileY := getTileIndex(mouseX, mouseY)

	coords := coordinates{
		x: tileX * tileSize,
		y: tileY * tileSize,
	}

	// stream delete
	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			return
		}
		delete(worldTiles, coords)
		return
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
	// worldTiles[coords] = ebiten.NewImage(tileSize, tileSize)
	worldTiles[coords] = ebiten.NewImageFromImage(spriteSheet.SubImage(image.Rectangle{
		image.Point{X: selectedSpriteCoords.x, Y: selectedSpriteCoords.y},
		image.Point{X: selectedSpriteCoords.x + 16, Y: selectedSpriteCoords.y + 16},
	}))
}
