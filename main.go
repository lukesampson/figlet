package main

import (
	"fmt"
	"os"
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
	
	_, err = readFont(fontpath)
	if err != nil {
		fmt.Println(err); os.Exit(1)
	}

	message := os.Args[1]

	for _, c:= range message {
		fmt.Println(string(c))
	}

}