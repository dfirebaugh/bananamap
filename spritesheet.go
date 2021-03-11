package main

import (
	"image/color"
	"log"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	spriteSheet          *ebiten.Image
	loadedSpriteSheet    *ebiten.Image
	spriteCursor         *ebiten.Image
	spriteCursorOp       *ebiten.DrawImageOptions
	selectedSpriteCoords coordinates
)

func initSpriteSheet() {
	var err error
	spriteCursorOp = &ebiten.DrawImageOptions{}
	spriteCursor = ebiten.NewImage(16, 16)
	spriteSheet = ebiten.NewImage(spriteSheetWidth, spriteSheetHeight)
	loadedSpriteSheet, _, err = ebitenutil.NewImageFromFile("resources/images/tiles.png")
	if err != nil {
		log.Fatal(err)
	}
}

func drawSpriteSheet(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(0), float64(canvasHeight-spriteSheetHeight))
	background.Fill(color.Black)
	screen.DrawImage(background, op)
	screen.DrawImage(loadedSpriteSheet, op)
	screen.DrawImage(spriteSheet, op)
	spriteSheet.Clear()
	spriteCursor.Fill(color.White)
	spriteSheet.DrawImage(spriteCursor, spriteCursorOp)
}

func getSpriteIndex(mouseX, mouseY int) (int, int) {
	mapX := mouseX
	mapY := mouseY - (screenHeight - spriteSheetHeight)
	return mapX / (16), mapY / (16)
}

func spriteSheetClick() {
	x, y := getSpriteIndex(ebiten.CursorPosition())
	selectedSpriteCoords.x = x
	selectedSpriteCoords.y = y

	// fmt.Println(iX, iY)
	spriteCursorOp.GeoM.Reset()
	spriteCursorOp.GeoM.Translate(float64(selectedSpriteCoords.x*16), float64(selectedSpriteCoords.y*16))
}
