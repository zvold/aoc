package hand

import (
	"fmt"
	"log"
	"strings"
)

type handstrength int

const (
	highCard handstrength = iota
	onePair
	twoPair
	threeOfKind
	fullHouse
	fourOfKind
	fiveOfKind
)

func (h handstrength) String() string {
	return []string{
		"high card",
		"one pair",
		"two pair",
		"3 of a kind",
		"full house",
		"4 of a kind",
		"5 of a kind",
	}[h]
}

var (
	// Normal card ordering for task 1.
	kindStrengths = map[byte]int{
		'2': 1,
		'3': 2,
		'4': 3,
		'5': 4,
		'6': 5,
		'7': 6,
		'8': 7,
		'9': 8,
		'T': 9,
		'J': 10,
		'Q': 11,
		'K': 12,
		'A': 13,
	}

	// Jokers have the smallest value on their own.
	jokerStrengths = map[byte]int{
		'J': 0,
		'2': 1,
		'3': 2,
		'4': 3,
		'5': 4,
		'6': 5,
		'7': 6,
		'8': 7,
		'9': 8,
		'T': 9,
		'Q': 11,
		'K': 12,
		'A': 13,
	}
)

type Hand struct {
	Cards         string
	Strength      handstrength
	JokerStrength handstrength
}

func (l Hand) String() string {
	return fmt.Sprintf("[%s] (%s) (%s)",
		l.Cards, l.Strength, l.JokerStrength)
}
func (l Hand) CompareStrength(r Hand) bool {
	if l.Strength != r.Strength {
		return l.Strength > r.Strength
	}
	for i := 0; i < 5; i++ {
		if l.Cards[i] != r.Cards[i] {
			return kindStrengths[l.Cards[i]] > kindStrengths[r.Cards[i]]
		}
	}
	return true
}
func (l Hand) CompareJokerStrength(r Hand) bool {
	if l.JokerStrength != r.JokerStrength {
		return l.JokerStrength > r.JokerStrength
	}
	for i := 0; i < 5; i++ {
		if l.Cards[i] != r.Cards[i] {
			return jokerStrengths[l.Cards[i]] > jokerStrengths[r.Cards[i]]
		}
	}
	return true
}

func NewHand(cards string) Hand {
	if len(cards) != 5 {
		log.Fatalf("Invalid hand string: %s", cards)
	}
	return Hand{cards, getStrength(cards), getJokerStrength(cards)}
}

// Calculates the strength of a given hand.
func getStrength(cards string) handstrength {
	if len(cards) != 5 {
		log.Fatalf("Invalid hand string: %s", cards)
	}

	m := make(map[rune]int)
	for _, v := range cards {
		m[v]++
	}

	switch len(m) {
	case 1:
		return fiveOfKind
	case 2:
		if maxValue(m) == 4 {
			return fourOfKind
		} else {
			return fullHouse
		}
	case 3:
		if maxValue(m) == 3 {
			return threeOfKind
		} else {
			return twoPair
		}
	case 4:
		return onePair
	case 5:
		return highCard
	default:
		panic("Hand with more than 5 cards.")
	}
}

// Calculates the max strength when 'J' can pretend to be any card.
func getJokerStrength(cards string) handstrength {
	best := getStrength(cards)
	if !strings.Contains(cards, "J") {
		return best
	}

	// Replace 'J' with different kinds and see if it improves the strength.
	for v := range kindStrengths {
		if v == 'J' {
			continue
		}
		cards2 := strings.ReplaceAll(cards, "J", string(v))
		strength := getStrength(cards2)
		if strength > best {
			best = strength
		}
	}
	return best
}

func maxValue(m map[rune]int) (result int) {
	for _, v := range m {
		if v > result {
			result = v
		}
	}
	return
}
