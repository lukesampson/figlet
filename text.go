package main

type figText struct {
	art [][]rune
	text string
}

func newFigText(height int) *figText {
	return &figText { art: make([][]rune, height) }
}

func (ft *figText) width() int {
	return len(ft.art[0])
}

func (ft *figText) height() int {
	return len(ft.art)
}

func (ft *figText) splitAt(index int) (*figText, *figText) {
	if !(index < ft.width()) {
		panic("split index is out of range")
	}

	h := ft.height()
	a, b := newFigText(h), newFigText(h)

	for i := 0; i < h; i++ {
		(*a).art[i] = (*ft).art[i][:index]
		(*b).art[i] = (*ft).art[i][index:]
	}

	return a, b
}

func (ft *figText) String() string {
	str := ""
	for _, line := range ft.art {
		str += string(line) + "\n"
	}
	return str
}

