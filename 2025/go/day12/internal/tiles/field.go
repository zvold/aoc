package tiles

import (
	"fmt"
	"math/bits"
)

// Represents a mutable rectangular field of specific dimensions.
type Field struct {
	data   []uint64
	width  byte
	height byte
	text   string // Cached string representation, cleared on mutations.
}

func (f *Field) String() string {
	if len(f.text) != 0 {
		return f.text
	}

	for j := range f.height {
		b := make([]byte, f.width+1)
		b[f.width] = '\n'
		for i := range f.width {
			if f.data[j]&(1<<i) != 0 {
				b[i] = 'o'
			} else {
				b[i] = '.'
			}
		}
		f.text += string(b)
	}
	return f.text
}

func (f *Field) Width() int {
	return int(f.width)
}

func (f *Field) Height() int {
	return int(f.height)
}

// Returns how many cells are set.
func (f *Field) Count() int {
	r := 0
	for _, row := range f.data {
		r += bits.OnesCount64(row)
	}
	return r
}

func CreateField(w, h int) *Field {
	if w < 0 || h < 0 || w > 64 || h > 64 {
		panic("Unsupported field size.")
	}
	return &Field{
		data:   make([]uint64, h),
		width:  byte(w),
		height: byte(h),
		text:   "",
	}
}

// Mutates the field by placing the tile at (x, y).
func (f *Field) Place(x, y int, tile *Tile) {
	if !f.Fits(x, y, tile) {
		panic(fmt.Sprintf("Tile doesn't fit at position %d, %d:\n%s\n%s", x, y, tile, f))
	}
	for j := range 3 {
		f.data[y+j] ^= (uint64(tile.mask[j]) << x)
	}
	f.text = ""
}

// Mutates the field by removing the tile at (x, y).
func (f *Field) Remove(x, y int, tile *Tile) {
	if !f.Present(x, y, tile) {
		panic(fmt.Sprintf("Tile not present, cannot be removed from %d, %d:\n%s\n%s", x, y, tile, f))
	}
	for j := range 3 {
		f.data[y+j] ^= (uint64(tile.mask[j]) << x)
	}
	f.text = ""
}

// Returns true if tile fits at position (x, y).
func (f *Field) Fits(x, y int, tile *Tile) bool {
	// This assumes any 3×3 tile a has 3×3 bounding box.
	// This is true for tiles we're interested in.
	if x < 0 || x >= f.Width()-2 || y < 0 || y >= f.Height()-2 {
		return false
	}
	return f.overlap(x, y, tile) == 0
}

// Returns true if the tile is present at position (x, y).
func (f *Field) Present(x, y int, tile *Tile) bool {
	// This assumes any 3×3 tile a has 3×3 bounding box.
	// This is true for tiles we're interested in.
	if x < 0 || x >= f.Width()-2 || y < 0 || y >= f.Height()-2 {
		return false
	}
	return f.overlap(x, y, tile) == tile.Count()
}

// Returns count of overlapping cells if the tile is above (x, y).
func (f *Field) overlap(x, y int, tile *Tile) int {
	c := 0
	for j := range 3 {
		c += bits.OnesCount64(f.data[y+j] & (uint64(tile.mask[j]) << x))
	}
	return c
}
