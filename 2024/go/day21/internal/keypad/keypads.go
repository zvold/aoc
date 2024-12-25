package keypad

import (
	"log"
	"math"
	"strings"

	u "github.com/zvold/aoc/2023/go/util"
)

var (
	NumKeypad Keypad = Keypad{values: "789456123#0A"}
	DirKeypad Keypad = Keypad{values: "#^A<v>"}
)

type Keypad struct {
	values string
}

// Returns the physical configuration of this particular keypad.
func (k Keypad) Key2Pos(b byte) u.Pos {
	i := strings.IndexByte(k.values, b)
	if i == -1 || b == '#' {
		log.Fatalf("Invalid key %c", b)
	}
	return u.Pos{X: i % 3, Y: i / 3}
}

// Returns the key in the position 'p' of the keypad.
func (k Keypad) pos2key(p u.Pos) string {
	i := p.Y*3 + p.X
	return k.values[i : i+1]
}

// Returns all possible sequences of Moves for getting from 'p' to 'p2'.
func (k Keypad) Moves(p, p2 u.Pos) map[string]bool {
	r := make(map[string]bool, 0)

	vertical, horizontal := "", ""
	if p2.X > p.X {
		vertical = strings.Repeat(">", p2.X-p.X)
	} else {
		vertical = strings.Repeat("<", p.X-p2.X)
	}
	if p2.Y > p.Y {
		horizontal = strings.Repeat("v", p2.Y-p.Y)
	} else {
		horizontal = strings.Repeat("^", p.Y-p2.Y)
	}

	if k.safe(p, vertical+horizontal) {
		r[vertical+horizontal] = true
	}
	if k.safe(p, horizontal+vertical) {
		r[horizontal+vertical] = true
	}
	return r
}

// Returns 'true' if sequence of moves never moves through '#'.
func (k Keypad) safe(p u.Pos, s string) bool {
	for _, v := range []byte(s) {
		switch v {
		case '<':
			p = p.Move(u.W)
		case '>':
			p = p.Move(u.E)
		case '^':
			p = p.Move(u.N)
		case 'v':
			p = p.Move(u.S)
		case 'A':
			continue
		default:
			log.Fatalf("Unexpected command %c", v)
		}
		if k.pos2key(p) == "#" {
			return false
		}
	}
	return true
}

func (k Keypad) Translate(s string) map[string]bool {
	return k.Translate2(k.Key2Pos('A'), s)
}

// Translates a sequence 's' of buttons that need to be pressed, starting from position 'p', into a
// set of sequences of moves like '<A^A>^^AvvvA'.
func (k Keypad) Translate2(p u.Pos, s string) map[string]bool {
	// Destination position.
	p2 := k.Key2Pos(s[0])

	// All possible move sequences for getting from p to p2.
	paths := k.Moves(p, p2)

	result := make(map[string]bool, 0)

	if len(s) == 1 {
		for path := range paths {
			result[path+"A"] = true
		}
	} else {
		// Append 'A' + every possible path from the recursive call.
		for path := range paths {
			for s := range k.Translate2(p2, s[1:]) {
				result[path+"A"+s] = true
			}
		}
	}

	return result
}

// Translates a sequence of moves and actuations into the sequence of keys.
func (k Keypad) Forward(s string) string {
	var r string
	p := k.Key2Pos('A') // Initial position.
	for _, v := range []byte(s) {
		switch v {
		case '<':
			p = p.Move(u.W)
		case '>':
			p = p.Move(u.E)
		case '^':
			p = p.Move(u.N)
		case 'v':
			p = p.Move(u.S)
		case 'A':
			r += k.pos2key(p)
		default:
			log.Fatalf("Unexpected command %c", v)
		}
	}
	return r
}

// Some explanation, as this was not very trivial, and I won't remember what's going on here.
// Consider the sequence human -> r1 -> r2 -> ... -> rN -> keypad
//                        L0      L1    L2    ...    LN
//
// First thing to realize is that once a key is pressed on the final keypad, all previous robots
// and the human are pointing to 'A'. This it because robot rN just received 'actuate' command,
// meaning robot rN-1 just pressed 'A'. Meaning, rN-1 just received 'actuate' command, so rN-2 just
// pressed 'A'...
//
// So, every keypress on the final keypad is pressed independently (the whole chain of previous
// robots are starting from the default position).
//
// Let's say the final robot has just pressed '3' in the sequence '379A'. It's pointing to key '3'.
// What is the sequence of moves it needs to do to press '7' from this position? There are two:
//   <<^^A
//   ^^<<A
// It's unclear which is best, so we'll check the "cost" of both.
//
// So our robot rN needs to make the moves < < ^ ^ A. What does this mean for robot rN-1?
// It means rN-1 needs to move its arm 'A' -> '<' then 'actuate', then make the move '<' -> '<'
// and actuate, and so on:
//   rN:     <     <     ^     ^    A
//  rN-1:  A→< A <→< A <→^ A ^→^ A      (this is implemented by calling Cost2 function)
//
// How does rN-1 move its arm from say 'A'->'v' and actuate it? Well there are two ways:
//   <vA
//   v<A
// (this is the same problem as before but for a lower-level robot...)
//
// The last level robot is controlled via the sequence directly by the human, so e.g. when it needs
// to move say 'A'->'<' and actuate, that would be just 'v<<A' (cost 4).
//
// In other words:
//   costLn("<<^^A") =
//   costLn-1(A→<)    + costLn-1(<→<) + ... costLn-1(^→A) =
//   costLn-1("v<<A") + costLn-1("A") + ... costLn-1(">A") =
//   costLn-2(A→v) + ...

func (k Keypad) Cost(level int, s string) uint64 {
	//fmt.Printf("level %d, path '%s'\n", level, s)
	var r uint64 = 0
	s = "A" + s // All robots are on 'A' key after the final robot outputted the last key.
	for i := range len(s) - 1 {
		r += k.Cost2(level-1, k.Key2Pos(s[i]), k.Key2Pos(s[i+1]))
	}

	return r
}

type costkey struct {
	level int
	p     u.Pos
	p2    u.Pos
}

var cache map[costkey]uint64 = make(map[costkey]uint64, 0)

func (k Keypad) Cost2(level int, p, p2 u.Pos) uint64 {
	if level == 0 {
		//fmt.Printf("\tlevel 0, %s -> %s, cost %d\n", k.pos2key(p), k.pos2key(p2), p.Manhattan(p2)+1)
		return uint64(p.Manhattan(p2) + 1) // +1 because we also press 'A' in the end.
	} else {
		if v, ok := cache[costkey{level: level, p: p, p2: p2}]; ok {
			return v
		}
		paths := k.Moves(p, p2)
		var min uint64 = math.MaxUint64
		for p := range paths {
			c := uint64(k.Cost(level, p+"A"))
			if c < min {
				min = c
			}
		}
		cache[costkey{level: level, p: p, p2: p2}] = min
		return min
	}
}
