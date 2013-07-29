package main

import (
	"strings"
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

func smushem(lch rune, rch rune, smushmode int) (rune, bool) {
	if lch == ' ' { return rch, true }
	if rch == ' ' { return lch, true }

	if smushmode & SMSmush == 0 { return 0, false }

	if smushmode & 63 == 0 {
		// This is smushing by universal overlapping.
		
	}
	return 0, false
}

// returns true if the word could be added to the line
func addWord(f font, word string, line []string) bool {
	return false
}

// Gets the next line that will fit in allowed width
func nextLine(f font, msg string, width int) ([]string, string) {
	line := make([]string, f.header.charheight)
	words := strings.Split(msg, " ")
	for i, word := range words {
		if addWord(f, word, line) {
			msg = strings.Join(words[i+1:], " ")
		} else if i == 0 { // word longer than line
			panic("forced word break not implemented")
		} else {
			break 
		}
	}
	return line, msg
}


func getLines(f font, msg string, width int) [][]string {
	lines := make([][]string, 0, 1) // make room for at least one line
	for len(msg) > 0 {
		var line []string
		line, msg = nextLine(f, msg, width)
		lines = append(lines, line)
	}
	return lines
}