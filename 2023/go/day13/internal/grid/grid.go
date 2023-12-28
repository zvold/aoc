package grid

import (
	"hash/fnv"
	"log"
)

type Grid struct {
	lines  []string
	hashes []uint32
	// When transposed=true, the rows are stored as columns in the 'lines' above.
	transposed bool
	// All smudge-equal rows (when transposed=false), or columns (when true).
	smudgeMap map[pair]bool
}

// Specifies the indexes of "smudge-equal" rows (or columns).
type pair struct {
	i, j int
}

// Dir specifies the direction of found reflection.
type Dir int

const (
	H Dir = iota
	V
)

func NewGrid(lines []string) Grid {
	hashes := make([]uint32, 0, len(lines))
	for _, v := range lines {
		hashes = append(hashes, hashs(v))
	}
	return Grid{lines, hashes, false, constructSmudgeMap(lines, H)}
}

func (g Grid) transpose() Grid {
	hashes := make([]uint32, 0, len(g.lines))
	tmp := make([]byte, 0, len(g.lines))
	for i := 0; i < len(g.lines[0]); i++ {
		for j := 0; j < len(g.lines); j++ {
			tmp = append(tmp, g.lines[j][i])
		}
		hashes = append(hashes, hashb(tmp))
		tmp = tmp[:0] // Keep the underlying array allocated.
	}
	return Grid{g.lines, hashes, true, constructSmudgeMap(g.lines, V)}
}

func constructSmudgeMap(lines []string, dir Dir) map[pair]bool {
	result := make(map[pair]bool)
	switch dir {
	case H:
		for i := 0; i < len(lines); i++ {
			for j := i + 1; j < len(lines); j++ {
				if smudgeEqual(lines[i], lines[j]) {
					result[pair{i, j}] = true
					result[pair{j, i}] = true
				}
			}
		}
	case V:
		for i := 0; i < len(lines[0]); i++ {
			for j := i + 1; j < len(lines[0]); j++ {
				// 'i' and 'j' are column indexes.
				if smudgeEqualColumns(lines, i, j) {
					result[pair{i, j}] = true
					result[pair{j, i}] = true
				}
			}
		}
	}
	return result
}

// Column version of smudgeEqual() function.
func smudgeEqualColumns(lines []string, c1 int, c2 int) bool {
	c := 0
	for _, line := range lines {
		if line[c1] != line[c2] {
			c++
		}
		if c > 1 {
			return false
		}
	}
	return c == 1 // Columns that are equal w/o fixing are not smudge-equal.
}

// Returns true if two strings can be made equal by fixing exactly 1 smudge.
func smudgeEqual(s1 string, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}
	c := 0
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			c++
		}
		if c > 1 {
			return false
		}
	}
	return c == 1 // Lines that are equal w/o fixing are not smudge-equal.
}

// FindReflection returns the first found reflection in the grid.
// If smudged is set, exactly one smudge must be fixed before reflection is found.
func (g Grid) FindReflection(smudged bool) (Dir, int) {
	r := g.findRowReflection(smudged)
	if r != -1 {
		return H, r
	} else {
		v := g.transpose().findRowReflection(smudged)
		if v == -1 {
			log.Fatal("Vertical reflection must exist if horizontal doesn't.")
		}
		return V, v
	}
}

// Returns the first row number after which reflection occurs.
func (g Grid) findRowReflection(smudged bool) int {
	for j := 0; j < len(g.hashes)-1; j++ {
		if g.checkReflection(j, smudged) {
			return j
		}
	}
	return -1
}

// Checks if the grid is reflected around axis b/w j and j+1.
func (g Grid) checkReflection(j int, smudged bool) bool {
	fixes := 0
	k := j + 1
	for j >= 0 && k < len(g.hashes) {
		switch smudged {
		case true:
			if g.hashes[j] != g.hashes[k] {
				// Rows are not equal by themselves:
				if fixes >= 1 {
					// We can only fix one smudge.
					return false
				} else {
					if g.smudgeMap[pair{j, k}] {
						// Rows are smudge-equal, use the allowed fix and continue.
						fixes++
					} else {
						return false
					}
				}
			}
		case false:
			if g.hashes[j] != g.hashes[k] {
				return false
			}
		}
		j--
		k++
	}

	if smudged {
		// For smudge-reflection, we must fix exactly one smudge.
		return fixes == 1
	} else {
		return true
	}
}

func hashs(s string) uint32 {
	return hashb([]byte(s))
}

func hashb(s []byte) uint32 {
	h := fnv.New32a()
	_, err := h.Write(s)
	if err != nil {
		return 0
	}
	return h.Sum32()
}
