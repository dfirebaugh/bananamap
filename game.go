package main

import (
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type Screen interface {
	UpdateGrid()
	DrawCanvas(*ebiten.Image)
	DrawSpriteSheet(*ebiten.Image)
	GetLevel() *Level
}

type Game struct {
	Canvas Screen
	Width  int
	Height int
}

func NewGame(screen Screen) *Game {
	return &Game{
		Canvas: screen,
		Width:  ScreenWidth,
		Height: ScreenHeight,
	}
}

func (g *Game) Update() error {
	g.Canvas.UpdateGrid()
	inputs(g)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Canvas.DrawCanvas(screen)
	g.Canvas.DrawSpriteSheet(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.Width, g.Height
}
