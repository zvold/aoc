package tiles

import (
	"fmt"
	"slices"
)

// Represents all rotations (including mirrors) of a 3×3 tile.
type MetaTile struct {
	tiles []*Tile
}

func (m *MetaTile) String() string {
	r := ""
	for j := range 3 {
		r += "|"
		for _, t := range m.tiles {
			r += fmt.Sprintf(" %s |", t.String()[4*j:4*j+3])
		}
		r += "\n"
	}
	return r
}

func (m *MetaTile) Size() int {
	return len(m.tiles)
}

func (m *MetaTile) Get(i int) *Tile {
	return m.tiles[i]
}

func CreateMetaTile(tile *Tile) *MetaTile {
	variants := make([]*Tile, 0)
	for range 4 {
		variants = append(variants, tile)
		variants = append(variants, tile.Mirror())
		tile = tile.Rotate()
	}

	metatile := MetaTile{tiles: make([]*Tile, 0)}
	for _, v := range variants {
		if slices.IndexFunc(metatile.tiles, func(t *Tile) bool {
			return t.Equals(v)
		}) == -1 {
			metatile.tiles = append(metatile.tiles, v)
		}
	}

	return &metatile
}
