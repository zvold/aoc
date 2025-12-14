package tiles

import "math/bits"

// Describes a 3x3 tile.
type Tile struct {
	mask [3]uint8
	text string
}

func (t *Tile) String() string {
	return t.text
}

func (t *Tile) Equals(t2 *Tile) bool {
	return t.mask == t2.mask
}

func (t *Tile) Count() int {
	r := 0
	for _, row := range t.mask {
		r += bits.OnesCount8(row)
	}
	return r
}

// Expects a string in the form "..#..##.#".
func CreateTile(s string) *Tile {
	if len(s) > 9 {
		panic("Tiles that are not 3×3 are not supported.")
	}
	r := Tile{}
	for i := range 9 {
		if s[i] == '#' {
			r.mask[i/3] |= 1 << (i % 3)
			r.text += "o"
		} else {
			r.text += "."
		}
		if i%3 == 2 {
			r.text += "\n"
		}
	}
	return &r
}

// Returns a tile rotated 90° clockwise.
func (t *Tile) Rotate() *Tile {
	b := []byte(".........")
	for i := range 3 {
		for j := range 3 {
			if t.mask[j]&(1<<i) != 0 {
				b[3*i+(2-j)] = '#'
			}
		}
	}
	return CreateTile(string(b))
}

// Returns a tile mirrored around vertical axis.
func (t *Tile) Mirror() *Tile {
	b := []byte(".........")
	for i := range 3 {
		for j := range 3 {
			if t.mask[j]&(1<<i) != 0 {
				b[3*j+(2-i)] = '#'
			}
		}
	}
	return CreateTile(string(b))
}
