package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type SpriteSheet interface {
	SetSelection([]Cursor)
	GetSelection() []Cursor
	SpriteIndicatorTranslate(Cursor)
	SpriteSheetClick()
	DrawSpriteSheet(*ebiten.Image)
}

type gameScreen struct {
	Width             int
	Height            int
	spriteSize        int
	spriteSheetHeight int
	spriteSheetWidth  int
	loadedSpriteSheet *ebiten.Image
	canvas            *ebiten.Image
	spriteSheet       SpriteSheet
	startX, startY    int
	offsetX, offsetY  float64
	isCollidable      bool
	gridLines         []Line
	lvl               *Level
	currentLayer      int
	tileSize          int
	ngridLine         int
}

func NewScreen(spriteSheet SpriteSheet, loadedSpriteSheet *ebiten.Image) *gameScreen {
	s := &gameScreen{
		spriteSheet:       spriteSheet,
		loadedSpriteSheet: loadedSpriteSheet,
		isCollidable:      false,
		currentLayer:      0,
		tileSize:          32,
		ngridLine:         128,
		offsetX:           0,
		offsetY:           0,
		Height:            CanvasHeight,
		Width:             CanvasWidth,
		spriteSheetHeight: SpriteSheetHeight,
		spriteSheetWidth:  SpriteSheetWidth,
		spriteSize:        SpriteSize,
	}
	s.canvas = ebiten.NewImage(s.Width, s.Height)
	go func() {
		s.gridLines = s.initGrid()
	}()

	s.lvl = Create(Source{
		Img:      s.loadedSpriteSheet,
		TileSize: 16,
	},
		s.Width, s.Height, s.tileSize, 1)
	return s
}

func (s *gameScreen) DrawSpriteSheet(screen *ebiten.Image) {
	s.spriteSheet.DrawSpriteSheet(screen)
}

func (s *gameScreen) initGrid() []Line {
	var lines []Line
	for n := 0.0; int(n) <= s.ngridLine; n++ {
		horizontal := NewLine(s.Width, 1, 0, float64(int(n)*s.tileSize))
		vertical := NewLine(1, s.Height, float64(int(n)*s.tileSize), 0)

		lines = append(lines, *horizontal)
		lines = append(lines, *vertical)
	}
	return lines
}

func (s *gameScreen) UpdateGrid() {
	s.updateGridPan()
	s.canvasInputs()
}

func (s *gameScreen) updateGridPan() {
	mouseX, mouseY := ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
		s.offsetX = s.offsetX + (float64(mouseX) - float64(s.startX))
		s.offsetY = s.offsetY + (float64(mouseY) - float64(s.startY))
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		s.offsetX = 0
		s.offsetY = 0
	}
	s.startX, s.startY = ebiten.CursorPosition()
}

func (s *gameScreen) drawGrid(lines []Line) {
	for _, g := range lines {
		g.Draw(s.canvas)
	}
}

func (s *gameScreen) DrawCanvas(screen *ebiten.Image) {
	// Get the x, y position of the cursor from the CursorPosition() function
	x, y := ebiten.CursorPosition()

	// Display the information with "X: xx, Y: xx" format
	ebitenutil.DebugPrint(screen, fmt.Sprintf("X: %d, Y: %d", x, y))
	s.drawCanvas(screen)
}

func (s *gameScreen) drawCanvas(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	s.canvas.Fill(color.NRGBA{0x00, 0x40, 0x80, 0xff})
	s.lvl.Draw(s.canvas)
	s.drawGrid(s.gridLines)
	op.GeoM.Translate(s.offsetX, s.offsetY)
	screen.DrawImage(s.canvas, op)
}

func (s *gameScreen) getTileIndex(mouseX, mouseY int) (int, int) {
	mapX := mouseX - int(s.offsetX)
	mapY := mouseY - int(s.offsetY)
	return mapX / s.tileSize, mapY / s.tileSize
}

func (s *gameScreen) canvasInputs() {
	selection := s.spriteSheet.GetSelection()
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		selection[0].Coords.X--
		s.spriteSheet.SpriteIndicatorTranslate(selection[0])
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		selection[0].Coords.X++
		s.spriteSheet.SpriteIndicatorTranslate(selection[0])
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		selection[0].Coords.Y--
		s.spriteSheet.SpriteIndicatorTranslate(selection[0])
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		selection[0].Coords.Y++
		s.spriteSheet.SpriteIndicatorTranslate(selection[0])
	}
	s.spriteSheet.SetSelection(selection)

	// toggle
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mouseX, mouseY := ebiten.CursorPosition()
		// if the click is in the spritesheet / toolbox area
		if mouseY > s.Height-s.spriteSheetHeight {
			s.spriteSheet.SpriteSheetClick()
			return
		}

		tileX, tileY := s.getTileIndex(mouseX, mouseY)

		coords := Coordinates{
			X: tileX,
			Y: tileY,
		}

		index := CoordsToIndex(coords, s.Width/s.tileSize, s.tileSize)
		if index >= (s.Width/s.tileSize)*(s.Height/s.tileSize) {
			return
		}
		if index < 0 {
			return
		}

		for _, sel := range selection {
			s.lvl.UpdateTile(Coordinates{
				X: coords.X + sel.Coords.X - selection[0].Coords.X,
				Y: coords.Y + sel.Coords.Y - selection[0].Coords.Y,
			}, 0, sel.Coords, s.spriteSize, s.loadedSpriteSheet)
		}
		s.lvl.ExportJSON()
		return
	}

	mouseX, mouseY := ebiten.CursorPosition()

	tileX, tileY := s.getTileIndex(mouseX, mouseY)

	coords := Coordinates{
		X: tileX,
		Y: tileY,
	}

	// stream delete
	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		// if the click is in the spritesheet / toolbox area
		if mouseY > s.Height-s.spriteSheetHeight {
			s.spriteSheet.SpriteSheetClick()

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
	if mouseY > s.Height-s.spriteSheetHeight {
		s.spriteSheet.SpriteSheetClick()
		return
	}

	// stream draw
	s.lvl.UpdateTile(coords, 0, selection[0].Coords, s.spriteSize, s.loadedSpriteSheet)
}

func (s *gameScreen) GetLevel() *Level {
	return s.lvl
}
