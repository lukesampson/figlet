package main

import "testing"

func TestSmushem_LChrEmpty_Always_Returns_Rchr(t *testing.T) {
	lchr := ' '
	rchrs := []rune { 'a', '!', '$' }
	for _, rchr := range rchrs {
		if x:= smushem(lchr, rchr); x != rchr {
			t.Errorf("smushem(%v, %v) = %v, want %v", lchr, rchr, x, rchr)
		}
	}
}