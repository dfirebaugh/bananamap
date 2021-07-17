package main

import (
	_ "image/png"
	"log"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func main() {
	loadedSpriteSheet, _, err := ebitenutil.NewImageFromFile("resources/images/tiles.png")
	if err != nil {
		log.Fatal(err)
	}

	spriteSheet := NewSpriteSheet(loadedSpriteSheet)
	gameScreen := NewScreen(spriteSheet, loadedSpriteSheet)
	game := NewGame(gameScreen)

	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("bananamap")
	ebiten.SetWindowResizable(true)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
