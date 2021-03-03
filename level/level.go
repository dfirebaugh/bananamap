// Package level Pass in the path to your png
//   file that was created with bananamap and
//   the `Load` function will return
package level

// Tile is a single tile in the level.
type Tile struct {
	// Index references where the tile exists
	// in the source PNG
	Index        int
	isCollidable bool
}

// Level represents a map level that
//  that has been configured
type Level struct {
	Height uint
	Width  uint
	Layers [][]Tile
}
