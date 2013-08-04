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

func printUsage() {
	fmt.Println("Usage: figlet [ -lcrhR ] [ -f fontfile ]")
	fmt.Println("              [ -w outputwidth ] [ -m smushmode ]")
	fmt.Println("              [ message ]")
	fmt.Println()
	fmt.Println("Show available fonts:")
	fmt.Println("       figlet -list")
}

func printHelp() {
	printUsage()
	fmt.Println()
	fmt.Println("For more info see https://github.com/lukesampson/figlet")
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

func listFonts() {
	fontsdir, err := findFonts()
	if err != nil {
		fmt.Println(err); os.Exit(1)
	}

	fmt.Printf("Fonts in %v:\n", fontsdir)

	for _, fontname := range fontNames(fontsdir) {
		fmt.Printf("%v:\n", fontname)
		fpath, _ := findFont(fontsdir, fontname)
		f, err := readFont(fpath)
		if err != nil {
			fmt.Println(err)
		}
		s := (&f).settings()

		printMsg(fontname, f, 80, s, "left")
		fmt.Println()
	}
}

func main() {
	// options
	fontname := flag.String("f", defaultFont, "name of font to use")
	reverse := flag.Bool("R", false, "reverse output")
	alignRight := flag.Bool("r", false, "right-align output")
	alignCenter := flag.Bool("c", false, "center-align output")
	outputWidth := flag.Int("w", 80, "output width")
	list := flag.Bool("list", false, "list available fonts")
	help := flag.Bool("h", false, "show help")
	flag.Parse()

	if *list {
		listFonts(); os.Exit(0)
	}

	if *help {
		printHelp(); os.Exit(0)
	}

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

	s := f.settings()
	if *reverse {
		s.rtol = !s.rtol
	}

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