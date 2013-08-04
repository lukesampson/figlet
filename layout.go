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
func smushem(lch rune, rch rune, s settings) rune {
	if lch == ' ' { return rch }
	if rch == ' ' { return lch }

	//fmt.Printf("smushem %q + %q\n", lch, rch)

	if s.smushmode & SMSmush == 0 { // smush not enabled
		return 0
	}

	if s.smushmode & SMKern == 0 { // smush but not kern
		// This is smushing by universal overlapping

		// ensure overlapping preference to visible chars (spaces handled already)
		if lch == s.hardblank { return rch }
		if rch == s.hardblank { return lch }

		// ensure dominant char overlaps, depending on right-to-left parameter
		if s.rtol { return lch }
		return rch
	}

	if s.smushmode & SMHardBlank == SMHardBlank {
		if lch == s.hardblank && rch == s.hardblank { return s.hardblank }
	}

	if lch == s.hardblank || rch == s.hardblank { return 0 }

	if s.smushmode & SMEqual == SMEqual {
		if lch == rch { return lch }
	}

	if s.smushmode & SMLowLine == SMLowLine {
		if lch == '_' && strings.ContainsRune("|/\\[]{}()<>", rch) { return rch }
		if rch == '_' && strings.ContainsRune("|/\\[]{}()<>", lch) { return lch }
	}

	if s.smushmode & SMHierarchy == SMHierarchy {
		hrchy := []string { "|", "/\\", "[]", "{}", "()", "<>" } // low -> high precedence
		inHrchy := func(low rune, high rune, i int) bool {
			return strings.ContainsRune(hrchy[i], low) && strings.ContainsRune(strings.Join(hrchy[i+1:], ""), high)
		}
		for i, _ := range hrchy {
			if inHrchy(lch, rch, i) { return rch }
			if inHrchy(rch, lch, i) { return lch }
		}
	}

	if s.smushmode & SMPair == SMPair {
		if lch=='[' && rch==']' { return '|' }
		if rch=='[' && lch==']' { return '|' }
		if lch=='{' && rch=='}' { return '|' }
		if rch=='{' && lch=='}' { return '|' }
		if lch=='(' && rch==')' { return '|' }
		if rch=='(' && lch==')' { return '|' }
	}

	if s.smushmode & SMBigX == SMBigX {
		if lch == '/' && rch == '\\' { return '|' }
		if lch == '\\' && rch == '/' { return 'Y' }
		if lch == '>' && rch == '<' { return 'X' }
	}
	return 0
}

// smushamt returns the maximum amount that the character can be smushed
// into the line.
func smushamt(char *figText, line *figText, s settings) int {
	//fmt.Println(line)
	//fmt.Println(char)
	if s.smushmode & (SMSmush | SMKern) == 0 {
		return 0;
  	}

  	empty := func (ch rune) bool {
		return ch == 0 || ch == ' '
	}

	maxsmush := char.width()
	for row := 0; row < char.height(); row++ {
		var left, right []rune
		if s.rtol {
			left, right = (*char).art[row], (*line).art[row]
		} else {
			left, right = (*line).art[row], (*char).art[row]
		}

		// find number of empty chars for left and right
		var i, j int
		for i = 0; i < len(left) && empty(left[len(left) - 1 - i]); i++ { }
		for j = 0; j < len(right) && empty(right[j]); j++ { }

		// the amount of smushing possible just by removing empty spaces
		rowsmush := j + i
		//fmt.Printf("j: %v, i: %v, rowsmush: %v\n", j, i, rowsmush)

		if i < len(left) && j < len(right) {
			// see if we can smush it even further
			lch := left[len(left) - 1 - i]
			rch := right[j]
			if !empty(lch) && !empty(rch) {
				if smushem(lch, rch, s) != 0 {
					//fmt.Printf("managed to smush %q + %q (%v)\n", lch, rch, s.smushmode)
					rowsmush++
				}
			}
		}

		if rowsmush < maxsmush { maxsmush = rowsmush }
	}

	//fmt.Printf("maxsmush: %v\n", maxsmush)
	return maxsmush;
}

type settings struct {
	smushmode int
	hardblank rune
	rtol bool
}

// Adds the given character onto the end of the given line.
func addChar(char *figText, line *figText, s settings) figText {
	smushamount := smushamt(char, line, s)
	return smushChar(char, line, smushamount, s)
}

func smushChar(char *figText, line *figText, amount int, s settings) figText {
	var result *figText
	if s.rtol {
		result = char.copy()
	} else {
		result = line.copy()
	}

	linelen := result.width()

	//fmt.Println("smushamount", amount)

	for row := 0; row < char.height(); row++ {
		left, right := &(*result).art[row], &(*char).art[row]
		if s.rtol {
			right = &(*line).art[row]
		}

		for k := 0; k < amount; k++ {
			column := linelen - amount + k
			if column < 0 { column = 0 }

			rch := (*right)[k]
			
			if column >= len(*left) {
				if rch == ' ' { continue } // absorb space
				*left = append(*left, ' ')
			}

			lch := (*left)[column]
			smushed := smushem(lch, rch, s)			

			//fmt.Printf("row %v, col %v, lch %q, rch %q, smushed %q\n", row, column, lch, rch, smushed)
			(*left)[column] = smushed

		}
		*left = append(*left, (*right)[amount:]...)
	}

	return *result
}

// gets the font entry for the given character, or the 'missing'
// character if the font doesn't contain this character
func getChar(c rune, f font) *figText {
	 l, ok := f.chars[c]
	 if !ok {
		l = f.chars[0]
	 }
	 return &figText { text: string(c), art: l }
}

func getWord(w string, f font, s settings) *figText {
	word := newFigText(f.header.charheight)
	for _, c := range w {
		c := getChar(c, f)
		*word = addChar(c, word, s)
	}
	(*word).text = w
	return word
}

func getWords(msg string, f font, s settings) []figText {
	words := make([]figText, 0)
	for _, word := range strings.Split(msg, " ") {
		words = append(words, *getWord(word, f, s))
	}
	return words
}

func breakWord(word *figText, maxwidth int, f font, s settings) (*figText, *figText) {
	h := word.height()
	a, b := word, newFigText(h)

	text := (*word).text

	for i := len(text) - 1; i > 0 && a.width() > maxwidth; i-- {
		a = getWord(text[:i], f, s)
		b = getWord(text[i:], f, s)
	}

	return a, b
}

func getLines(msg string, f font, maxwidth int, s settings) []figText {
	lines := make([]figText, 1)
	words := getWords(msg, f, s)

	// make empty first line
	lines[0] = *newFigText(f.header.charheight)

	linenum := 0
	for i, word := range words {
		// add a space between words
		if i > 0 {
			// don't smush space
			lineWithSpace := smushChar(getChar(' ', f), &lines[linenum], 0, s)
			//fmt.Println(&lineWithSpace)
			if lineWithSpace.width() <= maxwidth {
				lines[linenum] = lineWithSpace
			}
		}

		// check if we need to wrap
		if lines[linenum].width() + word.width() > maxwidth { // need to wrap
			lines = append(lines, figText { art: make([][]rune, f.header.charheight) })

			if word.width() > maxwidth {
				a, b := breakWord(&word, maxwidth - lines[linenum].width(), f, s)

				lines[linenum] = addChar(a, &lines[linenum], s)
				word = *b
			}

			linenum++
		}

		lines[linenum] = addChar(&word, &lines[linenum], s)
	}

	return lines
}