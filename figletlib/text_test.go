package figletlib

import (
	"fmt"
	"testing"
)

func TestCopy(t *testing.T) {
	a := &FigText{
		text: "test",
		art: [][]rune{
			[]rune("ABC"),
			[]rune("DEF"),
		},
	}

	b := a.copy()

	astr, bstr := fmt.Sprint(a), fmt.Sprint(b)
	if bstr != astr {
		t.Errorf("copy:\n%v\nwas not equal to original:\n%v\n", bstr, astr)
	}

	if x := (*b).text; x != (*a).text {
		t.Errorf("copied text was %v, want %v", x, (*a).text)
	}

	(*a).text = "changed"
	if x := (*b).text; x != "test" {
		t.Errorf("text for copy was updated to %v\n", x)
	}

	(*a).art[1][1] = 'X'

	if fmt.Sprint(b) != bstr {
		t.Errorf("b.art changed when a.art changed")
	}
}
