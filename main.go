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

func printLines(lines [][][]rune) {
	for _, line := range lines {
		for _, subline := range line {
			for _, outchar := range subline {
				fmt.Printf("%c", outchar)
			}
			fmt.Println()
		}
	}
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

	msg := strings.Join(os.Args[1:], " ")

	printLines(getLines(msg, f, 80))

}