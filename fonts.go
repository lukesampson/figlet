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
	if p, err := ctx.Import(pkgName, "", build.FindOnly); err == nil {
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

type fontHeader struct {
	hardBlank string
}

func parseHeader(header string) (fontHeader, error) {
	h := fontHeader {}

	reMagicNum := regexp.MustCompile(`^[ft]lf2.`)

	if !reMagicNum.MatchString(header) {
		return h, fmt.Errorf("invalid font header: %v", header)
	}

	headerParts := strings.Split(reMagicNum.ReplaceAllString(header, ""), " ")
	h.hardBlank = headerParts[0]

	headerNums := make([]int32, len(headerParts)-1)
	for i, s := range headerParts[1:] {
		num, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return h, fmt.Errorf("invalid font header: %v: %v", header, err)
		}
		headerNums[i] = int32(num)
	}

	return h, nil
}

func loadFont(file string) (string, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil { return "", err }

	lines := strings.Split(string(bytes), "\n")

	header, err := parseHeader(lines[0])
	if err != nil { return "", err }

	fmt.Printf("%v\n", header)

	return "", nil
}