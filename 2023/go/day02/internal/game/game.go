package game

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	gameIdRegex = regexp.MustCompile(`Game (\d+): (.*)`)
	colorRegex  = regexp.MustCompile(`\b(\d+) (green|blue|red)\b`)
)

type Cubes struct {
	Red   int
	Green int
	Blue  int
}

func (c *Cubes) String() string {
	return fmt.Sprintf("[r: %d, g: %d, b: %d]", c.Red, c.Green, c.Blue)
}
func (c *Cubes) Power() int {
	return c.Red * c.Green * c.Blue
}

type Game struct {
	Id   int
	Sets []*Cubes
}

func (g *Game) String() string {
	return fmt.Sprintf("[id: %d, sets: %s", g.Id, g.Sets)
}
func (g *Game) MinCubes() *Cubes {
	r := new(Cubes)
	for _, s := range g.Sets {
		if s.Red > r.Red {
			r.Red = s.Red
		}
		if s.Green > r.Green {
			r.Green = s.Green
		}
		if s.Blue > r.Blue {
			r.Blue = s.Blue
		}
	}
	return r
}

func ParseGame(s string) (*Game, error) {
	groups := gameIdRegex.FindStringSubmatch(s)
	if len(groups) < 3 {
		return nil, fmt.Errorf("cannot parse: %s", s)
	}

	result := new(Game)

	var err error
	result.Id, err = strconv.Atoi(groups[1])
	if err != nil {
		return nil, err
	}

	result.Sets, err = parseCubeSets(groups[2])
	if err != nil {
		return nil, err
	}

	return result, nil
}

func parseCubeSets(s string) ([]*Cubes, error) {
	var result []*Cubes

	parts := strings.Split(s, ";")
	if len(parts) == 0 {
		return result, nil
	}

	for _, p := range parts {
		c, err := parseCubeSet(p)
		if err != nil {
			return nil, err
		}
		result = append(result, c)
	}

	return result, nil
}

func parseCubeSet(s string) (*Cubes, error) {
	result := new(Cubes)

	matches := colorRegex.FindAllStringSubmatch(s, -1)

	var err error
	for _, v := range matches {
		switch v[2] {
		case "green":
			result.Green, err = strconv.Atoi(v[1])
		case "red":
			result.Red, err = strconv.Atoi(v[1])
		case "blue":
			result.Blue, err = strconv.Atoi(v[1])
		}
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
