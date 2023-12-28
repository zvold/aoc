package card

import "testing"

func TestCard_Count(t *testing.T) {
	c := ParseCard("Card 4: 41 92 73 84 69 | 59 84 76 51 58  5 54 83")

	want := 1
	got := c.Count()
	if got != want {
		t.Errorf("Card.Count() should return %v, returned %v", want, got)
	}
}

func TestCard_Points(t *testing.T) {
	c := ParseCard("Card 2: 13 32 20 16 61 | 61 30 68 82 17 32 24 19")

	want := 2
	got := c.Points()
	if got != want {
		t.Errorf("Card.Points() should return %v, returned %v", want, got)
	}
}
