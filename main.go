package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"flag"
)

const (
	pkgName = "github.com/lukesampson/figlet"
	defaultFont = "standard"
)

func printusage() {
	fmt.Println("Usage: figlet [ -lcrhvR ] [ -f fontfile ]")
	fmt.Println("              [ -w outputwidth ] [ -m smushmode ]")
	fmt.Println("              [ message ]")
}

func printLines(lines []figText, hardblank rune, maxwidth int, align string) {
	padleft := func(linelen int) {
		switch align {
		case "right":
			fmt.Print(strings.Repeat(" ", maxwidth - linelen))
		case "center":
			fmt.Print(strings.Repeat(" ", (maxwidth - linelen) / 2))
		}
	}

	for _, line := range lines {
		for _, subline := range line.art {
			padleft(len(subline))
			for _, outchar := range subline {
				if outchar == hardblank { outchar = ' '}
				fmt.Printf("%c", outchar)
			}
			if len(subline) < maxwidth && align != "right" {
				fmt.Println()
			}
		}
	}
}

func printMsg(msg string, f font, maxwidth int, s settings, align string) {
	lines := getLines(msg, f, maxwidth, s)
	printLines(lines, s.hardblank, maxwidth, align)
}

func main() {
	// options
	fontname := flag.String("f", defaultFont, "name of font to use")
	reverse := flag.Bool("R", false, "reverse output")
	alignRight := flag.Bool("r", false, "right-align output")
	alignCenter := flag.Bool("c", false, "center-align output")
	outputWidth := flag.Int("w", 80, "output width")
	flag.Parse()

	var align string
	if *alignRight {
		align = "right"
	} else if *alignCenter {
		align = "center"
	}

	f, err := getFont(*fontname)
	if err != nil {
		fmt.Println(err); os.Exit(1)
	}

	msg := strings.Join(flag.Args(), " ")

	s := settings {
		smushmode: f.header.smush2,
		hardblank: '$',
		rtol: *reverse }

	maxwidth := *outputWidth

	if(msg == "") {
		reader := bufio.NewReader(os.Stdin)
		for {
			msg, err = reader.ReadString('\n')
			if err != nil {
				if err == io.EOF { os.Exit(0) }
				msg = ""
			}
			printMsg(msg, f, maxwidth, s, align)
		}
	}
	printMsg(msg, f, maxwidth, s, align)
}