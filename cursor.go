package main

import (
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type Cursor struct {
	Coords Coordinates
	IMG    *ebiten.Image
	OP     *ebiten.DrawImageOptions
}
