package main

import (
	"fmt"
	_ "image/png"
	"log"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth       = 640
	screenHeight      = 640
	spriteSheetWidth  = 400
	spriteSheetHeight = 224
	canvasHeight      = screenHeight
	canvasWidth       = screenWidth * 5
)

var (
	background *ebiten.Image
)

type Game struct {
}

func init() {
	initCanvas()
	initSpriteSheet()

	var err error

	background = ebiten.NewImage(screenWidth, canvasHeight)
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error {
	updateGrid()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Get the x, y position of the cursor from the CursorPosition() function
	x, y := ebiten.CursorPosition()

	// Display the information with "X: xx, Y: xx" format
	ebitenutil.DebugPrint(screen, fmt.Sprintf("X: %d, Y: %d", x, y))
	drawCanvas(screen)
	drawSpriteSheet(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("bananamap")
	ebiten.SetWindowResizable(true)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
