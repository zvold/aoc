package util

// Direction is the enum with all possible directions.
type Direction int

const (
	N Direction = iota
	E
	S
	W
)

var dirStrings = []string{N: "↑", E: "→", S: "↓", W: "←"}

func (d Direction) String() string {
	return dirStrings[d]
}

func (d Direction) IsOpposite(d2 Direction) bool {
	return d.Opposite() == d2
}

func (d Direction) Opposite() Direction {
	return (d + 2) % 4
}

var Shifts = map[Direction]Pos{
	N: {0, -1},
	E: {+1, 0},
	S: {0, +1},
	W: {-1, 0},
}

type Pos struct {
	X, Y int
}

func (p Pos) Move(d Direction) Pos {
	return Pos{p.X + Shifts[d].X, p.Y + Shifts[d].Y}
}

func (p Pos) Manhattan(p2 Pos) int {
	return Abs(p.X-p2.X) + Abs(p.Y-p2.Y)
}
