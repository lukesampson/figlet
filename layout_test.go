package main

import "testing"

func Test_smush_with_lchr_empty_always_returns_rchr(t *testing.T) {
	rchrs := []rune { 'a', '!', '$' }
	for _, rchr := range rchrs {
		testSmushemAllSmushModes(t, ' ', rchr, rchr)
	}
}

func Test_smush_with_rchr_empty_always_returns_lchr(t *testing.T) {
	lchrs := []rune { 'a', '!', '$' }
	for _, lchr := range lchrs {
		testSmushemAllSmushModes(t, lchr, ' ', lchr)
	}
}

func Test_smush_with_smush_not_set_returns_false(t *testing.T) {
	lchr, rchr, smushmode := '|', '|', 0
	if x, ok := smushem(lchr, rchr, smushmode); ok != false {
		t.Errorf("smushem(%q, %q, %v) = (%q, %v), want (?, false)", lchr, rchr, smushmode, x, ok)
	}
}

func Test_smush_universal(t *testing.T) {
	//lchr, rchr, smushmode := ''
}

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

func testSmushemAllSmushModes(t *testing.T, lchr rune, rchr rune, expect rune) {
	for _, smushmode := range smushModes() {
		if x, ok := smushem(lchr, rchr, smushmode); x != expect {
			t.Errorf("smushem(%q, %q, %v) = (%q, %v), want %q", lchr, rchr, smushmode, x, ok, expect)
		}
	}
}