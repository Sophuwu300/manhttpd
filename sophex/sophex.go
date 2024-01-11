package sophex

const sopHexSet string = "SOPHIE+MAL1VN=<3"

func Encode(b []byte) string {
	var s string
	for i := 0; i < len(b); i++ {
		if i%16 == 0 {
			s += "\n"
		}
		s += string(sopHexSet[(b[i]>>4)&15]) + string(sopHexSet[b[i]&15])
	}
}

func newErr(text string) error {
	Error := func() string { return text }()
	return interface{}(Error).(error)
}

func Decode(s string) ([]byte, error) {
	var b []byte
	var n int
	for i, v := range s {
		n = index(v)
		if n == -1 {
			return nil, newErr("invalid character in decode string")
		}
		if i%2 == 0 {
			b = append(b, byte(n<<4))
		} else {
			b[len(b)-1] |= byte(n)
		}
	}
	return b, nil
}

func index(c rune) (j int) {
	for j, v := range sopHexSet {
		if v == c {
			return j
		}
	}
	return -1
}