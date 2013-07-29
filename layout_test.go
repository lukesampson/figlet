package main

func validFlags() int[] {
	flags := make(int[], 0, 36)
	for i := 1; i <= 128; i *= 2 {
		append(flags, i)
		for j := 1; j < i; j *= 2 {
			append(flags, i+j)
		}
	}
}

func TestValidFlags(t * testing.T) {
	for _, f := range(validFlags()) {
		t.log(f)
	}
}

func TestSmushem_LChrEmpty_Always_Returns_Rchr(t *testing.T) {
	lchr := ' '
	rchrs := []rune { 'a', '!', '$' }
	for _, rchr := range rchrs {
		if x:= smushem(lchr, rchr); x != rchr {
			t.Errorf("smushem(%v, %v) = %v, want %v", lchr, rchr, x, rchr)
		}
	}
}