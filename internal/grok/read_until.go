package grok

import (
	"bytes"
	"fmt"
)

// Matcher which reads bytes until it matches a target
type readUntil struct {
	targetByte byte
}

func (r readUntil) Match(pkt []byte) ([]byte, int, error) {
	indexOf := bytes.IndexByte(pkt, r.targetByte)
	if indexOf == -1 {
		return nil, 0, fmt.Errorf("Could not find terminating byte %02x", r.targetByte)
	}

	bytes := make([]byte, indexOf+1)
	copy(bytes, pkt[:indexOf+1])
	return bytes, indexOf + 1, nil
}
