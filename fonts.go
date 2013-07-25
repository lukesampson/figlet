package main

import (
	"fmt"
	"go/build"
	"path/filepath"
	"os"
	"errors"
	"io/ioutil"
	//"regexp"
	"strings"
	"strconv"
)

// smush modes
const (
	SMEqual = 1
	SMLowLine = 2
	SMHierarchy = 4
	SMPair = 8
	SMBigX = 16
	SMHardBlank = 32
	SMKern = 64
	SMSmush = 128
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
	if !strings.HasSuffix(font, ".flf") { font += ".flf" }
	path := filepath.Join(dir, font)
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}
	return "", fmt.Errorf("couldn't find font %v in %v", font, dir)
}

type fontHeader struct {
	hardBlank string
	height int32
	baseLine int32
	maxLength int32
	smush int32
	commentLines int32
	printDirection int32
	smush2 int32
}

func parseHeader(header string) (fontHeader, error) {
	h := fontHeader {}

	magic_num := "flf2a"
	if !strings.HasPrefix(header, magic_num) {
		return h, fmt.Errorf("invalid font header: %v", header)
	}

	headerParts := strings.Split(header[len(magic_num):], " ")
	h.hardBlank = headerParts[0]

	nums := make([]int32, len(headerParts)-1)
	for i, s := range headerParts[1:] {
		num, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return h, fmt.Errorf("invalid font header: %v: %v", header, err)
		}
		nums[i] = int32(num)
	}

	h.height = nums[0]
	h.baseLine = nums[1]
	h.maxLength = nums[2]
	h.smush = nums[3]
	h.commentLines = nums[4]

	// these are optional for backwards compatibility
	if len(nums) > 5 { h.printDirection = nums[5] }
	if len(nums) > 6 { h.smush2 = nums[6] }

	// if no smush2, decode smush into smush2
	if len(nums) < 7 {
		if h.smush == 0 {
			h.smush2 = SMKern
		} else if h.smush < 0 {
			h.smush2 = 0
		} else {
			h.smush2 = (h.smush & 31) | SMSmush
		}
	}

	return h, nil
}

type font struct {
	header fontHeader
	comment string
	chars [][]string
}

func loadFont(file string) (font, error) {
	f := font {}

	bytes, err := ioutil.ReadFile(file)
	if err != nil { return f, err }

	lines := strings.Split(string(bytes), "\n")

	f.header, err = parseHeader(lines[0])
	if err != nil { return f, err }

	f.comment = strings.Join(lines[1:f.header.commentLines+1], "\n")

	//height := f.header.height
	for i := 32; i < 128; i++ {
		//charlines := make([]string, height)
	}

	fmt.Println(f)
	
	return f, nil
}