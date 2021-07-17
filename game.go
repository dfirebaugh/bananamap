package main

import (
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type Screen interface {
	UpdateGrid()
	DrawCanvas(*ebiten.Image)
	DrawSpriteSheet(*ebiten.Image)
}

type Game struct {
	screen Screen
	Width  int
	Height int
}

func NewGame(screen Screen) *Game {
	return &Game{
		screen: screen,
		Width:  ScreenWidth,
		Height: ScreenHeight,
	}
}

func (g *Game) Update() error {
	g.screen.UpdateGrid()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.screen.DrawCanvas(screen)
	g.screen.DrawSpriteSheet(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.Width, g.Height
}
