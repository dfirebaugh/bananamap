package main

import (
	"fmt"
	"image/color"
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
	spriteSheet *ebiten.Image
	background  *ebiten.Image
)

type Game struct {
}

func init() {
	initCanvas()
	var err error
	spriteSheet, _, err = ebitenutil.NewImageFromFile("resources/images/tiles.png")
	if err != nil {
		log.Fatal(err)
	}
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

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(0), float64(canvasHeight-spriteSheetHeight))
	background.Fill(color.Black)
	screen.DrawImage(background, op)
	screen.DrawImage(spriteSheet, op)
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
