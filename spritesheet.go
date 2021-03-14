package main

import (
	"bananamap/level"
	"image/color"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type cursor struct {
	coords level.Coordinates
	img    *ebiten.Image
	op     *ebiten.DrawImageOptions
}

var (
	spriteSize        = 16
	spriteSheet       *ebiten.Image
	loadedSpriteSheet *ebiten.Image
	// selectedSpriteCoords level.Coordinates
	selection []cursor
)

func initSpriteSheet() {
	spriteSheet = ebiten.NewImage(spriteSheetWidth, spriteSheetHeight)

	newCursor := cursor{
		coords: level.Coordinates{
			X: 0,
			Y: 0,
		},
		img: ebiten.NewImage(spriteSize, spriteSize),
		op:  &ebiten.DrawImageOptions{},
	}

	selection = make([]cursor, 1)

	selection[0] = newCursor
}

func drawSpriteSheet(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(0), float64(canvasHeight-spriteSheetHeight))
	background.Fill(color.Black)
	screen.DrawImage(background, op)
	screen.DrawImage(loadedSpriteSheet, op)
	screen.DrawImage(spriteSheet, op)
	spriteSheet.Clear()
	for _, c := range selection {
		c.img.Fill(color.White)
		spriteSheet.DrawImage(c.img, c.op)
	}
}

func getSpriteIndex(mouseX, mouseY int) (int, int) {
	mapX := mouseX
	mapY := mouseY - (screenHeight - spriteSheetHeight)
	return mapX / (spriteSize), mapY / (spriteSize)
}

func spriteSheetClick() {
	x, y := getSpriteIndex(ebiten.CursorPosition())

	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		// add tiles to the selection multiple tiles to be added to the selection

		newCursor := cursor{
			coords: level.Coordinates{
				X: x,
				Y: y,
			},
			op:  &ebiten.DrawImageOptions{},
			img: ebiten.NewImage(spriteSize, spriteSize),
		}
		selection = append(selection, newCursor)

		for _, s := range selection {
			spriteIndicatorTranslate(s)
		}
		return
	}
	// reset selection to one element
	selection = selection[:1]

	selection[0].coords.X = x
	selection[0].coords.Y = y

	for _, s := range selection {
		spriteIndicatorTranslate(s)
	}
}

func spriteIndicatorTranslate(c cursor) {
	c.op.GeoM.Reset()
	c.op.GeoM.Translate(float64(c.coords.X*spriteSize), float64(c.coords.Y*spriteSize))
}
