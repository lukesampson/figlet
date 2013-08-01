package main

import (
	"errors"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Ä Ö Ü ä ö ü ß
var deutsch = []rune { 196, 214, 220, 228, 246, 252, 223 };

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
	hardblank string
	charheight int
	baseLine int
	maxlen int
	smush int
	commentLines int
	printDirection int
	smush2 int
}

func readHeader(header string) (fontHeader, error) {
	h := fontHeader {}

	magic_num := "flf2a"
	if !strings.HasPrefix(header, magic_num) {
		return h, fmt.Errorf("invalid font header: %v", header)
	}

	headerParts := strings.Split(header[len(magic_num):], " ")
	h.hardblank = headerParts[0]

	nums := make([]int, len(headerParts)-1)
	for i, s := range headerParts[1:] {
		num, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return h, fmt.Errorf("invalid font header: %v: %v", header, err)
		}
		nums[i] = int(num)
	}

	h.charheight = nums[0]
	h.baseLine = nums[1]
	h.maxlen = nums[2]
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

func readFontChar(lines []string, currline int, height int) [][]rune {
	char := make([][]rune, height)
	for row := 0; row < height; row++ {
		line := lines[currline+row]
		
		k := len(line) - 1

		// remove any trailing whitespace after end char
		ws := regexp.MustCompile(`\s`)
		for k > 0 && ws.MatchString(line[k:k+1]) { k-- }

		if k > 0 {
			// remove end marks
			endchar := line[k]
			for k > 0 && line[k] == endchar { k-- }
		}

		char[row] = []rune(line[:k+1])
	}

	return char
}

type font struct {
	header fontHeader
	comment string
	chars map[rune] [][]rune
}

func readFont(file string) (font, error) {
	f := font {}

	bytes, err := ioutil.ReadFile(file)
	if err != nil { return f, err }

	lines := strings.Split(string(bytes), "\n")

	f.header, err = readHeader(lines[0])
	if err != nil { return f, err }

	f.comment = strings.Join(lines[1:f.header.commentLines+1], "\n")

	f.chars = make(map[rune] [][]rune)
	charheight := int(f.header.charheight)
	currline := int(f.header.commentLines)+1
	
	// allocate 0, the 'missing' character
	f.chars[0] = make([][]rune, charheight)

	// standard ASCII characters
	for ord := ' '; ord <= '~'; ord++ {
		f.chars[ord] = readFontChar(lines, currline, charheight)
		currline += charheight
	}

	// 7 german characters
	for i := 0; i < 7; i++ {
		f.chars[deutsch[i]] = readFontChar(lines, currline, charheight)
		currline += charheight
	}

	// code-tagged characters
	for currline < len(lines) {
		var code int;
		_, err := fmt.Sscan(lines[currline], &code)
		if err != nil { break }
		currline++
		f.chars[rune(code)] = readFontChar(lines, currline, charheight)

		currline += charheight
	}

	
	return f, nil
}

func getFont(name string) (font, error) {
	fontsdir, err := findFonts()
	if err != nil {
		return font { }, err
	}

	fontname := defaultFont
	fontpath, err := findFont(fontsdir, fontname)
	if err != nil {
		return font { }, err
	}
	
	f, err := readFont(fontpath)
	if err != nil {
		return font { }, err
	}

	return f, nil
}