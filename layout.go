package main

import (
	"strings"
	//"fmt"
)

// smush modes
const (
	SMEqual = 1
	SMLowLine = 2
	SMHierarchy = 4
	SMPair = 8
	SMBigX = 16
	SMHardBlank = 32
	SMKern = 64
	SMSmush = 128
)

// gets the font entry for the given character, or the 'missing'
// character if the font doesn't contain this character
func getChar(c rune, f font) []string {
	 l, ok := f.chars[c]
	 if !ok {
		l = f.chars[0]
	 }
	 return l
}

// Given 2 characters, attempts to smush them into 1, according to
// smushmode.  Returns smushed character or '\0' if no smushing can be
// done.

// smushmode values are sum of following (all values smush blanks):
// 1: Smush equal chars (not hardblanks)
// 2: Smush '_' with any char in hierarchy below
// 4: hierarchy: "|", "/\", "[]", "{}", "()", "<>"
//    Each class in hier. can be replaced by later class.
// 8: [ + ] -> |, { + } -> |, ( + ) -> |
// 16: / + \ -> X, > + < -> X (only in that order)
// 32: hardblank + hardblank -> hardblank
func smushem(lch rune, rch rune, mode int, hardblank rune, rtol bool) rune {
	if lch == ' ' { return rch }
	if rch == ' ' { return lch }

	if mode & SMSmush == 0 { // smush not enabled
		return 0
	}

	if mode & SMKern == 0 { // smush but not kern
		// This is smushing by universal overlapping

		// ensure overlapping preference to visible chars (spaces handled already)
		if lch == hardblank { return rch }
		if rch == hardblank { return lch }

		// ensure dominant char overlaps, depending on right-to-left parameter
		if rtol { return lch }
		return rch
	}

	if mode & SMHardBlank == SMHardBlank {
		if lch == hardblank && rch == hardblank { return hardblank }
	}

	if lch == hardblank || rch == hardblank { return 0 }

	if mode & SMEqual == SMEqual {
		if lch == rch { return lch }
	}

	if mode & SMLowLine == SMLowLine {
		if lch == '_' && strings.ContainsRune("|/\\[]{}()<>", rch) { return rch }
		if rch == '_' && strings.ContainsRune("|/\\[]{}()<>", lch) { return lch }
	}

	if mode & SMHierarchy == SMHierarchy {
		hrchy := []string { "|", "/\\", "[]", "{}", "()", "<>" } // low -> high precedence
		inHrchy := func(low rune, high rune, i int) bool {
			return strings.ContainsRune(hrchy[i], low) && strings.ContainsRune(strings.Join(hrchy[i+1:], ""), high)
		}
		for i, _ := range hrchy {
			if inHrchy(lch, rch, i) { return rch }
			if inHrchy(rch, lch, i) { return lch }
		}
	}

	if mode & SMPair == SMPair {
		if lch=='[' && rch==']' { return '|' }
		if rch=='[' && lch==']' { return '|' }
		if lch=='{' && rch=='}' { return '|' }
		if rch=='{' && lch=='}' { return '|' }
		if lch=='(' && rch==')' { return '|' }
		if rch=='(' && lch==')' { return '|' }
	}

	if mode & SMBigX == SMBigX {
		if lch == '/' && rch == '\\' { return '|' }
		if lch == '\\' && rch == '/' { return 'Y' }
		if lch == '>' && rch == '<' { return 'X' }
	}
	return 0
}

// smushamt returns the maximum amount that the character can be smushed
// into the line.
func smushamt(char []string, line []string, smushmode int, hardblank rune, rtol bool) int {
	if (smushmode & (SMSmush | SMKern)) == 0 {
		return 0;
  	}

  	charwidth := len(char[0])
  	charheight := len(char)

  	empty := func (ch rune) bool {
		return ch == 0 || ch == ' '
	}

	maxsmush := charwidth
	for row := 0; row < charheight; row++ {
		var left, right []rune
		if rtol {
			left, right = []rune(char[row]), []rune(line[row])
		} else {
			left, right = []rune(line[row]), []rune(char[row])
		}

		// find first non-empty index in left and right
		var i, j int
		for i = len(left) - 1; i >= 0 && empty(left[i]); i-- { }
		for j = 0; j < len(right) && empty(right[j]); j++ { }

		// the amount of smushing possible just by removing empty spaces
		rowsmush := j + len(left) - i + 1

		// see if we can smush it further
		lch := left[i]
		rch := right[j]
		if !empty(lch) && !empty(rch) {
			if smushem(lch, rch, smushmode, hardblank, rtol) != 0 { rowsmush++ }
		}

		if rowsmush < maxsmush { maxsmush = rowsmush }
	}

	return maxsmush;
}

// Attempts to add the given character onto the end of the given line.
// Returns true if this succeeded, false otherwise.
func addChar(c rune, linep *[]string, maxwidth int, f font, smushmode int, hardblank rune, rtol bool) bool {
	line := *linep
	char := getChar(c, f)
	smushamount := smushamt(char, line, smushmode, hardblank, rtol)

	linelen := len(line[0])
	charheight, charwidth := len(char), len(char[0])

	if linelen + charwidth - smushamount > maxwidth { return false }

	for row := 0; row < charheight; row++ {
		if rtol { panic ("right-to-left not implemented") }
		for k := 0; k < smushamount; k++ {
			column := linelen - smushamount + k
			if column < 0 { column = 0 }

			lch, rch := rune(line[row][column]), rune(char[row][k])
			smushed := smushem(lch, rch, smushmode, hardblank, rtol)
			line[row] = line[row][:column] + string(smushed)
		}
		line[row] += char[row][smushamount:]
	}

	return true
}

// delete this!
type figWord struct {
	art []string
	text string
}

func getWord(w string, f font) []string {
	word := make([]string, f.header.charheight)
	for _, c := range w {
		// todo: addchar func
		char := getChar(c, f)
		for i, charline := range char {
			word[i] += charline
		}
	}

	return word
}

func getWords(msg string, f font) []figWord {
	words := make([]figWord, 0)
	for _, word := range strings.Split(msg, " ") {
		words = append(words, figWord { text: word, art: getWord(word, f) })
	}
	return words
}


func getLines(msg string, f font, width int) [][]string {
	lines := make([][]string, 1) // make room for at least one line
	words := getWords(msg, f)

	// kludge: add first line
	lines[0] = make([]string, f.header.charheight)

	// smoodge everything together for testing
	for _, word := range words {
		for j, wordline := range word.art {
			lines[0][j] += wordline
		}
	}

	return lines
}