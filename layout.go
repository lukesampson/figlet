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
func smushamt(char []string, line []string, smushmode int, rtol bool) int {
	return 0
}

func getWords(f font, msg string) [][]string {
	strings.Split(msg, " ")
	return nil
}


func getLines(f font, msg string, width int) [][]string {
	lines := make([][]string, 0, 1) // make room for at least one line
	words := getWords(f, msg)

	// smoodge everything together for testing
	for _, word := range words {
		for j, wordline := range word {
			lines[j] = append(lines[j], wordline)
		}
	}

	return lines
}