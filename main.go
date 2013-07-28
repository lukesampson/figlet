package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	pkgName = "github.com/lukesampson/figlet"
	defaultFont = "standard"
)

// Globals affected by command line options
var (
	deutschflag bool
	justification int
	paragraphflag int
	right2left bool
	multibyte int
	cmdinput bool
	smushmode int
)

func printusage() {
	fmt.Println("Usage: figlet [ -cklnoprstvxDELNRSWX ] [ -d fontdirectory ]")
	fmt.Println("              [ -f fontfile ] [ -m smushmode ] [ -w outputwidth ]")
	fmt.Println("              [ -C controlfile ] [ -I infocode ] [ message ]")
}

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

func main() {
	fontsdir, err := findFonts()
	if err != nil {
		fmt.Println(err); os.Exit(1)
	}

	fontname := defaultFont
	fontpath, err := findFont(fontsdir, fontname)
	if err != nil {
		fmt.Println(err); os.Exit(1)
	}
	
	f, err := readFont(fontpath)
	if err != nil {
		fmt.Println(err); os.Exit(1)
	}

	msg := os.Args[1]

	for _, c:= range msg {
		fmt.Println(string(c))
		l := getChar(c, f)

		for _, line := range l {
			fmt.Println(line)
		}	
	}

}