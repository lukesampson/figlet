package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/lukesampson/figlet/figletlib"
)

const (
	defaultFont = "standard"
)

func printUsage() {
	fmt.Println("Usage: figlet [ -lcrhR ] [ -f fontfile ] [ -I infocode ]")
	fmt.Println("              [ -w outputwidth ] [ -m smushmode ]")
	fmt.Println("              [ message ]")
}

func printHelp() {
	printUsage()
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("For more info see https://github.com/lukesampson/figlet")
}

func printVersion(fontsdir string) {
	fmt.Println("Figlet version: go-1.0")
	fmt.Printf("Fonts: %v\n", fontsdir)
}

func printInfoCode(infocode int, infodata []string) {
	fmt.Println(infodata[infocode])
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
	infoCode := flag.Int("I", -1, "infocode")
	infoCode2 := flag.Bool("I2", false, "show default font directory")
	infoCode3 := flag.Bool("I3", false, "show default font")
	infoCode4 := flag.Bool("I4", false, "show output width")
	infoCode5 := flag.Bool("I5", false, "show supported font formats")
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
		listFonts(fontsdir)
		os.Exit(0)
	}

	if *help {
		printHelp()
		os.Exit(0)
	}

	if *version {
		printVersion(fontsdir)
		os.Exit(0)
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

	ic := *infoCode

	if *infoCode2 {
		ic = 2
	} else if *infoCode3 {
		ic = 3
	} else if *infoCode4 {
		ic = 4
	} else if *infoCode5 {
		ic = 5
	}

	if ic > 1 && ic < 6 {
		outputWidthString := strconv.Itoa(*outputWidth)
		var infoData = []string{2: fontsdir, 3: *fontname, 4: outputWidthString, 5: "flf2"}
		printInfoCode(ic, infoData)
		os.Exit(0)
	} else if ic != -1 {
		fmt.Println("ERROR: invalid infocode", ic)
		os.Exit(1)
	}

	maxwidth := *outputWidth
	if msg == "" {
		reader := bufio.NewReader(os.Stdin)
		for {
			msg, err = reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					os.Exit(0)
				}
				msg = ""
			}
			figletlib.PrintMsg(msg, f, maxwidth, s, align)
		}
	}
	figletlib.PrintMsg(msg, f, maxwidth, s, align)
}
