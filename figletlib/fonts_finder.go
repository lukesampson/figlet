
package figletlib

import (
  "errors"
  "fmt"
  "go/build"
  "io/ioutil"
  "os"
  "path/filepath"
  "strings"
)

const (
  pkgName = "github.com/lukesampson/figlet"
)

func FindFontsDirectory() (string, error) {
  // try <bindir>/fonts
  bin := os.Args[0]
  if !filepath.IsAbs(bin) {
    return "", fmt.Errorf("find fonts: bin path %v is not absolute", bin)
  }
  bindir := filepath.Dir(bin)
  fonts := filepath.Join(bindir, "figletlib", "fonts")
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

func FontNames(dir string) []string {
  names := make([]string, 0)

  fis, err := ioutil.ReadDir(dir)
  if err != nil {
    fmt.Println(err); os.Exit(1)
  }

  for _, fi := range(fis) {
    name := fi.Name()
    if strings.HasSuffix(name, ".flf") {
      names = append(names, strings.TrimSuffix(name, ".flf"))
    }
  }

  return names
}

func FindFont(dir string, font string) (string, error) {
  if !strings.HasSuffix(font, ".flf") { font += ".flf" }
  path := filepath.Join(dir, font)
  if _, err := os.Stat(path); err == nil {
    return path, nil
  }
  return "", fmt.Errorf("couldn't find font %v in %v", font, dir)
}

func GetFont(name string) (Font, error) {
  fontsdir, err := FindFontsDirectory()
  if err != nil {
    return Font{}, err
  }

  fontpath, err := FindFont(fontsdir, name)
  if err != nil {
    return Font{}, err
  }

  f, err := ReadFont(fontpath)
  if err != nil {
    return Font{}, err
  }

  return f, nil
}
