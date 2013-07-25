package main

import (
	"fmt"
	"go/build"
	"path/filepath"
	"os"
	"errors"
	"io/ioutil"
	"regexp"
	"strings"
	"strconv"
)

const (
	basePkg = "github.com/lukesampson/figlet"
	defaultFont = "standard"
)

func findFonts() (string, error) {
	// try <bindir>/fonts
	bin := os.Args[0]
	if !filepath.IsAbs(bin) {
		return "", fmt.Errorf("find fonts: bin path %v is not absolute", bin)
	}
	bindir := filepath.Dir(bin)
	fonts := filepath.Join(bindir, "fonts")
	if _, err := os.Stat(fonts); err == nil {
		return fonts, nil
	}

	// try src directory
	ctx := build.Default
	if p, err := ctx.Import(basePkg, "", build.FindOnly); err == nil {
		fonts := filepath.Join(p.Dir, "fonts")
		if _, err := os.Stat(fonts); err == nil {
			return fonts, nil
		}
	}

	return "", errors.New("couldn't find fonts directory")
}

func findFont(dir string, font string) (string, error) {
	for _, ext := range([]string { ".flf", ".tlf" }) {
		path := filepath.Join(dir, fmt.Sprintf("%v%v", font, ext))
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("couldn't find font %v in %v", font, dir)
}

func loadFont(file string) (string, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil { return "", err }

	lines := strings.Split(string(bytes), "\n")

	header := lines[0]
	reMagicNum := regexp.MustCompile(`^[ft]lf2.`)

	if !reMagicNum.MatchString(header) {
		return "", fmt.Errorf("%v isn't a valid figlet font", file)
	}

	headerParts := strings.Split(reMagicNum.ReplaceAllString(header, ""), " ")
	//hardBlank := headerParts[0]

	headerNums := make([]int32, len(headerParts)-1)
	for i, s := range headerParts[1:] {
		num, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return "", fmt.Errorf("invalid font header: %v: %v", header, err)
		}
		headerNums[i] = int32(num)
	}

	fmt.Printf("header nums: %v (%v)\n", headerNums, len(headerNums))

	return "", nil
}

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
	
	_, err = loadFont(fontpath)
	if err != nil {
		fmt.Println(err); os.Exit(1)
	}
}