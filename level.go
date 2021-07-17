package main

import (
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"

	"github.com/hajimehoshi/ebiten/v2"
)

type Coordinates struct {
	X     int
	Y     int
	Index int
}

// Tile is a single tile in the level.
type Tile struct {
	// Index references where the tile exists
	// in the source PNG
	Index        int  `json:"index"`
	IsCollidable bool `json:"isCollidable"`
}

type Layer struct {
	tileImages map[Coordinates]*ebiten.Image
	Tiles      []Tile
}

type Source struct {
	Img      *ebiten.Image
	TileSize int
}

// Level represents a map level that
//  that has been configured
type Level struct {
	Height      uint    `json:"height"`
	Width       uint    `json:"width"`
	TileSize    uint    `json:"tileSize"`
	Layers      []Layer `json:"layers"`
	sourceImage Source
}

func Create(sourceImg Source, width, height, tileSize, numLayers int) *Level {
	l := new(Level)
	l.Width = uint(width)
	l.Height = uint(height)
	l.TileSize = uint(tileSize)
	l.Layers = make([]Layer, numLayers)
	l.sourceImage = sourceImg

	for i := 0; i < numLayers; i++ {
		l.Layers[i].Tiles = make([]Tile, (width/tileSize)*(height/tileSize))
		l.Layers[i].tileImages = make(map[Coordinates]*ebiten.Image)
	}
	return l
}

func (l Level) Draw(canvas *ebiten.Image) {
	for coords, tile := range l.Layers[0].tileImages {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(2, 2)
		op.GeoM.Translate(float64(coords.X*int(l.TileSize)), float64(coords.Y*int(l.TileSize)))
		canvas.DrawImage(tile, op)
	}
}

func CoordsToIndex(coords Coordinates, tilesInRow, tSize int) int {
	return (coords.X + (coords.Y * tilesInRow)) / tSize
}

func (l Level) UpdateTile(coords Coordinates, layer int, selectionCoords Coordinates, size int, sourceImage *ebiten.Image) {
	// map the index in the Level to the subImage index in the PNG
	indexInPNG := CoordsToIndex(selectionCoords, l.sourceImage.Img.Bounds().Max.X/l.sourceImage.TileSize, 1)
	indexInLevel := CoordsToIndex(coords, int(l.Width)/int(l.TileSize), int(l.TileSize))

	if indexInLevel > 2000 {
		println(indexInLevel)
		return
	}

	if indexInLevel < 0 {
		println(indexInLevel)
		return
	}

	l.Layers[layer].Tiles[indexInLevel].Index = indexInPNG

	l.Layers[layer].tileImages[coords] = ebiten.NewImageFromImage(sourceImage.SubImage(image.Rectangle{
		image.Point{X: selectionCoords.X * size, Y: selectionCoords.Y * size},
		image.Point{X: (selectionCoords.X * size) + size, Y: (selectionCoords.Y * size) + size},
	}))
}

func (l Level) Export() {
	b, err := json.MarshalIndent(l, "", " ")
	if err != nil {
		fmt.Println(err)
		return
	}

	_ = ioutil.WriteFile("test.json", b, 0644)
}

// JSONToTileMap returns an int slice of tiles.  This represents which subimages from the spritesheet
//  to represent and where they exist in the level.
//  provide the JSON and which layer to return
func JSONToTileMap(b []byte, layer int) ([]int, error) {
	var tileMap []int
	var lvlJSON Level
	err := json.Unmarshal(b, &lvlJSON)
	if err != nil {
		fmt.Println(err)
		return tileMap, err
	}

	for _, t := range lvlJSON.Layers[layer].Tiles {
		tileMap = append(tileMap, t.Index)
	}

	return tileMap, nil
}

// JSONToCollideMap returns an bool slice of collidable tiles.  This represents which tiles have collision enabled.
//  provide the JSON and which layer to return
func JSONToCollideMap(b []byte, layer int) ([]bool, error) {
	var collideMap []bool
	var lvlJSON Level
	err := json.Unmarshal(b, &lvlJSON)
	if err != nil {
		fmt.Println(err)
		return collideMap, err
	}

	for _, t := range lvlJSON.Layers[layer].Tiles {
		collideMap = append(collideMap, t.IsCollidable)
	}

	return collideMap, nil
}
