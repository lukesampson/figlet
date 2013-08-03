package main

import (
	"fmt"
	"os"
	"strings"
	"flag"
)

const (
	pkgName = "github.com/lukesampson/figlet"
	defaultFont = "standard"
)

func printusage() {
	fmt.Println("Usage: figlet [ -cklnoprstvxDELNRSWX ] [ -d fontdirectory ]")
	fmt.Println("              [ -f fontfile ] [ -m smushmode ] [ -w outputwidth ]")
	fmt.Println("              [ -C controlfile ] [ -I infocode ] [ message ]")
}

func printLines(lines []figText, hardblank rune) {
	for _, line := range lines {
		for _, subline := range line.art {
			for _, outchar := range subline {
				if outchar == hardblank { outchar = ' '}
				fmt.Printf("%c", outchar)
			}
			fmt.Println()
		}
	}
}

func main() {
	fontname := flag.String("f", defaultFont, "name of the font")
	flag.Parse()

	f, err := getFont(*fontname)
	if err != nil {
		fmt.Println(err); os.Exit(1)
	}

	msg := strings.Join(flag.Args(), " ")

	s := settings {
		smushmode: SMKern + SMSmush + SMEqual + SMLowLine + SMHierarchy + SMPair,
		hardblank: '$',
		rtol: false }

	lines := getLines(msg, f, 80, s)
	printLines(lines, s.hardblank)

}