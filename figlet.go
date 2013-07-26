package main

import (
	"fmt"
	"os"
)

const (
	pkgName = "github.com/lukesampson/figlet"
	defaultFont = "standard"
)

func main() {
	//fmt.Println(os.Args[0])

	fontsdir, err := findFonts()
	if err != nil {
		fmt.Println(err); os.Exit(1)
	}

	fontname := defaultFont
	fontpath, err := findFont(fontsdir, fontname)
	if err != nil {
		fmt.Println(err); os.Exit(1)
	}
	
	font, err := readFont(fontpath)
	if err != nil {
		fmt.Println(err); os.Exit(1)
	}

	//fmt.Println(font.chars['a'])
	fmt.Println(font)

}