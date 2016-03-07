package figletlib

import (
	"testing"
)

func Test_smush_with_lch_empty_always_returns_rch(t *testing.T) {
	rchs := []rune{'a', '!', '$'}
	for _, rch := range rchs {
		testSmushemAllSmushModes(t, ' ', rch, rch)
	}
}

func Test_smush_with_rch_empty_always_returns_lch(t *testing.T) {
	lchs := []rune{'a', '!', '$'}
	for _, lch := range lchs {
		testSmushemAllSmushModes(t, lch, ' ', lch)
	}
}

func Test_smush_with_smush_not_set_returns_null(t *testing.T) {
	lch, rch := '|', '|'
	if x := smushem(lch, rch, testSettings(0)); x != 0 {
		t.Errorf("smushem(%q, %q, %v) = %q, want %q", lch, rch, 0, x, 0)
	}
}

func Test_smush_universal(t *testing.T) {
	// smush mode of SMSmush but not SMKern is universal smushing
	lch, rch := '|', '$'
	s := testSettings(SMSmush)

	if x := smushem(lch, rch, s); x != lch {
		t.Errorf("should return lch when rch is hardblank, returned %q", x)
	}

	lch, rch = rch, lch // swap
	if x := smushem(lch, rch, s); x != rch {
		t.Errorf("should return rch when lch is hardblank, returned %q", x)
	}

	lch, rch = 'l', 'r'
	s.rtol = true
	if x := smushem(lch, rch, s); x != lch {
		t.Errorf("should return lch when right2left, returned %q", x)
	}

	s.rtol = false
	if x := smushem(lch, rch, s); x != rch {
		t.Errorf("should return rch when !right2left, returned %q", x)
	}

	testSmush('|', '$', s.smushmode, '|', t)
}

func Test_smush_combines_2_hardblanks_when_SMHardBlank(t *testing.T) {
	s := testSettings(SMSmush + SMKern + SMHardBlank)

	if x := smushem('$', '$', s); x != '$' {
		t.Errorf("should smush 2 hardblanks to 1, returned %q", x)
	}
}

func Test_smush_doesnt_combine_any_hardblank_when_not_SMHardBlank(t *testing.T) {
	s := testSettings(SMSmush + SMKern)

	if x := smushem('$', '|', s); x != 0 {
		t.Errorf("returned %q", x)
	}
}

func Test_smush_equal(t *testing.T) {
	if x := smushem('x', 'x', testSettings(SMSmush+SMKern+SMEqual)); x != 'x' {
		t.Errorf("expected 'x', returned %q", x)
	}
}

func Test_smush_lowline(t *testing.T) {
	replacements := "|/\\[]{}()<>"
	for _, r := range replacements {
		testSmushLowLine('_', r, r, t)
		testSmushLowLine(r, '_', r, t)
	}
}

func Test_smush_heirarchy(t *testing.T) {
	testSmushHierarchy('|', '\\', '\\', t)
	testSmushHierarchy('}', '|', '}', t)

	testSmushHierarchy('/', '>', '>', t)
	testSmushHierarchy('{', '\\', '{', t)

	testSmushHierarchy('[', '(', '(', t)
	testSmushHierarchy('>', ']', '>', t)

	testSmushHierarchy('}', ')', ')', t)
	testSmushHierarchy('<', '{', '<', t)

	testSmushHierarchy('(', '<', '<', t)
	testSmushHierarchy('>', '(', '>', t)
}

func Test_smush_pairs(t *testing.T) {
	testSmushPair('[', ']', '|', t)
	testSmushPair(']', '[', '|', t)

	testSmushPair('(', ')', '|', t)
	testSmushPair(')', '(', '|', t)

	testSmushPair('{', '}', '|', t)
	testSmushPair('}', '{', '|', t)
}

func Test_smush_bigX(t *testing.T) {
	mode := SMKern + SMSmush + SMBigX
	testSmush('/', '\\', mode, '|', t)
	testSmush('\\', '/', mode, 'Y', t)
	testSmush('>', '<', mode, 'X', t)
}

func Test_smush_hardblank(t *testing.T) {
	mode := SMKern + SMSmush + SMHardBlank
	testSmush('$', '$', mode, '$', t)
	testSmush(' ', '$', mode, '$', t)
}

func Test_smushamt(t *testing.T) {
	testSmushamtLine("|_ ", "  _", 3, t)
	testSmushamtLine("", "__", 0, t)
}

func Test_addChar(t *testing.T) {
	testAddCharLine("", "", "", t)
	testAddCharLine("", "__", "__", t)
	testAddCharLine("", " __", "__", t)
	testAddCharLine("|_ ", "  _", "|__", t)
	testAddCharLine("|_ ", "   _", "|__", t)
}

func Test_smushChar(t *testing.T) {
	testSmushCharLine("/ ", "  / ", 4, "/ ", t)
	testSmushCharLine("", " _", 1, "_", t)
}

func testSmushamtLine(l string, c string, want int, t *testing.T) {
	line := newFigText(1)
	char := newFigText(1)

	(*line).art[0] = []rune(l)
	(*char).art[0] = []rune(c)

	s := testSettings(SMKern + SMSmush + SMEqual + SMLowLine + SMHierarchy + SMPair)

	if x := smushamt(char, line, s); x != want {
		t.Errorf("smushamt = %v, want %v", x, want)
	}
}

func testAddCharLine(l string, c string, expect string, t *testing.T) {
	line := newFigText(1)
	char := newFigText(1)

	(*line).art[0] = []rune(l)
	(*char).art[0] = []rune(c)

	s := testSettings(SMKern + SMSmush + SMEqual + SMLowLine + SMHierarchy + SMPair)

	result := addChar(char, line, s)

	if string(result.art[0]) != expect {
		t.Errorf("addChar %q + %q made %q, expected %q", l, c, string((*line).art[0]), expect)
	}
}

func testSmushCharLine(l string, c string, amount int, expect string, t *testing.T) {
	if amount > len(c) {
		t.Errorf("amount %v is > character length %v", amount, len(c))
		return
	}
	line := newFigText(1)
	char := newFigText(1)

	(*line).art[0] = []rune(l)
	(*char).art[0] = []rune(c)

	s := testSettings(SMKern + SMSmush + SMEqual + SMLowLine + SMHierarchy + SMPair)

	result := smushChar(char, line, amount, s)

	if string(result.art[0]) != expect {
		t.Errorf("smushChar %q + %q made %q, expected %q", l, c, string((*line).art[0]), expect)
	}
}

func testSettings(smushmode int) Settings {
	return Settings{
		smushmode: smushmode,
		hardblank: '$',
		rtol:      false,
	}
}

func testSmushLowLine(l rune, r rune, expect rune, t *testing.T) {
	testSmush(l, r, SMKern+SMSmush+SMLowLine, expect, t)
}
func testSmushHierarchy(l rune, r rune, expect rune, t *testing.T) {
	testSmush(l, r, SMKern+SMSmush+SMHierarchy, expect, t)
}
func testSmushPair(l rune, r rune, expect rune, t *testing.T) {
	testSmush(l, r, SMKern+SMSmush+SMPair, expect, t)
}
func testSmush(l rune, r rune, mode int, expect rune, t *testing.T) {
	if x := smushem(l, r, testSettings(mode)); x != expect {
		t.Errorf("smush %q + %q => %q, want %q", l, r, x, expect)
	}
}

func validSmushModes() []int {
	modes := make([]int, 0, 36+1)
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
		if x := smushem(lch, rch, testSettings(smushmode)); x != expect {
			t.Errorf("smushem(%q, %q, %v...) = %q, want %q", lch, rch, smushmode, x, expect)
		}
	}
}
