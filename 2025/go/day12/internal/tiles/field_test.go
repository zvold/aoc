package tiles

import (
	"fmt"
	"math/rand/v2"
	"testing"
)

func Example_createField_1() {
	field := CreateField(7, 3)
	fmt.Println(field.String())

	// Output:
	// .......
	// .......
	// .......
}

func Test_fits_1(t *testing.T) {
	field := CreateField(7, 3)
	tile := createRandomTile()

	if field.Fits(-1, 0, tile) {
		t.Fatalf("Boundary check X test failed.")
	}
	if field.Fits(0, -1, tile) {
		t.Fatalf("Boundary check Y test failed.")
	}
	if field.Fits(5, 0, tile) {
		t.Fatalf("Boundary check X test failed.")
	}
	if field.Fits(0, 1, tile) {
		t.Fatalf("Boundary check Y test failed.")
	}
}

func Example_place_1() {
	field := CreateField(4, 6)
	tile := CreateTile("##..##..#")
	field.Place(1, 1, tile)
	field.Place(0, 2, tile)
	fmt.Println(field.String())

	// Output:
	// ....
	// .oo.
	// oooo
	// .ooo
	// ..o.
	// ....
}

func Test_fits_2(t *testing.T) {
	field := CreateField(4, 5)
	tile := CreateTile("........#")
	field.Place(1, 0, tile)

	tile = tile.Mirror().Rotate().Rotate()

	if field.Fits(1, 2, tile) {
		t.Fatalf("Field shouldn't fit the tile at (1, 2):\n%s\n%s", tile, field)
	}
}

func Test_place_2(t *testing.T) {
	for range 100 {
		field := CreateField(20, 30)
		for range 1000 { // Place 1000 random tiles at random positions.
			tile := createRandomTile()
			x, y := rand.IntN(field.Width()), rand.IntN(field.Height())
			if field.Fits(x, y, tile) {
				c := field.Count()
				field.Place(x, y, tile)
				if field.Count() != c+tile.Count() {
					t.Fatalf("Wrong cell count after placing\n%s\n%s", tile, field)
				}
			}
		}
	}
}

func Test_remove_1(t *testing.T) {
	type op struct {
		x, y int
		tile *Tile
	}

	for range 100 {
		field := CreateField(30, 20)
		if field.Count() != 0 {
			t.Fatalf("Just created field is not empty\n%s", field)
		}
		history := make([]op, 0)
		for range 1000 { // Place 1000 random tiles at random positions.
			tile := createRandomTile()
			x, y := rand.IntN(field.Width()), rand.IntN(field.Height())

			if field.Fits(x, y, tile) {
				field.Place(x, y, tile)
				history = append(history, op{x: x, y: y, tile: tile}) // Remember all placed tiles.
				if !field.Present(x, y, tile) {
					t.Fatalf("The tile was just placed, should be present\n%s\n%s", tile, field)
				}
			}
		}
		// Now remove all tiles.
		for _, h := range history {
			field.Remove(h.x, h.y, h.tile)
		}
		if field.Count() != 0 {
			t.Fatalf("Non-empty field after removing all tiles\n%s", field)
		}
	}
}

func Test_remove_2(t *testing.T) {
	field := CreateField(15, 25)
	for range 1000 {
		tile := createRandomTile()
		x, y := rand.IntN(field.Width()-2), rand.IntN(field.Height()-2)
		field.Place(x, y, tile)
		if field.Count() != tile.Count() {
			t.Fatalf("Wrong field.Count() after placing tile\n%s\n%s", tile, field)
		}

		field.Remove(x, y, tile)
		if field.Count() != 0 {
			t.Fatalf("Wrong field.Count() after removing tile\n%s\n%s", tile, field)
		}
	}
}
