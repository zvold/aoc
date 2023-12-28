package game

import "testing"

func TestCubes_Power(t *testing.T) {
	cubes := Cubes{2, 5, 7}
	want := 70
	got := cubes.Power()
	if got != want {
		t.Errorf("Cubes.Power() for '%v' should return %v, returned %v", cubes, want, got)
	}
}

func TestGame_MinCubes(t *testing.T) {
	g, err := ParseGame("Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red")
	if err != nil {
		t.Fatal(err)
	}

	want := &Cubes{20, 13, 6}
	got := g.MinCubes()
	if *got != *want {
		t.Errorf("Game.MinCubes() should return %v, returned %v", want, got)
	}
}

func TestParseGame(t *testing.T) {
	g, err := ParseGame("Game 18: 6 red, 4 green, 7 blue; 2 red, 3 green, 12 blue; 3 red, 6 blue, 6 green; " +
		"9 red, 10 blue; 6 green, 4 blue, 2 red; 12 red, 12 blue, 9 green")
	if err != nil {
		t.Fatal(err)
	}

	want := "[id: 18, sets: [[r: 6, g: 4, b: 7] [r: 2, g: 3, b: 12] [r: 3, g: 6, b: 6] " +
		"[r: 9, g: 0, b: 10] [r: 2, g: 6, b: 4] [r: 12, g: 9, b: 12]]"
	got := g.String()
	if got != want {
		t.Errorf("Game.String() should return %v, returned %v", want, got)
	}
}
