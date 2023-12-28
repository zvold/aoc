package card

import (
	"log"
	"regexp"
	"strconv"
)

var cardRe = regexp.MustCompile(`Card\s+(\d+): (.*) \| (.*)`)
var numRe = regexp.MustCompile(`\b\d+\b`)

type Card struct {
	Id      int
	Winning map[int]bool
	Present map[int]bool
}

func (c *Card) Count() int {
	count := 0
	for k := range c.Present {
		if c.Winning[k] {
			count++
		}
	}
	return count
}

func (c *Card) Points() int {
	count := c.Count()
	if count == 0 {
		return 0
	} else {
		return 1 << (count - 1)
	}
}

func ParseCard(s string) *Card {
	groups := cardRe.FindStringSubmatch(s)
	if len(groups) < 4 {
		log.Fatalf("cannot parse card: '%s'", s)
	}
	var err error
	var c Card
	c.Id, err = strconv.Atoi(groups[1])
	if err != nil {
		log.Fatalf("invalid card id: %s", groups[1])
	}

	c.Winning, err = ParseSet(groups[2])
	if err != nil {
		log.Fatalf("cannot parse winning set: '%s'", groups[2])
	}

	c.Present, err = ParseSet(groups[3])
	if err != nil {
		log.Fatalf("cannot parse present set: '%s'", groups[3])
	}

	return &c
}

func ParseSet(s string) (map[int]bool, error) {
	m := make(map[int]bool)

	for _, r := range numRe.FindAllStringIndex(s, -1) {
		n, err := strconv.Atoi(s[r[0]:r[1]])
		if err != nil {
			return nil, err
		}
		m[n] = true
	}

	return m, nil
}
