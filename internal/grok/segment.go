package grok

// A segment is a parsed out chunk of a packet.
type Segment struct {
	Name   string
	Offset int
	Length int
	Bytes  []byte
	Skip   bool
}
