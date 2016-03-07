package figletlib

import (
	"fmt"
	//"testing"
)

func outputMode(mode int) {
	modemap := make(map[string]int)
	modemap["SMEqual"] = SMEqual
	modemap["SMLowLine"] = SMLowLine
	modemap["SMHierarchy"] = SMHierarchy
	modemap["SMPair"] = SMPair
	modemap["SMBigX"] = SMBigX
	modemap["SMHardBlank"] = SMHardBlank
	modemap["SMKern"] = SMKern
	modemap["SMSmush"] = SMSmush

	for name := range modemap {
		m := modemap[name]
		if (mode & m) != 0 {
			fmt.Println(name)
		} else {
			fmt.Println("!", name)
		}
	}
}
