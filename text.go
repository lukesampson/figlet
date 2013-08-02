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

