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

func (ft *figText) String() string {
	str := ""
	for _, line := range ft.art {
		str += string(line) + "\n"
	}
	return str
}

func (ft *figText) copy() *figText {
	copied := newFigText(ft.height())

	(*copied).text = (*ft).text
	for i := 0; i < ft.height(); i++ {
		width := ft.width()
		(*copied).art[i] = make([]rune, width)
		for j := 0; j < width; j++ {
			(*copied).art[i][j] = (*ft).art[i][j]
		}
	}

	return copied
}

