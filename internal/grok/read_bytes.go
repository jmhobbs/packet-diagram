package grok

// Matcher which reads exactly a number of bytes
type readBytes struct {
	length int
}

func (r readBytes) Match(pkt []byte) ([]byte, int, error) {
	bytes := make([]byte, r.length)
	copy(bytes, pkt[:r.length])
	return bytes, r.length, nil
}
