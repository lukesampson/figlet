package main

import "testing"

func validSmushModes() []int {
	modes := make([]int, 0, 36 + 1)
	modes = append(modes, 0)
	for i := 1; i <= 128; i *= 2 {
		modes = append(modes, i)
		for j := 1; j < i; j *= 2 {
			modes = append(modes, i+j)
		}
	}
	return modes
}

func smushModes() []int {
	modes := make([]int, 8+1)
	modes[0] = 0
	for i := uint(0); i < 8; i++ {
		modes[i+1] = 1 << i
	}
	return modes
}

func TestSmushem_LChrEmpty_Always_Returns_Rchr(t *testing.T) {
	lchr := ' '
	rchrs := []rune { 'a', '!', '$' }
	for _, rchr := range rchrs {
		if x, ok := smushem(lchr, rchr, 0); ok && x != rchr {
			t.Errorf("smushem(%v, %v) = %v, want %v", lchr, rchr, x, rchr)
		}
	}
}