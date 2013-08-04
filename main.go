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
	fmt.Println("Usage: figlet [ -cklnoprstvxDELNRSWX ] [ -d fontdirectory ]")
	fmt.Println("              [ -f fontfile ] [ -m smushmode ] [ -w outputwidth ]")
	fmt.Println("              [ -I infocode ] [ message ]")
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

func main() {
	fontname := flag.String("f", defaultFont, "use this font")
	rtol := flag.Bool("R", false, "reverse output")
	alignRight := flag.Bool("r", false, "right-align output")
	alignCenter := flag.Bool("c", false, "center-align output")

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
		rtol: *rtol }

	maxwidth := 80

	if(msg == "") {
		reader := bufio.NewReader(os.Stdin)
		for {
			msg, err = reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					os.Exit(0)
				}
				msg = ""
			}
			lines := getLines(msg, f, maxwidth, s)
			printLines(lines, s.hardblank, maxwidth, align)
		}
	}

	lines := getLines(msg, f, maxwidth, s)
	printLines(lines, s.hardblank, maxwidth, align)

}