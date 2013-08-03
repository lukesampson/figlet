package main

import (
	"testing"
	"fmt"
)

func TestCopy(t *testing.T) {
	a := &figText {
		text: "test",
		art: [][]rune {
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
}