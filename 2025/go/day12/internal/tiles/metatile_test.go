package tiles

import "fmt"

func Example_createMetaTile_1() {
	tile := CreateTile("##..##..#")
	metatile := CreateMetaTile(tile)
	fmt.Println(metatile.String())

	// Output:
	// | oo. | .oo | ..o | o.. |
	// | .oo | oo. | .oo | oo. |
	// | ..o | o.. | oo. | .oo |
}

func Example_createMetaTile_2() {
	tile := CreateTile("#.#####.#")
	metatile := CreateMetaTile(tile)
	fmt.Println(metatile.String())

	// Output:
	// | o.o | ooo |
	// | ooo | .o. |
	// | o.o | ooo |
}

func Example_createMetaTile_3() {
	tile := CreateTile("##.#.##..")
	metatile := CreateMetaTile(tile)
	fmt.Println(metatile.String())

	// Output:
	// | oo. | .oo | ooo | ooo | ..o | o.. | .o. | .o. |
	// | o.o | o.o | ..o | o.. | o.o | o.o | o.. | ..o |
	// | o.. | ..o | .o. | .o. | .oo | oo. | ooo | ooo |
}
