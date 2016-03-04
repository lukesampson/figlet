package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"flag"

	"./figletlib"
)

const (
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

func printVersion() {
	fmt.Println("Figlet version: go-1.0")
	dir, err := figletlib.FindFonts()
	if err != nil {
		dir = fmt.Sprintf("ERROR: couldn't find fonts: %v", err)
	}
	fmt.Printf("Fonts: %v\n", dir)
}

func printLines(lines []figletlib.FigText, hardblank rune, maxwidth int, align string) {
	padleft := func(linelen int) {
		switch align {
		case "right":
			fmt.Print(strings.Repeat(" ", maxwidth - linelen))
		case "center":
			fmt.Print(strings.Repeat(" ", (maxwidth - linelen) / 2))
		}
	}

	for _, line := range lines {
		for _, subline := range line.Art() {
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

func printMsg(msg string, f figletlib.Font, maxwidth int, s figletlib.Settings, align string) {
	lines := figletlib.GetLines(msg, f, maxwidth, s)
	printLines(lines, s.HardBlank(), maxwidth, align)
}

func listFonts() {
	fontsdir, err := figletlib.FindFonts()
	if err != nil {
		fmt.Println(err); os.Exit(1)
	}

	fmt.Printf("Fonts in %v:\n", fontsdir)

	for _, fontname := range figletlib.FontNames(fontsdir) {
		fmt.Printf("%v:\n", fontname)
		fpath, _ := figletlib.FindFont(fontsdir, fontname)
		f, err := figletlib.ReadFont(fpath)
		if err != nil {
			fmt.Println(err)
		}
		s := (&f).Settings()

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
	version := flag.Bool("v", false, "show version info")
	flag.Parse()

	if *list {
		listFonts(); os.Exit(0)
	}

	if *help {
		printHelp(); os.Exit(0)
	}

	if *version {
		printVersion(); os.Exit(0)
	}

	var align string
	if *alignRight {
		align = "right"
	} else if *alignCenter {
		align = "center"
	}

	f, err := figletlib.GetFont(*fontname)
	if err != nil {
		fmt.Println(err); os.Exit(1)
	}

	msg := strings.Join(flag.Args(), " ")

	s := f.Settings()
	if *reverse {
		s.SetRtoL(true)
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