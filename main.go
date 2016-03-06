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

func printVersion(fontsdir string) {
	fmt.Println("Figlet version: go-1.0")
	fmt.Printf("Fonts: %v\n", fontsdir)
}

func listFonts(fontsdir string) {
	fmt.Printf("Fonts in %v:\n", fontsdir)
	fonts, _ := figletlib.FontNamesInDir(fontsdir)
	for _, fontname := range fonts {
		fmt.Printf("%v:\n", fontname)
		f, err := figletlib.GetFontByName(fontsdir, fontname)
		if err != nil {
			fmt.Println(err)
		}

		s := f.Settings()
		figletlib.PrintMsg(fontname, f, 80, s, "left")
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
	fontsDirectory := flag.String("d", "", "fonts directory")
	flag.Parse()

	fontsdir := *fontsDirectory
	if fontsdir == "" {
		fontsdir = figletlib.GuessFontsDirectory()
		if fontsdir == "" {
			fmt.Println("ERROR: couldn't find fonts directory, specify -d")
			os.Exit(1)
		}
	}

	if *list {
		listFonts(fontsdir); os.Exit(0)
	}

	if *help {
		printHelp(); os.Exit(0)
	}

	if *version {
		printVersion(fontsdir); os.Exit(0)
	}

	var align string
	if *alignRight {
		align = "right"
	} else if *alignCenter {
		align = "center"
	}

	f, err := figletlib.GetFontByName(fontsdir, *fontname)
	if err != nil {
		fmt.Println("ERROR: couldn't find font", *fontname, "in dir", fontsdir)
		os.Exit(1)
	}

	msg := strings.Join(flag.Args(), " ")

	s := f.Settings()
	if *reverse {
		s.SetRtoL(true)
	}

	maxwidth := *outputWidth
	if msg == "" {
		reader := bufio.NewReader(os.Stdin)
		for {
			msg, err = reader.ReadString('\n')
			if err != nil {
				if err == io.EOF { os.Exit(0) }
				msg = ""
			}
			figletlib.PrintMsg(msg, f, maxwidth, s, align)
		}
	}
	figletlib.PrintMsg(msg, f, maxwidth, s, align)
}