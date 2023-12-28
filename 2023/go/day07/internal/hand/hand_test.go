package hand

import "testing"

func TestHand_CompareStrength(t *testing.T) {
	h1 := NewHand("KK677")
	h2 := NewHand("KTJJT")

	if !h1.CompareStrength(h2) {
		t.Errorf("Hand.CompareStrength(): hand rank should be {%v} > {%v}", h1, h2)
	}
}

func TestHand_CompareJokerStrength(t *testing.T) {
	h1 := NewHand("T55J5")
	h2 := NewHand("KTJJT")

	if h1.CompareJokerStrength(h2) {
		t.Errorf("Hand.CompareJokerStrength(): hand rank should be {%v} < {%v}", h1, h2)
	}
}
