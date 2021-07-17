package main

import (
	"image/color"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type spriteSheet struct {
	Height            int
	Width             int
	spriteSize        int
	canvasHeight      int
	canvasWidth       int
	screenHeight      int
	selection         []Cursor
	background        *ebiten.Image
	spriteSheet       *ebiten.Image
	loadedSpriteSheet *ebiten.Image
}

func NewSpriteSheet(loadedSpriteSheet *ebiten.Image) *spriteSheet {
	ss := &spriteSheet{
		spriteSize:        SpriteSize,
		loadedSpriteSheet: loadedSpriteSheet,
		spriteSheet:       ebiten.NewImage(SpriteSheetWidth, SpriteSheetWidth),
		background:        ebiten.NewImage(ScreenWidth, ScreenHeight),
		canvasHeight:      CanvasHeight,
		canvasWidth:       CanvasWidth,
		screenHeight:      ScreenHeight,
		Height:            SpriteSheetHeight,
		Width:             SpriteSheetWidth,
	}

	newCursor := Cursor{
		Coords: Coordinates{
			X: 0,
			Y: 0,
		},
		IMG: ebiten.NewImage(SpriteSize, SpriteSize),
		OP:  &ebiten.DrawImageOptions{},
	}

	ss.selection = make([]Cursor, 1)

	ss.selection[0] = newCursor
	return ss
}

func (ss *spriteSheet) GetSelection() []Cursor {
	return ss.selection
}

func (ss *spriteSheet) SetSelection(selection []Cursor) {
	ss.selection = selection
}

func (ss *spriteSheet) DrawSpriteSheet(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(0), float64(ss.canvasHeight-ss.Height))
	ss.background.Fill(color.Black)
	screen.DrawImage(ss.background, op)
	screen.DrawImage(ss.loadedSpriteSheet, op)
	screen.DrawImage(ss.spriteSheet, op)
	ss.spriteSheet.Clear()
	for _, c := range ss.selection {
		c.IMG.Fill(color.White)
		ss.spriteSheet.DrawImage(c.IMG, c.OP)
	}
}

func (ss *spriteSheet) getSpriteIndex(mouseX, mouseY int) (int, int) {
	mapX := mouseX
	mapY := mouseY - (ss.screenHeight - ss.Height)
	return mapX / (ss.spriteSize), mapY / (ss.spriteSize)
}

func (ss *spriteSheet) SpriteSheetClick() {
	x, y := ss.getSpriteIndex(ebiten.CursorPosition())

	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		// add tiles to the selection multiple tiles to be added to the selection

		newCursor := Cursor{
			Coords: Coordinates{
				X: x,
				Y: y,
			},
			OP:  &ebiten.DrawImageOptions{},
			IMG: ebiten.NewImage(ss.spriteSize, ss.spriteSize),
		}
		ss.selection = append(ss.selection, newCursor)

		for _, s := range ss.selection {
			ss.SpriteIndicatorTranslate(s)
		}
		return
	}
	// reset selection to one element
	ss.selection = ss.selection[:1]

	ss.selection[0].Coords.X = x
	ss.selection[0].Coords.Y = y

	for _, s := range ss.selection {
		ss.SpriteIndicatorTranslate(s)
	}
}

func (ss *spriteSheet) SpriteIndicatorTranslate(c Cursor) {
	c.OP.GeoM.Reset()
	c.OP.GeoM.Translate(float64(c.Coords.X*ss.spriteSize), float64(c.Coords.Y*ss.spriteSize))
}
