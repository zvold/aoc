package util

import "golang.org/x/exp/constraints"

// Abs returns the absolute value of its argument.
func Abs[T constraints.Integer](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

// Gcd returns GCD of two integers.
func Gcd(i, j int) int {
	a := max(i, j)
	b := min(i, j)
	if a == 0 {
		return b
	} else if b == 0 {
		return a
	}
	return Gcd(a%b, b)
}

// Gcd2 returns GCD of a slice of integers.
func Gcd2(values ...int) int {
	if len(values) == 1 {
		return values[0]
	}
	if len(values) == 2 {
		return Gcd(values[0], values[1])
	}
	return Gcd2(append(values[2:], Gcd(values[0], values[1]))...)
}

// Lcm returns LCM of two integers.
func Lcm(i, j int) int {
	return Abs(max(i, j)) / Gcd(i, j) * Abs(min(i, j))
}

// Lcm2 returns LCM of a slice of integers.
func Lcm2(values ...int) (r int) {
	if len(values) == 1 {
		return values[0]
	}
	if len(values) == 2 {
		return Lcm(values[0], values[1])
	}
	return Lcm2(append(values[2:], Lcm(values[0], values[1]))...)
}
