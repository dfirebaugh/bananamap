package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

// Line is used as a fake line primitive
type Line struct {
	img  *ebiten.Image
	geoM *ebiten.GeoM
}

// NewLine makes a new line for use.
func NewLine(width, height int, x, y float64) *Line {
	img := ebiten.NewImage(width, height)
	img.Fill(colornames.Black)
	geoM := &ebiten.GeoM{}
	geoM.Translate(x, y)
	return &Line{img: img, geoM: geoM}
}

func NewLinex(startX, startY int, endX, endY float64) *Line {
	img := ebiten.NewImage(int(endX), int(endY))
	img.Fill(colornames.Black)
	geoM := &ebiten.GeoM{}
	geoM.Translate(float64(startX), float64(startY))
	return &Line{img: img, geoM: geoM}
}

// Draw renders the line to the screen.
func (l *Line) Draw(targetImage *ebiten.Image) {
	targetImage.DrawImage(l.img, &ebiten.DrawImageOptions{GeoM: *l.geoM})
}
