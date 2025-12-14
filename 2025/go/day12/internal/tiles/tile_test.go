package tiles

import (
	"fmt"
	"math/rand/v2"
	"testing"
)

func createRandomTile() *Tile {
	b := []byte(".........")
	for i := range 9 {
		if rand.IntN(100) < 50 {
			b[i] = '#'
		}
	}
	return CreateTile(string(b))
}

func Example_createTile_1() {
	tile := CreateTile("..#.#####")
	fmt.Println(tile.String())

	// Output:
	// ..o
	// .oo
	// ooo
}

func Example_rotate_1() {
	tile := CreateTile("##..##..#")
	fmt.Println(tile.String())
	fmt.Println(tile.Rotate().String())

	// Output:
	// oo.
	// .oo
	// ..o
	//
	// ..o
	// .oo
	// oo.
}

func Example_mirror_1() {
	tile := CreateTile("##..#####")
	fmt.Println(tile.String())
	fmt.Println(tile.Mirror().String())

	// Output:
	// oo.
	// .oo
	// ooo
	//
	// .oo
	// oo.
	// ooo
}

func Test_equal_1(t *testing.T) {
	tile1 := CreateTile(".##.#####")
	tile2 := CreateTile("#.#####.#")
	if tile1.Equals(tile2) {
		t.Fatalf("Tiles shouldn't be equal:\n%s!=\n%s", tile1, tile2)
	}
}

func Test_mirror_1(t *testing.T) {
	for range 1000 {
		t1 := createRandomTile()
		t2 := t1.Mirror().Mirror()
		if !t1.Equals(t2) {
			t.Fatalf("Double-mirrored tile should be equal:\n%s==\n%s", t1, t2)
		}
	}
}

func Test_rotate_1(t *testing.T) {
	for range 1000 {
		t1 := createRandomTile()
		t2 := t1.Rotate().Rotate().Rotate().Rotate()
		if !t1.Equals(t2) {
			t.Fatalf("360°-rotated tile should be equal:\n%s==\n%s", t1, t2)
		}
	}
}

func Test_symmetrical(t *testing.T) {
	tile1 := CreateTile("#.#####.#")
	tile2 := tile1.Mirror()
	if !tile1.Equals(tile2) {
		t.Fatalf("Mirrored tile should be equal:\n%s==\n%s", tile1, tile2)
	}

	tile3 := tile1.Rotate().Rotate()
	if !tile1.Equals(tile3) {
		t.Fatalf("180°-rotated tile should be equal:\n%s==\n%s", tile1, tile3)
	}
}
