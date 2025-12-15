package util

import "testing"

func TestGcd_1(t *testing.T) {
	i := 252
	j := 105
	want := 21
	got := Gcd(i, j)
	if want != got {
		t.Fatalf("Gcd(%d, %d) should return %d, returned %d", i, j, want, got)
	}
}

func TestGcd2_1(t *testing.T) {
	v := []int{252, 105, 21 * 17}
	want := 21
	got := Gcd2(v...)
	if want != got {
		t.Fatalf("Gcd2(%v) should return %d, returned %d", v, want, got)
	}
}

func TestGcd2_2(t *testing.T) {
	v := []int{2794, 2910, 3012, 3050}
	want := 2
	got := Gcd2(v...)
	if want != got {
		t.Fatalf("Gcd2(%v) should return %d, returned %d", v, want, got)
	}
}

func TestLcm_1(t *testing.T) {
	i := 2793
	j := 2917
	want := 8147181
	got := Lcm(i, j)
	if want != got {
		t.Fatalf("Lcm(%d, %d) should return %d, returned %d", i, j, want, got)
	}
}

func TestLcm_2(t *testing.T) {
	i := 3013
	j := 3051
	want := 9192663
	got := Lcm(i, j)
	if want != got {
		t.Fatalf("Lcm(%d, %d) should return %d, returned %d", i, j, want, got)
	}
}

func TestLcm2_1(t *testing.T) {
	v := []int{279, 297, 313, 301}
	want := 867419091
	got := Lcm2(v...)
	if want != got {
		t.Fatalf("Lcm2(%v) should return %d, returned %d", v, want, got)
	}
}
