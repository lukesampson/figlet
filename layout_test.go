package main

import "testing"

func Test_smush_with_lch_empty_always_returns_rch(t *testing.T) {
	rchs := []rune { 'a', '!', '$' }
	for _, rch := range rchs {
		testSmushemAllSmushModes(t, ' ', rch, rch)
	}
}

func Test_smush_with_rch_empty_always_returns_lch(t *testing.T) {
	lchs := []rune { 'a', '!', '$' }
	for _, lch := range lchs {
		testSmushemAllSmushModes(t, lch, ' ', lch)
	}
}

func Test_smush_with_smush_not_set_returns_null(t *testing.T) {
	lch, rch, smushmode := '|', '|', 0
	if x := smushem(lch, rch, smushmode, '$', false); x != 0 {
		t.Errorf("smushem(%q, %q, %v) = %q, want %q", lch, rch, smushmode, x, 0)
	}
}

func Test_smush_universal(t *testing.T) {
	// smush mode of SMSmush but not SMKern is universal smushing
	lch, rch, mode, hblank, rtol := '|', '$', SMSmush, '$', false

	if x := smushem(lch, rch, mode, hblank, rtol); x != lch {
		t.Errorf("should return lch when rch is hardblank, returned %q", x)
	}

	lch, rch = rch, lch // swap
	if x := smushem(lch, rch, mode, hblank, rtol); x != rch {
		t.Errorf("should return rch when lch is hardblank, returned %q", x)
	}

	lch, rch = 'l', 'r'
	if x := smushem(lch, rch, mode, hblank, true); x != lch {
		t.Errorf("should return lch when right2left, returned %q", x)
	}

	if x := smushem(lch, rch, mode, hblank, false); x != rch {
		t.Errorf("should return rch when !right2left, returned %q", x)
	}
}

func Test_smush_combines_2_hardblanks_when_SMHardBlank(t *testing.T) {
	mode := SMSmush + SMKern + SMHardBlank

	if x := smushem('$', '$', mode, '$', false); x != '$' {
		t.Errorf("should smush 2 hardblanks to 1, returned %q", x)
	}
}

func Test_smush_doesnt_combine_any_hardblank_when_not_SMHardBlank(t *testing.T) {
	mode := SMSmush + SMKern

	if x := smushem('$', '|', mode, '$', false); x != 0 {
		t.Errorf("returned %q", x)
	}
}

func Test_smush_equal(t *testing.T) {
	if x := smushem('x', 'x', SMSmush + SMKern + SMEqual, '$', false); x != 'x' {
		t.Errorf("expected 'x', returned %q", x)
	}
}

func Test_smush_lowline(t *testing.T) {
	low, mode, hblank, rtol := '_', SMSmush + SMKern + SMLowLine, '$', false
	replacements := "|/\\[]{}()<>"
	for _, r := range replacements {
		if x := smushem(low, r, mode, hblank, rtol); x != r {
			t.Errorf("smush %q with %q should return %q, returned %q", low, r, r, x)
		}

		if x := smushem(r, low, mode, hblank, rtol); x != r {
			t.Errorf("smush %q with %q should return %q, returned %q", r, low, r, x)
		}
	}
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

func testSmushemAllSmushModes(t *testing.T, lch rune, rch rune, expect rune) {
	for _, smushmode := range smushModes() {
		if x := smushem(lch, rch, smushmode, '$', false); x != expect {
			t.Errorf("smushem(%q, %q, %v...) = %q, want %q", lch, rch, smushmode, x, expect)
		}
	}
}