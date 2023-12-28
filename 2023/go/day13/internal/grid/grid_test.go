package grid

import "testing"

func TestGrid_FindReflection_1(t *testing.T) {
	g := NewGrid([]string{
		"###.##..#",
		"#####....",
		"#####....",
		"##....##.",
		"####.#..#",
		"...#.....",
		".#.######",
	})

	wantDir, wantI := V, 6
	gotDir, gotI := g.FindReflection(false)
	if wantDir != gotDir {
		t.Errorf("grid.FindReflection() should return (%v, ...), returned %v", wantDir, gotDir)
	}
	if wantI != gotI {
		t.Errorf("grid.FindReflection() should return (..., %v), returned %v", wantI, gotI)
	}
}

func TestGrid_FindReflection_2(t *testing.T) {
	g := NewGrid([]string{
		"###.##..#",
		"#####....",
		"#####....",
		"##....##.",
		"####.#..#",
		"...#.....",
		".#.######",
	})

	wantDir, wantI := V, 0
	gotDir, gotI := g.FindReflection(true)
	if wantDir != gotDir {
		t.Errorf("grid.FindReflection() should return (%v, ...), returned %v", wantDir, gotDir)
	}
	if wantI != gotI {
		t.Errorf("grid.FindReflection() should return (..., %v), returned %v", wantI, gotI)
	}
}

func TestGrid_FindReflection_3(t *testing.T) {
	g := NewGrid([]string{
		"#.##..##.",
		"..#.##.#.",
		"##......#", // < Smudge-equal reflection here.
		"##......#",
		"..#.##.#.",
		"..##..##.",
		"#.#.##.#.",
	})

	wantDir, wantI := H, 2
	gotDir, gotI := g.FindReflection(true)
	if wantDir != gotDir {
		t.Errorf("grid.FindReflection() should return (%v, ...), returned %v", wantDir, gotDir)
	}
	if wantI != gotI {
		t.Errorf("grid.FindReflection() should return (..., %v), returned %v", wantI, gotI)
	}
}
